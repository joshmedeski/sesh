package namer

import (
	"strings"
)

func determineBareWorktreePath(out string) string {
	for line := range strings.SplitSeq(strings.TrimSpace(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			return ""
		}

		if strings.Contains(line, "(bare)") {
			parts := strings.Fields(line)
			if len(parts) >= 2 && parts[1] == "(bare)" {
				path := parts[0]
				// TODO: move `.bare` folder naming convention to configuration?
				if strings.HasSuffix(path, "/.bare") {
					trimmedPath := strings.TrimSuffix(path, "/.bare")
					return trimmedPath
				} else {
					return parts[0]
				}
			}
		}
	}
	return ""
}

func getBareWorktreePath(n *RealNamer, path string) (bool, string, error) {
	isGit, list, err := n.git.WorktreeList(path)
	if err != nil {
		return false, "", nil
	}
	if !isGit {
		return false, "", nil
	}
	barePath := determineBareWorktreePath(list)
	if barePath == "" {
		return false, "", nil
	}
	return true, barePath, nil
}

// Gets the name from a git bare repository
func gitBareName(n *RealNamer, path string) (string, error) {
	isGit, barePath, err := getBareWorktreePath(n, path)
	if err != nil {
		return "", err
	}
	if isGit && barePath != "" {
		relativePath := strings.TrimPrefix(path, barePath)
		baseDir := n.pathwrap.Base(barePath)
		name := baseDir + relativePath
		return name, nil
	}
	return "", nil
}
