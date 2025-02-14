package connector

import (
	"github.com/joshmedeski/sesh/v2/model"
)

func dirStrategy(c *RealConnector, name string) (model.Connection, error) {
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
			// TODO: what is the best name for this? "dir" isn't technically a source
			// it's not used in any list command
			Src:  "dir",
			Name: nameFromPath,
			Path: absPath,
		},
	}, nil
}
