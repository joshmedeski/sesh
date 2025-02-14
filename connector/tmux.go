package connector

import "github.com/joshmedeski/sesh/v2/model"

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
		c.startup.Exec(connection.Session)
	}
	return c.tmux.SwitchOrAttach(connection.Session.Name, opts)
}
