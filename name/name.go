package name

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/tmux"
)

func convertToValidName(name string) string {
	validName := strings.ReplaceAll(name, ".", "_")
	validName = strings.ReplaceAll(validName, ":", "_")
	return validName
}

func nameFromPath(result string) string {
	name := ""
	if path.IsAbs(result) {
		gitName := nameFromGit(result)
		if gitName != "" {
			name = gitName
		} else {
			name = filepath.Base(result)
		}
	}
	return name
}

func nameFromGit(result string) string {
	gitRootPath := git.RootPath(result)
	if gitRootPath == "" {
		return ""
	}
	root := ""
	base := ""
	gitWorktreePath := git.WorktreePath(result)
	if gitWorktreePath != "" {
		root = gitWorktreePath
		base = filepath.Base(gitWorktreePath)
	} else {
		root = gitRootPath
		base = filepath.Base(gitRootPath)
	}
	relativePath := strings.TrimPrefix(result, root)
	nameFromGit := base + relativePath
	return nameFromGit
}

func DetermineName(choice string, path string) string {
	session, _ := tmux.FindSession(choice)
	if session != nil {
		return session.Name
	}

	// TODO: parent directory config option detection
	pathName := nameFromPath(path)
	if pathName != "" {
		return convertToValidName(pathName)
	}

	return convertToValidName(choice)
}
