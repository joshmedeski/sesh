package session

import (
	"fmt"
	"joshmedeski/sesh/dir"
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

func DetermineName(entry string) string {
	name := entry
	fullPath, err := dir.FullPath(entry)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if path.IsAbs(fullPath) {
		// TODO: parent directory feature flag detection
		baseName := filepath.Base(fullPath)
		gitRootPath := git.RootPath(fullPath)
		if gitRootPath != "" {
			gitRootBaseName := filepath.Base(gitRootPath)
			relativePath := strings.TrimPrefix(fullPath, gitRootPath)
			name = gitRootBaseName + relativePath
		} else {
			name = baseName
		}
	}
	return convertToValidName(name)
}
