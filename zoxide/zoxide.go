package zoxide

import (
	"github.com/Wingsdh/cc-sesh/v2/model"
	"github.com/Wingsdh/cc-sesh/v2/shell"
)

type Zoxide interface {
	ListResults() ([]*model.ZoxideResult, error)
	Add(path string) error
	Query(path string) (*model.ZoxideResult, error)
}

type RealZoxide struct {
	shell shell.Shell
}

func NewZoxide(shell shell.Shell) Zoxide {
	return &RealZoxide{shell}
}
