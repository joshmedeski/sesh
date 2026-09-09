package tmuxinator

import (
	"errors"

	"github.com/joshmedeski/sesh/v2/model"
)

// Disabled is the tmuxinator to wire in when the multiplexer isn't tmux.
// tmuxinator can only start tmux sessions, so on any other backend its configs
// would list but never connect; an empty list keeps them out of the picker and
// out of the connector's strategy chain.
type Disabled struct{}

func NewDisabled() Tmuxinator { return Disabled{} }

func (Disabled) List() ([]*model.TmuxinatorConfig, error) { return nil, nil }

func (Disabled) Start(string) (string, error) {
	return "", errors.New("tmuxinator only works with tmux")
}
