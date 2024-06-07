package connector

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
)

func establishConfigConnection(c *RealConnector, name string, opts model.ConnectOpts) (string, error) {
	session, exists := c.lister.FindConfigSession(name)
	if !exists {
		return "", nil
	}
	if session.Path != "" {
		return "", fmt.Errorf("found config session '%s' has no path", name)
	}
	// TODO: run startup command or startup script
	c.tmux.NewSession(session.Name, session.Path)
	c.zoxide.Add(session.Path)
	return c.tmux.SwitchOrAttach(name, opts)
}

func establishDirConnection(c *RealConnector, name string, _ model.ConnectOpts) (string, error) {
	isDir, absPath := c.dir.Dir(name)
	if !isDir {
		return "", nil
	}
	// TODO: get session name from directory
	// c.tmux.NewSession(session.Name, absPath)
	// c.zoxide.Add(session.Path)
	// return switchOrAttach(c, name, opts)
	return absPath, nil
}

func establishZoxideConnection(c *RealConnector, name string, _ model.ConnectOpts) (string, error) {
	isDir, absPath := c.dir.Dir(name)
	if !isDir {
		return "", nil
	}
	// TODO: get session name from directory
	// c.tmux.NewSession(session.Name, absPath)
	// c.zoxide.Add(session.Path)
	// return switchOrAttach(c, name, opts)
	return absPath, nil
}

// TODO: send to logging (local txt file?)
func (c *RealConnector) Connect(name string, opts model.ConnectOpts) (string, error) {
	// TODO: make it configurable to change the order of connection establishments?
	// ["tmux", "config", "dir", "zoxide"]
	// TODO: make it configurable to disable certain strategies (including flags for optimized fzf commands)
	// sesh connect --config (sesh list --config | fzf)
	strategies := []func(*RealConnector, string) (model.Connection, error){
		tmuxStrategy,
		configStrategy,
		dirStrategy,
		// establishZoxideConnection,
	}

	for _, strategy := range strategies {
		if connection, err := strategy(c, name); err != nil {
			return "", fmt.Errorf("failed to establish connection: %w", err)
		} else if connection.Found {
			// TODO: allow CLI flag to disable zoxide and overwrite all settings?
			// sesh connect --ignore-zoxide "dotfiles"
			if connection.AddToZoxide {
				c.zoxide.Add(connection.Session.Path)
			}
			if connection.New {
				c.tmux.NewSession(connection.Session.Name, connection.Session.Path)
			}
			// TODO: configure the ability to create a session in a detached way (like update)
			// TODO: configure the ability to create a popup instead of switching
			return c.tmux.SwitchOrAttach(connection.Session.Name, opts)
		}
	}

	return "", fmt.Errorf("no connection found for '%s'", name)
}
