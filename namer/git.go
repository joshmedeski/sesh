package namer

import (
	"strings"
)

// determineGitRootPath parses `git worktree list --porcelain` output and returns
// the main working tree path. The first entry (block of lines up to a blank line)
// describes the main working tree: `worktree <path>` followed by attributes like
// `HEAD <sha>`, `branch <ref>`, or `bare`. Porcelain preserves paths containing
// spaces and gives a dedicated `bare` token instead of a `(bare)` substring.
func determineGitRootPath(out string) string {
	var path string
	var isBare bool
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			if path != "" {
				break
			}
			continue
		}
		if p, ok := strings.CutPrefix(line, "worktree "); ok {
			path = p
			continue
		}
		if line == "bare" {
			isBare = true
		}
	}
	if path == "" {
		return ""
	}
	if isBare {
		if strings.HasSuffix(path, "/.bare") {
			return strings.TrimSuffix(path, "/.bare")
		}
		if strings.HasSuffix(path, "/.git") {
			return strings.TrimSuffix(path, "/.git")
		}
	}
	return path
}

func getGitRootPath(n *RealNamer, path string) (bool, string, error) {
	isGit, list, err := n.git.WorktreeList(path)
	if err != nil {
		return false, "", nil
	}
	if !isGit {
		return false, "", nil
	}
	rootPath := determineGitRootPath(list)
	if rootPath == "" {
		return false, "", nil
	}
	return true, rootPath, nil
}

// gitAnchor returns the path used as the naming anchor for git strategies.
// When git_namer_use_worktree_root is true, the current worktree's top-level
// directory (`git rev-parse --show-toplevel`) is used so sibling worktrees
// (e.g. `git worktree add ../feat`) get named after themselves rather than
// the main clone. Otherwise the main working tree root is used.
func gitAnchor(n *RealNamer, path string) (bool, string, error) {
	if n.config.GitNamerUseWorktreeRoot {
		isGit, topLevel, err := n.git.ShowTopLevel(path)
		if err != nil {
			return false, "", nil
		}
		if !isGit || topLevel == "" {
			return false, "", nil
		}
		return true, topLevel, nil
	}
	return getGitRootPath(n, path)
}

// anchorRepoName returns the repo-name segment for a git anchor path, applying
// git_dir_length when set (mirrors dir_length behavior, but scoped to git).
func anchorRepoName(n *RealNamer, anchor string) string {
	if n.config.GitDirLength > 1 {
		return lastNComponents(anchor, n.config.GitDirLength)
	}
	return n.pathwrap.Base(anchor)
}

// Names a session as <repoName>/<relativePath> where repoName is the basename
// of the anchor (main working tree or current worktree top-level, depending on
// config) and relativePath is the input path's offset from that anchor.
func gitName(n *RealNamer, path string) (string, error) {
	isGit, anchor, err := gitAnchor(n, path)
	if err != nil {
		return "", err
	}
	if !isGit || anchor == "" {
		return "", nil
	}
	repoName := anchorRepoName(n, anchor)
	relativePath := strings.TrimPrefix(path, anchor)
	return repoName + relativePath, nil
}

// Names the session root as <repoName>/<relativeWorktreePath> where the relative
// path is derived from the current worktree's top-level directory. This collapses
// nested subdirectories to their containing worktree (main tree or linked worktree).
func gitRootName(n *RealNamer, path string) (string, error) {
	if n.config.GitNamerUseWorktreeRoot {
		isGit, topLevel, err := n.git.ShowTopLevel(path)
		if err != nil {
			return "", err
		}
		if !isGit || topLevel == "" {
			return "", nil
		}
		return anchorRepoName(n, topLevel), nil
	}

	isGit, rootPath, err := getGitRootPath(n, path)
	if err != nil {
		return "", err
	}
	if !isGit || rootPath == "" {
		return "", nil
	}
	repoName := anchorRepoName(n, rootPath)
	_, topLevel, _ := n.git.ShowTopLevel(path)
	if topLevel == "" {
		return repoName, nil
	}
	relativePath := strings.TrimPrefix(topLevel, rootPath)
	return repoName + relativePath, nil
}
