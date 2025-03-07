package connector

import (
	"github.com/joshmedeski/sesh/v2/model"
)

func configWildcardStrategy(c *RealConnector, name string) (model.Connection, error) {
	wildcard, exists := c.lister.FindConfigWildcard(name)
	if !exists {
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
			Src:                   "config",
			Path:                  absPath,
			Name:                  nameFromPath,
			StartupCommand:        wildcard.StartupCommand,
			PreviewCommand:        wildcard.PreviewCommand,
			DisableStartupCommand: wildcard.DisableStartupCommand,
		},
	}, nil
}
