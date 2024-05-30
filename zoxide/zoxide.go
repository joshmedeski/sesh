package zoxide

import (
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/shell"
)

type Zoxide interface {
	ListResults() ([]*model.ZoxideResult, error)
	Add(path string) error
}

type RealZoxide struct {
	shell shell.Shell
}

func NewZoxide(shell shell.Shell) Zoxide {
	return &RealZoxide{shell}
}
