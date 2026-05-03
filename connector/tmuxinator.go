package connector

import (
	"fmt"

	"github.com/Wingsdh/cc-sesh/v2/model"
)

func tmuxinatorStrategy(c *RealConnector, name string) (model.Connection, error) {
	session, exists := c.lister.FindTmuxinatorConfig(name)
	if !exists {
		return model.Connection{Found: false}, nil
	}

	return model.Connection{
		Found:       true,
		Session:     session,
		New:         true,
		AddToZoxide: false,
	}, nil
}

func connectToTmuxinator(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	if _, err := c.tmuxinator.Start(connection.Session.Name); err != nil {
		return "", fmt.Errorf("failed to start tmuxinator session: %w", err)
	}
	return c.tmux.SwitchOrAttach(connection.Session.Name, opts)
}
