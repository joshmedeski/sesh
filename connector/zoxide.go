package connector

import "github.com/joshmedeski/sesh/model"

func zoxideToTmuxName(c *RealConnector, path string) (string, error) {
	fullPath, err := c.home.ExpandHome(path)
	if err != nil {
		return "", err
	}
	nameFromPath, err := c.namer.FromPath(fullPath)
	if err != nil {
		return "", err
	}
	return nameFromPath, nil
}

func zoxideStrategy(c *RealConnector, path string) (model.Connection, error) {
	session, exists := c.lister.FindZoxideSession(path)
	if !exists {
		return model.Connection{Found: false}, nil
	}
	nameFromPath, err := c.namer.FromPath(session.Path)
	if err != nil {
		return model.Connection{}, err
	}
	session.Name = nameFromPath
	return model.Connection{
		Found:       true,
		Session:     session,
		New:         true,
		AddToZoxide: true,
	}, nil
}
