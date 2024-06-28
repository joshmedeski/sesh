package namer

import (
	"strings"

	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/pathwrap"
)

type Namer interface {
	// Names a sesh session from a given path
	FromPath(path string) (string, error)
}

type RealNamer struct {
	pathwrap pathwrap.Path
	git      git.Git
}

func NewNamer(pathwrap pathwrap.Path, git git.Git) Namer {
	return &RealNamer{
		pathwrap: pathwrap,
		git:      git,
	}
}

func (n *RealNamer) FromPath(path string) (string, error) {
	isGit, topLevelDir, _ := n.git.ShowTopLevel(path)
	if isGit && topLevelDir != "" {
		relativePath := strings.TrimPrefix(path, topLevelDir)
		baseDir := n.pathwrap.Base(topLevelDir)
		name := baseDir + relativePath
		return name, nil
	}
	return n.pathwrap.Base(path), nil
}
