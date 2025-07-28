package namer

import (
	"fmt"
	"strings"
)

func gitBareRootName(n *RealNamer, path string) (string, error) {
	isGit, commonDir, _ := n.git.GitCommonDir(path)
	if isGit && strings.HasSuffix(commonDir, "/.bare") {
		topLevelDir := strings.TrimSuffix(commonDir, "/.bare")
		name, err := n.home.ShortenHome(topLevelDir)
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
	isGit, commonDir, _ := n.git.GitCommonDir(path)
	if isGit && strings.HasSuffix(commonDir, "/.bare") {
		topLevelDir := strings.TrimSuffix(commonDir, "/.bare")

		// Use dirName logic to respect dir_length for git bare repos
		repoName, err := dirName(n, topLevelDir, n.config.DirLength)
		if err != nil {
			return "", err
		}

		relativePath := strings.TrimPrefix(path, topLevelDir)
		name := repoName + relativePath
		return name, nil
	} else {
		return "", nil
	}
}

func gitRootName(n *RealNamer, path string) (string, error) {
	isGit, topLevelDir, _ := n.git.ShowTopLevel(path)
	if isGit && topLevelDir != "" {
		// Use dirName logic to respect dir_length for git root names
		repoName, err := dirName(n, topLevelDir, n.config.DirLength)
		if err != nil {
			return "", err
		}
		return repoName, nil
	} else {
		return "", nil
	}
}

// Gets the name from a git repository
func gitName(n *RealNamer, path string) (string, error) {
	isGit, topLevelDir, _ := n.git.ShowTopLevel(path)
	if isGit && topLevelDir != "" {
		// Use dirName logic to respect dir_length for git repos
		repoName, err := dirName(n, topLevelDir, n.config.DirLength)
		if err != nil {
			return "", err
		}

		relativePath := strings.TrimPrefix(path, topLevelDir)
		name := repoName + relativePath
		return name, nil
	} else {
		return "", nil
	}
}
