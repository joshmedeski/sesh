package connector

import (
	"github.com/joshmedeski/sesh/v2/model"
)

func configStrategy(c *RealConnector, name string) (model.Connection, error) {
	config, exists := c.lister.FindConfigSession(name)
	if !exists {
		return model.Connection{Found: false}, nil
	}

	return model.Connection{
		Found:       true,
		Session:     config,
		New:         true,
		AddToZoxide: true,
	}, nil
}
