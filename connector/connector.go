package connector

import (
	"github.com/joshmedeski/sesh/dir"
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/namer"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

type Connector interface {
	Connect(name string, opts model.ConnectOpts) (string, error)
}

type RealConnector struct {
	config model.Config
	dir    dir.Dir
	home   home.Home
	lister lister.Lister
	namer  namer.Namer
	tmux   tmux.Tmux
	zoxide zoxide.Zoxide
}

func NewConnector(
	config model.Config,
	dir dir.Dir,
	home home.Home,
	lister lister.Lister,
	namer namer.Namer,
	tmux tmux.Tmux,
	zoxide zoxide.Zoxide,
) Connector {
	return &RealConnector{
		config,
		dir,
		home,
		lister,
		namer,
		tmux,
		zoxide,
	}
}
