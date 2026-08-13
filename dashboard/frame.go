package dashboard

import (
	"strings"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
)

// framePane is a single pane inside the shared frame.
type framePane struct {
	title   string
	content string
	width   int
	focused bool
}

// borderText renders s in the border color.
func borderText(s string) string {
	return lipgloss.NewStyle().Foreground(colorBorder).Render(s)
}

// renderFrame draws one shared frame around an ordered list of panes:
//
//	┌─ title ─┬─ title ─┬─ title ─┐
//	│ content │ content │ content │
//	└─────────┴─────────┴─────────┘
//
// Panes are separated by junction characters (┬ ┴ on the border lines, │ on
// content rows). Content height inside the frame is height - 2.
func renderFrame(panes []framePane, height int) string {
	if len(panes) == 0 {
		return ""
	}

	innerHeight := height - 2
	if innerHeight < 1 {
		innerHeight = 1
	}

	var b strings.Builder

	// Top border with per-pane titles.
	b.WriteString(borderText("┌"))
	for i, p := range panes {
		if i > 0 {
			b.WriteString(borderText("┬"))
		}
		b.WriteString(frameSegment(p.title, p.width, p.focused))
	}
	b.WriteString(borderText("┐"))
	b.WriteString("\n")

	// Content rows.
	for row := 0; row < innerHeight; row++ {
		b.WriteString(borderText("│"))
		for _, p := range panes {
			b.WriteString(padWidth(paneLine(p.content, row), p.width))
			b.WriteString(borderText("│"))
		}
		b.WriteString("\n")
	}

	// Bottom border.
	b.WriteString(borderText("└"))
	for i, p := range panes {
		if i > 0 {
			b.WriteString(borderText("┴"))
		}
		b.WriteString(borderText(strings.Repeat("─", p.width)))
	}
	b.WriteString(borderText("┘"))

	return b.String()
}

// frameSegment builds a top-border segment of the given width with the pane
// title embedded as "─ title ───". The focused pane's title is accent+bold,
// others are dimmed. Titles too long for the segment are truncated with "…";
// segments narrower than 4 cells render as plain "─".
func frameSegment(title string, width int, focused bool) string {
	if width < 4 {
		if width < 1 {
			width = 1
		}
		return borderText(strings.Repeat("─", width))
	}

	style := dimmedStyle()
	if focused {
		style = accentStyle()
	}

	maxTitle := width - 4
	if maxTitle < 1 {
		maxTitle = 1
	}
	if utf8.RuneCountInString(title) > maxTitle {
		title = truncateRight(title, maxTitle)
	}

	filler := width - 3 - utf8.RuneCountInString(title)
	if filler < 1 {
		filler = 1
	}

	return borderText("─") + style.Render(" "+title+" ") + borderText(strings.Repeat("─", filler))
}

// paneLine returns the row-th line of a multi-line content string, or "" when
// out of range.
func paneLine(content string, row int) string {
	lines := strings.Split(content, "\n")
	if row < 0 || row >= len(lines) {
		return ""
	}
	return lines[row]
}

// padWidth truncates s to at most width display columns (ANSI-aware) and then
// pads it to exactly width columns. It never wraps: content wider than the pane
// is truncated rather than flowed onto the next line.
func padWidth(s string, width int) string {
	if width < 1 {
		width = 1
	}
	// Truncate first (MaxWidth alone does not word-wrap when Width is unset).
	truncated := lipgloss.NewStyle().MaxWidth(width).Render(s)
	// Then pad to the exact width (no wrapping since already within width).
	return lipgloss.NewStyle().Width(width).Render(truncated)
}
