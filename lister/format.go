package lister

import (
	"github.com/joshmedeski/sesh/v2/formatter"
	"github.com/joshmedeski/sesh/v2/model"
)

// activeWindowNameFormat returns only the active window's name for each tmux
// session. Inactive windows render as an empty string and are ignored by
// ListAllWindowNames' parser.
const activeWindowNameFormat = "#{?window_active,#{window_name},}"

func needsFormatting(opts ListOptions) bool {
	return !opts.Json && (opts.Icons || opts.FormatSet)
}

// Format applies display-only formatting to a session list through the injected
// formatter dependency. The formatter returns a copy so domain/cache data stays raw.
func (l *RealLister) Format(sessions model.SeshSessions, opts ListOptions) (model.SeshSessions, error) {
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
