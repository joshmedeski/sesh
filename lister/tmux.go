package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
)

func tmuxKey(name string) string {
	return fmt.Sprintf("tmux:%s", name)
}

func listTmuxSessions(t tmux.Tmux) (model.SeshSessionMap, error) {
	tmuxSessions, err := t.ListSessions()
	if err != nil {
		return nil, fmt.Errorf("couldn't list tmux sessions: %q", err)
	}
	sessions := make(model.SeshSessionMap)
	for _, session := range tmuxSessions {
		key := tmuxKey(session.Name)
		sessions[key] = model.SeshSession{
			Src: "tmux",
			// TODO: prepend icon if configured
			Name:     session.Name,
			Path:     session.Path,
			Attached: session.Attached,
			Windows:  session.Windows,
		}
	}
	return sessions, nil
}

func (l *RealLister) FindTmuxSession(name string) (model.SeshSession, bool) {
	sessions, err := listTmuxSessions(l.tmux)
	if err != nil {
		return model.SeshSession{}, false
	}
	if session, exists := sessions[name]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
