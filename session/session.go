package session

import (
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

type Session interface {
	List(opts ListOptions) ([]model.SeshSession, error)
}

type RealSession struct {
	tmux   tmux.Tmux
	zoxide zoxide.Zoxide
}

func NewSession(tmux tmux.Tmux, zoxide zoxide.Zoxide) Session {
	return &RealSession{tmux, zoxide}
}
