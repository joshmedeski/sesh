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
		for _, window := range connection.Session.WindowConfigs {
			if ret, err := c.tmux.NewWindow(window.Path); err != nil {
				return ret, err
			}

			if ret, err := c.tmux.SendKeys(connection.Session.Name, window.StartupScript); err != nil {
				return ret, err
			}
		}
		c.tmux.NextWindow()
	}
	return c.tmux.SwitchOrAttach(connection.Session.Name, opts)
}
