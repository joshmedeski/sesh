package startup

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type Startup interface {
	Exec(session model.SeshSession) (string, error)
	ResolveCommand(session model.SeshSession) (string, error)
	WrapForShell(command string) string
}

type RealStartup struct {
	os       oswrap.Os
	lister   lister.Lister
	tmux     tmux.Tmux
	config   model.Config
	home     home.Home
	replacer replacer.Replacer
}

func NewStartup(
	os oswrap.Os, config model.Config, lister lister.Lister, tmux tmux.Tmux, home home.Home, replacer replacer.Replacer,
) Startup {
	return &RealStartup{os, lister, tmux, config, home, replacer}
}

// ResolveCommand walks the strategy chain (per-session config → wildcard →
// default) and returns the first non-empty startup command without executing
// it. Returns "" when no strategy applies.
func (s *RealStartup) ResolveCommand(session model.SeshSession) (string, error) {
	strategies := []func(*RealStartup, model.SeshSession) (string, error){
		configStrategy,
		configWildcardStartupStrategy,
		defaultConfigStrategy,
	}
	for _, strategy := range strategies {
		command, err := strategy(s, session)
		if err != nil {
			return "", fmt.Errorf("failed to determine startup command: %w", err)
		}
		if command != "" {
			return command, nil
		}
	}
	return "", nil
}

// WrapForShell returns a shell-command string suitable for passing as the
// trailing positional argument to `tmux new-session` / `tmux new-window`.
// The resulting pane runs $SHELL interactively and executes command as part
// of pane creation, avoiding send-keys races without re-running shell init a
// second time after the command exits.
// Returns "" for empty input so callers can detect "no command".
func (s *RealStartup) WrapForShell(command string) string {
	if command == "" {
		return ""
	}
	shellPath := s.os.Getenv("SHELL")
	if shellPath == "" {
		shellPath = "/bin/sh"
	}
	return posixSingleQuote(shellPath) + " -i -c " + posixSingleQuote(command)
}

func (s *RealStartup) Exec(session model.SeshSession) (string, error) {
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

		// Inject the window's startup_script as its initial shell-command so
		// it runs reliably regardless of shell-init speed (issue #188).
		wrapped := s.WrapForShell(windowConfig.StartupScript)
		if ret, err := s.tmux.NewWindowInSession(windowConfig.Name, windowConfig.Path, session.Name, wrapped); err != nil {
			return ret, err
		}
	}
	if len(session.WindowNames) > 0 {
		if _, err := s.tmux.SelectWindow(session.Name + ":^"); err != nil {
			return "", err
		}
	}

	// The main startup_command is injected into `tmux new-session` by the
	// connector (before Exec runs). Here we only resolve it for logging.
	command, err := s.ResolveCommand(session)
	if err != nil {
		return "", err
	}
	if command != "" {
		return fmt.Sprintf("resolved startup command: %s", command), nil
	}
	return "", nil
}
