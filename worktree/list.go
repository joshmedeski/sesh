package worktree

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"time"

	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/joshmedeski/sesh/v2/model"
)

// IssueCacheName and IssueCacheVersion identify the namespace holding fetched
// issue titles. Bump the version whenever forge.Issue changes shape: the new
// version starts from an empty file and refetches, which is always safe because
// the values are only ever a cache of what the forge already knows.
const (
	IssueCacheName    = "forge-issues"
	IssueCacheVersion = 2

	// IssueCacheTTL is how long a fetched title is trusted. Titles change
	// rarely, so this is deliberately long: the cost of a stale title is a
	// slightly wrong label, while the cost of a short TTL is a network round
	// trip on a listing the user expects to be instant.
	// This is reused for PRs as well.
	IssueCacheTTL = 24 * time.Hour
	// IssueCacheMissingTTL is shorter because an issue that did not resolve
	// may simply have been created since, or moved.
	// This is reused for PRs as well.
	IssueCacheMissingTTL = time.Hour
)

type WorktreeNumber struct {
	Number       int
	IsExplicitPr bool
}

// List returns every worktree under the configured root, each paired with its
// issue title. Titles come from the cache when they are still fresh; anything
// missing or expired is fetched in a single batched request.
func (w *RealWorktree) List(opts model.WorktreeListOpts) ([]model.WorktreeEntry, error) {
	cfg, err := w.resolveListConfig(opts)
	if err != nil {
		return nil, err
	}
	// Caught here rather than left to the fetch, which treats failures as
	// transient and degrades to bare numbers. A block with no repo is a config
	// mistake that no amount of retrying fixes, so it gets a real error.
	if cfg.Repo == "" {
		return nil, fmt.Errorf(
			"the [[worktree]] block for path %q is missing `repo = \"org/repo\"`", cfg.Path)
	}

	repoPath, err := w.home.ExpandPath(cfg.Path)
	if err != nil {
		return nil, err
	}
	root := w.worktreeRoot(cfg, repoPath)

	f, err := w.ForgeFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	numbers, err := w.worktreeNumbers(f, root)
	if err != nil {
		return nil, err
	}

	issues, err := w.resolveIssues(cfg, f, numbers, opts.Refresh)
	if err != nil {
		return nil, err
	}

	entries := make([]model.WorktreeEntry, 0, len(numbers))
	for _, number := range numbers {
		name := f.WorktreeName(number.Number, number.IsExplicitPr)
		refType := model.RefTypeIssue
		if number.IsExplicitPr {
			refType = model.RefTypePr
		}
		entry := model.WorktreeEntry{
			Number:  number.Number,
			RefType: refType,
			Name:    name,
			Path:    w.path.Join(root, name),
			Icon:    cfg.Icon,
		}
		if issue, ok := issues[name]; ok {
			entry.Title = issue.Title
			entry.State = string(issue.State)
		}
		entry.Created = w.createdAt(entry.Path)
		entries = append(entries, entry)
	}
	return entries, nil
}

// resolveListConfig picks the [[worktree]] block to list. --repo matches on the
// repo name, --path on the resolved worktree root, and neither falls back to
// detecting the repo from the current directory.
func (w *RealWorktree) resolveListConfig(opts model.WorktreeListOpts) (model.WorktreeConfig, error) {
	if opts.Repo != "" {
		return w.resolveConfig(opts.Repo)
	}
	if opts.Path == "" {
		return w.resolveConfig("")
	}

	wanted, err := w.home.ExpandPath(opts.Path)
	if err != nil {
		return model.WorktreeConfig{}, err
	}
	var known []string
	for _, cfg := range w.config.WorktreeConfigs {
		repoPath, err := w.home.ExpandPath(cfg.Path)
		if err != nil {
			continue
		}
		root := w.worktreeRoot(cfg, repoPath)
		if w.sameDir(root, wanted) {
			return cfg, nil
		}
		known = append(known, root)
	}
	return model.WorktreeConfig{}, fmt.Errorf(
		"no [[worktree]] config has a worktree root at %q (known roots: %v)", wanted, known)
}

// worktreeNumbers lists the issue numbers that have a worktree directory. Only
// numerically named directories count — that is the naming Connect creates, and
// it keeps unrelated files in the root from showing up as entries.
func (w *RealWorktree) worktreeNumbers(f forge.Forge, root string) ([]*WorktreeNumber, error) {
	dirEntries, err := w.os.ReadDir(root)
	if err != nil {
		// An absent root just means no worktrees have been created yet.
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading worktree root %q: %w", root, err)
	}
	var numbers []*WorktreeNumber
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}
		number, isExplicitPr, err := f.WorktreeNumber(dirEntry.Name())
		if err != nil {
			continue
		}
		numbers = append(numbers, &WorktreeNumber{Number: number, IsExplicitPr: isExplicitPr})
	}
	sort.Slice(numbers, func(i, j int) bool {
		iIsExplicitPr := numbers[i].IsExplicitPr
		jIsExplicitPr := numbers[j].IsExplicitPr

		if iIsExplicitPr == jIsExplicitPr {
			return numbers[i].Number < numbers[j].Number
		}
		return !iIsExplicitPr && jIsExplicitPr
	})
	return numbers, nil
}

// resolveIssues returns the issues and explicit pull requests for given
// worktree numbers using the provided forge driver, reading the cache first and
// batching everything stale into one request. Fetch failures are logged and
// swallowed: a missing title degrades the listing to a bare number, which is
// far better than failing it outright.
//
// refresh counts every number as stale, so the whole listing is refetched and
// the cache written back over.
func (w *RealWorktree) resolveIssues(
	cfg model.WorktreeConfig,
	f forge.Forge,
	numbers []*WorktreeNumber,
	refresh bool,
) (map[string]forge.Issue, error) {
	if len(numbers) == 0 {
		return nil, nil
	}

	host, err := forge.Host(cfg)
	if err != nil {
		return nil, err
	}
	repo := cfg.Repo
	entries := w.issues.Load()
	issues := make(map[string]forge.Issue, len(numbers))

	keys := make([]string, 0, len(numbers))
	byKey := make(map[string]*WorktreeNumber, len(numbers))
	for _, number := range numbers {
		key := f.IssueKey(host, repo, number.Number, number.IsExplicitPr)
		keys = append(keys, key)
		byKey[key] = number
		// Stale values are used as-is here; the refetch below overwrites them
		// when it succeeds, so a rate-limited run still shows real titles.
		if issue, found, _ := w.issues.Lookup(entries, key); found {
			issues[f.WorktreeName(number.Number, number.IsExplicitPr)] = issue
		}
	}

	staleKeys := keys
	if !refresh {
		staleKeys = w.issues.Stale(entries, keys)
	}
	if len(staleKeys) == 0 {
		return issues, nil
	}

	staleIssueNumbers := make([]int, 0, len(staleKeys))
	stalePrNumbers := make([]int, 0, len(staleKeys))
	for _, key := range staleKeys {
		number := byKey[key]
		if number.IsExplicitPr {
			stalePrNumbers = append(stalePrNumbers, number.Number)
		} else {
			staleIssueNumbers = append(staleIssueNumbers, number.Number)
		}
	}

	foundIssues, missingIssues, err := f.Issues(host, repo, staleIssueNumbers)
	if err != nil {
		slog.Debug("worktree list: could not fetch issues", "repo", repo, "error", err)
		return issues, nil
	}
	foundPrs, missingPrs, err := f.PrsAsIssues(host, repo, stalePrNumbers)
	if err != nil {
		slog.Debug("worktree list: could not fetch pull requests", "repo", repo, "error", err)
		return issues, nil
	}

	for number, issue := range foundIssues {
		issues[f.WorktreeName(number, false)] = issue
		entries.Put(f.IssueKey(host, repo, number, false), issue)
	}
	for _, number := range missingIssues {
		// Negatively cached so a worktree for a deleted issue stops being
		// refetched on every listing.
		entries.PutMissing(f.IssueKey(host, repo, number, false))
	}
	for number, issue := range foundPrs {
		issues[f.WorktreeName(number, true)] = issue
		entries.Put(f.IssueKey(host, repo, number, true), issue)
	}
	for _, number := range missingPrs {
		// Negatively cached so a worktree for a deleted pull request stops being
		// refetched on every listing.
		entries.PutMissing(f.IssueKey(host, repo, number, true))
	}

	if err := w.issues.Save(entries); err != nil {
		slog.Debug("worktree list: could not save issue cache", "error", err)
	}
	return issues, nil
}

func (w *RealWorktree) createdAt(dir string) time.Time {
	for _, marker := range []string{".git", ".jj"} {
		if info, err := w.os.Stat(w.path.Join(dir, marker)); err == nil {
			return info.ModTime()
		}
	}
	return time.Time{}
}
