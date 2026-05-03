package connector

import "github.com/Wingsdh/cc-sesh/v2/model"

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
		// Resolve the startup command BEFORE creating the session so we can
		// inject it as the pane's initial shell-command. This eliminates the
		// send-keys race with slow shell init
		var rawCmd string
		if opts.Command != "" {
			rawCmd = opts.Command
		} else {
			resolved, err := c.startup.ResolveCommand(connection.Session)
			if err != nil {
				return "", err
			}
			rawCmd = resolved
		}
		shellCmd := c.startup.WrapForShell(rawCmd)
		if _, err := c.tmux.NewSession(connection.Session.Name, connection.Session.Path, shellCmd); err != nil {
			return "", err
		}
		if opts.Command == "" {
			if _, err := c.startup.Exec(connection.Session); err != nil {
				return "", err
			}
		}
	}
	return c.tmux.SwitchOrAttach(connection.Session.Name, opts)
}
