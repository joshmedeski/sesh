package namer

import "strings"

// Gets the name from a git bare repository
func gitBareName(n *RealNamer, path string) (string, error) {
	var name string
	isGit, commonDir, err := n.git.GitCommonDir(path)
	if err != nil {
		return "", err
	}
	if isGit && strings.HasSuffix(commonDir, "/.bare") {
		topLevelDir := strings.TrimSuffix(commonDir, "/.bare")
		relativePath := strings.TrimPrefix(path, topLevelDir)
		baseDir := n.pathwrap.Base(topLevelDir)
		name = baseDir + relativePath
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
