package namer

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/pathwrap"
)

type Namer interface {
	// Names a sesh session from a given path
	Name(path string) (string, error)
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

func convertToValidName(name string) string {
	validName := strings.ReplaceAll(name, ".", "_")
	validName = strings.ReplaceAll(validName, ":", "_")
	return validName
}

func (n *RealNamer) Name(path string) (string, error) {
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
