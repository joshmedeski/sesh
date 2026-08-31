package connector

import (
	"errors"
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

// targetSession resolves the session a window or pane connect should land in.
// An empty name means the attached session, which is what makes `sesh window
// connect notes` work with no flags from inside tmux.
func (c *RealConnector) targetSession(name string) (model.SeshSession, error) {
	if name != "" {
		return c.ensureSession(name)
	}
	session, exists := c.lister.GetAttachedTmuxSession()
	if !exists {
		return model.SeshSession{}, errors.New("not inside a tmux session, use --target to specify one")
	}
	return session, nil
}

// ensureSession resolves a name to a running tmux session, starting it when it
// isn't running yet. Unlike Connect it never moves the client, so a background
// window or pane connect can target a session that doesn't exist yet without
// yanking the user out of the one they're in.
func (c *RealConnector) ensureSession(name string) (model.SeshSession, error) {
	connection, err := c.Resolve(name)
	if err != nil {
		return model.SeshSession{}, err
	}
	if !connection.Found {
		return model.SeshSession{}, fmt.Errorf("no session found for '%s'", name)
	}
	if connection.Session.Src == "tmux-pane" {
		return model.SeshSession{}, fmt.Errorf("'%s' is a pane, not a session", name)
	}
	if connection.AddToZoxide {
		c.zoxide.Add(connection.Session.Path)
	}
	if !connection.New {
		return connection.Session, nil
	}
	if connection.Session.Src == "tmuxinator" {
		if _, err := c.tmuxinator.Start(connection.Session.Name); err != nil {
			return model.SeshSession{}, fmt.Errorf("failed to start tmuxinator session: %w", err)
		}
		return connection.Session, nil
	}
	if _, err := c.tmux.NewSession(connection.Session.Name, connection.Session.Path); err != nil {
		return model.SeshSession{}, fmt.Errorf("failed to create tmux session: %w", err)
	}
	c.startup.Exec(connection.Session)
	return connection.Session, nil
}

// focusSession brings the client to a session, honouring the same
// switch-vs-attach rules as Connect.
func (c *RealConnector) focusSession(sessionName string, doSwitch bool) (string, error) {
	// When invoked from outside tmux (e.g. osascript/browser dispatcher),
	// `switch-client -t` has no current client to act on. Find the active
	// client, switch it explicitly, and bring the terminal to the front.
	if doSwitch && !c.tmux.IsAttached() {
		return c.connectDetached(sessionName)
	}
	return c.tmux.SwitchOrAttach(sessionName, model.ConnectOpts{Switch: doSwitch})
}
