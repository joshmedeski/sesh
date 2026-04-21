package startup

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/wezterm"
)

type Startup interface {
	Exec(session model.SeshSession) (string, error)
}

type RealStartup struct {
	lister   lister.Lister
	tmux     tmux.Tmux
	wezterm  wezterm.Wezterm
	config   model.Config
	home     home.Home
	replacer replacer.Replacer
}

func NewStartup(
	config model.Config, lister lister.Lister, tmux tmux.Tmux, wezterm wezterm.Wezterm, home home.Home, replacer replacer.Replacer,
) Startup {
	return &RealStartup{lister, tmux, wezterm, config, home, replacer}
}

func (s *RealStartup) terminal(sessionName string) Terminal {
	if s.config.Backend == "wezterm" {
		return NewWeztermTerminal(s.wezterm, sessionName)
	}
	return NewTmuxTerminal(s.tmux)
}

func (s *RealStartup) Exec(session model.SeshSession) (string, error) {
	strategies := []func(*RealStartup, model.SeshSession) (string, error){
		configStrategy,
		configWildcardStartupStrategy,
		defaultConfigStrategy,
	}

	term := s.terminal(session.Name)

	windows := make(model.SeshWindowMap)
	for _, window := range s.config.WindowConfigs {
		key := lister.ConfigKey(window.Name)
		var path string = ""
		var err error = nil
		if window.Path != "" {
			path, err = s.home.ExpandHome(window.Path)
			if err != nil {
				return "", fmt.Errorf("couldn't expand home: %q", err)
			}
		}

		windows[key] = model.WindowConfig{
			Name:          window.Name,
			Path:          path,
			StartupScript: window.StartupScript,
		}
	}

	for _, window := range session.WindowNames {
		windowConfig, ok := windows[lister.ConfigKey(window)]
		if !ok {
			return "", fmt.Errorf("window %s is not defined in config", window)
		}
		if windowConfig.Path == "" {
			path, err := s.home.ExpandHome(session.Path)
			if err != nil {
				return "", fmt.Errorf("couldn't expand home: %q", err)
			}
			windowConfig.Path = path
		}

		// create the new window
		if ret, err := term.CreateWindow(windowConfig.Path, windowConfig.Name); err != nil {
			return ret, err
		}
		if ret, err := term.SendCommand(session.Name, windowConfig.StartupScript); err != nil {
			return ret, err
		}
	}
	term.FocusFirstWindow()

	for _, strategy := range strategies {
		if command, err := strategy(s, session); err != nil {
			return "", fmt.Errorf("failed to determine startup command: %w", err)
		} else if command != "" {
			term.SendCommand(session.Name, command)
			return fmt.Sprintf("executing startup command: %s", command), nil
		}
	}

	return "", nil // no command to run
}
