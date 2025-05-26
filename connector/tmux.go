package connector

import (
	"regexp"
	
	"github.com/joshmedeski/sesh/v2/model"
)

func tmuxStrategy(c *RealConnector, name string) (model.Connection, error) {
	// Check if this is a marked window format: session:window_name(window_number)
	re := regexp.MustCompile(`^(.+):(.+)\((\d+)\)$`)
	if matches := re.FindStringSubmatch(name); len(matches) == 4 {
		sessionName := matches[1]
		windowNumber := matches[3]
		
		// Look for the actual session
		session, exists := c.lister.FindTmuxSession(sessionName)
		if !exists {
			return model.Connection{Found: false}, nil
		}
		
		// Modify the session name to include the window number for connection
		session.Name = sessionName + ":" + windowNumber
		
		return model.Connection{
			Found:       true,
			Session:     session,
			New:         false,
			AddToZoxide: true,
		}, nil
	}
	
	// Normal session lookup
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
