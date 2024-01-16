package session

import (
	"path"

	"github.com/joshmedeski/sesh/dir"
	"github.com/joshmedeski/sesh/tmux"
)

func DeterminePath(choice string) (string, error) {
	fullPath := dir.FullPath(choice)
	if path.IsAbs(fullPath) {
		return fullPath, nil
	}

	if tmux.IsSession(fullPath) {
		return fullPath, nil
	}
	// TODO: if not absolute path, get zoxide results
	// TODO: get zoxide result if not path and tmux session doesn't exist
	return fullPath, nil
}
