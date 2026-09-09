package connector

import (
	"github.com/joshmedeski/sesh/v2/model"
)

func configWildcardStrategy(c *RealConnector, name string) (model.Connection, error) {
	path, err := c.home.ExpandPath(name)
	if err != nil {
		return model.Connection{}, err
	}

	isDir, absPath := c.dir.Dir(path)
	if !isDir {
		return model.Connection{Found: false}, nil
	}

	wc, found := c.lister.FindConfigWildcard(absPath)
	if !found {
		return model.Connection{Found: false}, nil
	}

	nameFromPath, err := c.namer.Name(absPath)
	if err != nil {
		return model.Connection{}, err
	}

	// A wildcard match describes how to create the session, not that it is
	// missing. Without this check the create path runs again on every reconnect,
	// sending the startup command into a session that is already up.
	if existing, ok := c.lister.FindTmuxSessionByBase(nameFromPath); ok {
		return model.Connection{
			Found:       true,
			New:         false,
			AddToZoxide: true,
			Session:     existing,
		}, nil
	}

	return model.Connection{
		Found:       true,
		New:         true,
		AddToZoxide: true,
		Session: model.SeshSession{
			Src:                   "config_wildcard",
			Name:                  nameFromPath,
			Path:                  absPath,
			WindowNames:           wc.Windows,
			DisableStartupCommand: wc.DisableStartCommand,
		},
	}, nil
}
