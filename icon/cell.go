package icon

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

func (glyphs Glyphs) Source(src string) (string, color.Color) {
	if g, ok := glyphs[src]; ok {
		var code int
		switch {
		case g.ColorCode >= 90 && g.ColorCode <= 97:
			code = g.ColorCode - 82
		case g.ColorCode >= 30 && g.ColorCode <= 37:
			code = g.ColorCode - 30
		default:
			code = g.ColorCode
		}
		return g.Icon + " ", lipgloss.ANSIColor(code)
	}
	return "? ", lipgloss.ANSIColor(8)
}

// ColumnWidth is the display width an icon column needs: one cell for the
// single-width source glyphs, more when a configured icon is wider (emoji are
// two cells). It is measured across the whole config rather than the visible
// rows so the column keeps one width.
func ColumnWidth(config model.Config) int {
	width := 1
	for _, session := range config.SessionConfigs {
		width = max(width, Width(session.Icon))
	}
	for _, wildcard := range config.WildcardConfigs {
		width = max(width, Width(wildcard.Icon))
	}
	for _, source := range config.Icons {
		width = max(width, Width(source.Icon))
	}
	return width
}

// Width is the display width of a configured icon. Trailing spaces are the
// user's own width override and are not counted.
func Width(icn string) int {
	return lipgloss.Width(strings.TrimRight(icn, " "))
}

// Cell renders an icon column cell of width+1 cells: the configured icon, or
// the colored source glyph when there is none, followed by the gap before the
// next column. A configured icon is left unstyled.
func (glyphs Glyphs) Cell(custom, src string, width int) string {
	if custom == "" {
		glyph, clr := glyphs.Source(src)
		cell := lipgloss.NewStyle().Foreground(clr).Render(glyph)
		return cell + pad(width+1-lipgloss.Width(cell))
	}
	return custom + " " + pad(width+1-lipgloss.Width(strings.TrimRight(custom, " ")+" "))
}

func pad(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}
