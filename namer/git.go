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

// Names a session as <repoName>/<relativePath> where repoName is the basename
// of the main working tree and relativePath is the input path's offset from it.
func gitName(n *RealNamer, path string) (string, error) {
	isGit, rootPath, err := getGitRootPath(n, path)
	if err != nil {
		return "", err
	}
	if !isGit || rootPath == "" {
		return "", nil
	}
	repoName := n.pathwrap.Base(rootPath)
	relativePath := strings.TrimPrefix(path, rootPath)
	return repoName + relativePath, nil
}

// Names the session root as <repoName>/<relativeWorktreePath> where the relative
// path is derived from the current worktree's top-level directory. This collapses
// nested subdirectories to their containing worktree (main tree or linked worktree).
func gitRootName(n *RealNamer, path string) (string, error) {
	isGit, rootPath, err := getGitRootPath(n, path)
	if err != nil {
		return "", err
	}
	if !isGit || rootPath == "" {
		return "", nil
	}
	repoName := n.pathwrap.Base(rootPath)
	_, topLevel, _ := n.git.ShowTopLevel(path)
	if topLevel == "" {
		return repoName, nil
	}
	relativePath := strings.TrimPrefix(topLevel, rootPath)
	return repoName + relativePath, nil
}
