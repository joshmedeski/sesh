package lister

import (
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

type Lister interface {
	List(opts ListOptions) ([]model.SeshSession, error)
}

type RealLister struct {
	config model.Config
	home   home.Home
	tmux   tmux.Tmux
	zoxide zoxide.Zoxide
}

func NewLister(config model.Config, home home.Home, tmux tmux.Tmux, zoxide zoxide.Zoxide) Lister {
	return &RealLister{config, home, tmux, zoxide}
}
