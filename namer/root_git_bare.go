package namer

import "strings"

// Gets the root name from a git bare repository
// Returns the bare repo name plus the worktree directory (e.g., "myrepo/main")
func gitBareRootName(n *RealNamer, path string) (string, error) {
	isGit, barePath, err := getBareWorktreePath(n, path)
	if err != nil {
		return "", err
	}
	if isGit && barePath != "" {
		repoName, err := dirName(n, barePath)
		if err != nil {
			return "", err
		}

		// Get the worktree's top-level directory
		_, topLevelDir, _ := n.git.ShowTopLevel(path)
		if topLevelDir != "" {
			relativePath := strings.TrimPrefix(topLevelDir, barePath)
			return repoName + relativePath, nil
		}
		return repoName, nil
	}
	return "", nil
}
