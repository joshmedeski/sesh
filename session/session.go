package session

import (
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
)

type Session interface {
	List(opts ListOptions) ([]model.SeshSession, error)
}

type RealSession struct {
	tmux tmux.Tmux
}

func NewSession(tmux tmux.Tmux) Session {
	return &RealSession{tmux}
}
