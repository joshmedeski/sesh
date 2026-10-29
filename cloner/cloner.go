package cloner

import (
	"fmt"
	"os"
	"strings"

	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/jj"
	"github.com/joshmedeski/sesh/v2/model"
)

type Cloner interface {
	// Clones a git repository
	Clone(opts model.GitCloneOptions) (string, error)
}

type RealCloner struct {
	connector connector.Connector
	git       git.Git
	jj        jj.Jj
}

func NewCloner(connector connector.Connector, git git.Git, jj jj.Jj) Cloner {
	return &RealCloner{
		connector: connector,
		git:       git,
		jj:        jj,
	}
}

func (c *RealCloner) Clone(opts model.GitCloneOptions) (string, error) {
	path := getPath(opts)
	switch opts.Vcs {
	case "jj":
		if err := c.jj.Clone(opts.Repo, path, opts.GitFlags...); err != nil {
			return "", err
		}
	case "", "git":
		if _, err := c.git.Clone(opts.Repo, opts.CmdDir, opts.Dir, opts.GitFlags...); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported vcs %q: use git or jj", opts.Vcs)
	}

	newOpts := model.ConnectOpts{}
	if _, err := c.connector.Connect(path, newOpts); err != nil {
		return "", err
	}

	return "", nil

}

func getPath(opts model.GitCloneOptions) string {
	var path string
	if opts.CmdDir != "" {
		path = opts.CmdDir
	} else {
		path, _ = os.Getwd()
	}

	if opts.Dir != "" {
		path = path + "/" + opts.Dir
	} else {
		repoName := getRepoName(opts.Repo)
		path = path + "/" + repoName
	}
	return path
}

func getRepoName(url string) string {
	parts := strings.Split(url, "/")
	lastPart := parts[len(parts)-1]
	repoName := strings.TrimSuffix(lastPart, ".git")
	return repoName
}
