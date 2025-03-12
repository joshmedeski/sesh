package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func configKey(name string) string {
	return fmt.Sprintf("config:%s", name)
}

func listConfig(l *RealLister) (model.SeshSessions, error) {
	windows := make(model.SeshWindowMap)
	for _, window := range l.config.WindowConfigs {
		key := configKey(window.Name)
		path, err := l.home.ExpandHome(window.Path)
		if err != nil {
			return model.SeshSessions{}, fmt.Errorf("couldn't expand home: %q", err)
		}

		if window.StartupScript != "" && window.DisableStartScript {
			return model.SeshSessions{}, fmt.Errorf("startup_script and disable_start_script are mutually exclusive")
		}

		windows[key] = model.WindowConfig{
			Name:               window.Name,
			Path:               path,
			StartupScript:      window.StartupScript,
			DisableStartScript: window.DisableStartScript,
		}
	}

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

			if session.StartupCommand != "" && session.DisableStartCommand {
				return model.SeshSessions{}, fmt.Errorf("startup_command and disable_start_command are mutually exclusive")
			}

			windowConfigs := make([]model.WindowConfig, len(session.Windows))
			for _, window := range session.Windows {
				windowConfig, ok := windows[configKey(window)]
				if !ok {
					return model.SeshSessions{}, fmt.Errorf("window %s is not defined in config", window)
				}
				windowConfigs = append(windowConfigs, windowConfig)
			}

			directory[key] = model.SeshSession{
				Src:                   "config",
				Name:                  session.Name,
				Path:                  path,
				StartupCommand:        session.StartupCommand,
				PreviewCommand:        session.PreviewCommand,
				DisableStartupCommand: session.DisableStartCommand,
				Tmuxinator:            session.Tmuxinator,
				WindowConfigs:         windowConfigs,
			}
		}
	}
	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) FindConfigSession(name string) (model.SeshSession, bool) {
	key := configKey(name)
	sessions, _ := listConfig(l)
	if session, exists := sessions.Directory[key]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
