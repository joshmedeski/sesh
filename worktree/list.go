package worktree

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"strconv"
	"time"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/model"
)

// IssueCacheName and IssueCacheVersion identify the namespace holding fetched
// issue titles. Bump the version whenever github.Issue changes shape: the new
// version starts from an empty file and refetches, which is always safe because
// the values are only ever a cache of what GitHub already knows.
const (
	IssueCacheName    = "github-issues"
	IssueCacheVersion = 1

	// IssueCacheTTL is how long a fetched title is trusted. Titles change
	// rarely, so this is deliberately long: the cost of a stale title is a
	// slightly wrong label, while the cost of a short TTL is a network round
	// trip on a listing the user expects to be instant.
	IssueCacheTTL = 24 * time.Hour
	// IssueCacheMissingTTL is shorter because an issue that did not resolve
	// may simply have been created since, or moved.
	IssueCacheMissingTTL = time.Hour
)

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

	numbers, err := w.worktreeNumbers(root)
	if err != nil {
		return nil, err
	}

	issues := w.resolveIssues(cfg.Repo, numbers)

	entries := make([]model.WorktreeEntry, 0, len(numbers))
	for _, number := range numbers {
		entry := model.WorktreeEntry{
			Number: number,
			Path:   w.path.Join(root, strconv.Itoa(number)),
		}
		if issue, ok := issues[number]; ok {
			entry.Title = issue.Title
			entry.State = issue.State
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// resolveListConfig picks the [[worktree]] block to list. --repo matches on the
// GitHub name, --path on the resolved worktree root, and neither falls back to
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
func (w *RealWorktree) worktreeNumbers(root string) ([]int, error) {
	dirEntries, err := w.os.ReadDir(root)
	if err != nil {
		// An absent root just means no worktrees have been created yet.
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading worktree root %q: %w", root, err)
	}

	var numbers []int
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}
		number, err := strconv.Atoi(dirEntry.Name())
		if err != nil {
			continue
		}
		numbers = append(numbers, number)
	}
	sort.Ints(numbers)
	return numbers, nil
}

// resolveIssues returns the issue for each number, reading the cache first and
// batching everything stale into one request. Fetch failures are logged and
// swallowed: a missing title degrades the listing to a bare number, which is
// far better than failing it outright.
func (w *RealWorktree) resolveIssues(repo string, numbers []int) map[int]github.Issue {
	if len(numbers) == 0 {
		return nil
	}

	entries := w.issues.Load()
	issues := make(map[int]github.Issue, len(numbers))

	keys := make([]string, 0, len(numbers))
	byKey := make(map[string]int, len(numbers))
	for _, number := range numbers {
		key := issueKey(repo, number)
		keys = append(keys, key)
		byKey[key] = number
		// Stale values are used as-is here; the refetch below overwrites them
		// when it succeeds, so a rate-limited run still shows real titles.
		if issue, found, _ := w.issues.Lookup(entries, key); found {
			issues[number] = issue
		}
	}

	staleKeys := w.issues.Stale(entries, keys)
	if len(staleKeys) == 0 {
		return issues
	}

	staleNumbers := make([]int, 0, len(staleKeys))
	for _, key := range staleKeys {
		staleNumbers = append(staleNumbers, byKey[key])
	}

	found, missing, err := w.github.Issues(repo, staleNumbers)
	if err != nil {
		slog.Debug("worktree list: could not fetch issues", "repo", repo, "error", err)
		return issues
	}

	for number, issue := range found {
		issues[number] = issue
		entries.Put(issueKey(repo, number), issue)
	}
	for _, number := range missing {
		// Negatively cached so a worktree for a deleted issue stops being
		// refetched on every listing.
		entries.PutMissing(issueKey(repo, number))
	}

	if err := w.issues.Save(entries); err != nil {
		slog.Debug("worktree list: could not save issue cache", "error", err)
	}
	return issues
}

// issueKey namespaces cache keys by repo so worktrees for different repos with
// overlapping numbers cannot collide.
func issueKey(repo string, number int) string {
	return repo + "#" + strconv.Itoa(number)
}
