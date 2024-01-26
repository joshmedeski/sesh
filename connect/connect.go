package connect

import (
	"fmt"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/session"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

func Connect(choice string, alwaysSwitch bool, command string, config *config.Config) error {
	session, err := session.Determine(choice, config)
	if err != nil {
		return fmt.Errorf("unable to connect to %q: %w", choice, err)
	}
	if err = zoxide.Add(session.Path); err != nil {
		return fmt.Errorf("unable to connect to %q: %w", choice, err)
	}
	return tmux.Connect(tmux.TmuxSession{
		Name: session.Name,
		Path: session.Path,
	}, alwaysSwitch, command, session.Path, config)
}
