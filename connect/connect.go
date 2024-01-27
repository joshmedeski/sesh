package connect

import (
	"fmt"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

func Connect(
	choice string,
	alwaysSwitch bool,
	command string,
	config *config.Config,
) error {
	cmd, err := tmux.NewCommand(tmux.Options{})
	if err != nil {
		return fmt.Errorf("unable to configure the tmux command: %w", err)
	}
	s, err := cmd.SessionByName(choice)
	if err != nil {
		return err
	}
	if err = zoxide.Add(s.Path); err != nil {
		return fmt.Errorf("unable to connect to %q: %w", choice, err)
	}

	return tmux.Connect(s, alwaysSwitch, command, s.Path, config)
}
