package session

import (
	"fmt"
	"os"
	"path"

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
	if path.IsAbs(fullPath) {
		return fullPath, nil
	}

	if tmux.IsSession(fullPath) {
		return fullPath, nil
	}

	zoxideResult, err := zoxide.Query(fullPath)
	if err != nil {
		fmt.Println("Couldn't query zoxide", err)
		os.Exit(1)
	}
	if zoxideResult != nil {
		return zoxideResult.Path, nil
	}

	return fullPath, nil
}
