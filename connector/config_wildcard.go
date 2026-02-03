package connector

import (
	"github.com/joshmedeski/sesh/v2/model"
)

func configWildcardStrategy(c *RealConnector, name string) (model.Connection, error) {
	_, found := c.lister.FindConfigWildcard(name)
	if !found {
		return model.Connection{Found: false}, nil
	}

	path, err := c.home.ExpandHome(name)
	if err != nil {
		return model.Connection{}, err
	}

	isDir, absPath := c.dir.Dir(path)
	if !isDir {
		return model.Connection{Found: false}, nil
	}

	nameFromPath, err := c.namer.Name(absPath)
	if err != nil {
		return model.Connection{}, err
	}

	return model.Connection{
		Found:       true,
		New:         true,
		AddToZoxide: true,
		Session: model.SeshSession{
			Src:  "config_wildcard",
			Name: nameFromPath,
			Path: absPath,
		},
	}, nil
}
