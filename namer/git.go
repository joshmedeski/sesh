package namer

import (
	"fmt"
	"strings"
)

func gitBareRootName(n *RealNamer, path string) (string, error) {
	isGit, worktreeList, _ := n.git.WorkTreeList(path)
	if isGit && worktreeList != "" {
		bareRoot, err := parseBareFromWorkTreeList(worktreeList)
		if err != nil {
			return "", err
		}
		name, err := n.home.ShortenHome(bareRoot)
		if err != nil {
			return "", fmt.Errorf("couldn't shorten path: %q", err)
		}
		return name, nil
	} else {
		return "", nil
	}
}

// Gets the name from a git bare repository
func gitBareName(n *RealNamer, path string) (string, error) {
	isGit, worktreeList, _ := n.git.WorkTreeList(path)
	if isGit && worktreeList != "" {
		bareRoot, err := parseBareFromWorkTreeList(worktreeList)
		if err != nil {
			return "", err
		}
		relativePath := strings.TrimPrefix(path, bareRoot)
		baseDir := n.pathwrap.Base(bareRoot)
		name := baseDir + relativePath
		return name, nil
	} else {
		return "", nil
	}
}

func gitRootName(n *RealNamer, path string) (string, error) {
	isGit, topLevelDir, _ := n.git.ShowTopLevel(path)
	if isGit && topLevelDir != "" {
		name, err := n.home.ShortenHome(topLevelDir)
		if err != nil {
			return "", fmt.Errorf("couldn't shorten path: %q", err)
		}
		return name, nil
	} else {
		return "", nil
	}
}

// Gets the name from a git repository
func gitName(n *RealNamer, path string) (string, error) {
	isGit, topLevelDir, _ := n.git.ShowTopLevel(path)
	if isGit && topLevelDir != "" {
		relativePath := strings.TrimPrefix(path, topLevelDir)
		baseDir := n.pathwrap.Base(topLevelDir)
		name := baseDir + relativePath
		return name, nil
	} else {
		return "", nil
	}
}
