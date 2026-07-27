package picker

import (
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
)

// WildcardFinder resolves a path to the [[wildcard]] block that matches it. The
// picker borrows the lister's matching so an icon covers exactly the paths that
// startup_command and preview_command already do, first match in config order
// included.
type WildcardFinder interface {
	FindConfigWildcard(path string) (model.WildcardConfig, bool)
}

// buildIconResolver indexes the icons declared in config into a lookup for the
// picker. It returns nil when nothing declares one, so the common case costs the
// TUI nothing per row.
//
// The most specific match wins: an exact [[session]] name, then a [[session]]
// path, then a [[wildcard]] pattern. Sessions are indexed by path as well as by
// name because the same directory is often listed by another source under a
// derived name — a zoxide entry for a configured session's path still gets its
// icon.
func buildIconResolver(config model.Config, h home.Home, wildcards WildcardFinder) IconFunc {
	byName := make(map[string]string)
	byPath := make(map[string]string)
	for _, session := range config.SessionConfigs {
		if session.Icon == "" {
			continue
		}
		if session.Name != "" {
			byName[session.Name] = session.Icon
		}
		if session.Path == "" {
			continue
		}
		// An unexpandable path is skipped rather than fatal: the icon is
		// cosmetic, and the lister already reports a broken path.
		if path, err := h.ExpandPath(session.Path); err == nil {
			byPath[filepath.Clean(path)] = session.Icon
		}
	}

	hasWildcardIcon := false
	for _, wildcard := range config.WildcardConfigs {
		if wildcard.Icon != "" {
			hasWildcardIcon = true
			break
		}
	}

	if len(byName) == 0 && len(byPath) == 0 && !hasWildcardIcon {
		return nil
	}

	return func(session model.SeshSession) string {
		if icn, ok := byName[session.Name]; ok {
			return icn
		}
		if session.Path == "" {
			return ""
		}
		if icn, ok := byPath[filepath.Clean(session.Path)]; ok {
			return icn
		}
		if !hasWildcardIcon || wildcards == nil {
			return ""
		}
		// A first match that declares no icon still wins, matching how the rest
		// of the wildcard config resolves: the session falls back to its source
		// glyph rather than searching on for a later pattern with an icon.
		if wildcard, ok := wildcards.FindConfigWildcard(session.Path); ok {
			return wildcard.Icon
		}
		return ""
	}
}

// iconColWidth is the display width the picker's icon column needs: one cell for
// the single-width source glyphs, more when a configured icon is wider (emoji
// are two cells).
//
// It is measured across the whole config rather than the visible rows so the
// column keeps one width — scrolling past an emoji can't shift the names, and a
// row with a narrow icon lines up with one that has a wide icon.
//
// Trailing spaces are ignored, so padding an icon by hand to match how the
// terminal draws it doesn't widen the column for every other row. See
// Model.iconCell.
func iconColWidth(config model.Config) int {
	width := 1
	for _, session := range config.SessionConfigs {
		width = max(width, iconWidth(session.Icon))
	}
	for _, wildcard := range config.WildcardConfigs {
		width = max(width, iconWidth(wildcard.Icon))
	}
	return width
}

func iconWidth(icn string) int {
	return lipgloss.Width(strings.TrimRight(icn, " "))
}
