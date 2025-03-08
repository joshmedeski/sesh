package connector

import (
	"github.com/joshmedeski/sesh/v2/model"
)

func tmuxinatorStrategy(c *RealConnector, name string) (model.Connection, error) {
	session, exists := c.lister.FindTmuxinatorConfig(name)
	if !exists {
		return model.Connection{Found: false}, nil
	}

	return model.Connection{
		Found:       true,
		Session:     session,
		New:         true,
		AddToZoxide: false,
	}, nil
}

func connectToTmuxinator(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	return c.tmuxinator.Start(connection.Session.Name)
}
