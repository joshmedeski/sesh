package session

import (
	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

type Session interface {
	List(opts ListOptions) ([]model.SeshSession, error)
}

type RealSession struct {
	config config.Config
	home   home.Home
	tmux   tmux.Tmux
	zoxide zoxide.Zoxide
}

func NewSession(config config.Config, home home.Home, tmux tmux.Tmux, zoxide zoxide.Zoxide) Session {
	return &RealSession{config, home, tmux, zoxide}
}
