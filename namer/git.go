package namer

import (
	"strings"
)

// determineGitRootPath parses the first non-empty line of `git worktree list`
// output and returns the main working tree path. For bare repos (first line
// contains "(bare)"), it trims a trailing "/.bare" or "/.git" suffix so the name
// reflects the parent directory rather than the bare-storage folder.
func determineGitRootPath(out string) string {
	out = strings.TrimSpace(out)
	if out == "" {
		return ""
	}
	firstLine := out
	if idx := strings.Index(out, "\n"); idx >= 0 {
		firstLine = out[:idx]
	}
	firstLine = strings.TrimSpace(firstLine)
	if firstLine == "" {
		return ""
	}
	parts := strings.Fields(firstLine)
	if len(parts) == 0 {
		return ""
	}
	path := parts[0]
	if strings.Contains(firstLine, "(bare)") {
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
