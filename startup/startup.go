package startup

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/mux"
	"github.com/joshmedeski/sesh/v2/replacer"
)

type Startup interface {
	Exec(session model.SeshSession) (string, error)
}

type RealStartup struct {
	lister   lister.Lister
	mux      mux.Mux
	config   model.Config
	home     home.Home
	replacer replacer.Replacer
}

func NewStartup(
	config model.Config, lister lister.Lister, m mux.Mux, home home.Home, replacer replacer.Replacer,
) Startup {
	return &RealStartup{lister, m, config, home, replacer}
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
		var path string = ""
		var err error = nil
		if window.Path != "" {
			path, err = s.home.ExpandPath(window.Path)
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
			path, err := s.home.ExpandPath(session.Path)
			if err != nil {
				return "", fmt.Errorf("couldn't expand home: %q", err)
			}
			windowConfig.Path = path
		}

		// create the new window
		if ret, err := s.mux.NewWindowInSession(model.TmuxWindowOpts{
			Name:          windowConfig.Name,
			StartDir:      windowConfig.Path,
			TargetSession: session.Name,
		}); err != nil {
			return ret, err
		}
		if ret, err := s.mux.SendKeys(session.Name, windowConfig.StartupScript); err != nil {
			return ret, err
		}
	}
	s.mux.NextWindowInSession(session.Name)

	for _, strategy := range strategies {
		if command, err := strategy(s, session); err != nil {
			return "", fmt.Errorf("failed to determine startup command: %w", err)
		} else if command != "" {
			s.mux.SendKeys(session.Name, command)
			return fmt.Sprintf("executing startup command: %s", command), nil
		}
	}

	return "", nil // no command to run
}
