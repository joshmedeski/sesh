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

	// When invoked from outside tmux (e.g. osascript/browser dispatcher),
	// `switch-client -t` has no current client to act on. Find the active
	// client, switch it explicitly, and bring the terminal to the front.
	if opts.Switch && !c.tmux.IsAttached() {
		return c.connectDetached(connection.Session.Name)
	}

	return c.tmux.SwitchOrAttach(connection.Session.Name, opts)
}

func (c *RealConnector) connectDetached(sessionName string) (string, error) {
	client := c.tmux.ResolveClient()
	if client != "" {
		c.tmux.SwitchClientTarget(client, sessionName)
	}
	c.focuser.Activate(c.config.Terminal)
	if client != "" {
		return fmt.Sprintf("switching to tmux session: %s", sessionName), nil
	}
	return fmt.Sprintf("connected to tmux session: %s", sessionName), nil
}
