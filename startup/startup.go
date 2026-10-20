package startup

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type Startup interface {
	Exec(session model.SeshSession) (string, error)
}

type RealStartup struct {
	lister   lister.Lister
	tmux     tmux.Tmux
	config   model.Config
	home     home.Home
	replacer replacer.Replacer
}

func NewStartup(
	config model.Config, lister lister.Lister, tmux tmux.Tmux, home home.Home, replacer replacer.Replacer,
) Startup {
	return &RealStartup{lister, tmux, config, home, replacer}
}

func (s *RealStartup) Exec(session model.SeshSession) (string, error) {
	strategies := []func(*RealStartup, model.SeshSession) (string, error){
		configStrategy,
		configWildcardStartupStrategy,
		defaultConfigStrategy,
	}

	windows := make(model.SeshWindowMap)
	for _, window := range s.config.WindowConfigs {
		key := lister.ConfigKey(window.Name)
		if window.Path != "" {
			path, err := s.home.ExpandPath(window.Path)
			if err != nil {
				return "", fmt.Errorf("couldn't expand home: %q", err)
			}
			window.Path = path
		}
		windows[key] = window
	}

	panes := make(model.SeshPaneMap)
	for _, pane := range s.config.PaneConfigs {
		key := lister.ConfigKey(pane.Name)
		panes[key] = pane
	}

	for _, windowName := range session.WindowNames {
		if ret, err := s.createWindow(session, windowName, windows, panes); err != nil {
			return ret, err
		}
	}
	_, _ = s.tmux.NextWindowInSession(session.Name)

	for _, strategy := range strategies {
		if command, err := strategy(s, session); err != nil {
			return "", fmt.Errorf("failed to determine startup command: %w", err)
		} else if command != "" {
			_, _ = s.tmux.SendKeys(session.Name+":", command)
			return fmt.Sprintf("executing startup command: %s", command), nil
		}
	}

	return "", nil // no command to run
}

func (s *RealStartup) createWindow(session model.SeshSession, name string, windows model.SeshWindowMap, panes model.SeshPaneMap) (string, error) {
	window, ok := windows[lister.ConfigKey(name)]
	if !ok {
		return "", fmt.Errorf("window %s is not defined in config", name)
	}
	if window.Path == "" {
		path, err := s.home.ExpandPath(session.Path)
		if err != nil {
			return "", fmt.Errorf("couldn't expand home: %q", err)
		}
		window.Path = path
	}

	if ret, err := s.tmux.NewWindowInSession(model.TmuxWindowOpts{
		Name: window.Name, StartDir: window.Path, TargetSession: session.Name,
	}); err != nil {
		return ret, err
	}
	target := session.Name + ":" + window.Name
	if ret, err := s.tmux.SendKeys(target, window.StartupScript); err != nil {
		return ret, err
	}
	return s.createPanes(target, window, panes)
}

func (s *RealStartup) createPanes(target string, window model.WindowConfig, panes model.SeshPaneMap) (string, error) {
	for _, paneName := range window.Panes {
		pane, ok := panes[lister.ConfigKey(paneName)]
		if !ok {
			return "", fmt.Errorf("pane %s is not defined in config", paneName)
		}
		// TODO: implement focus?
		if _, err := s.tmux.NewPaneInWindow(model.TmuxPaneOpts{
			TargetWindow: target,
			StartDir:     window.Path,
			Command:      pane.StartupCommand,
			Layout:       window.Layout,
		}); err != nil {
			return "", err
		}
	}
	return "", nil
}
