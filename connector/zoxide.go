package connector

import "github.com/joshmedeski/sesh/v2/model"

func zoxideToTmuxName(c *RealConnector, path string) (string, error) {
	fullPath, err := c.home.ExpandHome(path)
	if err != nil {
		return "", err
	}
	name, err := c.namer.Name(fullPath)
	if err != nil {
		return "", err
	}
	return name, nil
}

func zoxideStrategy(c *RealConnector, path string) (model.Connection, error) {
	session, exists := c.lister.FindZoxideSession(path)
	if !exists {
		return model.Connection{Found: false}, nil
	}
	name, err := c.namer.Name(session.Path)
	if err != nil {
		return model.Connection{}, err
	}
	session.Name = name
	return model.Connection{
		Found:       true,
		Session:     session,
		New:         true,
		AddToZoxide: true,
	}, nil
}
