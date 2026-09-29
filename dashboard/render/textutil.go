// textutil.go
package render

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/x/ansi"
)

// formatAge renders a compact relative age ("2h"/"3d"/"4mo") for a session's
// last attached time, or "" when lastAttached is nil/zero.
func formatAge(lastAttached *time.Time) string {
	if lastAttached == nil || lastAttached.IsZero() {
		return ""
	}
	since := time.Since(*lastAttached)
	if since < 0 {
		return ""
	}
	switch {
	case since < time.Hour:
		m := int(since.Minutes())
		if m < 1 {
			m = 1
		}
		return fmt.Sprintf("%dm", m)
	case since < 24*time.Hour:
		return fmt.Sprintf("%dh", int(since.Hours()))
	case since < 30*24*time.Hour:
		return fmt.Sprintf("%dd", int(since.Hours()/24))
	default:
		return fmt.Sprintf("%dmo", int(since.Hours()/(24*30)))
	}
}

// CollapseHome replaces the home directory prefix of path with "~".
func CollapseHome(path, homeDir string) string {
	if homeDir == "" {
		return path
	}
	if path == homeDir {
		return "~"
	}
	if strings.HasPrefix(path, homeDir+string(filepath.Separator)) {
		return "~" + path[len(homeDir):]
	}
	return path
}

// Paren wraps s in parentheses when non-empty, otherwise returns "".
func Paren(s string) string {
	if s == "" {
		return ""
	}
	return "(" + s + ")"
}

// TruncateRight truncates s to maxRunes runes, appending "…" when truncated.
func TruncateRight(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	if maxRunes <= 1 {
		return "…"
	}
	return string(runes[:maxRunes-1]) + "…"
}

// TruncateRightANSI truncates an ANSI-styled string to maxWidth visible cells,
// appending "…" when truncated. Escape sequences are preserved so embedded
// colours (e.g. the git-status parts) survive truncation.
func TruncateRightANSI(s string, maxWidth int) string {
	return ansi.Truncate(s, maxWidth, "…")
}

// truncateDirLeft truncates a directory path from the left, preferring to keep
// the last two path segments (so the directory name and its parent stay
// visible). Falls back to a "…"-prefixed tail when segments are too long.
func truncateDirLeft(dir string, maxRunes int) string {
	if utf8.RuneCountInString(dir) <= maxRunes {
		return dir
	}
	if kept := lastTwoSegments(dir); utf8.RuneCountInString(kept) <= maxRunes {
		return kept
	}
	if maxRunes <= 1 {
		return "…"
	}
	runes := []rune(dir)
	return "…" + string(runes[len(runes)-(maxRunes-1):])
}

// lastTwoSegments keeps the leading "~/" (if present) plus the final two path
// segments of dir.
func lastTwoSegments(dir string) string {
	prefix := ""
	rest := dir
	if strings.HasPrefix(rest, "~/") {
		prefix = "~/"
		rest = rest[2:]
	} else if rest == "~" {
		return "~"
	}
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) <= 2 {
		return prefix + strings.Join(parts, "/")
	}
	return prefix + strings.Join(parts[len(parts)-2:], "/")
}
