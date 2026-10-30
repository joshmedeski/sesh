package startup

import (
	"cmp"
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
		windows[lister.ConfigKey(window.Name)] = window
	}

	panes := make(model.SeshPaneMap)
	for _, pane := range s.config.PaneConfigs {
		key := lister.ConfigKey(pane.Name)
		panes[key] = pane
	}

	for _, windowName := range session.WindowNames {
		if err := s.createWindow(session, windowName, windows, panes); err != nil {
			return "", err
		}
	}
	_, _ = s.tmux.NextWindowInSession(session.Name)

	var message string
	for _, strategy := range strategies {
		if command, err := strategy(s, session); err != nil {
			return "", fmt.Errorf("failed to determine startup command: %w", err)
		} else if command != "" {
			if _, err := s.tmux.SendKeys(session.Name+":", command); err != nil {
				return "", fmt.Errorf("failed to send startup command: %w", err)
			}
			message = fmt.Sprintf("executing startup command: %s", command)
			break
		}
	}

	// panes come after the command so it lands in the main pane
	if firstPanes := s.firstWindowPanes(session); len(firstPanes.Panes) > 0 {
		path, err := s.home.ExpandPath(session.Path)
		if err != nil {
			return "", fmt.Errorf("first window: couldn't expand home: %q", err)
		}
		// "^" is tmux's lowest-index window
		if err := s.createPanes(session.Name+":^", path, firstPanes, panes); err != nil {
			return "", fmt.Errorf("first window: %w", err)
		}
	}

	return message, nil
}

func (s *RealStartup) firstWindowPanes(session model.SeshSession) model.WindowPanes {
	for _, sessionConfig := range s.config.SessionConfigs {
		if lister.ConfigKey(sessionConfig.Name) == lister.ConfigKey(session.Name) && len(sessionConfig.Panes) > 0 {
			return sessionConfig.WindowPanes
		}
	}
	if wildcard, ok := s.lister.FindConfigWildcard(session.Path); ok && len(wildcard.Panes) > 0 {
		return wildcard.WindowPanes
	}
	return s.config.DefaultSessionConfig.WindowPanes
}

func (s *RealStartup) createWindow(session model.SeshSession, name string, windows model.SeshWindowMap, panes model.SeshPaneMap) error {
	window, ok := windows[lister.ConfigKey(name)]
	if !ok {
		return fmt.Errorf("window %s is not defined in config", name)
	}
	// session path is the fallback for windows without their own
	path, err := s.home.ExpandPath(cmp.Or(window.Path, session.Path))
	if err != nil {
		return fmt.Errorf("window %s: couldn't expand home: %q", name, err)
	}

	// "session:index" target; window names may contain dots or digits
	target, err := s.tmux.NewWindowInSession(model.TmuxWindowOpts{
		Name: window.Name, StartDir: path, TargetSession: session.Name,
	})
	if err != nil {
		return fmt.Errorf("window %s: %w", name, err)
	}
	if _, err := s.tmux.SendKeys(target, window.StartupScript); err != nil {
		return fmt.Errorf("window %s: %w", name, err)
	}
	if err := s.createPanes(target, path, window.WindowPanes, panes); err != nil {
		return fmt.Errorf("window %s: %w", name, err)
	}
	return nil
}

func (s *RealStartup) createPanes(target string, startDir string, windowPanes model.WindowPanes, panes model.SeshPaneMap) error {
	var focusID, zoomID, prevID string
	for _, paneName := range windowPanes.Panes {
		pane, ok := panes[lister.ConfigKey(paneName)]
		if !ok {
			return fmt.Errorf("pane %s is not defined in config", paneName)
		}
		// split the previous pane so panes keep config order
		paneID, err := s.tmux.NewPaneInWindow(model.TmuxPaneOpts{
			TargetWindow: target,
			TargetPane:   prevID,
			StartDir:     startDir,
		})
		if err != nil {
			return fmt.Errorf("pane %s: %w", paneName, err)
		}
		if pane.StartupCommand != "" {
			if _, err := s.tmux.SendKeys(paneID, pane.StartupCommand); err != nil {
				return fmt.Errorf("pane %s: run startup command: %w", paneName, err)
			}
		}
		prevID = paneID
		if pane.Focus && focusID == "" {
			focusID = paneID
		}
		if pane.Zoom && zoomID == "" {
			zoomID = paneID
		}
	}
	if windowPanes.Layout != "" {
		if _, err := s.tmux.SelectLayout(target, windowPanes.Layout); err != nil {
			return fmt.Errorf("apply layout %s: %w", windowPanes.Layout, err)
		}
	}
	if focusID != "" {
		if _, err := s.tmux.SelectPaneByID(focusID); err != nil {
			return fmt.Errorf("focus pane: %w", err)
		}
	}
	zoomTarget := zoomID
	if zoomTarget == "" && windowPanes.Zoom && len(windowPanes.Panes) > 0 {
		zoomTarget = cmp.Or(focusID, target)
	}
	if zoomTarget != "" {
		if _, err := s.tmux.ZoomPane(zoomTarget); err != nil {
			return fmt.Errorf("zoom pane: %w", err)
		}
	}
	return nil
}
