package session

import (
	"joshmedeski/sesh/dir"
	"joshmedeski/sesh/tmux"
	"path"
)

func DeterminPath(choice string) (string, error) {
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
