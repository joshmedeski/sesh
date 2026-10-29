package namer

import "strings"

func jjName(n *RealNamer, path string) (string, error) {
	isJj, workspace, repo := n.jj.Roots(path)
	if !isJj {
		return "", nil
	}
	anchor := repo
	if n.config.GitNamerUseWorktreeRoot {
		anchor = workspace
	}
	return anchorRepoName(n, anchor) + strings.TrimPrefix(n.pathwrap.ToSlash(path), anchor), nil
}

func jjRootName(n *RealNamer, path string) (string, error) {
	isJj, workspace, repo := n.jj.Roots(path)
	if !isJj {
		return "", nil
	}
	if n.config.GitNamerUseWorktreeRoot {
		return anchorRepoName(n, workspace), nil
	}
	return anchorRepoName(n, repo) + strings.TrimPrefix(workspace, repo), nil
}
