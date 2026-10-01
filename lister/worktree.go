package lister

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
)

type worktreeRoot struct {
	prefix string
	config model.WorktreeConfig
}

func worktreeRoots(config model.Config, h home.Home) []worktreeRoot {
	var roots []worktreeRoot
	for _, worktree := range config.WorktreeConfigs {
		if worktree.Path == "" {
			continue
		}
		if repoPath, err := h.ExpandPath(worktree.Path); err == nil {
			roots = append(roots, worktreeRoot{worktree.Root(repoPath) + string(filepath.Separator), worktree})
		}
	}
	return roots
}

func (r worktreeRoot) number(path string) (int, bool) {
	rest, ok := strings.CutPrefix(path, r.prefix)
	if !ok {
		return 0, false
	}
	dir, _, _ := strings.Cut(rest, string(filepath.Separator))
	number, err := strconv.Atoi(dir)
	return number, err == nil
}

func worktreeResolver(config model.Config, h home.Home, issues *cache.Namespace[github.Issue]) func(path string) *model.WorktreeInfo {
	roots := worktreeRoots(config, h)
	if len(roots) == 0 {
		return nil
	}

	var entries cache.Entries[github.Issue]
	return func(path string) *model.WorktreeInfo {
		if path == "" {
			return nil
		}
		path = filepath.Clean(path)
		for _, root := range roots {
			number, ok := root.number(path)
			if !ok {
				continue
			}
			info := &model.WorktreeInfo{Repo: root.config.Repo, Number: number}
			if issues == nil || info.Repo == "" {
				return info
			}
			if entries == nil {
				entries = issues.Load()
			}
			if issue, found, _ := issues.Lookup(entries, github.IssueKey(info.Repo, number)); found {
				info.Title, info.State = issue.Title, issue.State
			}
			return info
		}
		return nil
	}
}
