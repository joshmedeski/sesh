package connector

import (
	"github.com/Wingsdh/cc-sesh/v2/dir"
	"github.com/Wingsdh/cc-sesh/v2/home"
	"github.com/Wingsdh/cc-sesh/v2/lister"
	"github.com/Wingsdh/cc-sesh/v2/model"
	"github.com/Wingsdh/cc-sesh/v2/namer"
	"github.com/Wingsdh/cc-sesh/v2/startup"
	"github.com/Wingsdh/cc-sesh/v2/tmux"
	"github.com/Wingsdh/cc-sesh/v2/tmuxinator"
	"github.com/Wingsdh/cc-sesh/v2/zoxide"
)

type Connector interface {
	Connect(name string, opts model.ConnectOpts) (string, error)
}

type RealConnector struct {
	config     model.Config
	dir        dir.Dir
	home       home.Home
	lister     lister.Lister
	namer      namer.Namer
	startup    startup.Startup
	tmux       tmux.Tmux
	zoxide     zoxide.Zoxide
	tmuxinator tmuxinator.Tmuxinator
}

func NewConnector(
	config model.Config,
	dir dir.Dir,
	home home.Home,
	lister lister.Lister,
	namer namer.Namer,
	startup startup.Startup,
	tmux tmux.Tmux,
	zoxide zoxide.Zoxide,
	tmuxinator tmuxinator.Tmuxinator,
) Connector {
	return &RealConnector{
		config,
		dir,
		home,
		lister,
		namer,
		startup,
		tmux,
		zoxide,
		tmuxinator,
	}
}
