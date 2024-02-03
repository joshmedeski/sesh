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

func Connect(
	choice string,
	alwaysSwitch bool,
	command string,
	config *config.Config,
) error {
	t, err := tmux.NewCommand()
	if err != nil {
		return err
	}

	var errorStack []error
	isActiveSession := true
	s, err := t.GetSession(choice)
	if err != nil {
		isActiveSession = false
		errorStack = append(errorStack, err)
	}
	sessionName, sessionPath := s.Name(), s.Path()
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
		sessionName = filepath.Base(p)
		sessionPath = p
	}

	if err = zoxide.Add(s.Path()); err != nil {
		return fmt.Errorf("unable to connect to %q: %w", choice, err)
	}

	return t.Connect(sessionName, alwaysSwitch, command, sessionPath, config)
}
