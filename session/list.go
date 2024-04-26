package session

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
)

type ListOptions struct {
	Config       bool
	HideAttached bool
	Icons        bool
	Json         bool
	Tmux         bool
	Zoxide       bool
}

func listTmuxSessions(t tmux.Tmux) ([]model.SeshSession, error) {
	tmuxSessions, err := t.ListSessions()
	if err != nil {
		return nil, fmt.Errorf("couldn't list tmux sessions: %q", err)
	}
	sessions := make([]model.SeshSession, len(tmuxSessions))
	for i, session := range tmuxSessions {
		sessions[i] = model.SeshSession{
			Src:      "tmux",
			Name:     session.Name,
			Path:     session.Path,
			Attached: session.Attached,
			Windows:  session.Windows,
		}
	}
	return sessions, nil
}

func (s *RealSession) List(opts ListOptions) ([]model.SeshSession, error) {
	return listTmuxSessions(s.tmux)
}
