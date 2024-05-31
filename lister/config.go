package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
)

func listConfigSessions(c model.Config) model.SeshSessions {
	orderedIndex := make([]string, 0)
	directory := make(model.SeshSessionMap)
	for _, session := range c.SessionConfigs {
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
	}
}

func (l *RealLister) FindConfigSession(name string) (model.SeshSession, bool) {
	sessions := listConfigSessions(l.config)
	if session, exists := sessions.Directory[name]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
