package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
)

func listTmuxSessions(t tmux.Tmux) (model.SeshSessions, error) {
	tmuxSessions, err := t.ListSessions()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list tmux sessions: %q", err)
	}
	numOfSessions := len(tmuxSessions)
	orderedIndex := make([]string, numOfSessions)
	directory := make(model.SeshSessionMap)
	for _, session := range tmuxSessions {
		key := fmt.Sprintf("tmux:%s", session.Name)
		orderedIndex = append(orderedIndex, key)
		directory[key] = model.SeshSession{
			Src: "tmux",
			// TODO: prepend icon if configured
			Name:     session.Name,
			Path:     session.Path,
			Attached: session.Attached,
			Windows:  session.Windows,
		}
	}
	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) FindTmuxSession(name string) (model.SeshSession, bool) {
	sessions, err := listTmuxSessions(l.tmux)
	if err != nil {
		return model.SeshSession{}, false
	}
	if session, exists := sessions.Directory[name]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
