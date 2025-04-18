package namer

import (
	"errors"
	"fmt"
	"strings"
)

func gitBareRootName(n *RealNamer, path string) (string, error) {
	isGit, bareRoot, err := gitBareRootFromWorkTreeList(n, path)
	if err != nil {
		return "", err
	}

	if isGit && bareRoot != "" {
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
	isGit, bareRoot, err := gitBareRootFromWorkTreeList(n, path)
	if err != nil {
		return "", err
	}

	if isGit && bareRoot != "" {
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

func gitBareRootFromWorkTreeList(n *RealNamer, path string) (bool, string, error) {
	isGit, out, err := n.git.WorkTreeList(path)
	if err != nil {
		return false, "", err
	}

	fields := strings.Fields(out)
	if len(fields) == 0 {
		return false, "", errors.New("error parsing git worktree fields")
	}

	bareRoot := fields[0]
	if isGit && strings.HasSuffix(bareRoot, "/.bare") {
		bareRoot = strings.TrimSuffix(bareRoot, "/.bare")
	}

	return true, bareRoot, nil
}
