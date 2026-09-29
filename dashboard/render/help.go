package render

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// RenderHelp renders the keybinding reference for the active page, padded to
// width x height. sortable adds the sort bind for pages whose list can be
// re-sorted.
func RenderHelp(page, width, height int, sortable bool) string {
	binds := []keybind{
		{"tab / shift+tab", "next / previous tab"},
		{"j/k ↑/↓", "move"},
		{"enter", "open"},
		{"/", "filter"},
		{"r", "refresh"},
	}
	if sortable {
		binds = append(binds, keybind{"s", "cycle sort"})
	}
	if page == 0 {
		binds = append(binds,
			keybind{"ctrl+h / ctrl+l", "focus pane left / right"},
			keybind{"ctrl+j / ctrl+k", "focus pane below / above"},
			keybind{"1-9", "focus pane"},
			keybind{"ctrl+d", "kill tmux session"},
			keybind{"click", "focus pane and select row"},
		)
	}
	binds = append(binds,
		keybind{"?", "close help"},
		keybind{"q / esc", "quit"},
	)

	keyWidth := 0
	for _, b := range binds {
		keyWidth = max(keyWidth, lipgloss.Width(b.key))
	}
	lines := []string{"", "  " + accentStyle().Render("Keybindings"), ""}
	for _, b := range binds {
		lines = append(lines, "  "+accentStyle().Width(keyWidth).Render(b.key)+"  "+TextStyle().Render(b.label))
	}
	return lipgloss.NewStyle().Width(width).Height(height).MaxHeight(height).Render(strings.Join(lines, "\n"))
}
