package dashboard

import (
	"fmt"
	"strings"
)

func renderFooter(width int) string {
	if width < 10 {
		return ""
	}
	controls := "  j/k Navigate  |  Enter Attach  |  t Toggle  |  q Exit"
	return fmt.Sprintf("├─ CONTROLS %s┤\n%s\n", strings.Repeat("─", width-14), controls)
}

func renderHeader(title string, totalItems int, width int) string {
	if width < 10 {
		return title
	}
	right := fmt.Sprintf("Total: %d", totalItems)
	return fmt.Sprintf("  %s  %s", title, right)
}
