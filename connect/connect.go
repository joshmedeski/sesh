package connect

import (
	"joshmedeski/sesh/session"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
)

func Connect(choice string, alwaysSwitch bool) error {
	session := session.Determine(choice)
	zoxide.Add(session.Path)
	tmux.Connect(tmux.TmuxSession{
		Name: session.Name,
		Path: session.Path,
	}, alwaysSwitch)
	return nil
}
