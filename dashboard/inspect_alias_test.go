package dashboard

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestInspectAliasRowSelected(t *testing.T) {
	row := renderOpenRow(100, true, false, "wallpaper", "wp", 0, 1, "~/d", "", "", nil, nil)
	fmt.Printf("SELECTED ROW:\n%s\n\n", row)
	fmt.Printf("STRIPPED:\n%s\n\n", ansi.Strip(row))
	fmt.Printf("CONTAINS 48;5;8 BEFORE NAME: %v\n", stringsContainsInOrder(row, "48;5;8m", "wallpaper"))
	fmt.Printf("CONTAINS 48;5;8 AFTER WP: %v\n", stringsContainsInOrder(row, "wp", "48;5;8m"))
}

func stringsContainsInOrder(s, a, b string) bool {
	ia := strings.Index(s, a)
	ib := strings.Index(s, b)
	return ia >= 0 && ib > ia
}
