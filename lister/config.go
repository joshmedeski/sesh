package lister

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/model"
)

func configKey(name string) string {
	return fmt.Sprintf("config:%s", name)
}

func listConfig(l *RealLister) (model.SeshSessions, error) {
	activeSessions, _ := listTmux(l)

	orderedIndex := make([]string, 0)
	directory := make(model.SeshSessionMap)
	for _, session := range l.config.SessionConfigs {
		if session.Name != "" {
			key := configKey(session.Name)
			orderedIndex = append(orderedIndex, key)
			path, err := l.home.ExpandHome(session.Path)
			if err != nil {
				return model.SeshSessions{}, fmt.Errorf("couldn't expand home: %q", err)
			}
			// check if session is attached
			isAttached := 0
			tmuxKey := strings.Replace(key, "config:", "tmux:", 1)
			tmuxSession := activeSessions.Directory[tmuxKey]
			if tmuxSession != (model.SeshSession{}) {
				isAttached = tmuxSession.Attached
			}
			directory[key] = model.SeshSession{
				Src:            "config",
				Attached:       isAttached,
				Name:           session.Name,
				Path:           path,
				StartupCommand: session.StartupCommand,
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
	key := configKey(name)
	if session, exists := sessions.Directory[key]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
