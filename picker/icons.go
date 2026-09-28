package picker

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

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
