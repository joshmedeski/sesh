package cloner

import (
	"github.com/joshmedeski/sesh/connector"
	"github.com/joshmedeski/sesh/git"
)

type Cloner interface {
	// Clones a git repository
	Clone(path string) (string, error)
}

type RealCloner struct {
	connector connector.Connector
	git       git.Git
}

func NewCloner(connector connector.Connector, git git.Git) Cloner {
	return &RealCloner{
		connector: connector,
		git:       git,
	}
}

func (c *RealCloner) Clone(path string) (string, error) {
	// TODO: clone
	// TODO: get name of directory
	// TODO: connect to that directory
	return "", nil
}
