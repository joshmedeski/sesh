package connect

import (
	"fmt"
	"joshmedeski/sesh/dir"
	"joshmedeski/sesh/session"
	"joshmedeski/sesh/tmux"
	"os"
)

func Connect(choice string) error {
	fullPath, err := dir.FullPath(choice)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	sessionName := session.DetermineName(choice)
	session := tmux.TmuxSession{
		Name:           sessionName,
		StartDirectory: fullPath,
	}
	tmux.Connect(session)
	return nil
}
