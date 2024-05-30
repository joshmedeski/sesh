package connector

import "github.com/joshmedeski/sesh/model"

func tmuxStrategy(c *RealConnector, name string) (model.Connection, error) {
	// TODO: find by name or by path?
	session, exists := c.lister.FindTmuxSession(name)
	if !exists {
		return model.Connection{
			Found: false,
		}, nil
	}

	// TODO: make zoxide add configurable

	return model.Connection{
		Found:       true,
		Session:     session,
		New:         false,
		AddToZoxide: true,
		// Switch: true
	}, nil
}
