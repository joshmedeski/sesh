package startup

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type Startup interface {
	Exec(session model.SeshSession) (string, error)
}

type RealStartup struct {
	lister lister.Lister
	tmux   tmux.Tmux
	config model.Config
}

func NewStartup(config model.Config, lister lister.Lister, tmux tmux.Tmux) Startup {
	return &RealStartup{lister, tmux, config}
}

func (s *RealStartup) Exec(session model.SeshSession) (string, error) {
	strategies := []func(*RealStartup, model.SeshSession) (string, error){
		configStrategy,
		defaultConfigStrategy,
	}

	for _, strategy := range strategies {
		if command, err := strategy(s, session); err != nil {
			return "", fmt.Errorf("failed to determine startup command: %w", err)
		} else if command != "" {
			s.tmux.SendKeys(session.Name, command)
			return fmt.Sprintf("executing startup command: %s", command), nil
		}
	}
	return "", nil // no command to run
}
