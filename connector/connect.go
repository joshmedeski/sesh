package connector

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

// TODO: send to logging (local txt file?)
func (c *RealConnector) Connect(name string, opts model.ConnectOpts) (string, error) {
	// TODO: make it configurable to change the order of connection establishments?
	// ["tmux", "config", "dir", "zoxide"]
	// TODO: make it configurable to disable certain strategies (including flags for optimized fzf commands)
	// sesh connect --config (sesh list --config | fzf)
	strategies := []func(*RealConnector, string) (model.Connection, error){
		tmuxStrategy,
		tmuxinatorStrategy,
		configStrategy,
		dirStrategy,
		zoxideStrategy,
	}

	connectStrategy := map[string]func(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error){
		"tmux":       connectToTmux,
		"tmuxinator": connectToTmuxinator,
		"config":     connectToTmux,
		"dir":        connectToTmux,
		"zoxide":     connectToTmux,
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
			return connectStrategy[connection.Session.Src](c, connection, opts)
		}
	}

	return "", fmt.Errorf("no connection found for '%s'", name)
}
