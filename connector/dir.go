package connector

import "github.com/joshmedeski/sesh/model"

func dirStrategy(c *RealConnector, name string) (model.Connection, error) {
	path, err := c.home.ExpandHome(name)
	if err != nil {
		return model.Connection{}, err
	}
	isDir, absPath := c.dir.Dir(path)
	if !isDir {
		return model.Connection{Found: false}, nil
	}
	return model.Connection{
		Found:       true,
		New:         true,
		AddToZoxide: true,
		Session: model.SeshSession{
			Src:  "zoxide",
			Name: name,
			Path: absPath,
		},
	}, nil
}
