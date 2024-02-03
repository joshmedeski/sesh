package session

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/joshmedeski/sesh/dir"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

func DeterminePath(choice string) (string, error) {
	if choice == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return cwd, nil
	}
	fullPath := dir.FullPath(choice)

	realPath, err := filepath.EvalSymlinks(choice)
	if err == nil && path.IsAbs(realPath) {
		return realPath, nil
	}

	if path.IsAbs(fullPath) {
		return fullPath, nil
	}

	isSession, sessionPath := tmux.IsSession(fullPath)
	if isSession && sessionPath != "" {
		return sessionPath, nil
	}

	zoxideResult, err := zoxide.Query(fullPath)
	if err != nil {
		fmt.Println("Couldn't query zoxide", err)
		os.Exit(1)
	}
	if zoxideResult.Path != "" {
		return zoxideResult.Path, nil
	}

	return fullPath, nil
}
