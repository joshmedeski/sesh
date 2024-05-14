package connector

import (
	"fmt"

	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/model"
)

func switchOrAttach(c *RealConnector, name string, opts model.ConnectOpts) (string, error) {
	if opts.Switch || c.tmux.IsAttached() {
		if _, err := c.tmux.SwitchClient(name); err != nil {
			return "", fmt.Errorf("failed to switch to tmux session: %w", err)
		} else {
			return fmt.Sprintf("switching to existing tmux session: %s", name), nil
		}
	} else {
		if _, err := c.tmux.AttachSession(name); err != nil {
			return "", fmt.Errorf("failed to attach to tmux session: %w", err)
		} else {
			return fmt.Sprintf("attaching to existing tmux session: %s", name), nil
		}
	}
}

func establishTmuxConnection(c *RealConnector, name string, opts model.ConnectOpts) (string, error) {
	session, exists := c.lister.FindTmuxSession(name)
	if !exists {
		return "", nil
	}
	return switchOrAttach(c, session.Name, opts)
}

func establishConfigConnection(c *RealConnector, name string, opts model.ConnectOpts) (string, error) {
	sessions, err := c.lister.List(lister.ListOptions{Config: true})
	if err != nil {
		return "", err
	}
	for _, session := range sessions {
		if session.Name == name {
			if session.Path != "" {
				return "", fmt.Errorf("found config session '%s' has no path", name)
			}
			c.tmux.NewSession(session.Name, session.Path)
			switchOrAttach(c, name, opts)
		}
	}
	return "", nil // no tmux connection was established
}

func (c *RealConnector) Connect(name string, opts model.ConnectOpts) (string, error) {
	if tmuxConnected, err := establishTmuxConnection(c, name, opts); err != nil {
		// TODO: send to logging (local txt file?)
		return "", fmt.Errorf("failed to establish tmux connection: %w", err)
	} else if tmuxConnected != "" {
		return tmuxConnected, nil
	}

	// TODO: if name is config session, create session from config

	// TODO: if name is directory, create session from directory

	// TODO: if name matches zoxide result, create session from result
	return "connect", nil
}
