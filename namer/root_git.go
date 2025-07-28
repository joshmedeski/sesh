package namer

func gitRootName(n *RealNamer, path string) (string, error) {
	isGit, topLevelDir, _ := n.git.ShowTopLevel(path)
	if isGit && topLevelDir != "" {
		name, err := dirName(n, topLevelDir)
		if err != nil {
			return "", err
		}
		return name, nil
	} else {
		return "", nil
	}
}
