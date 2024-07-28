package connect

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/icons"
	"github.com/joshmedeski/sesh/session"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

func Connect(
	choice string,
	alwaysSwitch bool,
	command string,
	config *config.Config,
) error {
	if strings.HasPrefix(choice, icons.TmuxIcon) || strings.HasPrefix(choice, icons.ZoxideIcon) || strings.HasPrefix(choice, icons.ConfigIcon) {
		choice = choice[4:]
	}

	session, err := session.Determine(choice, config)
	if err != nil {
		return fmt.Errorf("unable to connect to %q: %w", choice, err)
	}

	if err = zoxide.Add(session.Path); err != nil {
		return fmt.Errorf("unable to connect to %q: %w", choice, err)
	}
	return tmux.Connect(tmux.TmuxSession{
		Name:     session.Name,
		Path:     session.Path,
		PathList: session.PathList,
	}, alwaysSwitch, command, session.Path, config)
}
