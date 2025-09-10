package namer

import "fmt"

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
