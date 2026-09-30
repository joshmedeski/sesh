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

type columnKind int

const (
	fixedColumn columnKind = iota
	fillColumn
	branchColumnKind
)

type columnDef struct {
	kind     columnKind
	width    int
	minWidth int
	build    func(width int) Col
}

func layoutColumns(width int, order []string, defs map[string]columnDef, branchLongest int) []Col {
	picked := make([]columnDef, 0, len(order))
	fixed, fills := 0, 0
	for _, id := range order {
		def, ok := defs[id]
		if !ok || width < def.minWidth {
			continue
		}
		picked = append(picked, def)
		switch def.kind {
		case fixedColumn:
			fixed += def.width
		case fillColumn:
			fills++
		}
	}

	flex := width - 2 - fixed - max(len(picked)-1, 0)
	branchWidth := 0
	for _, def := range picked {
		if def.kind != branchColumnKind {
			continue
		}
		if fills > 0 {
			branchWidth = branchColumn(flex, branchLongest)
		} else {
			branchWidth = max(min(max(branchLongest, len("BRANCH")), flex), 1)
		}
	}
	fillSpace := max(flex-branchWidth, fills)
	fillWidth, fillExtra := fillSpace/max(fills, 1), fillSpace%max(fills, 1)

	cols := make([]Col, 0, len(picked))
	for _, def := range picked {
		w := def.width
		switch def.kind {
		case fillColumn:
			w = fillWidth
			if fills--; fills == 0 {
				w += fillExtra
			}
		case branchColumnKind:
			w = branchWidth
		}
		col := def.build(w)
		col.Width = w
		cols = append(cols, col)
	}
	return cols
}

// OpenColumns are the columns a sessions list supports, in their default
// order.
var OpenColumns = []string{"icon", "attached", "title", "alias", "directory", "git_branch", "git_status", "age", "alerts"}

// WorktreeColumns are the columns a worktree list supports, in their default
// order.
var WorktreeColumns = []string{"icon", "ghi_number", "ghi_state", "ghi_title", "git_branch", "git_status", "age"}

// IssueColumns are the GitHub issue columns a sessions list can add, for
// sessions inside a [[worktree]] checkout.
var IssueColumns = []string{"ghi_number", "ghi_state", "ghi_title"}

// Issue is the GitHub issue a row's worktree maps to.
type Issue struct {
	Number int
	Title  string
	State  string
}

func issueColumnDefs(issue Issue, selected bool) map[string]columnDef {
	number := ""
	if issue.Number > 0 {
		number = "#" + strconv.Itoa(issue.Number)
	}
	numberStyle := DimmedStyle()
	titleStyle := TextStyle()
	if issue.State == "CLOSED" {
		titleStyle = DimmedStyle()
	}
	if selected {
		numberStyle = TextStyle()
		titleStyle = TextStyle()
	}
	return map[string]columnDef{
		"ghi_number": {width: 7, build: func(int) Col {
			return Col{Title: "ISSUE", Text: number, Style: numberStyle}
		}},
		"ghi_state": {width: 6, build: func(int) Col {
			return Col{Title: "STATE", Text: strings.ToLower(issue.State), Style: lipgloss.NewStyle().Foreground(issueStateColors[issue.State])}
		}},
		"ghi_title": {kind: fillColumn, build: func(w int) Col {
			return Col{Title: "TITLE", Text: TruncateRight(issue.Title, w), Style: titleStyle}
		}},
	}
}

var configuredColumns = []string{"icon", "running", "title", "alias", "directory", "git_branch", "git_status"}

func orDefault(columns, defaults []string) []string {
	if len(columns) == 0 {
		return defaults
	}
	return columns
}

// renderOpenRow renders a sessions list row with the default columns:
// marker(2) | icon | att(2) | name(18) | alias(longest) | dir(fill) |
// branch(longest) | status(12) | age(5, last attached) | alerts(2).
// Progressive drop: <90 cols drop status+age+alerts, <70 drop branch+att.
func renderOpenRow(width int, selected, current bool, name, alias string, attached, windows int, dir, branch, status string, lastAttached *time.Time, alerts []string) string {
	return RenderOpenRowFocused(width, nil, lipgloss.Width(branch), AliasColumn(alias), IconCol("", "tmux", 1, selected), selected, current, true, name, alias, attached, windows, dir, branch, status, lastAttached, alerts, Issue{})
}

// RenderOpenRowFocused is renderOpenRow with explicit columns (nil for the
// defaults) and focused flag, so unfocused panes render a dimmed selection
// highlight.
func RenderOpenRowFocused(width int, columns []string, branchCol, aliasCol int, iconCol Col, selected, current, focused bool, name, alias string, attached, windows int, dir, branch, status string, lastAttached *time.Time, alerts []string, issue Issue) string {
	cols := openCols(width, columns, branchCol, aliasCol, iconCol, selected, current, focused, name, alias, attached, dir, branch, status, lastAttached, alerts, issue)
	return RenderRow(RowMarker(selected, focused), cols, selected, focused)
}

// RenderOpenHeader renders a sessions list's column titles.
func RenderOpenHeader(width int, columns []string, branchCol, aliasCol, iconWidth int) string {
	return RenderColumnTitles(openCols(width, columns, branchCol, aliasCol, Col{Width: iconWidth}, false, false, true, "", "", 0, "", "", "", nil, nil, Issue{}))
}

func openCols(width int, columns []string, branchCol, aliasCol int, iconCol Col, selected, current, focused bool, name, alias string, attached int, dir, branch, status string, lastAttached *time.Time, alerts []string, issue Issue) []Col {
	nameStyle := TextStyle()
	if current {
		nameStyle = accentStyle()
	}
	defs := issueColumnDefs(issue, selected)
	for id, def := range map[string]columnDef{
		"icon": {width: iconCol.Width, build: func(int) Col { return iconCol }},
		"attached": {width: 2, minWidth: 70, build: func(int) Col {
			if attached > 0 {
				return Col{Text: SuccessStyle().Render("●")}
			}
			return Col{}
		}},
		"title": {width: 18, build: func(w int) Col {
			return Col{Title: "NAME", Text: TruncateRight(name, w), Style: nameStyle}
		}},
		"directory": {kind: fillColumn, build: func(w int) Col {
			return Col{Title: "DIRECTORY", Text: truncateDirLeft(dir, w), Style: TextStyle()}
		}},
		"git_branch": {kind: branchColumnKind, minWidth: 70, build: func(w int) Col {
			return Col{Title: "BRANCH", Text: TruncateRight(branch, w), Style: BranchStyle(), Align: lipgloss.Left}
		}},
		"git_status": {width: 12, minWidth: 90, build: func(w int) Col {
			return Col{Title: "STATUS", Text: TruncateRightANSI(status, w), Style: BranchStyle()}
		}},
		"age": {width: 5, minWidth: 90, build: func(int) Col {
			return Col{Title: "AGE", Text: formatAge(lastAttached), Style: ageStyle(), Align: lipgloss.Left}
		}},
		"alerts": {width: 2, minWidth: 90, build: func(int) Col {
			if len(alerts) > 0 {
				return Col{Text: WarningStyle().Render("!")}
			}
			return Col{}
		}},
	} {
		defs[id] = def
	}
	if aliasCol > 0 {
		defs["alias"] = columnDef{width: aliasCol, build: func(int) Col {
			return Col{Title: "ALIAS", Text: aliasPill(alias, selected, focused)}
		}}
	}
	return layoutColumns(width, orDefault(columns, OpenColumns), defs, branchCol)
}

// renderConfiguredRow renders a Tab 2 (Configured) session row with columns:
// marker(2) | icon | state(2) | name(24) | alias(longest) | path(fill) |
// branch(16) | status(12). The branch drops below 70 cols.
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
	return RenderColumnTitles(configuredCols(width, aliasCol, Col{Width: iconWidth}, false, true, "", "", false, "", "", ""))
}

func configuredCols(width, aliasCol int, iconCol Col, selected, focused bool, name, alias string, running bool, path, branch, status string) []Col {
	if path == "" {
		path = "-"
	}
	defs := map[string]columnDef{
		"icon": {width: iconCol.Width, build: func(int) Col { return iconCol }},
		"running": {width: 2, build: func(int) Col {
			if running {
				return Col{Text: "●", Style: SuccessStyle()}
			}
			return Col{Text: "○", Style: DimmedStyle()}
		}},
		"title": {width: 24, build: func(w int) Col {
			return Col{Title: "NAME", Text: TruncateRight(name, w), Style: TextStyle()}
		}},
		"directory": {kind: fillColumn, build: func(w int) Col {
			return Col{Title: "DIRECTORY", Text: truncateDirLeft(path, w), Style: TextStyle()}
		}},
		"git_branch": {width: 16, minWidth: 70, build: func(w int) Col {
			return Col{Title: "BRANCH", Text: TruncateRight(branch, w), Style: BranchStyle()}
		}},
		"git_status": {width: 12, build: func(w int) Col {
			return Col{Title: "STATUS", Text: TruncateRightANSI(status, w)}
		}},
	}
	if aliasCol > 0 {
		defs["alias"] = columnDef{width: aliasCol, build: func(int) Col {
			return Col{Title: "ALIAS", Text: aliasPill(alias, selected, focused)}
		}}
	}
	return layoutColumns(width, configuredColumns, defs, 16)
}

func branchColumn(flex, longest int) int {
	return min(max(longest, 16), max(flex-16, 16))
}

// RenderWorktreeRowFocused renders a worktree row. The default columns are:
// marker(2) | icon | number(7) | state(6) | title(fill) | branch(longest) |
// status(12) | age(5, since created). The branch drops below 70 cols. Closed
// issues render their title dimmed.
func RenderWorktreeRowFocused(width int, columns []string, branchCol int, iconCol Col, selected, focused bool, number int, title, state, branch, status string, created *time.Time) string {
	cols := worktreeCols(width, columns, branchCol, iconCol, selected, Issue{Number: number, Title: title, State: state}, branch, status, created)
	return RenderRow(RowMarker(selected, focused), cols, selected, focused)
}

// RenderWorktreeHeader renders a worktree list's column titles.
func RenderWorktreeHeader(width int, columns []string, branchCol, iconWidth int) string {
	return RenderColumnTitles(worktreeCols(width, columns, branchCol, Col{Width: iconWidth}, false, Issue{}, "", "", nil))
}

var issueStateColors = map[string]lipgloss.ANSIColor{
	"OPEN":   lipgloss.ANSIColor(2),
	"MERGED": lipgloss.ANSIColor(5),
	"CLOSED": lipgloss.ANSIColor(1),
}

func worktreeCols(width int, columns []string, branchCol int, iconCol Col, selected bool, issue Issue, branch, status string, created *time.Time) []Col {
	defs := issueColumnDefs(issue, selected)
	defs["icon"] = columnDef{width: iconCol.Width, build: func(int) Col { return iconCol }}
	defs["git_branch"] = columnDef{kind: branchColumnKind, minWidth: 70, build: func(w int) Col {
		return Col{Title: "BRANCH", Text: TruncateRight(branch, w), Style: BranchStyle()}
	}}
	defs["git_status"] = columnDef{width: 12, build: func(w int) Col {
		return Col{Title: "STATUS", Text: TruncateRightANSI(status, w), Style: BranchStyle()}
	}}
	defs["age"] = columnDef{width: 5, build: func(int) Col {
		return Col{Title: "AGE", Text: formatAge(created), Style: ageStyle()}
	}}
	return layoutColumns(width, orDefault(columns, WorktreeColumns), defs, branchCol)
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
