package session

import (
	"fmt"
	"joshmedeski/sesh/git"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

type session struct {
	Type        string
	Value       string
	DisplayName string
}

type Srcs struct {
	Tmux   bool
	Zoxide bool
}

func checkAnyTrue(s interface{}) bool {
	val := reflect.ValueOf(s)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Ptr && !field.IsNil() && field.Elem().Bool() {
			return true
		}
	}
	return false
}

func Sessions(srcs Srcs) []string {
	var sessions []string
	anySrcs := checkAnyTrue(srcs)

	if !anySrcs || srcs.Tmux {
		tmuxSessions, err := tmux.Sessions()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		sessions = append(sessions, tmuxSessions...)
	}

	if !anySrcs || srcs.Zoxide {
		dirs, err := zoxide.Dirs()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		sessions = append(sessions, dirs...)
	}
	return sessions
}

func convertToValidName(name string) string {
	validName := strings.ReplaceAll(name, ".", "_")
	validName = strings.ReplaceAll(validName, ":", "_")
	return validName
}

// TODO: parent directory feature flag detection
func DetermineName(result string) string {
	name := result
	pathName := nameFromPath(result)
	if pathName != "" {
		name = pathName
	}
	return convertToValidName(name)
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
