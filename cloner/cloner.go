package cloner

import (
	// "fmt"

	"github.com/joshmedeski/sesh/connector"
	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/model"
)

type Cloner interface {
	// Clones a git repository
	Clone(opts model.GitCloneOptions) (string, error)
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

func (c *RealCloner) Clone(opts model.GitCloneOptions) (string, error) {
	if _, err := c.git.Clone(opts.Repo, opts.CmdDir, opts.Dir); err != nil {
		return "", err
	} else {
		return "", nil
	}

	// TODO: get name of directory
	// TODO: connect to that directory
}
