package tmux

import (
	"fmt"
	"log"

	"github.com/joshmedeski/sesh/config"
)

func Connect(
	s TmuxSession,
	alwaysSwitch bool,
	command string,
	sessionPath string,
	config *config.Config,
) error {
	session, _ := FindSession(s.Name)
	// TODO: load tmup if exists
	if session == nil {
		_, err := NewSession(s)
		if err != nil {
			return fmt.Errorf(
				"error when creating new tmux session %q: %w",
				s.Name,
				err,
			)
		}
		if command != "" {
			runPersistentCommand(s.Name, command)
		} else if startupScript := getStartupScript(sessionPath, config); startupScript != "" {
			err := execStartupScript(s.Name, startupScript)
			if err != nil {
				log.Fatal(err)
			}
		} else if startupCommand := getStartupCommand(sessionPath, config); startupCommand != "" {
			err := execStartupCommand(s.Name, startupCommand)
			if err != nil {
				log.Fatal(err)
			}
		} else if config.DefaultSessionConfig.StartupCommand != "" {
			err := execStartupCommand(s.Name, config.DefaultSessionConfig.StartupCommand)
			if err != nil {
				log.Fatal(err)
			}
		} else if config.DefaultSessionConfig.StartupScript != "" {
			err := execStartupScript(s.Name, config.DefaultSessionConfig.StartupScript)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	isAttached := isAttached()
	if isAttached || alwaysSwitch {
		switchSession(s.Name)
	} else {
		attachSession(s.Name)
	}
	return nil
}
