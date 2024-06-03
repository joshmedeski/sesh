package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
)

func listConfig(l *RealLister) (model.SeshSessions, error) {
	orderedIndex := make([]string, 0)
	directory := make(model.SeshSessionMap)
	for _, session := range l.config.SessionConfigs {
		if session.Name != "" {
			key := fmt.Sprintf("config:%s", session.Name)
			orderedIndex = append(orderedIndex, key)
			directory[key] = model.SeshSession{
				Src:  "config",
				Name: session.Name,
				Path: session.Path,
			}
		}
	}
	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) FindConfigSession(name string) (model.SeshSession, bool) {
	sessions, _ := listConfig(l)
	if session, exists := sessions.Directory[name]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
