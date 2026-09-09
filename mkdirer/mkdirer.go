package mkdirer

import (
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
)

type Mkdirer interface {
	// Creates a directory (relative or absolute path) if it doesn't already
	// exist, then connects to it as a session.
	Mkdir(path string, opts model.ConnectOpts) (string, error)
}

type RealMkdirer struct {
	os        oswrap.Os
	home      home.Home
	connector connector.Connector
}

func NewMkdirer(os oswrap.Os, home home.Home, connector connector.Connector) Mkdirer {
	return &RealMkdirer{
		os:        os,
		home:      home,
		connector: connector,
	}
}

func (m *RealMkdirer) Mkdir(path string, opts model.ConnectOpts) (string, error) {
	expandedPath, err := m.home.ExpandPath(path)
	if err != nil {
		return "", err
	}

	if err := m.os.MkdirAll(expandedPath, 0o755); err != nil {
		return "", err
	}

	return m.connector.Connect(expandedPath, opts)
}
