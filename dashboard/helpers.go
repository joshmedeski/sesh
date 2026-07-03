package dashboard

import "charm.land/lipgloss/v2"

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

// NewStyleBorder returns a new lipgloss style with the given parameters and border
func NewStyleBorder(width int, maxWidth int, height int, maxHeight int, color int, faint bool, padding []int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		MaxWidth(maxWidth).
		Height(height).
		MaxHeight(maxHeight).
		Foreground(lipgloss.ANSIColor(color)).
		Faint(faint).
		Padding(padding[0], padding[1], padding[2], padding[3]).
		Border(lipgloss.RoundedBorder(), true, true, true, true)
}
