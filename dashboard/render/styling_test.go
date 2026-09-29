package render

import (
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/git"
)

func firstResetIndex(s string) int {
	if i := strings.Index(s, "\x1b[0m"); i >= 0 {
		return i
	}
	return strings.Index(s, "\x1b[m")
}

// --- Row rendering ---

func TestRenderOpenRow_FullColumns(t *testing.T) {
	row := renderOpenRow(100, true, false, "mysession", "", 0, 3, "~/code/proj", "main", "+1 ~2", nil, nil)
	assert.Contains(t, row, "mysession")
	assert.Contains(t, row, "(main)")
	assert.Contains(t, row, "+1 ~2")
}

func TestRenderOpenRow_DropsStatusUnder90(t *testing.T) {
	row := renderOpenRow(80, false, false, "s", "", 0, 1, "~/d", "main", "+1", nil, nil)
	assert.Contains(t, row, "(main)")
	assert.NotContains(t, row, "+1")
}

func TestRenderOpenRow_DropsBranchUnder70(t *testing.T) {
	row := renderOpenRow(60, false, false, "foo", "", 0, 1, "~/d", "main", "+1", nil, nil)
	assert.NotContains(t, row, "(main)")
}

func TestRenderOpenRow_DropsWindowsUnder50(t *testing.T) {
	row := renderOpenRow(40, false, false, "foo", "", 0, 1, "~/d", "main", "+1", nil, nil)
	assert.NotContains(t, row, "1w")
}

func TestRenderOpenRow_AttachedIndicator(t *testing.T) {
	row := renderOpenRow(100, false, false, "s", "", 1, 1, "~/d", "", "", nil, nil)
	assert.Contains(t, row, "●")
	row = renderOpenRow(100, false, false, "s", "", 0, 1, "~/d", "", "", nil, nil)
	assert.NotContains(t, row, "●")
}

func TestRenderOpenRow_Age(t *testing.T) {
	twoH := time.Now().Add(-2 * time.Hour)
	row := renderOpenRow(100, false, false, "s", "", 0, 1, "~/d", "", "", &twoH, nil)
	assert.Contains(t, row, "2h")

	threeD := time.Now().Add(-3 * 24 * time.Hour)
	row = renderOpenRow(100, false, false, "s", "", 0, 1, "~/d", "", "", &threeD, nil)
	assert.Contains(t, row, "3d")

	fourMo := time.Now().Add(-4 * 30 * 24 * time.Hour)
	row = renderOpenRow(100, false, false, "s", "", 0, 1, "~/d", "", "", &fourMo, nil)
	assert.Contains(t, row, "4mo")

	// Age uses a dedicated light-gray foreground so it remains readable on the
	// cursor highlight background (which shares the dimmed gray colour).
	assert.Contains(t, row, "\x1b[38;5;7m")

	// Selected rows still show the age text with the same foreground.
	selected := renderOpenRow(100, true, false, "s", "", 0, 1, "~/d", "", "", &twoH, nil)
	assert.Contains(t, selected, "2h")
	assert.Contains(t, selected, "\x1b[38;5;7")

	// Nil/zero last attached → blank age.
	assert.Equal(t, "", formatAge(nil))
	zero := time.Time{}
	assert.Equal(t, "", formatAge(&zero))
}

func TestRenderOpenRow_Alerts(t *testing.T) {
	row := renderOpenRow(100, false, false, "s", "", 0, 1, "~/d", "", "", nil, []string{"bell"})
	assert.Contains(t, row, "!")
	row = renderOpenRow(100, false, false, "s", "", 0, 1, "~/d", "", "", nil, nil)
	assert.NotContains(t, row, "!")
}

func TestRenderOpenRow_CurrentHighlight(t *testing.T) {
	row := renderOpenRow(100, false, true, "mysession", "", 0, 1, "~/d", "", "", nil, nil)
	assert.Contains(t, row, "\x1b[1;38;5;14m") // bold cyan accent
}

func TestRenderConfiguredRow(t *testing.T) {
	row := renderConfiguredRow(100, false, "proj", "", true, "~/code/proj", "main", "+1")
	assert.Contains(t, row, "proj")
	assert.Contains(t, row, "●")
	assert.Contains(t, row, "~/code/proj")
	assert.Contains(t, row, "(main)")
}

func TestRenderConfiguredRowNotRunning(t *testing.T) {
	row := renderConfiguredRow(100, false, "proj", "", false, "", "main", "")
	assert.Contains(t, row, "○")
	assert.Contains(t, row, "-")
}

func TestRenderConfiguredRow_StartupCommandIndicator(t *testing.T) {
	// The startup-command "*" indicator was removed; the column stays blank.
	row := renderConfiguredRow(100, false, "proj", "make run", true, "~/code/proj", "", "")
	assert.NotContains(t, row, "*")
	assert.NotContains(t, row, "\x1b[38;5;11m") // no yellow

	row = renderConfiguredRow(100, false, "proj", "", true, "~/code/proj", "", "")
	assert.NotContains(t, row, "*")
}

func TestRenderConfiguredRow_DropsCmdAndBranchUnder70(t *testing.T) {
	row := renderConfiguredRow(60, false, "proj", "make run", true, "~/code/proj", "main", "")
	assert.NotContains(t, row, "*")
	assert.NotContains(t, row, "(main)")
}

// --- Git status formatting ---

func TestFormatGitStatusColored(t *testing.T) {
	got := FormatGitStatus(git.StatusSummary{Staged: 1, Unstaged: 2, Deleted: 3, Untracked: 4})
	// Each part carries its own distinct ANSI 256 foreground colour.
	assert.Contains(t, got, "\x1b[38;5;10m") // staged green
	assert.Contains(t, got, "\x1b[38;5;11m") // unstaged yellow
	assert.Contains(t, got, "\x1b[38;5;9m")  // deleted red
	assert.Contains(t, got, "\x1b[38;5;5m")  // untracked magenta
	// Visible text is unchanged.
	require.True(t, strings.Contains(got, "+1") && strings.Contains(got, "~2") &&
		strings.Contains(got, "-3") && strings.Contains(got, "!4"))
}

func TestFormatGitStatusOmitsZeroParts(t *testing.T) {
	got := FormatGitStatus(git.StatusSummary{Staged: 2, Untracked: 1})
	assert.NotContains(t, got, "~")
	assert.NotContains(t, got, "-")
	assert.Contains(t, got, "\x1b[38;5;10m")
	assert.Contains(t, got, "\x1b[38;5;5m")
}

func TestRenderOpenRow_StyledStatusKeepsWidth(t *testing.T) {
	styled := FormatGitStatus(git.StatusSummary{Staged: 1, Unstaged: 2})
	row := renderOpenRow(100, false, false, "mysession", "", 0, 3, "~/code/proj", "main", styled, nil, nil)
	// The status cell has no outer foreground, so the embedded colours remain.
	assert.Contains(t, row, "\x1b[38;5;10m")
	assert.Contains(t, row, "\x1b[38;5;11m")
	// Embedded escapes do not inflate the visible row width.
	plain := renderOpenRow(100, false, false, "mysession", "", 0, 3, "~/code/proj", "main", "+1 ~2", nil, nil)
	assert.Equal(t, lipgloss.Width(plain), lipgloss.Width(row))
}

func TestRenderConfiguredRow_StyledStatusKeepsWidth(t *testing.T) {
	styled := FormatGitStatus(git.StatusSummary{Deleted: 3, Untracked: 4})
	row := renderConfiguredRow(100, false, "proj", "", false, "~/code/proj", "main", styled)
	assert.Contains(t, row, "\x1b[38;5;9m")
	assert.Contains(t, row, "\x1b[38;5;5m")
	assert.Equal(t, 100, lipgloss.Width(row))
}

// TestGitStatusHoverBackgroundCoversAllParts verifies that when a row with a
// coloured git-status cell is highlighted, the cursor background is applied to
// every status glyph, not just the first one. The status string must avoid
// per-part reset sequences so the outer highlight background stays active
// across the whole cell.
func TestGitStatusHoverBackgroundCoversAllParts(t *testing.T) {
	status := FormatGitStatus(git.StatusSummary{Staged: 1, Unstaged: 2, Deleted: 3, Untracked: 4})
	cell := RenderRow("", []Col{{Text: status, Width: 12, Style: BranchStyle()}}, true, true)

	plain := strings.TrimRight(ansi.Strip(cell), " ")
	require.Equal(t, "+1 ~2 -3 !4", plain)

	// Each status glyph keeps its distinct foreground colour.
	for _, seq := range []string{"\x1b[38;5;10m", "\x1b[38;5;11m", "\x1b[38;5;9m", "\x1b[38;5;5m"} {
		assert.Contains(t, cell, seq)
	}

	// The highlight background starts before the first status glyph.
	greenIdx := strings.Index(cell, "\x1b[38;5;10m+1")
	require.GreaterOrEqual(t, greenIdx, 0)
	assert.Contains(t, cell[:greenIdx], ";48;5;8m")

	// No reset appears between the first status glyph and the final "!4", so
	// the cursor background stays active for every glyph.
	afterGreen := cell[greenIdx:]
	resetIdx := firstResetIndex(afterGreen)
	require.GreaterOrEqual(t, resetIdx, 0)
	lastTokenIdx := strings.Index(afterGreen, "!4")
	require.GreaterOrEqual(t, lastTokenIdx, 0)
	assert.Greater(t, resetIdx, lastTokenIdx)
}

// --- Alias-chip rendering ---

func TestRenderOpenRowWithAlias(t *testing.T) {
	row := renderOpenRow(100, false, false, "wallpaper", "wp", 0, 1, "~/d", "", "", nil, nil)
	assert.Contains(t, row, "wallpaper")
	// The alias is rendered as a rounded pill using powerline half circles,
	// matching the picker’s default alias chip style.
	assert.Contains(t, ansi.Strip(row), "\ue0b6wp\ue0b4")
	assert.NotContains(t, row, "[wp]")
	// The pill uses reverse video over the dashboard accent color so the
	// effective background is the theme accent (cyan 14).
	assert.Contains(t, row, "38;5;14m")
	assert.Contains(t, row, "\x1b[7")

	// Without an alias the chip is absent.
	noAlias := renderOpenRow(100, false, false, "wallpaper", "", 0, 1, "~/d", "", "", nil, nil)
	assert.NotContains(t, noAlias, "\ue0b6")
	assert.NotContains(t, noAlias, "[wp]")
}

// TestAliasRowHoverBackgroundCoversPillAndName verifies that when a row with
// an alias is selected, the alias pill keeps its own accent background while
// the cursor highlight continues unbroken around it. The pill and name must be
// emitted as one continuous ANSI run — no reset may interrupt the highlight
// between the pill’s left glyph and the end of the name.
func TestAliasRowHoverBackgroundCoversPillAndName(t *testing.T) {
	row := renderOpenRow(100, true, false, "wallpaper", "wp", 0, 1, "~/d", "", "", nil, nil)

	// The cursor background starts before the pill’s left glyph.
	glyphIdx := strings.Index(row, chipLeftGlyph)
	require.GreaterOrEqual(t, glyphIdx, 0)
	assert.Contains(t, row[:glyphIdx], ";48;5;8m")

	// The pill paints an explicit accent background around the alias text, so
	// the pill keeps its own fill instead of taking the cursor highlight.
	aliasIdx := strings.Index(row, "wp")
	require.GreaterOrEqual(t, aliasIdx, 0)
	assert.Contains(t, row[glyphIdx:aliasIdx], "\x1b[48;5;14m")
	// The label text contrasts with the accent fill.
	assert.Contains(t, row[glyphIdx:aliasIdx], "\x1b[38;5;0m")

	// The cursor background is restored after the pill, so the name that
	// follows stays on the row highlight.
	assert.Contains(t, row[aliasIdx:], "\x1b[48;5;8m")

	// No reset appears between the pill’s left glyph and the end of the name,
	// so the highlight stays active across the pill and the name.
	afterGlyph := row[glyphIdx:]
	nameEnd := strings.Index(afterGlyph, "wallpaper")
	require.GreaterOrEqual(t, nameEnd, 0)
	nameEnd += len("wallpaper")
	resetIdx := firstResetIndex(afterGlyph)
	require.GreaterOrEqual(t, resetIdx, 0)
	assert.GreaterOrEqual(t, resetIdx, nameEnd)

	// The rounded theme-aware pill is preserved: the alias reads as one shape
	// between the half circles, painted with the accent foreground.
	assert.Contains(t, ansi.Strip(row), "\ue0b6wp\ue0b4")
	assert.Contains(t, row, "38;5;14m")
}

// TestRenderOpenRow_NameColumnSlim locks in the slimmer Open-session name
// column: the name cell is 18 visible cells wide (was 20), so the directory
// column starts at marker(2) + att(2) + sep(1) + name(18) + sep(1) = 24.
func TestRenderOpenRow_NameColumnSlim(t *testing.T) {
	row := renderOpenRow(100, false, false, "short", "", 0, 1, "~/d", "", "", nil, nil)
	dirIdx := strings.Index(ansi.Strip(row), "~/d")
	require.GreaterOrEqual(t, dirIdx, 0)
	assert.Equal(t, 24, dirIdx)
}
