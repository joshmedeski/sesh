package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func tmuxinatorKey(name string) string {
	return fmt.Sprintf("tmuxinator:%s", name)
}

func listTmuxinator(l *RealLister) (model.SeshSessions, error) {
	tmuxinatorResults, err := l.tmuxinator.List()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list tmuxinator sessions: %q", err)
	}

	numTmuxinatorResults := len(tmuxinatorResults)
	orderedIndex := make([]string, numTmuxinatorResults)
	directory := make(model.SeshSessionMap)

	for i, session := range tmuxinatorResults {
		key := tmuxinatorKey(session.Name)
		orderedIndex[i] = key
		directory[key] = model.SeshSession{
			Src:  "tmuxinator",
			Name: session.Name,
		}
	}
	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) FindTmuxinatorConfig(name string) (model.SeshSession, bool) {
	sessions, _ := listTmuxinator(l)
	key := tmuxinatorKey(name)
	if session, exists := sessions.Directory[key]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
