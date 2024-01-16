package connect

import (
	"github.com/joshmedeski/sesh/session"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

func Connect(choice string, alwaysSwitch bool, command string) error {
	session := session.Determine(choice)
	zoxide.Add(session.Path)
	tmux.Connect(tmux.TmuxSession{
		Name: session.Name,
		Path: session.Path,
	}, alwaysSwitch, command)
	return nil
}
