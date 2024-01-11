package connect

import (
	"fmt"
	"joshmedeski/sesh/dir"
	"joshmedeski/sesh/session"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
	"os"
)

func Connect(choice string, alwaysSwitch bool) error {
	fullPath := dir.FullPath(choice)
	zoxide.Add(fullPath)
	sessionName := session.DetermineName(fullPath)
	if sessionName == "" {
		fmt.Println("Session couldn't be determined")
		os.Exit(1)
	}
	// TODO: get zoxide result if not path and tmux session doesn't exist
	session := tmux.TmuxSession{
		Name:           sessionName,
		StartDirectory: fullPath,
	}
	tmux.Connect(session, alwaysSwitch)
	return nil
}
