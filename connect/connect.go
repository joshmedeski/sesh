package connect

import (
	"joshmedeski/sesh/session"
	"joshmedeski/sesh/tmux"
)

func Connect(choice string) error {
	sessionName := session.DetermineName(choice)

	if tmux.IsSession(sessionName) {
		tmux.Connect(sessionName)
	}
	return nil

	// TODO: if starting with ~ then it's a dir
	// TODO: if dir, create new tmux session
	// 	print("is path")
	// } else {
	// TODO: if not, then it's a tmux session
	// TODO: if tmux session, then attach
	// 	print("is not path")
	// }
}
