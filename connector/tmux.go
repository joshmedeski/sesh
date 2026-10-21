package connector

import (
	"errors"
	"fmt"
	"log/slog"

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

// showStartupError surfaces startup errors in the tmux status line.
func (c *RealConnector) showStartupError(err error) {
	_, _ = c.tmux.DisplayMessage("sesh: " + err.Error())
}

func connectToTmux(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	var startupErr error
	if connection.New {
		// Everything below assumes the session now exists, so a failure here
		// makes the rest meaningless — and it is how "tmux isn't on the PATH"
		// first shows up when sesh is launched from a GUI app.
		if _, err := c.tmux.NewSession(connection.Session.Name, connection.Session.Path); err != nil {
			return "", fmt.Errorf("failed to create tmux session %q: %w", connection.Session.Name, err)
		}
		if opts.Command != "" {
			_, _ = c.tmux.SendKeys(connection.Session.Name+":", opts.Command)
		} else {
			_, startupErr = c.startup.Exec(connection.Session)
		}
	}

	message, err := c.focusSession(connection.Session.Name, opts.Switch)
	if startupErr == nil {
		return message, err
	}
	if err != nil {
		return message, errors.Join(err, startupErr)
	}
	c.showStartupError(startupErr)
	return message, fmt.Errorf("session %s startup failed: %w", connection.Session.Name, startupErr)
}

func (c *RealConnector) connectDetached(sessionName string) (string, error) {
	// Reaching tmux is checked separately from picking a client, because the two
	// failures need opposite handling and look identical from ResolveClient. A
	// GUI launcher (Leader Key, Raycast, an osascript binding) runs sesh with a
	// bare PATH that Homebrew isn't on, so every tmux call fails while
	// Activate's osascript still works — sesh brings the terminal forward and
	// silently changes nothing, which is indistinguishable from success.
	clients, err := c.tmux.ListClients()
	if err != nil {
		return "", fmt.Errorf("couldn't reach tmux to switch to '%s': %w", sessionName, err)
	}

	// Nothing attached, so there is no client to move. Bringing the terminal
	// forward is still the useful thing to do: the user lands where they can
	// attach.
	if len(clients) == 0 {
		_, _ = c.focuser.Activate(c.config.Terminal)
		return fmt.Sprintf("connected to tmux session: %s", sessionName), nil
	}

	client := c.tmux.ResolveClient()
	if client != "" {
		_, err = c.tmux.SwitchClientTarget(client, sessionName)
	} else {
		_, err = c.tmux.SwitchClient(sessionName)
	}

	// Worth focusing either way: on failure the user at least lands on the tmux
	// they need to sort out.
	_, _ = c.focuser.Activate(c.config.Terminal)

	if err != nil {
		slog.Error("failed to switch client", "client", client, "session", sessionName, "error", err)
		return "", fmt.Errorf("failed to switch to tmux session '%s': %w", sessionName, err)
	}
	return fmt.Sprintf("switching to tmux session: %s", sessionName), nil
}
