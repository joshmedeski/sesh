package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func ConfigKey(name string) string {
	return fmt.Sprintf("config:%s", name)
}

func listConfig(l *RealLister) (model.SeshSessions, error) {
	orderedIndex := make([]string, 0)
	directory := make(model.SeshSessionMap)
	for _, session := range l.config.SessionConfigs {
		if session.Name != "" {
			key := ConfigKey(session.Name)
			orderedIndex = append(orderedIndex, key)

			path, err := l.home.ExpandHome(session.Path)
			if err != nil {
				return model.SeshSessions{}, fmt.Errorf("couldn't expand home: %q", err)
			}

			description, err := l.home.ShortenHome(session.Path)
			if err != nil {
				return model.SeshSessions{}, fmt.Errorf("couldn't shorten path: %q", err)
			}

			if session.StartupCommand != "" && session.DisableStartCommand {
				return model.SeshSessions{}, fmt.Errorf("startup_command and disable_start_command are mutually exclusive")
			}

			directory[key] = model.SeshSession{
				Src:                   "config",
				Name:                  session.Name,
				Path:                  path,
				Description:           description,
				StartupCommand:        session.StartupCommand,
				PreviewCommand:        session.PreviewCommand,
				DisableStartupCommand: session.DisableStartCommand,
				Tmuxinator:            session.Tmuxinator,
				WindowNames:           session.Windows,
			}
		}
	}
	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) FindConfigSession(name string) (model.SeshSession, bool) {
	key := ConfigKey(name)
	sessions, _ := listConfig(l)
	if session, exists := sessions.Directory[key]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
