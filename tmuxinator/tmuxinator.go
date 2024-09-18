package tmuxinator

import (
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/shell"
)

type Tmuxinator interface {
	ListConfigs() ([]*model.TmuxinatorConfig, error)
	CreateSession(targetSession string) (string, error)
}

type RealTmuxinator struct {
	shell shell.Shell
}

func NewTmuxinator(shell shell.Shell) Tmuxinator {
	return &RealTmuxinator{shell}
}

func (t *RealTmuxinator) CreateSession(targetSession string) (string, error) {
	return t.shell.Cmd("tmuxinator", "start", targetSession)
}
