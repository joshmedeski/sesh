// rows.go
package render

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/joshmedeski/sesh/v2/icon"
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
	Title string
	Text  string
	Width int
	Style lipgloss.Style
	Align lipgloss.Position
}

// RenderColumnTitles renders the column titles of cols, aligned with RenderRow.
func RenderColumnTitles(cols []Col) string {
	titles := make([]Col, len(cols))
	for i, c := range cols {
		titles[i] = Col{Text: c.Title, Width: c.Width, Style: DimmedStyle().Bold(true)}
	}
	return RenderRow("  ", titles, false, true)
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
		rendered[i] = st.Render(ansi.Truncate(c.Text, c.Width, "…"))
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

const (
	pillLeftGlyph  = "\ue0b6"
	pillRightGlyph = "\ue0b4"
)

func aliasPill(alias string, selected, focused bool) string {
	if alias == "" {
		return ""
	}
	if !selected {
		fill := lipgloss.NewStyle().Foreground(colorText)
		return fill.Render(pillLeftGlyph) + fill.Reverse(true).Render(alias) + fill.Render(pillRightGlyph)
	}
	rowBg := colorHighlight
	if !focused {
		rowBg = colorHighlightDim
	}
	fg := ansiFg(colorText)
	return fg + pillLeftGlyph +
		ansiBg(colorText) + ansiFg(colorPillText) + alias +
		ansiBg(rowBg) + fg + pillRightGlyph
}

func ansiBg(c lipgloss.ANSIColor) string {
	return fmt.Sprintf("\x1b[48;5;%dm", int(c))
}

// AliasColumn returns the alias column width for the given aliases: wide
// enough for the longest alias and its header, or 0 when none are set.
func AliasColumn(aliases ...string) int {
	w := 0
	for _, a := range aliases {
		w = max(w, lipgloss.Width(a))
	}
	if w == 0 {
		return 0
	}
	return max(w+lipgloss.Width(pillLeftGlyph+pillRightGlyph), len("ALIAS"))
}

// renderOpenRow renders a Tab 1 (Open) session row with columns:
// marker(2) | icon | att(2) | name(18) | alias(longest) | dir(fill) |
// branch(longest) | status(12) | age(5, last attached) | alerts(2).
// Progressive drop: <90 cols drop status+age+alerts, <70 drop branch+att.
func renderOpenRow(width int, selected, current bool, name, alias string, attached, windows int, dir, branch, status string, lastAttached *time.Time, alerts []string) string {
	return RenderOpenRowFocused(width, lipgloss.Width(branch), AliasColumn(alias), IconCol("", "tmux", 1, selected), selected, current, true, name, alias, attached, windows, dir, branch, status, lastAttached, alerts)
}

// RenderOpenRowFocused is renderOpenRow with an explicit focused flag, so
// unfocused panes render a dimmed selection highlight.
func RenderOpenRowFocused(width, branchCol, aliasCol int, iconCol Col, selected, current, focused bool, name, alias string, attached, windows int, dir, branch, status string, lastAttached *time.Time, alerts []string) string {
	cols := openCols(width, branchCol, aliasCol, iconCol, selected, current, focused, name, alias, attached, dir, branch, status, lastAttached, alerts)
	return RenderRow(RowMarker(selected, focused), cols, selected, focused)
}

// RenderOpenHeader renders the Open tab column titles.
func RenderOpenHeader(width, branchCol, aliasCol, iconWidth int) string {
	return RenderColumnTitles(openCols(width, branchCol, aliasCol, Col{Width: iconWidth}, false, false, true, "", "", 0, "", "", "", nil, nil))
}

func openCols(width, branchCol, aliasCol int, iconCol Col, selected, current, focused bool, name, alias string, attached int, dir, branch, status string, lastAttached *time.Time, alerts []string) []Col {
	includeAtt := width >= 70
	includeBranch := width >= 70
	includeStatus := width >= 90
	includeAge := width >= 90
	includeAlerts := width >= 90

	fixed := 18 + iconCol.Width
	numCols := 3
	for _, c := range []struct {
		on    bool
		width int
	}{{includeAtt, 2}, {aliasCol > 0, aliasCol}, {includeBranch, 0}, {includeStatus, 12}, {includeAge, 5}, {includeAlerts, 2}} {
		if c.on {
			fixed += c.width
			numCols++
		}
	}

	flex := width - 2 - fixed - (numCols - 1)
	branchWidth := 0
	if includeBranch {
		branchWidth = branchColumn(flex, branchCol)
	}
	dirWidth := max(flex-branchWidth, 1)

	nameStyle := TextStyle()
	if current {
		nameStyle = accentStyle()
	}

	cols := []Col{iconCol}
	if includeAtt {
		attText := ""
		if attached > 0 {
			attText = SuccessStyle().Render("●")
		}
		cols = append(cols, Col{Text: attText, Width: 2})
	}
	cols = append(cols, Col{Title: "NAME", Text: TruncateRight(name, 18), Width: 18, Style: nameStyle})
	if aliasCol > 0 {
		cols = append(cols, Col{Title: "ALIAS", Text: aliasPill(alias, selected, focused), Width: aliasCol})
	}
	cols = append(cols, Col{Title: "DIRECTORY", Text: truncateDirLeft(dir, dirWidth), Width: dirWidth, Style: TextStyle()})
	if includeBranch {
		cols = append(cols, Col{Title: "BRANCH", Text: TruncateRight(branch, branchWidth), Width: branchWidth, Style: BranchStyle(), Align: lipgloss.Left})
	}
	if includeStatus {
		cols = append(cols, Col{Title: "STATUS", Text: TruncateRightANSI(status, 12), Width: 12, Style: BranchStyle()})
	}
	if includeAge {
		cols = append(cols, Col{Title: "AGE", Text: formatAge(lastAttached), Width: 5, Style: ageStyle(), Align: lipgloss.Left})
	}
	if includeAlerts {
		alertText := ""
		if len(alerts) > 0 {
			alertText = WarningStyle().Render("!")
		}
		cols = append(cols, Col{Text: alertText, Width: 2})
	}
	return cols
}

// renderConfiguredRow renders a Tab 2 (Configured) session row with columns:
// marker(2) | icon | cmd(2) | state(2) | name(24) | alias(longest) |
// path(fill) | branch(16) | status(12). The cmd column and branch drop
// together below 70 cols.
func renderConfiguredRow(width int, selected bool, name, alias, startupCommand string, running bool, path, branch, status string) string {
	return RenderConfiguredRowFocused(width, AliasColumn(alias), IconCol("", "config", 1, selected), selected, true, name, alias, startupCommand, running, path, branch, status)
}

// RenderConfiguredRowFocused is renderConfiguredRow with an explicit focused
// flag, so unfocused panes render a dimmed selection highlight.
func RenderConfiguredRowFocused(width, aliasCol int, iconCol Col, selected, focused bool, name, alias, startupCommand string, running bool, path, branch, status string) string {
	cols := configuredCols(width, aliasCol, iconCol, selected, focused, name, alias, running, path, branch, status)
	return RenderRow(RowMarker(selected, focused), cols, selected, focused)
}

// RenderConfiguredHeader renders the Configured tab column titles.
func RenderConfiguredHeader(width, aliasCol, iconWidth int) string {
	cols := configuredCols(width, aliasCol, Col{Width: iconWidth}, false, true, "", "", false, "", "", "")
	return RenderColumnTitles(cols)
}

func configuredCols(width, aliasCol int, iconCol Col, selected, focused bool, name, alias string, running bool, path, branch, status string) []Col {
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

	fixed := iconCol.Width + 24 + 2 + 12
	numCols := 5
	if aliasCol > 0 {
		fixed += aliasCol
		numCols++
	}
	if includeBranch {
		fixed += 2 + 16
		numCols += 2
	}
	pathWidth := max(width-2-fixed-(numCols-1), 1)

	cols := make([]Col, 0, numCols)
	cols = append(cols, iconCol)
	if includeBranch {
		cols = append(cols, Col{Width: 2})
	}
	cols = append(cols, Col{Text: stateText, Width: 2, Style: stateColStyle})
	cols = append(cols, Col{Title: "NAME", Text: TruncateRight(name, 24), Width: 24, Style: TextStyle()})
	if aliasCol > 0 {
		cols = append(cols, Col{Title: "ALIAS", Text: aliasPill(alias, selected, focused), Width: aliasCol})
	}
	cols = append(cols, Col{Title: "DIRECTORY", Text: truncateDirLeft(path, pathWidth), Width: pathWidth, Style: TextStyle()})
	if includeBranch {
		cols = append(cols, Col{Title: "BRANCH", Text: TruncateRight(branch, 16), Width: 16, Style: BranchStyle()})
	}
	cols = append(cols, Col{Title: "STATUS", Text: TruncateRightANSI(status, 12), Width: 12})
	return cols
}

func branchColumn(flex, longest int) int {
	return min(max(longest, 16), max(flex-16, 16))
}

// RenderWorktreeRowFocused renders a worktree tab row with columns:
// marker(2) | icon | number(7) | state(6) | title(fill) | branch(longest) |
// status(12) | age(5, since created). The branch drops below 70 cols. Closed
// issues render their title dimmed.
func RenderWorktreeRowFocused(width, branchCol int, iconCol Col, selected, focused bool, number int, title, state, branch, status string, created *time.Time) string {
	cols := worktreeCols(width, branchCol, iconCol, selected, "#"+strconv.Itoa(number), title, state, branch, status, created)
	return RenderRow(RowMarker(selected, focused), cols, selected, focused)
}

// RenderWorktreeHeader renders the worktree tab column titles.
func RenderWorktreeHeader(width, branchCol, iconWidth int) string {
	return RenderColumnTitles(worktreeCols(width, branchCol, Col{Width: iconWidth}, false, "", "", "", "", "", nil))
}

var issueStateColors = map[string]lipgloss.ANSIColor{
	"OPEN":   lipgloss.ANSIColor(2),
	"MERGED": lipgloss.ANSIColor(5),
	"CLOSED": lipgloss.ANSIColor(1),
}

func worktreeCols(width, branchCol int, iconCol Col, selected bool, number, title, state, branch, status string, created *time.Time) []Col {
	includeBranch := width >= 70

	fixed := iconCol.Width + 7 + 6 + 12 + 5
	numCols := 6
	if includeBranch {
		numCols++
	}
	flex := width - 2 - fixed - (numCols - 1)
	branchWidth := 0
	if includeBranch {
		branchWidth = branchColumn(flex, branchCol)
	}
	titleWidth := max(flex-branchWidth, 1)

	numberStyle := DimmedStyle()
	titleStyle := TextStyle()
	if state == "CLOSED" {
		titleStyle = DimmedStyle()
	}
	if selected {
		numberStyle = TextStyle()
		titleStyle = TextStyle()
	}

	cols := []Col{
		iconCol,
		{Title: "ISSUE", Text: number, Width: 7, Style: numberStyle},
		{Title: "STATE", Text: strings.ToLower(state), Width: 6, Style: lipgloss.NewStyle().Foreground(issueStateColors[state])},
		{Title: "TITLE", Text: TruncateRight(title, titleWidth), Width: titleWidth, Style: titleStyle},
	}
	if includeBranch {
		cols = append(cols, Col{Title: "BRANCH", Text: TruncateRight(branch, branchWidth), Width: branchWidth, Style: BranchStyle()})
	}
	cols = append(cols, Col{Title: "STATUS", Text: TruncateRightANSI(status, 12), Width: 12, Style: BranchStyle()})
	cols = append(cols, Col{Title: "AGE", Text: formatAge(created), Width: 5, Style: ageStyle()})
	return cols
}

// IconCol is the leading icon column of a list row: the configured icon, or
// the source glyph in its color when there is none.
func IconCol(custom, src string, width int, selected bool) Col {
	if custom != "" {
		return Col{Text: custom, Width: width}
	}
	glyph, clr := icon.SourceGlyph(src)
	if selected && clr == colorDimmed {
		clr = colorText
	}
	return Col{Text: strings.TrimRight(glyph, " "), Width: width, Style: lipgloss.NewStyle().Foreground(clr)}
}
