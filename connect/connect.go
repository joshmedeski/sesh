package connect

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

// Connect to an established tmux session if one with a name or path matching
// the 'choice' already exists, otherwise if the 'choice' is a valid file system
// path create and connect to a new tmux session with that name.
func Connect(
	choice string,
	alwaysSwitch bool,
	command string,
	config *config.Config,
) error {
	// Create a new tmux command.
	cmd, err := tmux.NewCommand(tmux.Options{})
	if err != nil {
		return fmt.Errorf("unable to configure the tmux command: %w", err)
	}

	// Check if the 'choice' is a valid tmux session name or path.
	var errorStack []error
	isActiveSession := true
	s, err := cmd.GetSession(choice)
	if err != nil {
		isActiveSession = false
		errorStack = append(errorStack, err)
	}
	if !isActiveSession {
		p, err := filepath.Abs(choice)
		if err != nil {
			errorStack = append(errorStack, err)
			p = choice
		}
		info, err := os.Stat(p)
		if err != nil {
			errorStack = append(errorStack, err)
			return fmt.Errorf(
				"unable to connect to %q: %w",
				choice,
				errors.Join(errorStack...),
			)
		}
		if !info.IsDir() {
			errorStack = append(
				errorStack,
				fmt.Errorf("%q found but is not a directory", p),
			)
			return errors.Join(errorStack...)
		}
		s = tmux.TmuxSession{
			Name:     filepath.Base(p),
			Path:     p,
			Attached: 0,
		}
	}

	// Add the path to zoxide.
	if err = zoxide.Add(s.Path); err != nil {
		return fmt.Errorf("unable to connect to %q: %w", choice, err)
	}

	// Connect to the tmux session.
	return tmux.Connect(s, alwaysSwitch, command, s.Path, config)
}
