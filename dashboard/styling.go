// styling.go
package dashboard

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// Shared color palette (see design spec).
var (
	colorAccent = lipgloss.ANSIColor(14) // accent (cyan)
	colorDimmed = lipgloss.ANSIColor(8)  // dimmed / border
	colorBorder = lipgloss.ANSIColor(8)
	colorText   = lipgloss.ANSIColor(15) // white text
	// colorPillText is the alias pill's label colour on a selected row. Black
	// mirrors the terminal's typical default background, so it stays legible
	// on the accent pill fill the same way the unselected chip's reverse-video
	// label does (its text is the terminal's own contrasting background).
	colorPillText  = lipgloss.ANSIColor(0)  // black
	colorAge       = lipgloss.ANSIColor(7)  // light gray (distinct from highlight bg)
	colorBranch    = lipgloss.ANSIColor(5)  // magenta
	colorStatus    = lipgloss.ANSIColor(10) // green
	colorHighlight = lipgloss.ANSIColor(8)

	// colorHighlightDim is the slightly dimmer selection background used when a
	// pane is not focused, so the selection stays visible without competing
	// with the focused pane's highlight.
	colorHighlightDim = lipgloss.ANSIColor(236)

	// git-status part colours (formatGitStatus).
	colorStaged    = lipgloss.ANSIColor(10) // green  "+N"
	colorUnstaged  = lipgloss.ANSIColor(11) // yellow "~N"
	colorDeleted   = lipgloss.ANSIColor(9)  // red    "-N"
	colorUntracked = lipgloss.ANSIColor(5)  // magenta "!N"
)

func accentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
}

// warningStyle is the yellow style used for alerts and the startup-command
// indicator.
func warningStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorUnstaged)
}

func dimmedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorDimmed)
}

// ageStyle is the neutral, muted foreground for the Open sessions age column.
// It deliberately avoids colorDimmed (the cursor background) so the age stays
// readable when a row is highlighted.
func ageStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorAge)
}

func textStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorText)
}

func branchStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorBranch)
}

func successStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorStatus)
}

// GroupNameRender is retained for backwards compatibility with the widget
// sections that render a single-line title bar.
func GroupNameRender(name string, width int) lipgloss.Style {
	return NewStyle(width, width, 1, 1, 15, false, []int{0, 0, 0, 0})
}

// cursorStyle is the shared selection highlight used by every list section: a
// subtle full-line background fill. Focused panes use the standard highlight;
// unfocused panes use a slightly dimmer fill.
func cursorStyle(focused bool) lipgloss.Style {
	if focused {
		return lipgloss.NewStyle().Background(colorHighlight)
	}
	return lipgloss.NewStyle().Background(colorHighlightDim)
}

// rowMarker returns the 2-column marker for a row: "▌ " (accent bold on the
// selection background) when selected, otherwise two spaces. The marker is
// dimmed when the pane is not focused.
func rowMarker(selected, focused bool) string {
	if !selected {
		return "  "
	}
	if focused {
		return accentStyle().Background(colorHighlight).Render("▌ ")
	}
	return dimmedStyle().Background(colorHighlightDim).Render("▌ ")
}

// col describes a single column in a list row.
type col struct {
	text  string
	width int
	style lipgloss.Style
	align lipgloss.Position
}

// renderRow joins columns with a single space and applies the shared cursor
// background to every cell (and separator) plus the marker when selected,
// producing a vim-like full-line highlight.
func renderRow(marker string, cols []col, selected, focused bool) string {
	bg := cursorStyle(focused)
	rendered := make([]string, len(cols))
	for i, c := range cols {
		st := c.style
		if selected {
			st = st.Inherit(bg)
		}
		st = st.Width(c.width)
		if c.align != 0 {
			st = st.Align(c.align)
		}
		rendered[i] = st.Render(c.text)
	}
	if selected {
		return marker + strings.Join(rendered, bg.Render(" "))
	}
	return marker + strings.Join(rendered, " ")
}

// renderSimpleRow renders a widget-section row (git, ssh, docker) with the
// shared cursor marker and a vim-like full-line highlight when selected. Cells
// are raw text with a foreground style and are concatenated directly (no
// separator is inserted), so any spacing or column padding must live in the
// cell text or its style; cell content is otherwise preserved exactly.
func renderSimpleRow(cells []col, selected, focused bool) string {
	bg := cursorStyle(focused)
	rendered := make([]string, len(cells))
	for i, c := range cells {
		st := c.style
		if selected {
			st = st.Inherit(bg)
		}
		rendered[i] = st.Render(c.text)
	}
	return rowMarker(selected, focused) + strings.Join(rendered, "")
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
	return renderOpenRowFocused(width, selected, current, true, name, alias, attached, windows, dir, branch, status, lastAttached, alerts)
}

// renderOpenRowFocused is renderOpenRow with an explicit focused flag, so
// unfocused panes render a dimmed selection highlight.
func renderOpenRowFocused(width int, selected, current, focused bool, name, alias string, attached, windows int, dir, branch, status string, lastAttached *time.Time, alerts []string) string {
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

	nameStyle := textStyle()
	if current {
		nameStyle = accentStyle()
	}

	cols := []col{}

	if includeAtt {
		attText := ""
		if attached > 0 {
			attText = successStyle().Render("●")
		}
		cols = append(cols, col{text: attText, width: 2})
	}

	chip := aliasChip(alias)
	nameBudget := 22
	if chip != "" {
		nameBudget -= lipgloss.Width(chip)
		if nameBudget < 1 {
			nameBudget = 1
		}
	}
	truncated := truncateRight(name, nameBudget)
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
	cols = append(cols, col{text: nameText, width: 18})

	// if includeWindows {
	// 	cols = append(cols, col{text: fmt.Sprintf("%2dw", windows), width: 7, style: textStyle(), align: lipgloss.Left})
	// }
	cols = append(cols, col{text: truncateDirLeft(dir, dirWidth), width: dirWidth, style: textStyle()})
	if includeBranch {
		cols = append(cols, col{text: truncateRight(paren(branch), 16), width: 14, style: branchStyle(), align: lipgloss.Left})
	}
	if includeStatus {
		cols = append(cols, col{text: truncateRightANSI(status, 12), width: 12, style: branchStyle()})
	}
	if includeAge {
		cols = append(cols, col{text: formatAge(lastAttached), width: 5, style: ageStyle(), align: lipgloss.Left})
	}
	if includeAlerts {
		alertText := ""
		if len(alerts) > 0 {
			alertText = warningStyle().Render("!")
		}
		cols = append(cols, col{text: alertText, width: 2})
	}

	return renderRow(rowMarker(selected, focused), cols, selected, focused)
}

// formatAge renders a compact relative age ("2h"/"3d"/"4mo") for a session's
// last attached time, or "" when lastAttached is nil/zero.
func formatAge(lastAttached *time.Time) string {
	if lastAttached == nil || lastAttached.IsZero() {
		return ""
	}
	since := time.Since(*lastAttached)
	if since < 0 {
		return ""
	}
	switch {
	case since < time.Hour:
		m := int(since.Minutes())
		if m < 1 {
			m = 1
		}
		return fmt.Sprintf("%dm", m)
	case since < 24*time.Hour:
		return fmt.Sprintf("%dh", int(since.Hours()))
	case since < 30*24*time.Hour:
		return fmt.Sprintf("%dd", int(since.Hours()/24))
	default:
		return fmt.Sprintf("%dmo", int(since.Hours()/(24*30)))
	}
}

// renderConfiguredRow renders a Tab 2 (Configured) session row with columns:
// marker(2) | cmd(2) | name(24) | state(2) | path(fill) | branch(16) |
// status(12). The cmd column ("*") and branch drop together below 70 cols.
func renderConfiguredRow(width int, selected bool, name, startupCommand string, running bool, path, branch, status string) string {
	return renderConfiguredRowFocused(width, selected, true, name, startupCommand, running, path, branch, status)
}

// renderConfiguredRowFocused is renderConfiguredRow with an explicit focused
// flag, so unfocused panes render a dimmed selection highlight.
func renderConfiguredRowFocused(width int, selected, focused bool, name, startupCommand string, running bool, path, branch, status string) string {
	includeBranch := width >= 70

	stateText := "○"
	stateColStyle := dimmedStyle()
	if running {
		stateText = "●"
		stateColStyle = successStyle()
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

	cols := make([]col, 0, numCols)
	if includeBranch {
		cmdText := ""
		// if startupCommand != "" {
		// 	cmdText = warningStyle().Render("*")
		// }
		cols = append(cols, col{text: cmdText, width: 2})
	}
	cols = append(cols, col{text: stateText, width: 2, style: stateColStyle})
	cols = append(cols, col{text: truncateRight(name, 24), width: 24, style: textStyle()})
	cols = append(cols, col{text: truncateDirLeft(path, pathWidth), width: pathWidth, style: textStyle()})
	if includeBranch {
		cols = append(cols, col{text: truncateRight(paren(branch), 16), width: 16, style: branchStyle()})
	}
	cols = append(cols, col{text: truncateRightANSI(status, 12), width: 12})

	return renderRow(rowMarker(selected, focused), cols, selected, focused)
}

// collapseHome replaces the home directory prefix of path with "~".
func collapseHome(path, homeDir string) string {
	if homeDir == "" {
		return path
	}
	if path == homeDir {
		return "~"
	}
	if strings.HasPrefix(path, homeDir+string(filepath.Separator)) {
		return "~" + path[len(homeDir):]
	}
	return path
}

// paren wraps s in parentheses when non-empty, otherwise returns "".
func paren(s string) string {
	if s == "" {
		return ""
	}
	return "(" + s + ")"
}

// truncateRight truncates s to maxRunes runes, appending "…" when truncated.
func truncateRight(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	if maxRunes <= 1 {
		return "…"
	}
	return string(runes[:maxRunes-1]) + "…"
}

// truncateRightANSI truncates an ANSI-styled string to maxWidth visible cells,
// appending "…" when truncated. Escape sequences are preserved so embedded
// colours (e.g. the git-status parts) survive truncation.
func truncateRightANSI(s string, maxWidth int) string {
	return ansi.Truncate(s, maxWidth, "…")
}

// truncateDirLeft truncates a directory path from the left, preferring to keep
// the last two path segments (so the directory name and its parent stay
// visible). Falls back to a "…"-prefixed tail when segments are too long.
func truncateDirLeft(dir string, maxRunes int) string {
	if utf8.RuneCountInString(dir) <= maxRunes {
		return dir
	}
	if kept := lastTwoSegments(dir); utf8.RuneCountInString(kept) <= maxRunes {
		return kept
	}
	if maxRunes <= 1 {
		return "…"
	}
	runes := []rune(dir)
	return "…" + string(runes[len(runes)-(maxRunes-1):])
}

// lastTwoSegments keeps the leading "~/" (if present) plus the final two path
// segments of dir.
func lastTwoSegments(dir string) string {
	prefix := ""
	rest := dir
	if strings.HasPrefix(rest, "~/") {
		prefix = "~/"
		rest = rest[2:]
	} else if rest == "~" {
		return "~"
	}
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) <= 2 {
		return prefix + strings.Join(parts, "/")
	}
	return prefix + strings.Join(parts[len(parts)-2:], "/")
}
