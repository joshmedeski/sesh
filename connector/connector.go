package connector

import (
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
)

type Connector interface {
	Connect(name string, opts ConnectOpts) (string, error)
}

type RealConnector struct {
	config model.Config
	home   home.Home
	lister lister.Lister
	tmux   tmux.Tmux
}

func NewConnector(config model.Config, home home.Home, lister lister.Lister, tmux tmux.Tmux) Connector {
	return &RealConnector{config, home, lister, tmux}
}
