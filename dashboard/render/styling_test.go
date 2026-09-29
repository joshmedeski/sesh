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
	"github.com/joshmedeski/sesh/v2/icon"
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
	assert.Contains(t, row, "main")
	assert.Contains(t, row, "+1 ~2")
}

func TestRenderOpenRow_DropsStatusUnder90(t *testing.T) {
	row := renderOpenRow(80, false, false, "s", "", 0, 1, "~/d", "main", "+1", nil, nil)
	assert.Contains(t, row, "main")
	assert.NotContains(t, row, "+1")
}

func TestRenderOpenRow_DropsBranchUnder70(t *testing.T) {
	row := renderOpenRow(60, false, false, "foo", "", 0, 1, "~/d", "main", "+1", nil, nil)
	assert.NotContains(t, row, "main")
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
	row := renderConfiguredRow(100, false, "proj", "", "", true, "~/code/proj", "main", "+1")
	assert.Contains(t, row, "proj")
	assert.Contains(t, row, "●")
	assert.Contains(t, row, "~/code/proj")
	assert.Contains(t, row, "main")
}

func TestRenderConfiguredRowNotRunning(t *testing.T) {
	row := renderConfiguredRow(100, false, "proj", "", "", false, "", "main", "")
	assert.Contains(t, row, "○")
	assert.Contains(t, row, "-")
}

func TestRenderConfiguredRow_StartupCommandIndicator(t *testing.T) {
	// The startup-command "*" indicator was removed; the column stays blank.
	row := renderConfiguredRow(100, false, "proj", "", "make run", true, "~/code/proj", "", "")
	assert.NotContains(t, row, "*")
	assert.NotContains(t, row, "\x1b[38;5;11m") // no yellow

	row = renderConfiguredRow(100, false, "proj", "", "", true, "~/code/proj", "", "")
	assert.NotContains(t, row, "*")
}

func TestRenderConfiguredRow_DropsCmdAndBranchUnder70(t *testing.T) {
	row := renderConfiguredRow(60, false, "proj", "", "make run", true, "~/code/proj", "main", "")
	assert.NotContains(t, row, "*")
	assert.NotContains(t, row, "main")
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
	row := renderConfiguredRow(100, false, "proj", "", "", false, "~/code/proj", "main", styled)
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

// --- Alias column ---

func cellIndex(s, sub string) int {
	i := strings.Index(s, sub)
	if i < 0 {
		return -1
	}
	return lipgloss.Width(s[:i])
}

func TestRenderOpenRowAliasPillAfterName(t *testing.T) {
	row := renderOpenRow(100, false, false, "wallpaper", "wp", 0, 1, "~/d", "", "", nil, nil)
	plain := ansi.Strip(row)
	assert.NotContains(t, row, "38;5;14m")
	assert.Contains(t, row, "38;5;15m")
	assert.Contains(t, plain, pillLeftGlyph+"wp"+pillRightGlyph)
	nameIdx := cellIndex(plain, "wallpaper")
	assert.Equal(t, nameIdx+18+1, cellIndex(plain, pillLeftGlyph))
	assert.Greater(t, cellIndex(plain, "~/d"), cellIndex(plain, pillRightGlyph))

	noAlias := ansi.Strip(renderOpenRow(100, false, false, "wallpaper", "", 0, 1, "~/d", "", "", nil, nil))
	assert.NotContains(t, noAlias, pillLeftGlyph)
	assert.Equal(t, nameIdx+18+1, cellIndex(noAlias, "~/d"))
}

func TestAliasPillSelectedKeepsRowHighlight(t *testing.T) {
	row := renderOpenRow(100, true, false, "wallpaper", "wp", 0, 1, "~/d", "", "", nil, nil)
	glyphIdx := strings.Index(row, pillLeftGlyph)
	aliasIdx := strings.Index(row, "wp")
	require.GreaterOrEqual(t, glyphIdx, 0)
	assert.Contains(t, row[glyphIdx:aliasIdx], "\x1b[48;5;15m")
	assert.Contains(t, row[glyphIdx:aliasIdx], "\x1b[38;5;0m")
	assert.Contains(t, row[aliasIdx:], "\x1b[48;5;8m")
}

func TestRenderConfiguredRowAliasPillAfterName(t *testing.T) {
	plain := ansi.Strip(renderConfiguredRow(100, false, "proj", "p", "", false, "~/code/proj", "", ""))
	nameIdx := cellIndex(plain, "proj")
	assert.Equal(t, nameIdx+24+1, cellIndex(plain, pillLeftGlyph+"p"+pillRightGlyph))
	assert.Greater(t, cellIndex(plain, "~/code/proj"), nameIdx+24+1)
}

// TestRenderOpenRow_NameColumnSlim locks in the slimmer Open-session name
// column: the name cell is 18 visible cells wide (was 20), so the directory
// column starts at marker(2) + icon(1) + sep(1) + att(2) + sep(1) + name(18)
// + sep(1) = 26.
func TestRenderOpenRow_NameColumnSlim(t *testing.T) {
	row := ansi.Strip(renderOpenRow(100, false, false, "short", "", 0, 1, "~/d", "", "", nil, nil))
	dirIdx := strings.Index(row, "~/d")
	require.GreaterOrEqual(t, dirIdx, 0)
	assert.Equal(t, 26, lipgloss.Width(row[:dirIdx]))
}

func TestRenderOpenRow_LongCellsDoNotWrap(t *testing.T) {
	for _, selected := range []bool{false, true} {
		row := renderOpenRow(120, selected, false, "re/sesh/extensions/sesh-worktree-cleanup", "sc", 1, 1, "~/c/re/sesh/extensions/sesh", "claude/telemetry-dashboard", "+1 ~113 -62 !4", nil, nil)
		assert.NotContains(t, row, "\n")
		assert.Equal(t, 120, lipgloss.Width(row))
	}
}

func TestRenderOpenRow_WideRowShowsFullBranch(t *testing.T) {
	row := renderOpenRow(200, false, false, "nu/w/10290", "", 1, 1, "~/c/nu/w/10290", "claude/telemetry-dashboard", "~2", nil, nil)
	assert.Contains(t, ansi.Strip(row), "claude/telemetry-dashboard")
	assert.Equal(t, 200, lipgloss.Width(row))
}

func TestRenderWorktreeRow_SelectedNumberUsesTextColor(t *testing.T) {
	row := RenderWorktreeRowFocused(120, 16, IconCol("", "config", 1, true), true, true, 7503, "closed issue", true, "main", "")
	selected := TextStyle().Inherit(cursorStyle(true))
	assert.Contains(t, row, selected.Width(7).Render("#7503"))
	assert.NotContains(t, row, DimmedStyle().Inherit(cursorStyle(true)).Width(7).Render("#7503"))
}

func TestIconCol_CustomIconElseSourceGlyph(t *testing.T) {
	assert.Equal(t, "🏠", IconCol("🏠", "tmux", 2, false).Text)

	glyph, clr := icon.SourceGlyph("tmux")
	col := IconCol("", "tmux", 2, false)
	assert.Equal(t, strings.TrimSpace(glyph), col.Text)
	assert.Equal(t, clr, col.Style.GetForeground())
}

func TestIconCol_SelectedDimmedGlyphUsesTextColor(t *testing.T) {
	assert.Equal(t, colorDimmed, IconCol("", "config", 1, false).Style.GetForeground())
	assert.Equal(t, colorText, IconCol("", "config", 1, true).Style.GetForeground())
}

func TestRenderOpenRow_IconIsFirstColumn(t *testing.T) {
	row := ansi.Strip(RenderOpenRowFocused(100, 0, 0, IconCol("🏠", "tmux", 2, false), false, false, true, "nutiliti", "", 1, 1, "~/c/nu", "main", "", nil, nil))
	assert.True(t, strings.HasPrefix(row, "  🏠"), row)
}

func assertTitlesAlign(t *testing.T, header, row string, titles map[string]string) {
	t.Helper()
	header, row = ansi.Strip(header), ansi.Strip(row)
	assert.Equal(t, lipgloss.Width(row), lipgloss.Width(header))
	for title, cell := range titles {
		assert.Equal(t, cellIndex(row, cell), cellIndex(header, title), title)
	}
}

func TestColumnTitlesAlignWithRows(t *testing.T) {
	icon := Col{Text: "x", Width: 1}
	alias := AliasColumn("ms")
	wide := map[string]string{"NAME": "mysession", "ALIAS": pillLeftGlyph, "DIRECTORY": "~/some/dir", "BRANCH": "main", "STATUS": "+1"}
	narrow := map[string]string{"NAME": "mysession", "ALIAS": pillLeftGlyph, "DIRECTORY": "~/some/dir"}
	for width, titles := range map[int]map[string]string{60: narrow, 80: wide, 120: wide} {
		if width < 90 {
			delete(titles, "STATUS")
		}
		assertTitlesAlign(t,
			RenderOpenHeader(width, 10, alias, 1),
			RenderOpenRowFocused(width, 10, alias, icon, false, false, true, "mysession", "ms", 0, 1, "~/some/dir", "main", "+1", nil, nil),
			titles)
	}

	assertTitlesAlign(t,
		RenderConfiguredHeader(100, alias, 1),
		RenderConfiguredRowFocused(100, alias, icon, false, true, "mysession", "ms", "", true, "~/some/dir", "main", "+1"),
		map[string]string{"NAME": "mysession", "ALIAS": pillLeftGlyph, "DIRECTORY": "~/some/dir", "BRANCH": "main", "STATUS": "+1"})

	assertTitlesAlign(t,
		RenderWorktreeHeader(100, 10, 1),
		RenderWorktreeRowFocused(100, 10, icon, false, true, 358, "tmux command updates", false, "jam/358", "+1"),
		map[string]string{"ISSUE": "#358", "TITLE": "tmux command updates", "BRANCH": "jam/358", "STATUS": "+1"})
}
