package namer

import (
	"strings"
)

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
