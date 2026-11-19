package home

import (
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/v2/oswrap"
)

type Home interface {
	ShortenHome(path string) (string, error)
	ExpandPath(path string) (string, error)
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

	if rel, err := filepath.Rel(home, path); err == nil && filepath.IsLocal(rel) {
		return filepath.Join("~", rel), nil
	}

	return path, nil
}

func (p *RealHome) ExpandPath(path string) (string, error) {
	expanded := p.os.ExpandEnv(path)
	home, err := p.os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(expanded, "~") {
		return strings.Replace(expanded, "~", home, 1), nil
	}
	return expanded, nil
}
