package startup

import (
	"cmp"
	"errors"
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

	// A failing window or pane never stops the rest of the startup.
	var errs []error
	for _, windowName := range session.WindowNames {
		if err := s.createWindow(session, windowName, windows, panes); err != nil {
			errs = append(errs, err)
		}
	}
	_, _ = s.tmux.NextWindowInSession(session.Name)

	for _, strategy := range strategies {
		command, err := strategy(s, session)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to determine startup command: %w", err))
			return "", errors.Join(errs...)
		}
		if command == "" {
			continue
		}
		if _, err := s.tmux.SendKeys(session.Name+":", command); err != nil {
			errs = append(errs, fmt.Errorf("failed to send startup command: %w", err))
			return "", errors.Join(errs...)
		}
		return fmt.Sprintf("executing startup command: %s", command), errors.Join(errs...)
	}

	return "", errors.Join(errs...) // no command to run
}

func (s *RealStartup) createWindow(session model.SeshSession, name string, windows model.SeshWindowMap, panes model.SeshPaneMap) error {
	window, ok := windows[lister.ConfigKey(name)]
	if !ok {
		return fmt.Errorf("window %s is not defined in config", name)
	}
	// session path is the fallback for windows without their own
	rawPath := cmp.Or(window.Path, session.Path)
	path, err := s.home.ExpandPath(rawPath)
	if err != nil {
		return fmt.Errorf("window %s: couldn't expand home: %q", name, err)
	}
	window.Path = path

	// "session:index" target; window names may contain dots or digits
	target, err := s.tmux.NewWindowInSession(model.TmuxWindowOpts{
		Name: window.Name, StartDir: window.Path, TargetSession: session.Name,
	})
	if err != nil {
		return fmt.Errorf("window %s: %w", name, err)
	}
	var errs []error
	if _, err := s.tmux.SendKeys(target, window.StartupScript); err != nil {
		errs = append(errs, fmt.Errorf("window %s: %w", name, err))
	}
	errs = append(errs, s.createPanes(target, window, panes))
	return errors.Join(errs...)
}

// createPanes keeps going past failing panes and returns all errors.
func (s *RealStartup) createPanes(target string, window model.WindowConfig, panes model.SeshPaneMap) error {
	var errs []error
	fail := func(err error) { errs = append(errs, fmt.Errorf("window %s: %w", window.Name, err)) }

	var focusID, zoomID, prevID string
	for _, paneName := range window.Panes {
		pane, ok := panes[lister.ConfigKey(paneName)]
		if !ok {
			fail(fmt.Errorf("pane %s is not defined in config", paneName))
			continue
		}
		// split the previous pane so panes keep config order
		paneID, err := s.tmux.NewPaneInWindow(model.TmuxPaneOpts{
			TargetWindow: target,
			TargetPane:   prevID,
			StartDir:     window.Path,
			Command:      pane.StartupCommand,
			Layout:       window.Layout,
		})
		if err != nil {
			fail(fmt.Errorf("pane %s: %w", paneName, err))
		}
		// a pane that exists despite the error (send-keys/layout failed) still counts
		if paneID == "" {
			continue
		}
		prevID = paneID
		if pane.Focus && focusID == "" {
			focusID = paneID
		}
		if pane.Zoom && zoomID == "" {
			zoomID = paneID
		}
	}
	if focusID != "" {
		if _, err := s.tmux.SelectPaneByID(focusID); err != nil {
			fail(fmt.Errorf("focus pane: %w", err))
		}
	}
	zoomTarget := zoomID
	if zoomTarget == "" && window.Zoom && prevID != "" {
		zoomTarget = cmp.Or(focusID, target)
	}
	if zoomTarget != "" {
		if _, err := s.tmux.ZoomPane(zoomTarget); err != nil {
			fail(fmt.Errorf("zoom pane: %w", err))
		}
	}
	return errors.Join(errs...)
}
