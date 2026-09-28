package lister

import (
	"cmp"
	"slices"

	"github.com/joshmedeski/sesh/v2/formatter"
	"github.com/joshmedeski/sesh/v2/model"
)

// activeWindowNameFormat returns only the active window's name for each tmux
// session. Inactive windows render as an empty string and are ignored by
// ListAllWindowNames' parser.
const activeWindowNameFormat = "#{?window_active,#{window_name},}"

func needsFormatting(opts ListOptions) bool {
	return opts.Json || opts.Icons || opts.FormatSet
}

// Format applies display-only formatting to a session list through the injected
// formatter dependency. The formatter returns a copy so domain/cache data stays raw.
func (l *RealLister) Format(sessions model.SeshSessions, opts ListOptions) (model.SeshSessions, error) {
	if opts.Json {
		return l.resolve(sessions), nil
	}
	if !needsFormatting(opts) {
		return sessions, nil
	}

	var activeWindowNames map[string]string
	if opts.FormatSet && l.formatter.UsesActiveWindowName(opts.Format) && hasTmuxSessions(sessions) {
		windowNames, err := l.tmux.ListAllWindowNames(activeWindowNameFormat)
		if err == nil {
			activeWindowNames = firstActiveWindowNameBySession(windowNames)
		}
	}

	var template *string
	if opts.FormatSet {
		template = &opts.Format
	}

	return l.formatter.Format(sessions, formatter.Options{
		Template:          template,
		Icons:             opts.Icons,
		NoColor:           opts.NoColor,
		IconExcludes:      opts.IconExcludes,
		ActiveWindowNames: activeWindowNames,
	})
}

func hasTmuxSessions(sessions model.SeshSessions) bool {
	for _, key := range sessions.OrderedIndex {
		if sessions.Directory[key].Src == "tmux" {
			return true
		}
	}
	return false
}

func firstActiveWindowNameBySession(windowNames map[string][]string) map[string]string {
	if len(windowNames) == 0 {
		return nil
	}

	active := make(map[string]string, len(windowNames))
	for session, names := range windowNames {
		if len(names) > 0 {
			active[session] = names[0]
		}
	}
	return active
}

// resolve fills in what `sesh list --json` reports beyond the raw source data:
// the icon the picker shows, the matching [[wildcard]], live tmux windows, and
// for sessions sesh creates on connect, the startup, preview, and windows that
// fall back from the session to its wildcard to [default_session].
func (l *RealLister) resolve(sessions model.SeshSessions) model.SeshSessions {
	resolveIcon := IconResolver(l.config, l.home, l)
	var liveWindows map[string][]model.TmuxWindow
	if hasTmuxSessions(sessions) {
		liveWindows, _ = l.tmux.ListAllWindows()
	}
	configByName := make(map[string]model.SessionConfig, len(l.config.SessionConfigs))
	for _, config := range l.config.SessionConfigs {
		configByName[config.Name] = config
	}

	directory := make(model.SeshSessionMap, len(sessions.Directory))
	for key, session := range sessions.Directory {
		if resolveIcon != nil {
			session.Icon = resolveIcon(session)
		}
		wildcard, _ := l.FindConfigWildcard(session.Path)
		session.Wildcard = wildcard.Pattern
		switch session.Src {
		case "tmux":
			session.TmuxWindows = liveWindows[session.Name]
			if config, ok := configByName[session.Name]; ok {
				session.Alias, session.AliasAutoConnect = config.Alias, config.AliasAutoConnect
			}
		case "config", "zoxide":
			session = l.applyFallbacks(session, wildcard)
		}
		directory[key] = session
	}
	return model.SeshSessions{OrderedIndex: sessions.OrderedIndex, Directory: directory}
}

func (l *RealLister) applyFallbacks(session model.SeshSession, wildcard model.WildcardConfig) model.SeshSession {
	defaults := l.config.DefaultSessionConfig
	if session.StartupCommand == "" && session.Tmuxinator == "" && !session.DisableStartupCommand {
		session.DisableStartupCommand = wildcard.DisableStartCommand
		if !session.DisableStartupCommand {
			session.StartupCommand = cmp.Or(wildcard.StartupCommand, defaults.StartupCommand)
		}
	}
	session.PreviewCommand = cmp.Or(session.PreviewCommand, wildcard.PreviewCommand, defaults.PreviewCommand)
	if len(session.WindowNames) == 0 {
		session.WindowNames = wildcard.Windows
	}
	if len(session.WindowNames) == 0 {
		session.WindowNames = defaults.Windows
	}
	session.WindowConfigs = l.windowConfigs(session)
	return session
}

func (l *RealLister) windowConfigs(session model.SeshSession) []model.WindowConfig {
	var configs []model.WindowConfig
	for _, name := range session.WindowNames {
		i := slices.IndexFunc(l.config.WindowConfigs, func(w model.WindowConfig) bool { return w.Name == name })
		if i < 0 {
			continue
		}
		window := l.config.WindowConfigs[i]
		window.Path = cmp.Or(window.Path, session.Path)
		if path, err := l.home.ExpandPath(window.Path); err == nil {
			window.Path = path
		}
		configs = append(configs, window)
	}
	return configs
}
