package dashboard

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func renderFooter(width int) string {
	// if width < 10 {
	// 	return ""
	// }
	controls := "\nj/k Navigate  |  Enter Attach  |  t Toggle  |  Ctrl+d Kill  |  q Exit"
	controlsStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8)).Faint(true)
	colStyle := lipgloss.NewStyle().Width(width)
	controls = colStyle.Render(controlsStyle.Render(controls))
	return fmt.Sprintf("%s", controls)
}

func renderHeader(title string, totalItems int, width int) string {
	if width < 10 {
		return title
	}
	right := fmt.Sprintf("Active sessions: %d", totalItems)
	sep := strings.Repeat(" ", max(width-(len(title)+len(right))-2, 0))
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8))
	headStyle := lipgloss.NewStyle().Width(width)
	header := headStyle.Render(headerStyle.Render(title + sep + " " + right))
	return fmt.Sprintf("%s", header)
}
