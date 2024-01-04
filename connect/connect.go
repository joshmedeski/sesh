package connect

import (
	"joshmedeski/sesh/session"
)

func Connect(choice string) error {
	sessionName := session.DetermineName(choice)
	print(sessionName)
	return nil

	// TODO: generate session name from path
	// TODO: if starting with ~ then it's a dir
	// TODO: if dir, create new tmux session
	// 	print("is path")
	// } else {
	// TODO: if not, then it's a tmux session
	// TODO: if tmux session, then attach
	// 	print("is not path")
	// }
}
