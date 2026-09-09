package picker

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/model"
)

func testWorktrees() []model.WorktreeEntry {
	return []model.WorktreeEntry{
		{Number: 9, Path: "/repo/.wk/9", Title: "Add an opt-in preview pane", State: "MERGED"},
		{Number: 409, Path: "/repo/.wk/409", Title: "Worktree support", State: "OPEN"},
		{Number: 423, Path: "/repo/.wk/423", Title: "Cap the fuzzy length penalty", State: "CLOSED"},
		{Number: 431, Path: "/repo/.wk/431"},
	}
}

func testWorktreeFetchFunc(entries []model.WorktreeEntry) WorktreeFetchFunc {
	return func(bool) ([]model.WorktreeEntry, error) { return entries, nil }
}

func testWorktreeOptions() WorktreeOptions {
	return WorktreeOptions{Prompt: "> ", Placeholder: "Filter worktrees..."}
}

// newTestWorktreeModel creates a model and simulates the async load completing,
// on a terminal wide enough that no title is clipped.
func newTestWorktreeModel() WorktreeModel {
	return newTestWorktreeModelWith(testWorktrees(), testWorktreeOptions())
}

func newTestWorktreeModelWith(entries []model.WorktreeEntry, opts WorktreeOptions) WorktreeModel {
	m := NewWorktreeModel(testWorktreeFetchFunc(entries), opts)
	result, _ := m.Update(worktreesLoadedMsg{entries: entries})
	result, _ = result.(WorktreeModel).Update(tea.WindowSizeMsg{Width: 120, Height: 24})
	return result.(WorktreeModel)
}

// worktreeView renders the model's frame with the ANSI styling stripped.
func worktreeView(m WorktreeModel) string {
	return ansi.Strip(m.View().Content)
}

// typeQuery feeds a query in one rune at a time, the way it would be typed.
func typeQuery(m WorktreeModel, query string) WorktreeModel {
	for _, r := range query {
		result, _ := m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
		m = result.(WorktreeModel)
	}
	return m
}

// filteredNumbers is the numbers of the rows currently shown, in order.
func filteredNumbers(m WorktreeModel) []int {
	numbers := make([]int, 0, len(m.filtered))
	for _, row := range m.filtered {
		numbers = append(numbers, row.item.entry.Number)
	}
	return numbers
}

func TestNewWorktreeModel(t *testing.T) {
	m := newTestWorktreeModel()
	assert.Len(t, m.allItems, 4)
	assert.Len(t, m.filtered, 4)
	assert.Equal(t, 0, m.cursor)
	assert.False(t, m.loading)
	assert.False(t, m.quit)

	_, picked := m.Chosen()
	assert.False(t, picked, "nothing is chosen until enter is pressed")
}

func TestNewWorktreeModel_StartsInLoadingState(t *testing.T) {
	m := NewWorktreeModel(testWorktreeFetchFunc(testWorktrees()), testWorktreeOptions())
	assert.True(t, m.loading)
	assert.Empty(t, m.allItems)
	assert.Contains(t, worktreeView(m), "Loading worktrees...")
}

func TestWorktreeModel_LoadError(t *testing.T) {
	m := NewWorktreeModel(nil, testWorktreeOptions())
	result, cmd := m.Update(worktreesLoadedMsg{err: errors.New("no [[worktree]] block")})

	assert.Error(t, result.(WorktreeModel).LoadErr())
	assert.NotNil(t, cmd, "a load error must quit the picker")
}

func TestWorktreeModel_EmptyList(t *testing.T) {
	m := newTestWorktreeModelWith(nil, testWorktreeOptions())
	assert.Contains(t, worktreeView(m), "No worktrees found")
}

func TestWorktreeView_ShowsNumberBadgeAndTitle(t *testing.T) {
	out := worktreeView(newTestWorktreeModel())

	assert.Contains(t, out, "409   OPEN    Worktree support")
	assert.Contains(t, out, "9   MERGED  Add an opt-in preview pane")
	assert.Contains(t, out, "423   CLOSED  Cap the fuzzy length penalty")
}

func TestWorktreeView_ShowsPillsWhenIconsAreEnabled(t *testing.T) {
	opts := testWorktreeOptions()
	opts.ShowIcons = true
	out := worktreeView(newTestWorktreeModelWith(testWorktrees(), opts))

	assert.Contains(t, out, "409  "+chipLeftGlyph+"OPEN"+chipRightGlyph+"   Worktree support")
	// The pill is the same width as the squared-off form, so the titles land in
	// the same column whichever form the badges take.
	assert.Equal(t, titleColumn(t, out, "Worktree support"),
		titleColumn(t, worktreeView(newTestWorktreeModel()), "Worktree support"))
}

// titleColumn is the display column a title starts in. It is counted in runes,
// not bytes, because the half circles rounding off a pill are multi-byte.
func titleColumn(t *testing.T, out, title string) int {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		if idx := strings.Index(line, title); idx >= 0 {
			return len([]rune(line[:idx]))
		}
	}
	t.Fatalf("expected a row containing %q", title)
	return 0
}

func TestWorktreeView_AlignsColumns(t *testing.T) {
	out := worktreeView(newTestWorktreeModel())

	// Every title starts in the same column, whatever the width of the number
	// and of the badge before it — including the row that has neither a state
	// nor a title.
	var titleColumns []int
	for _, title := range []string{"Worktree support", "Add an opt-in", "Cap the fuzzy"} {
		titleColumns = append(titleColumns, titleColumn(t, out, title))
	}
	require.Len(t, titleColumns, 3)
	assert.Equal(t, titleColumns[0], titleColumns[1], "titles must share a column")
	assert.Equal(t, titleColumns[1], titleColumns[2], "titles must share a column")

	// The numbers are right-aligned, so the single-digit one is indented to the
	// width of the widest.
	assert.Equal(t, 3, m9NumberWidth(t, out), "numbers must be right-aligned")
}

// m9NumberWidth is the column the single-digit number 9 ends at, which equals
// the width of the number column when the numbers are right-aligned.
func m9NumberWidth(t *testing.T, out string) int {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "MERGED") {
			return strings.Index(line, "9") + 1 - 2 // less the two-column prefix
		}
	}
	t.Fatal("expected a row for worktree 9")
	return 0
}

func TestWorktreeView_UnknownStateGetsNoBadge(t *testing.T) {
	out := worktreeView(newTestWorktreeModel())

	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "431") {
			assert.Equal(t, "431", strings.TrimSpace(line),
				"a worktree with no resolved issue shows its bare number")
			return
		}
	}
	t.Fatal("expected a row for worktree 431")
}

func TestStateBadge(t *testing.T) {
	for state, want := range map[string]string{
		"OPEN":   "OPEN",
		"CLOSED": "CLOSED",
		"MERGED": "MERGED",
		"open":   "OPEN",
	} {
		badge := stateBadge(state, true)
		assert.Contains(t, ansi.Strip(badge), want, "state %q", state)
		assert.Contains(t, badge, "\x1b[7", "state %q must be filled with reverse video", state)
	}

	assert.Equal(t, strings.TrimSpace(ansi.Strip(stateBadge("", true))), "",
		"an absent state renders as blanks")
}

func TestStateBadge_PillShape(t *testing.T) {
	badge := stateBadge("OPEN", true)

	assert.Contains(t, badge, chipLeftGlyph, "the pill is rounded off with half circles")
	assert.Contains(t, badge, chipRightGlyph)
	assert.Equal(t, chipLeftGlyph+"OPEN"+chipRightGlyph, strings.TrimSpace(ansi.Strip(badge)))

	// The half circles carry the fill color as a plain foreground, which is the
	// color the reversed label is filled with, so the pill reads as one shape.
	open, _ := badgeColor("OPEN")
	fill := lipgloss.NewStyle().Foreground(open)
	assert.Contains(t, badge, fill.Render(chipLeftGlyph))
	assert.Contains(t, badge, fill.Render(chipRightGlyph))
}

func TestStateBadge_SquaredOffWithoutNerdFonts(t *testing.T) {
	badge := stateBadge("OPEN", false)

	assert.NotContains(t, badge, chipLeftGlyph,
		"nerd font glyphs are only used when icons are enabled")
	assert.NotContains(t, badge, chipRightGlyph)
	assert.Equal(t, "OPEN", strings.TrimSpace(ansi.Strip(badge)))
	assert.Contains(t, badge, "\x1b[7", "the squared-off form is still filled")
}

func TestStateBadge_FixedColumnWidth(t *testing.T) {
	// Every badge occupies the same width, whatever its state and whichever form
	// it takes, so the titles after them line up.
	for _, showIcons := range []bool{true, false} {
		for _, state := range []string{"OPEN", "CLOSED", "MERGED", "", "SOMETHING_ELSE"} {
			assert.Equal(t, badgeCellWidth, len([]rune(ansi.Strip(stateBadge(state, showIcons)))),
				"badge column width for state %q (icons: %v)", state, showIcons)
		}
	}
}

func TestStateBadge_ColorsAreDistinct(t *testing.T) {
	open, _ := badgeColor("OPEN")
	closed, _ := badgeColor("CLOSED")
	merged, _ := badgeColor("MERGED")
	assert.NotEqual(t, open, closed)
	assert.NotEqual(t, open, merged)
	assert.NotEqual(t, closed, merged)

	_, badged := badgeColor("")
	assert.False(t, badged, "an absent state gets no badge")
	_, badged = badgeColor("DRAFT")
	assert.False(t, badged, "an unrecognized state gets no badge")
}

func TestWorktreeFilter_MatchesOnNumber(t *testing.T) {
	m := typeQuery(newTestWorktreeModel(), "409")

	assert.Equal(t, []int{409}, filteredNumbers(m))
}

func TestWorktreeFilter_MatchesOnTitle(t *testing.T) {
	m := typeQuery(newTestWorktreeModel(), "fuzzy")

	assert.Equal(t, []int{423}, filteredNumbers(m))
}

func TestWorktreeFilter_MatchesOnNumberAndTitleTogether(t *testing.T) {
	m := typeQuery(newTestWorktreeModel(), "409 worktree")

	assert.Equal(t, []int{409}, filteredNumbers(m),
		"a query spanning both columns must match as one string")
}

func TestWorktreeFilter_MatchesTitlelessRowByNumber(t *testing.T) {
	m := typeQuery(newTestWorktreeModel(), "431")

	assert.Equal(t, []int{431}, filteredNumbers(m),
		"a worktree whose title never resolved is still reachable by number")
}

func TestWorktreeFilter_EmptyQueryKeepsListOrder(t *testing.T) {
	m := typeQuery(newTestWorktreeModel(), "4")
	require.NotEqual(t, 4, len(m.filtered), "the query must narrow the list first")

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	m = result.(WorktreeModel)

	assert.Equal(t, []int{9, 409, 423, 431}, filteredNumbers(m),
		"an empty query shows every worktree, in the order it was listed")
}

func TestWorktreeFilter_NoMatch(t *testing.T) {
	m := typeQuery(newTestWorktreeModel(), "zzzzz")

	assert.Empty(t, m.filtered)
	// Enter on an empty list must not pick a phantom row.
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	_, picked := result.(WorktreeModel).Chosen()
	assert.False(t, picked)
}

func TestWorktreeFilter_ResetsCursor(t *testing.T) {
	result, _ := newTestWorktreeModel().Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m := result.(WorktreeModel)
	require.Equal(t, 1, m.cursor)

	m = typeQuery(m, "4")
	assert.Equal(t, 0, m.cursor, "filtering must put the cursor back on the first row")
	assert.Equal(t, 0, m.offset)
}

func TestWorktreeFilter_HighlightsMatchesInBothColumns(t *testing.T) {
	m := typeQuery(newTestWorktreeModel(), "409 worktree")
	require.Len(t, m.filtered, 1)

	row := m.filtered[0]
	assert.Equal(t, []int{0, 1, 2}, row.numberIndexes, "the number's own runes")
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7}, row.titleIndexes,
		"title indexes must be rebased past the number and its separating space")
}

func TestSplitMatchIndexes(t *testing.T) {
	// "409 Worktree": the space at index 3 belongs to neither column.
	numberIndexes, titleIndexes := splitMatchIndexes([]int{0, 2, 3, 4, 6}, 3)
	assert.Equal(t, []int{0, 2}, numberIndexes)
	assert.Equal(t, []int{0, 2}, titleIndexes)

	numberIndexes, titleIndexes = splitMatchIndexes(nil, 3)
	assert.Nil(t, numberIndexes)
	assert.Nil(t, titleIndexes)
}

func TestWorktreeModel_Enter_PicksHighlightedEntry(t *testing.T) {
	result, _ := newTestWorktreeModel().Update(tea.KeyPressMsg{Code: tea.KeyDown})
	result, cmd := result.(WorktreeModel).Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	entry, picked := result.(WorktreeModel).Chosen()
	assert.True(t, picked)
	assert.Equal(t, 409, entry.Number)
	assert.Equal(t, "/repo/.wk/409", entry.Path,
		"the whole entry is returned, so the caller keeps the path and state")
	assert.NotNil(t, cmd, "enter must quit the picker")
}

func TestWorktreeModel_Enter_WhileLoadingDoesNothing(t *testing.T) {
	m := NewWorktreeModel(testWorktreeFetchFunc(testWorktrees()), testWorktreeOptions())
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	_, picked := result.(WorktreeModel).Chosen()
	assert.False(t, picked)
	assert.Nil(t, cmd, "enter must not quit before the list has loaded")
}

func TestWorktreeModel_Quit(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{
		{Code: tea.KeyEscape},
		{Code: 'c', Mod: tea.ModCtrl},
	} {
		result, cmd := newTestWorktreeModel().Update(key)
		m := result.(WorktreeModel)

		assert.True(t, m.quit)
		_, picked := m.Chosen()
		assert.False(t, picked, "quitting must not pick the highlighted row")
		assert.NotNil(t, cmd)
	}
}

func TestWorktreeModel_CursorMovement(t *testing.T) {
	m := newTestWorktreeModel()

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Equal(t, 0, result.(WorktreeModel).cursor, "the cursor stops at the top")

	for i := 0; i < 10; i++ {
		result, _ = result.(WorktreeModel).Update(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	assert.Equal(t, 3, result.(WorktreeModel).cursor, "the cursor stops at the last row")

	result, _ = result.(WorktreeModel).Update(tea.KeyPressMsg{Code: 'k', Mod: tea.ModCtrl})
	assert.Equal(t, 2, result.(WorktreeModel).cursor, "ctrl+k moves up")

	result, _ = result.(WorktreeModel).Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	assert.Equal(t, 3, result.(WorktreeModel).cursor, "ctrl+n moves down")
}

func TestWorktreeModel_CursorOnEmptyListStaysValid(t *testing.T) {
	m := typeQuery(newTestWorktreeModel(), "zzzzz")
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})

	assert.Equal(t, 0, result.(WorktreeModel).cursor,
		"an empty list must not move the cursor off the end")
}

func TestWorktreeModel_ScrollsWithinVisibleRows(t *testing.T) {
	entries := make([]model.WorktreeEntry, 0, 20)
	for i := 1; i <= 20; i++ {
		entries = append(entries, model.WorktreeEntry{Number: i, Title: "issue", State: "OPEN"})
	}
	m := NewWorktreeModel(testWorktreeFetchFunc(entries), testWorktreeOptions())
	result, _ := m.Update(worktreesLoadedMsg{entries: entries})
	// Five rows fit: seven lines less the filter row and the blank under it.
	result, _ = result.(WorktreeModel).Update(tea.WindowSizeMsg{Width: 80, Height: 7})
	m = result.(WorktreeModel)
	require.Equal(t, 5, m.visibleCount())

	for i := 0; i < 6; i++ {
		result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		m = result.(WorktreeModel)
	}
	assert.Equal(t, 6, m.cursor)
	assert.Equal(t, 2, m.offset, "the list scrolls to keep the cursor visible")

	// The frame stays exactly the height of the terminal however far it scrolls.
	assert.Len(t, strings.Split(worktreeView(m), "\n"), 7)
}

func TestWorktreeView_ClipsLongTitles(t *testing.T) {
	long := strings.Repeat("really long title ", 10)
	entries := []model.WorktreeEntry{{Number: 1, Title: long, State: "OPEN"}}
	m := newTestWorktreeModelWith(entries, testWorktreeOptions())
	result, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	m = result.(WorktreeModel)

	out := worktreeView(m)
	// The worktree rows start after the filter row and the blank line under it.
	for _, line := range strings.Split(out, "\n")[headerLines:] {
		assert.LessOrEqual(t, len([]rune(line)), 40,
			"no row may exceed the terminal width, or it wraps and breaks the row count")
	}
	assert.Contains(t, out, "…", "a clipped title is marked as clipped")
	assert.Len(t, strings.Split(out, "\n"), 10, "the frame is the height of the terminal")
}

func TestWorktreeView_ClipsHighlightsWithTheTitle(t *testing.T) {
	entries := []model.WorktreeEntry{{Number: 1, Title: strings.Repeat("ab", 40), State: "OPEN"}}
	m := newTestWorktreeModelWith(entries, testWorktreeOptions())
	result, _ := m.Update(tea.WindowSizeMsg{Width: 30, Height: 10})
	m = typeQuery(result.(WorktreeModel), "ab")
	require.Len(t, m.filtered, 1)

	// Clipping past a match index must not index off the end of the title.
	title, indexes := m.clipTitle(entries[0].Title, []int{0, 1, 70}, 20)
	assert.Contains(t, title, "…")
	assert.Equal(t, []int{0, 1}, indexes, "matches past the cut are dropped")
}

func TestWorktreeView_NoClipBeforeTerminalSizeIsKnown(t *testing.T) {
	entries := []model.WorktreeEntry{{Number: 1, Title: strings.Repeat("x", 200), State: "OPEN"}}
	m := NewWorktreeModel(testWorktreeFetchFunc(entries), testWorktreeOptions())
	result, _ := m.Update(worktreesLoadedMsg{entries: entries})

	title, indexes := result.(WorktreeModel).clipTitle(entries[0].Title, []int{5}, 20)
	assert.Equal(t, entries[0].Title, title, "an unknown width gives nothing to clip against")
	assert.Equal(t, []int{5}, indexes)
}

func TestWorktreeView_CursorMarksHighlightedRow(t *testing.T) {
	m := newTestWorktreeModel()
	lines := strings.Split(worktreeView(m), "\n")

	assert.True(t, strings.HasPrefix(lines[2], "> "), "the first row starts highlighted")
	assert.True(t, strings.HasPrefix(lines[3], "  "), "other rows are not")
}

func TestWorktreeModel_CtrlR_RefetchesIgnoringTheCache(t *testing.T) {
	var refreshes []bool
	fetchFunc := func(refresh bool) ([]model.WorktreeEntry, error) {
		refreshes = append(refreshes, refresh)
		return testWorktrees(), nil
	}
	m := NewWorktreeModel(fetchFunc, testWorktreeOptions())
	result, _ := m.Update(worktreesLoadedMsg{entries: testWorktrees()})
	result, _ = result.(WorktreeModel).Update(tea.WindowSizeMsg{Width: 120, Height: 24})

	result, cmd := result.(WorktreeModel).Update(tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl})
	m = result.(WorktreeModel)
	require.NotNil(t, cmd, "ctrl+r must kick off a refetch")
	assert.True(t, m.refreshing)
	assert.Contains(t, worktreeView(m), "Refreshing worktrees...")
	// The rows already fetched stay on screen while the refetch runs.
	assert.Contains(t, worktreeView(m), "Worktree support")

	msg := cmd()
	assert.Equal(t, []bool{true}, refreshes, "the refetch must ignore the cache")

	result, _ = m.Update(msg)
	m = result.(WorktreeModel)
	assert.False(t, m.refreshing)
	assert.NotContains(t, worktreeView(m), "Refreshing worktrees...")
	assert.Len(t, m.filtered, 4)
}

func TestWorktreeModel_CtrlR_KeepsTheFilterAndTheCursor(t *testing.T) {
	m := typeQuery(newTestWorktreeModel(), "4")
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(WorktreeModel)
	require.Equal(t, 1, m.cursor)

	result, cmd := m.Update(tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl})
	result, _ = result.(WorktreeModel).Update(cmd())
	m = result.(WorktreeModel)

	assert.Equal(t, "4", m.filterInput.Value(), "a refresh must not clear the query")
	assert.Equal(t, []int{431, 409, 423}, filteredNumbers(m),
		"the query still narrows and ranks the refetched rows")
	assert.Equal(t, 1, m.cursor, "a refresh must not move the cursor")
}

func TestWorktreeModel_CtrlR_PullsTheCursorBackOntoAShrunkList(t *testing.T) {
	m := newTestWorktreeModel()
	for i := 0; i < 3; i++ {
		result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		m = result.(WorktreeModel)
	}
	require.Equal(t, 3, m.cursor)

	// The refetch comes back with the worktrees the cursor was sitting past
	// having been removed.
	result, _ := m.Update(worktreesLoadedMsg{entries: testWorktrees()[:2], refreshed: true})
	m = result.(WorktreeModel)

	assert.Equal(t, 1, m.cursor, "the cursor must land on the last remaining row")
	assert.Equal(t, 0, m.offset)
}

func TestWorktreeModel_CtrlR_IsIgnoredWhileAFetchIsInFlight(t *testing.T) {
	m := NewWorktreeModel(testWorktreeFetchFunc(testWorktrees()), testWorktreeOptions())
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl})
	assert.Nil(t, cmd, "the initial load is already fetching everything")

	m = newTestWorktreeModel()
	result, cmd := m.Update(tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl})
	require.NotNil(t, cmd)
	_, cmd = result.(WorktreeModel).Update(tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl})
	assert.Nil(t, cmd, "a second ctrl+r must not stack a second refetch")
}

func TestWorktreeModel_CtrlR_FailureKeepsTheList(t *testing.T) {
	m := newTestWorktreeModel()
	result, _ := m.Update(tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl})
	result, cmd := result.(WorktreeModel).Update(worktreesLoadedMsg{
		err: errors.New("gh rate limited"), refreshed: true,
	})
	m = result.(WorktreeModel)

	assert.Nil(t, cmd, "a failed refresh must not end the picker")
	assert.NoError(t, m.LoadErr())
	assert.Len(t, m.filtered, 4, "the rows fetched before the failure are still good")
	assert.Contains(t, worktreeView(m), "Refresh failed: gh rate limited")

	// The next refresh clears the report of the last one.
	result, _ = m.Update(tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl})
	assert.NotContains(t, worktreeView(result.(WorktreeModel)), "Refresh failed")
}

// autoRefreshOptions is the picker as `sesh worktree picker` builds it: cached
// rows first, then a refetch behind them.
func autoRefreshOptions() WorktreeOptions {
	opts := testWorktreeOptions()
	opts.AutoRefresh = true
	return opts
}

func TestWorktreeModel_AutoRefresh_CorrectsAStaleTitle(t *testing.T) {
	cached := []model.WorktreeEntry{
		{Number: 409, Path: "/repo/.wk/409", Title: "The old title", State: "OPEN"},
	}
	renamed := []model.WorktreeEntry{
		{Number: 409, Path: "/repo/.wk/409", Title: "The title I just edited", State: "OPEN"},
	}

	var refreshes []bool
	fetchFunc := func(refresh bool) ([]model.WorktreeEntry, error) {
		refreshes = append(refreshes, refresh)
		if refresh {
			return renamed, nil
		}
		return cached, nil
	}

	m := NewWorktreeModel(fetchFunc, autoRefreshOptions())
	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 24})
	result, cmd := result.(WorktreeModel).Update(worktreesLoadedMsg{entries: cached})
	m = result.(WorktreeModel)

	require.NotNil(t, cmd, "the cached rows must be followed by a refetch")
	assert.True(t, m.refreshing)
	assert.Contains(t, worktreeView(m), "The old title",
		"the cached title is on screen while the refetch runs")

	result, _ = m.Update(cmd())
	m = result.(WorktreeModel)

	assert.Equal(t, []bool{true}, refreshes, "the automatic refetch must ignore the cache")
	assert.False(t, m.refreshing)
	assert.Contains(t, worktreeView(m), "The title I just edited",
		"the edited title replaces the cached one without the user asking")
}

func TestWorktreeModel_AutoRefresh_RunsOnlyOnce(t *testing.T) {
	m := NewWorktreeModel(testWorktreeFetchFunc(testWorktrees()), autoRefreshOptions())
	result, cmd := m.Update(worktreesLoadedMsg{entries: testWorktrees()})
	require.NotNil(t, cmd)

	result, cmd = result.(WorktreeModel).Update(cmd())
	assert.Nil(t, cmd, "the refetch landing must not schedule another")

	// And a later ctrl+r is still available.
	_, cmd = result.(WorktreeModel).Update(tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl})
	assert.NotNil(t, cmd, "ctrl+r still forces a refetch after the automatic one")
}

func TestWorktreeModel_AutoRefresh_FailureIsSilent(t *testing.T) {
	m := NewWorktreeModel(testWorktreeFetchFunc(testWorktrees()), autoRefreshOptions())
	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 24})
	result, cmd := result.(WorktreeModel).Update(worktreesLoadedMsg{entries: testWorktrees()})
	require.NotNil(t, cmd)

	result, quitCmd := result.(WorktreeModel).Update(worktreesLoadedMsg{
		err: errors.New("dial tcp: no such host"), refreshed: true,
	})
	m = result.(WorktreeModel)

	assert.Nil(t, quitCmd, "a failed automatic refresh must not end the picker")
	assert.NoError(t, m.LoadErr())
	assert.Len(t, m.filtered, 4, "the cached rows it was trying to improve are still good")
	assert.NotContains(t, worktreeView(m), "Refresh failed",
		"nobody asked for this refresh, so its failure is not reported")

	// A ctrl+r the user did ask for still reports its failure.
	result, _ = m.Update(tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl})
	result, _ = result.(WorktreeModel).Update(worktreesLoadedMsg{
		err: errors.New("gh rate limited"), refreshed: true,
	})
	assert.Contains(t, worktreeView(result.(WorktreeModel)), "Refresh failed: gh rate limited")
}

func TestWorktreeModel_AutoRefresh_SkippedWhenThereAreNoWorktrees(t *testing.T) {
	m := NewWorktreeModel(testWorktreeFetchFunc(nil), autoRefreshOptions())
	_, cmd := m.Update(worktreesLoadedMsg{})
	assert.Nil(t, cmd, "there are no titles to refetch")
}

func TestWorktreeStatusLine_DoesNotWrap(t *testing.T) {
	m := newTestWorktreeModel()
	result, _ := m.Update(tea.WindowSizeMsg{Width: 20, Height: 24})
	result, _ = result.(WorktreeModel).Update(worktreesLoadedMsg{
		err: errors.New(strings.Repeat("very long error ", 10)), refreshed: true,
	})
	m = result.(WorktreeModel)

	// The status line is the second, under the filter row.
	status := strings.Split(worktreeView(m), "\n")[1]
	assert.LessOrEqual(t, len([]rune(status)), 20,
		"a long error must be clipped, or it wraps and breaks the row count")
	assert.Contains(t, status, "…")
	assert.Len(t, strings.Split(worktreeView(m), "\n"), 24,
		"the frame is still the height of the terminal")
}

func TestWorktreeOptions_QueryPrefillsFilter(t *testing.T) {
	opts := testWorktreeOptions()
	opts.Query = "fuzzy"
	m := newTestWorktreeModelWith(testWorktrees(), opts)

	assert.Equal(t, "fuzzy", m.filterInput.Value())
	assert.Equal(t, []int{423}, filteredNumbers(m),
		"a pre-filled query filters the list as soon as it loads")
}

func TestWorktreeOptions_QueryIsBackspaceable(t *testing.T) {
	opts := testWorktreeOptions()
	opts.Query = "409"
	m := newTestWorktreeModelWith(testWorktrees(), opts)

	for i := 0; i < 3; i++ {
		result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
		m = result.(WorktreeModel)
	}
	assert.Equal(t, "", m.filterInput.Value())
	assert.Len(t, m.filtered, 4, "the pre-filled query can be cleared like anything typed")
}
