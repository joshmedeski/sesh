package session

import (
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/path"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

type Session interface {
	List(opts ListOptions) ([]model.SeshSession, error)
}

type RealSession struct {
	path   path.Path
	tmux   tmux.Tmux
	zoxide zoxide.Zoxide
}

func NewSession(path path.Path, tmux tmux.Tmux, zoxide zoxide.Zoxide) Session {
	return &RealSession{path, tmux, zoxide}
}
