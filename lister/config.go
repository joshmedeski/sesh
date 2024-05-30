package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
)

func configKey(name string) string {
	return fmt.Sprintf("config:%s", name)
}

func listConfigSessions(c model.Config) model.SeshSessionMap {
	sessions := make(model.SeshSessionMap)
	for _, session := range c.SessionConfigs {
		if session.Name != "" {
			key := configKey(session.Name)
			sessions[key] = model.SeshSession{
				Src:  "config",
				Name: session.Name,
				Path: session.Path,
			}
		}
	}
	return sessions
}

func (l *RealLister) FindConfigSession(name string) (model.SeshSession, bool) {
	sessions := listConfigSessions(l.config)
	if session, exists := sessions[name]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
