package home

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/oswrap"
)

type Home interface {
	ShortenHome(path string) (string, error)
	ExpandHome(path string) (string, error)
}

type RealHome struct {
	os oswrap.Os
}

func NewHome(os oswrap.Os) Home {
	return &RealHome{os}
}

func (p *RealHome) ShortenHome(path string) (string, error) {
	home, err := p.os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1), nil
	}

	return path, nil
}

func (p *RealHome) ExpandHome(path string) (string, error) {
	home, err := p.os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(path, "~") {
		return strings.Replace(path, "~", home, 1), nil
	}
	return path, nil
}
