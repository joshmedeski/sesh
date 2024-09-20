package tmuxinator

import (
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/shell"
)

type Tmuxinator interface {
	List() ([]*model.TmuxinatorConfig, error)
	Start(targetSession string) (string, error)
}

type RealTmuxinator struct {
	shell shell.Shell
}

func NewTmuxinator(shell shell.Shell) Tmuxinator {
	return &RealTmuxinator{shell}
}

func (t *RealTmuxinator) Start(targetSession string) (string, error) {
	return t.shell.Cmd("tmuxinator", "start", targetSession)
}
