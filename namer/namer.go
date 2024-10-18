package namer

import (
	"fmt"

	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/pathwrap"
)

type Namer interface {
	// Names a sesh session from a given path
	Name(path string) (string, error)
	RootName(path string) (string, error)
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

func (n *RealNamer) Name(path string) (string, error) {
	path, err := n.pathwrap.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	strategies := []func(*RealNamer, string) (string, error){
		gitBareName,
		gitName,
		dirName,
	}
	for _, strategy := range strategies {
		name, err := strategy(n, path)
		if err != nil {
			return "", err
		}
		if name != "" {
			return convertToValidName(name), nil
		}
	}
	return "", fmt.Errorf("could not determine name from path: %s", path)
}

func (n *RealNamer) RootName(path string) (string, error) {
	path, err := n.pathwrap.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	strategies := []func(*RealNamer, string) (string, error){
		gitBareRootName,
		gitRootName,
		dirName,
	}
	for _, strategy := range strategies {
		name, err := strategy(n, path)
		if err != nil {
			return "", err
		}
		if name != "" {
			return convertToValidName(name), nil
		}
	}
	return "", fmt.Errorf("could not determine root name from path: %s", path)
}
