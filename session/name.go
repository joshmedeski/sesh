package session

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/git"
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

func DetermineName(result string, config *config.Config) string {
	name := result
	// TODO: parent directory config option detection
	pathName := nameFromPath(result)
	if pathName != "" {
		name = pathName
	}
	return convertToValidName(name)
}
