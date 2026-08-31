package connector

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func tmuxStrategy(c *RealConnector, name string) (model.Connection, error) {
	session, exists := c.lister.FindTmuxSession(name)
	if !exists {
		return model.Connection{Found: false}, nil
	}
	return model.Connection{
		Found:       true,
		Session:     session,
		New:         false,
		AddToZoxide: true,
	}, nil
}

func connectToTmux(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	if connection.New {
		c.tmux.NewSession(connection.Session.Name, connection.Session.Path)
		if opts.Command != "" {
			c.tmux.SendKeys(connection.Session.Name, opts.Command)
		} else {
			c.startup.Exec(connection.Session)
		}
	}

	return c.focusSession(connection.Session.Name, opts.Switch)
}

func (c *RealConnector) connectDetached(sessionName string) (string, error) {
	switched := false
	if clients, err := c.tmux.ListClients(); err == nil {
		for _, client := range clients {
			if client != "" {
				c.tmux.SwitchClientTarget(client, sessionName)
				switched = true
				break
			}
		}
	}
	c.focuser.Activate(c.config.Terminal)
	if switched {
		return fmt.Sprintf("switching to tmux session: %s", sessionName), nil
	}
	return fmt.Sprintf("connected to tmux session: %s", sessionName), nil
}
