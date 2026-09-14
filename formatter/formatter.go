package formatter

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/model"
)

// Options controls how sessions are rendered for display.
type Options struct {
	Template          *string
	Icons             bool
	NoColor           bool
	IconExcludes      []string
	ActiveWindowNames map[string]string
}

// Formatter renders sessions without coupling callers to icon or ANSI details.
type Formatter interface {
	Format(sessions model.SeshSessions, opts Options) (model.SeshSessions, error)
	UsesActiveWindowName(format string) bool
}

type RealFormatter struct {
	icon icon.Icon
}

func NewFormatter(icon icon.Icon) Formatter {
	return &RealFormatter{icon: icon}
}

var colorCodes = map[string]int{
	"black":          30,
	"red":            31,
	"green":          32,
	"yellow":         33,
	"blue":           34,
	"magenta":        35,
	"cyan":           36,
	"white":          37,
	"bright-black":   90,
	"gray":           90,
	"grey":           90,
	"bright-red":     91,
	"bright-green":   92,
	"bright-yellow":  93,
	"bright-blue":    94,
	"bright-magenta": 95,
	"bright-cyan":    96,
	"bright-white":   97,
}

func colorCode(color string) (int, bool) {
	if color == "" {
		return 0, true
	}
	code, ok := colorCodes[strings.ToLower(strings.TrimSpace(color))]
	return code, ok
}

func renderColors(format string, noColor bool) (string, error) {
	const prefix = "{fg:"

	var out strings.Builder
	for {
		start := strings.Index(format, prefix)
		if start == -1 {
			out.WriteString(format)
			break
		}

		out.WriteString(format[:start])
		format = format[start:]

		end := strings.IndexByte(format, '}')
		if end == -1 {
			return "", fmt.Errorf("unterminated foreground color token in format")
		}

		token := format[:end+1]
		color := strings.TrimSpace(token[len(prefix) : len(token)-1])
		code, ok := colorCode(color)
		if !ok || code == 0 {
			return "", fmt.Errorf("unsupported format color %q (use black, red, green, yellow, blue, magenta, cyan, white, or a bright-* variant)", color)
		}

		if !noColor {
			fmt.Fprintf(&out, "\033[%dm", code)
		}
		format = format[end+1:]
	}

	result := out.String()
	if noColor {
		return strings.ReplaceAll(result, "{/fg}", ""), nil
	}
	return strings.ReplaceAll(result, "{/fg}", "\033[39m"), nil
}

func sourceIconExcluded(excluded []string, source string) bool {
	for _, candidate := range excluded {
		if strings.EqualFold(strings.TrimSpace(candidate), source) {
			return true
		}
	}
	return false
}

func (f *RealFormatter) UsesActiveWindowName(format string) bool {
	return strings.Contains(format, "{active_window_name}") || strings.Contains(format, "{active_window_name_prefix}")
}

func (f *RealFormatter) Format(sessions model.SeshSessions, opts Options) (model.SeshSessions, error) {
	format := ""
	if opts.Template != nil {
		var err error
		format, err = renderColors(*opts.Template, opts.NoColor)
		if err != nil {
			return model.SeshSessions{}, err
		}
	}

	directory := make(model.SeshSessionMap, len(sessions.Directory))
	for key, session := range sessions.Directory {
		directory[key] = session
	}

	for _, key := range sessions.OrderedIndex {
		session, ok := directory[key]
		if !ok {
			continue
		}

		displayName := session.Name
		if opts.Icons && !sourceIconExcluded(opts.IconExcludes, session.Src) {
			if opts.NoColor {
				displayName = f.icon.AddIconNoColor(session)
			} else {
				displayName = f.icon.AddIcon(session)
			}
		}

		if opts.Template != nil {
			activeWindowName := ""
			if session.Src == "tmux" {
				activeWindowName = opts.ActiveWindowNames[session.Name]
			}

			activeWindowNamePrefix := ""
			if activeWindowName != "" {
				activeWindowNamePrefix = activeWindowName + " "
			}

			displayName = strings.NewReplacer(
				"{name}", displayName,
				"{session}", session.Name,
				"{source}", session.Src,
				"{path}", session.Path,
				"{active_window_name}", activeWindowName,
				"{active_window_name_prefix}", activeWindowNamePrefix,
			).Replace(format)
		}

		session.Name = displayName
		directory[key] = session
	}

	sessions.Directory = directory
	return sessions, nil
}
