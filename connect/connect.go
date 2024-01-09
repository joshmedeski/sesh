package connect

import (
	"joshmedeski/sesh/dir"
	"joshmedeski/sesh/session"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
)

func Connect(choice string) error {
	fullPath := dir.FullPath(choice)
	zoxide.Add(fullPath)
	sessionName := session.DetermineName(fullPath)
	// TODO: get zoxide result if not path and tmux session doesn't exist
	session := tmux.TmuxSession{
		Name:           sessionName,
		StartDirectory: fullPath,
	}
	tmux.Connect(session)
	return nil
}
