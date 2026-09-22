// styles.go
package render

import (
	"charm.land/lipgloss/v2"
)

// Shared color palette (see design spec).
var (
	colorAccent = lipgloss.ANSIColor(14) // accent (cyan)
	colorDimmed = lipgloss.ANSIColor(8)  // dimmed / border
	colorBorder = lipgloss.ANSIColor(8)
	colorText   = lipgloss.ANSIColor(15) // white text
	// colorPillText is the alias pill's label colour on a selected row. Black
	// mirrors the terminal's typical default background, so it stays legible
	// on the accent pill fill the same way the unselected chip's reverse-video
	// label does (its text is the terminal's own contrasting background).
	colorPillText  = lipgloss.ANSIColor(0)  // black
	colorAge       = lipgloss.ANSIColor(7)  // light gray (distinct from highlight bg)
	colorBranch    = lipgloss.ANSIColor(5)  // magenta
	colorStatus    = lipgloss.ANSIColor(10) // green
	colorHighlight = lipgloss.ANSIColor(8)

	// colorHighlightDim is the slightly dimmer selection background used when a
	// pane is not focused, so the selection stays visible without competing
	// with the focused pane's highlight.
	colorHighlightDim = lipgloss.ANSIColor(236)

	// git-status part colours (formatGitStatus).
	colorStaged    = lipgloss.ANSIColor(10) // green  "+N"
	colorUnstaged  = lipgloss.ANSIColor(11) // yellow "~N"
	ColorDeleted   = lipgloss.ANSIColor(9)  // red    "-N"
	colorUntracked = lipgloss.ANSIColor(5)  // magenta "!N"
)

func accentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
}

// WarningStyle is the yellow style used for alerts and the startup-command
// indicator.
func WarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorUnstaged)
}

// DimmedStyle is the muted foreground used for secondary text and unfocused
// elements.
func DimmedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorDimmed)
}

// ageStyle is the neutral, muted foreground for the Open sessions age column.
// It deliberately avoids colorDimmed (the cursor background) so the age stays
// readable when a row is highlighted.
func ageStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorAge)
}

// TextStyle is the default white foreground for row text.
func TextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorText)
}

// BranchStyle is the magenta foreground used for git branch names.
func BranchStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorBranch)
}

// SuccessStyle is the green foreground used for running/attached indicators.
func SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorStatus)
}

// GroupNameRender is retained for backwards compatibility with the widget
// sections that render a single-line title bar.
func GroupNameRender(name string, width int) lipgloss.Style {
	return NewStyle(width, width, 1, 1, 15, false, []int{0, 0, 0, 0})
}

// NewStyle returns a new lipgloss style with the given parameters
func NewStyle(width int, maxWidth int, height int, maxHeight int, color int, faint bool, padding []int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		MaxWidth(maxWidth).
		Height(height).
		MaxHeight(maxHeight).
		Foreground(lipgloss.ANSIColor(color)).
		Faint(faint).
		Padding(padding[0], padding[1], padding[2], padding[3])
}

// cursorStyle is the shared selection highlight used by every list section: a
// subtle full-line background fill. Focused panes use the standard highlight;
// unfocused panes use a slightly dimmer fill.
func cursorStyle(focused bool) lipgloss.Style {
	if focused {
		return lipgloss.NewStyle().Background(colorHighlight)
	}
	return lipgloss.NewStyle().Background(colorHighlightDim)
}
