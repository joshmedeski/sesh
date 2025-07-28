package namer

import (
	"strings"
)

// Gets the name from a git repository
func gitName(n *RealNamer, path string) (string, error) {
	isGit, topLevelDir, _ := n.git.ShowTopLevel(path)
	if isGit && topLevelDir != "" {
		repoName, err := dirName(n, topLevelDir)
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
