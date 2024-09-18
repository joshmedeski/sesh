package connector

import (
	"github.com/joshmedeski/sesh/model"
)

func tmuxinatorStrategy(c *RealConnector, name string) (model.Connection, error) {
	session, exists := c.lister.FindTmuxinatorConfig(name)
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
