package path

import (
	"strings"

	"github.com/joshmedeski/sesh/oswrap"
)

type Path interface {
	ShortenHome(path string) (string, error)
	ExpandHome(path string) (string, error)
}

type RealPath struct {
	os oswrap.Os
}

func NewPath(os oswrap.Os) Path {
	return &RealPath{os}
}

func (p *RealPath) ShortenHome(path string) (string, error) {
	home, err := p.os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1), nil
	}

	return path, nil
}

func (p *RealPath) ExpandHome(path string) (string, error) {
	home, err := p.os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(path, "~") {
		return strings.Replace(path, "~", home, 1), nil
	}
	return path, nil
}
