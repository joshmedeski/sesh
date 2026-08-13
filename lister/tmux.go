package lister

import (
	"fmt"
	"time"

	"github.com/joshmedeski/sesh/v2/model"
)

func tmuxKey(name string) string {
	return fmt.Sprintf("tmux:%s", name)
}

// tmuxToSesh maps a tmux session onto a SeshSession, copying the time fields
// nil-safely so the two models never alias mutable pointers.
func tmuxToSesh(session *model.TmuxSession) model.SeshSession {
	return model.SeshSession{
		Src:          "tmux",
		Name:         session.Name,
		Path:         session.Path,
		Attached:     session.Attached,
		Windows:      session.Windows,
		Created:      cloneTime(session.Created),
		LastAttached: cloneTime(session.LastAttached),
		Activity:     cloneTime(session.Activity),
		Alerts:       session.Alerts,
	}
}

// cloneTime returns a fresh copy of t, or nil when t is nil.
func cloneTime(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	c := *t
	return &c
}

func listTmux(l *RealLister) (model.SeshSessions, error) {
	tmuxSessions, err := l.tmux.ListSessions()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list tmux sessions: %q", err)
	}

	directory := make(map[string]model.SeshSession)
	orderedIndex := []string{}

	for _, session := range tmuxSessions {
		key := tmuxKey(session.Name)
		orderedIndex = append(orderedIndex, key)
		directory[key] = tmuxToSesh(session)
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

	filtered := sessions.OrderedIndex
	if len(l.config.Blacklist) > 0 {
		compiled := compileBlacklist(l.config.Blacklist)
		filtered = make([]string, 0, len(sessions.OrderedIndex))
		for _, index := range sessions.OrderedIndex {
			session := sessions.Directory[index]
			if !isBlacklisted(compiled, session.Name) {
				filtered = append(filtered, index)
			}
		}
	}

	if len(filtered) < 2 {
		return model.SeshSession{}, false
	}
	secondSessionIndex := filtered[1]
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
			return tmuxToSesh(session), true
		}
	}
	return model.SeshSession{}, false
}
