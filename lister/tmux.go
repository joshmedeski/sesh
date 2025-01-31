package lister

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/model"
)

func tmuxKey(name string) string {
	return fmt.Sprintf("tmux:%s", name)
}

func isBlacklisted(blacklist []string, name string) bool {
	for _, blacklistedName := range blacklist {
		if strings.EqualFold(blacklistedName, name) {
			return true
		}
	}
	return false
}

func listTmux(l *RealLister) (model.SeshSessions, error) {
	tmuxSessions, err := l.tmux.ListSessions()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list tmux sessions: %q", err)
	}
	numOfSessions := len(tmuxSessions)
	orderedIndex := make([]string, numOfSessions)
	directory := make(model.SeshSessionMap)
	for i, session := range tmuxSessions {
		key := tmuxKey(session.Name)
		orderedIndex[i] = key
		directory[key] = model.SeshSession{
			Src:      "tmux",
			Name:     session.Name,
			Path:     session.Path,
			Attached: session.Attached,
			Windows:  session.Windows,
		}
	}

	finalOrderedIndex := []string{}
	for _, key := range orderedIndex {
		if !isBlacklisted(l.config.Blacklist, directory[key].Name) {
			finalOrderedIndex = append(finalOrderedIndex, key)
		}
	}

	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: finalOrderedIndex,
	}, nil
}

func (l *RealLister) FindTmuxSession(name string) (model.SeshSession, bool) {
	sessions, err := listTmux(l)
	if err != nil {
		return model.SeshSession{}, false
	}
	key := tmuxKey(name)
	if session, exists := sessions.Directory[key]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}

func (l *RealLister) GetLastTmuxSession() (model.SeshSession, bool) {
	sessions, err := listTmux(l)
	if err != nil {
		return model.SeshSession{}, false
	}
	if len(sessions.OrderedIndex) < 2 {
		return model.SeshSession{}, false
	}
	secondSessionIndex := sessions.OrderedIndex[1]
	return sessions.Directory[secondSessionIndex], true
}

func (l *RealLister) GetAttachedTmuxSession() (model.SeshSession, bool) {
	return GetAttachedTmuxSession(l)
}

func GetAttachedTmuxSession(l *RealLister) (model.SeshSession, bool) {
	tmuxSessions, err := l.tmux.ListSessions()
	if err != nil {
		return model.SeshSession{}, false
	}
	for _, session := range tmuxSessions {
		if session.Attached != 0 {
			return model.SeshSession{
				Src:      "tmux",
				Name:     session.Name,
				Path:     session.Path,
				Attached: session.Attached,
				Windows:  session.Windows,
			}, true
		}
	}
	return model.SeshSession{}, false
}
