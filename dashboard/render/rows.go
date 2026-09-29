// rows.go
package render

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

// RowMarker returns the 2-column marker for a row: "▌ " (accent bold on the
// selection background) when selected, otherwise two spaces. The marker is
// dimmed when the pane is not focused.
func RowMarker(selected, focused bool) string {
	if !selected {
		return "  "
	}
	if focused {
		return accentStyle().Background(colorHighlight).Render("▌ ")
	}
	return DimmedStyle().Background(colorHighlightDim).Render("▌ ")
}

// Col describes a single column in a list row.
type Col struct {
	Text  string
	Width int
	Style lipgloss.Style
	Align lipgloss.Position
}

// RenderRow joins columns with a single space and applies the shared cursor
// background to every cell (and separator) plus the marker when selected,
// producing a vim-like full-line highlight.
func RenderRow(marker string, cols []Col, selected, focused bool) string {
	bg := cursorStyle(focused)
	rendered := make([]string, len(cols))
	for i, c := range cols {
		st := c.Style
		if selected {
			st = st.Inherit(bg)
		}
		st = st.Width(c.Width)
		if c.Align != 0 {
			st = st.Align(c.Align)
		}
		rendered[i] = st.Render(c.Text)
	}
	if selected {
		return marker + strings.Join(rendered, bg.Render(" "))
	}
	return marker + strings.Join(rendered, " ")
}

// RenderSimpleRow renders a widget-section row (git, ssh, docker) with the
// shared cursor marker and a vim-like full-line highlight when selected. Cells
// are raw text with a foreground style and are concatenated directly (no
// separator is inserted), so any spacing or column padding must live in the
// cell text or its style; cell content is otherwise preserved exactly.
func RenderSimpleRow(cells []Col, selected, focused bool) string {
	bg := cursorStyle(focused)
	rendered := make([]string, len(cells))
	for i, c := range cells {
		st := c.Style
		if selected {
			st = st.Inherit(bg)
		}
		rendered[i] = st.Render(c.Text)
	}
	return RowMarker(selected, focused) + strings.Join(rendered, "")
}

// Powerline half circles used to round off the alias chip, matching the
// picker’s default alias styling.
const (
	chipLeftGlyph  = "\ue0b6"
	chipRightGlyph = "\ue0b4"
)

// aliasChip renders a configured alias as a rounded pill. It returns an empty
// string when there is no alias. The chip uses reverse video over the dashboard
// accent color so the background is the theme accent and the foreground is the
// terminal’s own contrasting background color, matching the picker’s alias
// chip conventions. The half circles are painted with the same accent so the
// pill reads as one shape.
func aliasChip(alias string) string {
	if alias == "" {
		return ""
	}
	fill := lipgloss.NewStyle().Foreground(colorAccent)
	label := fill.Reverse(true).Render(alias)
	return fill.Render(chipLeftGlyph) + label + fill.Render(chipRightGlyph) + " "
}

// aliasChipSelected renders the alias pill for a selected row as a continuous
// ANSI run. The pill keeps an explicit accent background of its own rather than
// inheriting the cursor highlight, so the alias reads the same as its
// unselected chip. The label text is the contrasting terminal-background colour
// (colorPillText), while the powerline half circles stay accent-coloured over
// the cursor background, so the rounded edges read as one shape with the pill
// and the row highlight continues unbroken outside it. The cursor background is
// restored after the pill so the name that follows stays on the highlight.
func aliasChipSelected(alias string) string {
	accent := ansiFg(colorAccent)
	return accent + chipLeftGlyph +
		ansiBg(colorAccent) + ansiFg(colorPillText) + alias +
		ansiBg(colorHighlight) + accent + chipRightGlyph + " "
}

// ansiBg returns an ANSI 256-colour background sequence for c.
func ansiBg(c lipgloss.ANSIColor) string {
	return fmt.Sprintf("\x1b[48;5;%dm", int(c))
}

// namePrefix returns the raw ANSI prefix (bold + foreground) for the session
// name, without a trailing reset, so it can be embedded in a continuous ANSI
// run when the row is selected.
func namePrefix(current bool) string {
	if current {
		return "\x1b[1m" + ansiFg(colorAccent)
	}
	return ansiFg(colorText)
}

// renderOpenRow renders a Tab 1 (Open) session row with columns:
// marker(2) | alias+name(22) | att(2) | windows(5) | dir(fill) | branch(16) |
// status(12) | age(5, last attached) | alerts(2).
// Progressive drop: <90 cols drop status+age+alerts, <70 drop branch+att,
// <50 drop windows.
func renderOpenRow(width int, selected, current bool, name, alias string, attached, windows int, dir, branch, status string, lastAttached *time.Time, alerts []string) string {
	return RenderOpenRowFocused(width, selected, current, true, name, alias, attached, windows, dir, branch, status, lastAttached, alerts)
}

// RenderOpenRowFocused is renderOpenRow with an explicit focused flag, so
// unfocused panes render a dimmed selection highlight.
func RenderOpenRowFocused(width int, selected, current, focused bool, name, alias string, attached, windows int, dir, branch, status string, lastAttached *time.Time, alerts []string) string {
	includeWindows := width >= 50
	includeBranch := width >= 70
	includeAtt := width >= 70
	includeStatus := width >= 90
	includeAge := width >= 90
	includeAlerts := width >= 90

	fixed := 22
	if includeAtt {
		fixed += 2
	}
	if includeWindows {
		fixed += 5
	}
	if includeBranch {
		fixed += 16
	}
	if includeStatus {
		fixed += 12
	}
	if includeAge {
		fixed += 5
	}
	if includeAlerts {
		fixed += 2
	}

	numCols := 2 // name + dir
	if includeAtt {
		numCols++
	}
	if includeWindows {
		numCols++
	}
	if includeBranch {
		numCols++
	}
	if includeStatus {
		numCols++
	}
	if includeAge {
		numCols++
	}
	if includeAlerts {
		numCols++
	}

	dirWidth := max(width-2-fixed-(numCols-1), 1)

	nameStyle := TextStyle()
	if current {
		nameStyle = accentStyle()
	}

	cols := []Col{}

	if includeAtt {
		attText := ""
		if attached > 0 {
			attText = SuccessStyle().Render("●")
		}
		cols = append(cols, Col{Text: attText, Width: 2})
	}

	chip := aliasChip(alias)
	nameBudget := 22
	if chip != "" {
		nameBudget -= lipgloss.Width(chip)
		if nameBudget < 1 {
			nameBudget = 1
		}
	}
	truncated := TruncateRight(name, nameBudget)
	var nameText string
	switch {
	case selected && chip != "":
		// One continuous ANSI run so the cursor background survives the pill
		// and the name (no nested resets from pre-rendered text).
		nameText = aliasChipSelected(alias) + namePrefix(current) + truncated
	case selected:
		nameText = namePrefix(current) + truncated
	default:
		nameText = nameStyle.Render(truncated)
		if chip != "" {
			nameText = chip + nameText
		}
	}
	cols = append(cols, Col{Text: nameText, Width: 18})

	// if includeWindows {
	// 	cols = append(cols, Col{Text: fmt.Sprintf("%2dw", windows), Width: 7, Style: TextStyle(), Align: lipgloss.Left})
	// }
	cols = append(cols, Col{Text: truncateDirLeft(dir, dirWidth), Width: dirWidth, Style: TextStyle()})
	if includeBranch {
		cols = append(cols, Col{Text: TruncateRight(Paren(branch), 16), Width: 14, Style: BranchStyle(), Align: lipgloss.Left})
	}
	if includeStatus {
		cols = append(cols, Col{Text: TruncateRightANSI(status, 12), Width: 12, Style: BranchStyle()})
	}
	if includeAge {
		cols = append(cols, Col{Text: formatAge(lastAttached), Width: 5, Style: ageStyle(), Align: lipgloss.Left})
	}
	if includeAlerts {
		alertText := ""
		if len(alerts) > 0 {
			alertText = WarningStyle().Render("!")
		}
		cols = append(cols, Col{Text: alertText, Width: 2})
	}

	return RenderRow(RowMarker(selected, focused), cols, selected, focused)
}

// renderConfiguredRow renders a Tab 2 (Configured) session row with columns:
// marker(2) | cmd(2) | name(24) | state(2) | path(fill) | branch(16) |
// status(12). The cmd column ("*") and branch drop together below 70 cols.
func renderConfiguredRow(width int, selected bool, name, startupCommand string, running bool, path, branch, status string) string {
	return RenderConfiguredRowFocused(width, selected, true, name, startupCommand, running, path, branch, status)
}

// RenderConfiguredRowFocused is renderConfiguredRow with an explicit focused
// flag, so unfocused panes render a dimmed selection highlight.
func RenderConfiguredRowFocused(width int, selected, focused bool, name, startupCommand string, running bool, path, branch, status string) string {
	includeBranch := width >= 70

	stateText := "○"
	stateColStyle := DimmedStyle()
	if running {
		stateText = "●"
		stateColStyle = SuccessStyle()
	}

	if path == "" {
		path = "-"
	}

	fixed := 24 + 2 + 12 // name + state + status
	numCols := 4         // name + state + path + status
	if includeBranch {
		fixed += 2 + 16 // cmd + branch
		numCols += 2
	}
	pathWidth := width - 2 - fixed - (numCols - 1)
	if pathWidth < 1 {
		pathWidth = 1
	}

	cols := make([]Col, 0, numCols)
	if includeBranch {
		cmdText := ""
		// if startupCommand != "" {
		// 	cmdText = WarningStyle().Render("*")
		// }
		cols = append(cols, Col{Text: cmdText, Width: 2})
	}
	cols = append(cols, Col{Text: stateText, Width: 2, Style: stateColStyle})
	cols = append(cols, Col{Text: TruncateRight(name, 24), Width: 24, Style: TextStyle()})
	cols = append(cols, Col{Text: truncateDirLeft(path, pathWidth), Width: pathWidth, Style: TextStyle()})
	if includeBranch {
		cols = append(cols, Col{Text: TruncateRight(Paren(branch), 16), Width: 16, Style: BranchStyle()})
	}
	cols = append(cols, Col{Text: TruncateRightANSI(status, 12), Width: 12})

	return RenderRow(RowMarker(selected, focused), cols, selected, focused)
}
