package connector

import (
	"fmt"

	"github.com/joshmedeski/sesh/lister"
)

func establishMultiplexerConnection(c *RealConnector, name string, opts ConnectOpts) (string, error) {
	sessions, err := c.lister.List(lister.ListOptions{Tmux: true})
	if err != nil {
		return "", fmt.Errorf("determine tmux connection failed: %w", err)
	}
	for _, session := range sessions {
		if session.Name == name {
			// TODO: make this more robust (switch when applicable)
			c.tmux.AttachSession(name)
			return fmt.Sprintf("determine tmux connection succeeded: %s", name), nil
		}
	}
	return "", nil
}

type ConnectOpts struct {
	AlwaysSwitch bool
	Command      string
}

func (c *RealConnector) Connect(name string, opts ConnectOpts) (string, error) {
	if connected, err := establishMultiplexerConnection(c, name, opts); err != nil {
		// TODO: send to logging (local txt file?)
		return "", fmt.Errorf("failed to connect: %w", err)
	} else if connected != "" {
		return connected, nil
	}

	// TODO: if name is config session, create session from config
	// TODO: if name is directory, create session from directory
	// TODO: if name matches zoxide result, create session from result
	return "connect", nil
}
