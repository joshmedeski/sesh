package dashboard

import (
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

// tabTitles are the two permanent tab labels.
var tabTitles = []string{"Open", "Configured"}

// renderHeader renders the two-row header: a tab line (with an optional
// right-pinned "N active" count on wide terminals) followed by a full-width
// rule in the border color.
func renderHeader(activePage int, activeCount int, width int) string {
	var parts []string
	for i, t := range tabTitles {
		if i == activePage {
			parts = append(parts, accentStyle().Render(t))
		} else {
			parts = append(parts, dimmedStyle().Render(t))
		}
	}
	sep := lipgloss.NewStyle().Foreground(colorBorder).Render(" │ ")
	tabLine := strings.Join(parts, sep)

	right := ""
	if width >= 50 {
		right = dimmedStyle().Render(strconv.Itoa(activeCount) + " active")
	}

	spacer := strings.Repeat(" ", max(width-lipgloss.Width(tabLine)-lipgloss.Width(right), 0))
	row1 := lipgloss.NewStyle().Width(width).Render(tabLine + spacer + right)

	rule := lipgloss.NewStyle().Foreground(colorBorder).Render(strings.Repeat("─", width))

	return row1 + "\n" + rule
}

type keybind struct {
	key   string
	label string
}

// footerBinds returns the keybinds (with their right-pinned help bind) for a
// given active page. sortLabel is the current sort mode's display label, shown
// only on page 0.
func footerBinds(page int, sortLabel string) ([]keybind, keybind) {
	right := keybind{"?", "help"}
	if page == 0 {
		return []keybind{
			{"tab", "page"},
			{"j/k", "move"},
			{"enter", "open"},
			{"s", "sort:" + sortLabel},
			{"/", "filter"},
			{"r", "refresh"},
			{"ctrl+h/j/k/l", "panes"},
			{"ctrl+d", "kill"},
			{"q", "quit"},
		}, right
	}
	return []keybind{
		{"tab", "page"},
		{"j/k", "move"},
		{"enter", "open"},
		{"r", "refresh"},
		{"q", "quit"},
	}, right
}

// renderFooter renders the one-row footer for the active page. While the
// focused pane is filtering, a filter line replaces the binds entirely. Small
// terminals: <70 cols drop labels (keys only), <30 keep only tab j/k enter q ?.
// If the fully-labeled footer would overflow the available width, labels are
// dropped so the footer never wraps.
func renderFooter(page, width int, sortLabel string, filtering bool, query string) string {
	if filtering {
		return renderFilterFooter(query, width)
	}

	binds, right := footerBinds(page, sortLabel)

	if width < 30 {
		binds = []keybind{{"tab", ""}, {"j/k", ""}, {"enter", ""}, {"q", ""}}
		right = keybind{"?", ""}
		return renderFooterLine(binds, right, false, width)
	}

	showLabels := width >= 70
	if showLabels && !footerFits(binds, right, width) {
		showLabels = false
	}
	return renderFooterLine(binds, right, showLabels, width)
}

// renderFilterFooter renders the type-to-filter status line, replacing all
// normal binds (no right-pinned help).
func renderFilterFooter(query string, width int) string {
	parts := []string{
		accentStyle().Render("filter:") + " " + textStyle().Render(query),
		accentStyle().Render("esc") + " " + dimmedStyle().Render("clear"),
		accentStyle().Render("enter") + " " + dimmedStyle().Render("done"),
	}
	line := strings.Join(parts, lipgloss.NewStyle().Foreground(colorBorder).Render(" │ "))
	return lipgloss.NewStyle().Width(width).Render(line)
}

// footerFits reports whether the footer with labels would fit within width.
func footerFits(binds []keybind, right keybind, width int) bool {
	left := 0
	for i, b := range binds {
		if i > 0 {
			left += 3 // " │ "
		}
		left += len(b.key)
		if b.label != "" {
			left += 1 + len(b.label)
		}
	}
	rightW := len(right.key)
	if right.label != "" {
		rightW += 1 + len(right.label)
	}
	return left+rightW <= width
}

func renderFooterLine(binds []keybind, right keybind, showLabels bool, width int) string {
	format := func(b keybind) string {
		if showLabels && b.label != "" {
			return accentStyle().Render(b.key) + " " + dimmedStyle().Render(b.label)
		}
		return accentStyle().Render(b.key)
	}

	var parts []string
	for _, b := range binds {
		parts = append(parts, format(b))
	}
	left := strings.Join(parts, lipgloss.NewStyle().Foreground(colorBorder).Render(" │ "))
	rightStr := format(right)

	spacer := strings.Repeat(" ", max(width-lipgloss.Width(left)-lipgloss.Width(rightStr), 0))
	return lipgloss.NewStyle().Width(width).Render(left + spacer + rightStr)
}
