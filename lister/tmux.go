package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
)

func listTmux(l *RealLister) (model.SeshSessions, error) {
	tmuxSessions, err := l.tmux.ListSessions()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list tmux sessions: %q", err)
	}
	numOfSessions := len(tmuxSessions)
	orderedIndex := make([]string, numOfSessions)
	directory := make(model.SeshSessionMap)
	for i, session := range tmuxSessions {
		key := fmt.Sprintf("tmux:%s", session.Name)
		orderedIndex[i] = key
		directory[key] = model.SeshSession{
			Src:      "tmux",
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
	sessions, err := listTmux(l)
	if err != nil {
		return model.SeshSession{}, false
	}
	if session, exists := sessions.Directory[name]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
