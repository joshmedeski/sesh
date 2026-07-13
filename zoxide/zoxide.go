package zoxide

import (
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
)

type Zoxide interface {
	ListResults() ([]*model.ZoxideResult, error)
	Add(path string) error
	Query(path string) (*model.ZoxideResult, error)
}

// Default commands reproduce zoxide's behavior exactly, so an empty
// [frecency] config leaves the backend byte-identical to prior versions.
const (
	defaultListCommand  = "zoxide query --list --score"
	defaultQueryCommand = "zoxide query {}"
	defaultAddCommand   = "zoxide add {}"
)

type RealZoxide struct {
	shell        shell.Shell
	listCommand  string
	queryCommand string
	addCommand   string
}

// NewZoxide builds the frecency backend, falling back to the zoxide
// defaults for any command the frecency config leaves empty.
func NewZoxide(shell shell.Shell, frecency model.FrecencyConfig) Zoxide {
	return &RealZoxide{
		shell:        shell,
		listCommand:  orDefault(frecency.ListCommand, defaultListCommand),
		queryCommand: orDefault(frecency.QueryCommand, defaultQueryCommand),
		addCommand:   orDefault(frecency.AddCommand, defaultAddCommand),
	}
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
