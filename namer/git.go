package namer

import (
	"fmt"
	"strings"
)

func gitName(n *RealNamer, path string) (string, error) {
	isGit, topLevelDir, _ := n.git.GitRoot(path)
	if isGit && topLevelDir != "" {
		relativePath := strings.TrimPrefix(path, topLevelDir)
		baseDir := n.pathwrap.Base(topLevelDir)
		name := baseDir + relativePath
		return name, nil
	} else {
		return "", nil
	}
}

func gitRootName(n *RealNamer, path string) (string, error) {
	isGit, topLevelDir, _ := n.git.GitRoot(path)
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
