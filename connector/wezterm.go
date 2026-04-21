package connector

import "github.com/joshmedeski/sesh/v2/model"

func weztermStrategy(c *RealConnector, name string) (model.Connection, error) {
	session, exists := c.lister.FindWeztermWorkspace(name)
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

// connectToWezterm creates a WezTerm workspace via CLI. Note: the CLI cannot
// switch to a workspace -- the companion Lua plugin handles that step.
func connectToWezterm(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	if connection.New {
		if _, err := c.wezterm.SpawnWorkspace(connection.Session.Name, connection.Session.Path); err != nil {
			return "", err
		}
		if opts.Command != "" {
			// Get the pane list and send the command to the first pane of this workspace.
			panes, err := c.wezterm.ListAllPanes()
			if err == nil {
				for _, p := range panes {
					if p.Workspace == connection.Session.Name {
						c.wezterm.SendText(p.PaneID, opts.Command)
						break
					}
				}
			}
		} else {
			c.startup.Exec(connection.Session)
		}
	}
	return connection.Session.Name, nil
}
