package lister

import (
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
)

type worktreeRoot struct {
	prefix string
	config model.WorktreeConfig
	// Set once per root by worktreeResolver. forge stays nil when the config's driver
	// or host cannot be resolved, which costs the row its title but not its resolveRef.
	forge forge.Forge
	host  string
}

func worktreeRoots(config model.Config, h home.Home) []worktreeRoot {
	var roots []worktreeRoot
	for _, worktree := range config.WorktreeConfigs {
		if worktree.Path == "" {
			continue
		}
		if repoPath, err := h.ExpandPath(worktree.Path); err == nil {
			roots = append(roots, worktreeRoot{
				prefix: worktree.Root(repoPath) + string(filepath.Separator),
				config: worktree,
			})
		}
	}
	return roots
}

// resolveRef resolves the issue/PR reference for a worktree: its number, whether
// it explicitly refers to a pull request, and the forge-specific reference
// name. ok is false for a path that is not under this root and for a name the
// forge's driver rejects; either way the session is not treated as a worktree.
func (r worktreeRoot) resolveRef(path string) (number int, isExplicitPr bool, name string, ok bool) {
	rest, ok := strings.CutPrefix(path, r.prefix)
	if !ok {
		return 0, false, "", false
	}
	name, _, _ = strings.Cut(rest, string(filepath.Separator))
	// A root whose driver or host never resolved cannot know the forge's name
	// form, so only a plain number is recognized; the row keeps its number and
	// loses only its title.
	if r.forge == nil {
		number, err := strconv.Atoi(name)
		return number, false, "", err == nil
	}
	number, isExplicitPr, err := r.forge.WorktreeNumber(name)
	if err != nil {
		return 0, false, "", false
	}
	return number, isExplicitPr, name, true
}

func WorktreeResolver(config model.Config, h home.Home, forges forge.Resolver, issues *cache.Namespace[forge.Issue]) func(path string) *model.WorktreeInfo {
	roots := worktreeRoots(config, h)
	if len(roots) == 0 {
		return nil
	}

	// The driver and host are resolved once per root rather than per lookup (the closure
	// below runs for every session in a listing). A root that cannot be resolved keeps its
	// number and loses only its title, so one misconfigured block cannot fail a listing.
	for i := range roots {
		f, err := forges.FromConfig(roots[i].config)
		if err != nil {
			slog.Debug("worktree resolver: could not resolve forge", "repo", roots[i].config.Repo, "error", err)
			continue
		}
		host, err := forge.Host(roots[i].config)
		if err != nil {
			slog.Debug("worktree resolver: could not resolve host", "repo", roots[i].config.Repo, "error", err)
			continue
		}
		roots[i].forge, roots[i].host = f, host
	}

	var entries cache.Entries[forge.Issue]
	return func(path string) *model.WorktreeInfo {
		if path == "" {
			return nil
		}
		path = filepath.Clean(path)
		for _, root := range roots {
			number, isExplicitPr, name, ok := root.resolveRef(path)
			if !ok {
				continue
			}
			refType := model.RefTypeIssue
			if isExplicitPr {
				refType = model.RefTypePr
			}
			info := &model.WorktreeInfo{
				Repo:    root.config.Repo,
				Number:  number,
				Name:    name,
				RefType: refType,
			}
			if issues == nil || info.Repo == "" || root.forge == nil {
				return info
			}
			if entries == nil {
				entries = issues.Load()
			}
			// The key must be the driver's: the worktree command writes it with the same
			// method, and a key built anywhere else would silently miss. The flag is part
			// of it because forges like GitLab number issues and merge requests
			// independently, so "!101" and "101" are two different entries.
			key := root.forge.IssueKey(root.host, info.Repo, number, isExplicitPr)
			if issue, found, _ := issues.Lookup(entries, key); found {
				info.Title, info.State = issue.Title, string(issue.State)
			}
			return info
		}
		return nil
	}
}
