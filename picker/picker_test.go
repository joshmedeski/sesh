package picker

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/model"
)

func testSessions() model.SeshSessions {
	dir := model.SeshSessionMap{
		"s1": {Name: "my-project", Src: "tmux", Path: "/home/user/my-project"},
		"s2": {Name: "dotfiles", Src: "config", Path: "/home/user/dotfiles"},
		"s3": {Name: "~/code/app", Src: "zoxide", Path: "/home/user/code/app"},
		"s4": {Name: "rails-app", Src: "tmuxinator", Path: "/home/user/rails-app"},
		"s5": {Name: "notes", Src: "tmux", Path: "/home/user/notes"},
	}
	return model.SeshSessions{
		OrderedIndex: []string{"s1", "s2", "s3", "s4", "s5"},
		Directory:    dir,
	}
}

func testFetchFunc(sessions model.SeshSessions) FetchFunc {
	return func() (model.SeshSessions, error) {
		return sessions, nil
	}
}

func testOptions() Options {
	return Options{
		Prompt:                "> ",
		Placeholder:           "Filter sessions...",
		AliasAutoConnectDelay: 150 * time.Millisecond,
		AliasFilterPrefix:     model.DefaultAliasFilterPrefix,
	}
}

// testOptionsWith builds the default test options and lets a test tweak the
// handful of fields it cares about.
func testOptionsWith(fn func(*Options)) Options {
	opts := testOptions()
	fn(&opts)
	return opts
}

// newTestModel creates a model and simulates the async load completing.
func newTestModel() Model {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptions())
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	return result.(Model)
}

func TestNew(t *testing.T) {
	m := newTestModel()
	assert.Len(t, m.allItems, 5)
	assert.Len(t, m.filtered, 5)
	assert.Equal(t, 0, m.cursor)
	assert.Equal(t, "", m.chosen)
	assert.False(t, m.quit)
	assert.False(t, m.loading)
}

func TestNew_StartsInLoadingState(t *testing.T) {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptions())
	assert.True(t, m.loading)
	assert.Len(t, m.allItems, 0)
	assert.Len(t, m.filtered, 0)
}

func TestSrcIcon(t *testing.T) {
	for _, src := range []string{"tmux", "config", "zoxide", "tmuxinator"} {
		icn, clr := srcIcon(src)
		assert.NotEmpty(t, icn, "icon for %s should not be empty", src)
		assert.NotEqual(t, "? ", icn, "icon for %s should not be fallback", src)
		assert.NotNil(t, clr, "color for %s should not be nil", src)
	}

	icn, clr := srcIcon("other")
	assert.Equal(t, "? ", icn)
	assert.NotNil(t, clr)
}

func TestApplyFilter_EmptyPattern(t *testing.T) {
	m := newTestModel()
	assert.Len(t, m.filtered, 5)
	assert.Equal(t, "my-project", m.filtered[0].item.name)
}

func TestApplyFilter_WithPattern(t *testing.T) {
	m := newTestModel()
	m.filterInput.SetValue("dot")
	m.applyFilter()

	assert.Equal(t, 1, len(m.filtered))
	assert.Equal(t, "dotfiles", m.filtered[0].item.name)
	assert.Greater(t, len(m.filtered[0].matchedIndexes), 0)
}

func TestApplyFilter_FuzzyMatch(t *testing.T) {
	m := newTestModel()
	m.filterInput.SetValue("mp")
	m.applyFilter()

	found := false
	for _, f := range m.filtered {
		if f.item.name == "my-project" {
			found = true
			break
		}
	}
	assert.True(t, found, "fuzzy match should find 'my-project' for pattern 'mp'")
}

func TestApplyFilter_NoMatches(t *testing.T) {
	m := newTestModel()
	m.filterInput.SetValue("zzzzzzz")
	m.applyFilter()

	assert.Len(t, m.filtered, 0)
}

func TestCursorDown(t *testing.T) {
	m := newTestModel()
	m.height = 30

	m.cursorDown(1)
	assert.Equal(t, 1, m.cursor)

	m.cursorDown(1)
	assert.Equal(t, 2, m.cursor)
}

func TestCursorDown_ClampsAtEnd(t *testing.T) {
	m := newTestModel()
	m.height = 30

	m.cursorDown(100)
	assert.Equal(t, 4, m.cursor)
}

func TestCursorUp(t *testing.T) {
	m := newTestModel()
	m.height = 30
	m.cursor = 3

	m.cursorUp(1)
	assert.Equal(t, 2, m.cursor)
}

func TestCursorUp_ClampsAtZero(t *testing.T) {
	m := newTestModel()
	m.cursor = 0

	m.cursorUp(5)
	assert.Equal(t, 0, m.cursor)
}

func TestUpdate_Escape(t *testing.T) {
	m := newTestModel()
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	resultModel := result.(Model)

	assert.True(t, resultModel.Quit())
	assert.Equal(t, "", resultModel.Chosen())
}

func TestUpdate_CtrlC(t *testing.T) {
	m := newTestModel()
	result, _ := m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	resultModel := result.(Model)

	assert.True(t, resultModel.Quit())
}

func TestUpdate_Enter_SelectsSession(t *testing.T) {
	m := newTestModel()
	m.height = 30

	m.cursorDown(1)
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	resultModel := result.(Model)

	assert.False(t, resultModel.Quit())
	assert.Equal(t, "dotfiles", resultModel.Chosen())
}

func TestUpdate_Enter_ReturnsRawName(t *testing.T) {
	m := newTestModel()
	m.height = 30

	// Select the first item (tmux source "my-project")
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	resultModel := result.(Model)

	// Chosen should be the raw session name with no icon prefix
	assert.Equal(t, "my-project", resultModel.Chosen())
	assert.False(t, strings.HasPrefix(resultModel.Chosen(), "\033"), "Chosen() should not contain ANSI escape codes")
}

func TestUpdate_Enter_EmptyList(t *testing.T) {
	m := newTestModel()
	m.filterInput.SetValue("zzzzzzz")
	m.applyFilter()

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	resultModel := result.(Model)

	assert.Equal(t, "", resultModel.Chosen())
}

func TestUpdate_Enter_WhileLoading(t *testing.T) {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptions())
	assert.True(t, m.loading)

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	resultModel := result.(Model)

	assert.Equal(t, "", resultModel.Chosen(), "enter while loading should not select anything")
	assert.Nil(t, cmd, "enter while loading should not quit")
	assert.True(t, resultModel.loading, "should still be loading")
}

func TestUpdate_Escape_WhileLoading(t *testing.T) {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptions())

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	resultModel := result.(Model)

	assert.True(t, resultModel.Quit(), "escape while loading should quit")
}

func TestUpdate_SessionsLoaded(t *testing.T) {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptions())
	assert.True(t, m.loading)

	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	resultModel := result.(Model)

	assert.False(t, resultModel.loading)
	assert.Len(t, resultModel.allItems, 5)
	assert.Len(t, resultModel.filtered, 5)
	assert.Nil(t, resultModel.loadErr)
}

func TestUpdate_SessionsLoaded_WithPreTypedFilter(t *testing.T) {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptions())

	// Simulate typing "dot" before sessions arrive
	m.filterInput.SetValue("dot")

	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	resultModel := result.(Model)

	assert.False(t, resultModel.loading)
	assert.Len(t, resultModel.allItems, 5)
	assert.Equal(t, 1, len(resultModel.filtered), "pre-typed filter should be applied on load")
	assert.Equal(t, "dotfiles", resultModel.filtered[0].item.name)
}

func TestUpdate_SessionsLoadError(t *testing.T) {
	fetchErr := errors.New("zoxide not found")
	m := New(func() (model.SeshSessions, error) {
		return model.SeshSessions{}, fetchErr
	}, testOptions())

	result, _ := m.Update(sessionsLoadedMsg{err: fetchErr})
	resultModel := result.(Model)

	assert.Equal(t, fetchErr, resultModel.LoadErr())
}

func TestUpdate_ArrowDown(t *testing.T) {
	m := newTestModel()
	m.height = 30

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	resultModel := result.(Model)

	assert.Equal(t, 1, resultModel.cursor)
}

func TestUpdate_ArrowUp(t *testing.T) {
	m := newTestModel()
	m.height = 30
	m.cursor = 2

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	resultModel := result.(Model)

	assert.Equal(t, 1, resultModel.cursor)
}

func TestUpdate_WindowSize(t *testing.T) {
	m := newTestModel()

	result, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	resultModel := result.(Model)

	assert.Equal(t, 80, resultModel.width)
	assert.Equal(t, 24, resultModel.height)
}

func TestView_ReturnsNonEmpty(t *testing.T) {
	m := newTestModel()
	m.width = 60
	m.height = 24

	v := m.View()
	assert.NotZero(t, v)
}

func TestView_LoadingState(t *testing.T) {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptions())
	m.width = 60
	m.height = 24

	assert.True(t, m.Loading(), "model should be in loading state")
	// View should render without panicking
	v := m.View()
	assert.NotZero(t, v)
}

func TestHighlightMatches_NoIndexes(t *testing.T) {
	match := lipgloss.NewStyle().Bold(true)
	normal := lipgloss.NewStyle()
	result := highlightMatches("hello", nil, match, normal)
	assert.Contains(t, result, "hello")
}

func TestHighlightMatches_WithIndexes(t *testing.T) {
	match := lipgloss.NewStyle().Bold(true)
	normal := lipgloss.NewStyle()
	result := highlightMatches("hello", []int{0, 2}, match, normal)
	assert.NotEmpty(t, result)
}

func TestScrolling(t *testing.T) {
	m := newTestModel()
	// Short enough that the five test sessions can't all fit at once, which is
	// what makes the list scroll at all.
	m.height = 6

	visible := m.visibleCount()
	require.Less(t, visible, len(m.filtered), "the list must overflow to scroll")
	for i := 0; i < visible+2; i++ {
		m.cursorDown(1)
	}

	assert.Greater(t, m.offset, 0)
}

func TestHalfPageMovement(t *testing.T) {
	dir := make(model.SeshSessionMap)
	index := make([]string, 20)
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("s%d", i)
		dir[key] = model.SeshSession{Name: fmt.Sprintf("session-%d", i), Src: "tmux"}
		index[i] = key
	}
	sessions := model.SeshSessions{OrderedIndex: index, Directory: dir}
	m := New(testFetchFunc(sessions), testOptions())
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)
	m.height = 20

	half := m.visibleCount() / 2
	result, _ = m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	resultModel := result.(Model)
	assert.Equal(t, half, resultModel.cursor)
}

// newTestModelSeparatorAware creates a model with separator-aware matching enabled.
func newTestModelSeparatorAware() Model {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) { o.SeparatorAware = true }))
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	return result.(Model)
}

func TestNormalizeSeparators(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"my-project", "my project"},
		{"my_project", "my project"},
		{"~/code/app", "~ code app"},
		{"path\\to\\dir", "path to dir"},
		{"a-b_c/d\\e", "a b c d e"},
		{"no separators", "no separators"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalizeSeparators(tt.input))
		})
	}
}

func TestApplyFilter_SeparatorAware_SpaceMatchesDash(t *testing.T) {
	m := newTestModelSeparatorAware()
	m.filterInput.SetValue("my project")
	m.applyFilter()

	found := false
	for _, f := range m.filtered {
		if f.item.name == "my-project" {
			found = true
			break
		}
	}
	assert.True(t, found, "space in pattern should match dash in session name")
}

func TestApplyFilter_SeparatorAware_SpaceMatchesSlash(t *testing.T) {
	m := newTestModelSeparatorAware()
	m.filterInput.SetValue("code app")
	m.applyFilter()

	found := false
	for _, f := range m.filtered {
		if f.item.name == "~/code/app" {
			found = true
			break
		}
	}
	assert.True(t, found, "space in pattern should match slash in session name")
}

func TestApplyFilter_SeparatorAware_PatternNormalization(t *testing.T) {
	m := newTestModelSeparatorAware()
	m.filterInput.SetValue("my-project")
	m.applyFilter()

	found := false
	for _, f := range m.filtered {
		if f.item.name == "my-project" {
			found = true
			break
		}
	}
	assert.True(t, found, "dash in pattern should still match dash in session name")
}

func TestApplyFilter_SeparatorAware_Disabled(t *testing.T) {
	m := newTestModel()
	m.filterInput.SetValue("my project")
	m.applyFilter()

	for _, f := range m.filtered {
		assert.NotEqual(t, "my-project", f.item.name,
			"space should NOT match dash when separator-aware is disabled")
	}
}

// longNameSessions mirrors the shape of the bug the capped length penalty
// fixes: a tmux session whose name is long but starts with the query, competing
// with short zoxide paths that only contain it.
func longNameSessions() model.SeshSessions {
	dir := model.SeshSessionMap{
		"s1": {Name: "~/c/sesh", Src: "zoxide"},
		"s2": {Name: "~/c/re/sesh", Src: "zoxide"},
		"s3": {Name: "~/c/sesh/w/411", Src: "zoxide"},
		"s4": {Name: "sesh/w/423 — Add opt-in preview pane to the picker TUI", Src: "tmux"},
	}
	return model.SeshSessions{
		OrderedIndex: []string{"s1", "s2", "s3", "s4"},
		Directory:    dir,
	}
}

func newLongNameModel() Model {
	sessions := longNameSessions()
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) {
		o.SeparatorAware = true
	}))
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	return result.(Model)
}

func TestApplyFilter_LongNameMatchingAtStartRanksFirst(t *testing.T) {
	m := newLongNameModel()
	m.filterInput.SetValue("sesh")
	m.applyFilter()

	assert.Len(t, m.filtered, 4)
	assert.Equal(t, "sesh/w/423 — Add opt-in preview pane to the picker TUI",
		m.filtered[0].item.name,
		"a long name matching at its first character should outrank shorter partial matches")
}

func TestApplyFilter_LengthStillBreaksComparableMatches(t *testing.T) {
	m := newLongNameModel()
	m.filterInput.SetValue("c sesh")
	m.applyFilter()

	names := make([]string, 0, len(m.filtered))
	for _, f := range m.filtered {
		names = append(names, f.item.name)
	}
	assert.Equal(t, []string{"~/c/sesh", "~/c/sesh/w/411", "~/c/re/sesh"}, names,
		"among comparable matches the shorter name should still come first")
}

func TestApplyFilter_ExactNameStillBeatsALongerPrefixMatch(t *testing.T) {
	sessions := model.SeshSessions{
		OrderedIndex: []string{"s1", "s2"},
		Directory: model.SeshSessionMap{
			"s1": {Name: "sesh plus a very long tail of unmatched characters", Src: "tmux"},
			"s2": {Name: "sesh", Src: "tmux"},
		},
	}
	m := New(testFetchFunc(sessions), testOptions())
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)
	m.filterInput.SetValue("sesh")
	m.applyFilter()

	assert.Equal(t, "sesh", m.filtered[0].item.name,
		"capping the penalty should not let a long name overtake an exact match")
}

// sessionsWithWindows returns sessions carrying window names for display.
func sessionsWithWindows() model.SeshSessions {
	dir := model.SeshSessionMap{
		"s1": {Name: "sesh", Src: "tmux", WindowNames: []string{"editor", "server", "logs"}},
		"s2": {Name: "dotfiles", Src: "tmux", WindowNames: []string{"nvim", "shell"}},
		"s3": {Name: "scratch", Src: "tmux"},
	}
	return model.SeshSessions{
		OrderedIndex: []string{"s1", "s2", "s3"},
		Directory:    dir,
	}
}

// newTestModelWithWindows creates a loaded model with window display enabled.
func newTestModelWithWindows() Model {
	sessions := sessionsWithWindows()
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) { o.ShowWindows = true }))
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)
	m.width = 60
	m.height = 24
	return m
}

func TestView_ShowWindows(t *testing.T) {
	m := newTestModelWithWindows()
	out := fmt.Sprintf("%v", m.View())

	assert.Contains(t, out, "sesh")
	assert.Contains(t, out, "editor")
	assert.Contains(t, out, "server")
	assert.Contains(t, out, "nvim")
}

func TestView_ShowWindows_Disabled(t *testing.T) {
	sessions := sessionsWithWindows()
	m := New(testFetchFunc(sessions), testOptions())
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)
	m.width = 60
	m.height = 24

	out := fmt.Sprintf("%v", m.View())
	assert.Contains(t, out, "sesh")
	assert.NotContains(t, out, "editor")
}

func TestView_ShowWindows_NoWindowsRendersCleanly(t *testing.T) {
	m := newTestModelWithWindows()
	out := fmt.Sprintf("%v", m.View())

	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "scratch") {
			assert.Equal(t, "  scratch", strings.TrimRight(line, " \r"),
				"a session without window names should render exactly as before")
		}
	}
}

func TestView_ShowWindows_RowsFitContentWidth(t *testing.T) {
	dir := model.SeshSessionMap{
		"s1": {Name: "sesh", Src: "tmux", WindowNames: []string{
			"editor", "server", "logs", "database", "tests", "docs", "shell",
			"migrations", "worker", "queue", "metrics",
		}},
	}
	sessions := model.SeshSessions{OrderedIndex: []string{"s1"}, Directory: dir}
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) { o.ShowWindows = true }))
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)
	m.width = 60
	m.height = 24

	out := fmt.Sprintf("%v", m.View())
	var row string
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "sesh") {
			row = line
			break
		}
	}
	assert.NotEmpty(t, row)
	assert.LessOrEqual(t, lipgloss.Width(row), m.contentWidth(),
		"row must not exceed the content width")
	assert.Contains(t, row, "+", "elided windows should be summarized as +N")
}

func TestEnter_SelectsSessionNameOnly(t *testing.T) {
	m := newTestModelWithWindows()

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	resultModel := result.(Model)

	assert.Equal(t, "sesh", resultModel.Chosen(),
		"selection must be the session name, never window names")
}

func TestApplyFilter_DoesNotMatchWindowNames(t *testing.T) {
	m := newTestModelWithWindows()
	m.filterInput.SetValue("editor")
	m.applyFilter()

	assert.Empty(t, m.filtered, "typing a window name must not match its session")
}

func TestSessionItems_StringIgnoresWindowNames(t *testing.T) {
	items := buildItems(sessionsWithWindows(), false)
	for i := range items {
		assert.Equal(t, items[i].name, items.String(i),
			"fuzzy source must expose the session name only")
	}
}

// testAliases mirrors what buildAliases produces for a config with `wp` on
// wallpaper (auto-connect) and `dot` on dotfiles (no auto-connect).
func testAliases() map[string]Alias {
	return map[string]Alias{
		"wp":  {Alias: "wp", Target: "my-project", AutoConnect: true},
		"dot": {Alias: "dot", Target: "dotfiles"},
	}
}

// newAliasModel returns a loaded model with testAliases applied. The delay is
// kept tiny so tests that wait out a real tick stay fast.
func newAliasModel(tweak ...func(*Options)) Model {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) {
		o.Aliases = testAliases()
		o.AliasAutoConnectDelay = time.Millisecond
		for _, fn := range tweak {
			fn(o)
		}
	}))
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)
	m.width = 60
	m.height = 24
	return m
}

// typeFilter types a value one key at a time, exercising the same code path as
// a real user, and returns the model plus the command from the last keystroke.
func typeFilter(m Model, value string) (Model, tea.Cmd) {
	var cmd tea.Cmd
	for _, r := range value {
		var result tea.Model
		result, cmd = m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
		m = result.(Model)
	}
	return m, cmd
}

// aliasTick runs a command and digs out the auto-connect message it produces,
// returning nil when no auto-connect was scheduled. The text input always
// returns a cursor-blink command, so the command being non-nil says nothing on
// its own.
func aliasTick(cmd tea.Cmd) *aliasAutoConnectMsg {
	if cmd == nil {
		return nil
	}
	switch msg := cmd().(type) {
	case aliasAutoConnectMsg:
		return &msg
	case tea.BatchMsg:
		for _, batched := range msg {
			if found := aliasTick(batched); found != nil {
				return found
			}
		}
	}
	return nil
}

func TestBuildAliases(t *testing.T) {
	aliases := buildAliases([]model.SessionConfig{
		{Name: "wallpaper", Alias: "WP", AliasAutoConnect: true},
		{Name: "dotfiles", Alias: "dot"},
		{Name: "notes"},
		{Name: "", Alias: "orphan"},
	})

	assert.Len(t, aliases, 2, "sessions without an alias or name are skipped")
	assert.Equal(t, Alias{Alias: "WP", Target: "wallpaper", AutoConnect: true}, aliases["wp"],
		"aliases are keyed lowercased but keep their configured casing for display")
	assert.False(t, aliases["dot"].AutoConnect)
}

func TestAliasAutoConnectDelay(t *testing.T) {
	assert.Equal(t, 300*time.Millisecond, aliasAutoConnectDelay("300ms"))
	assert.Equal(t, 150*time.Millisecond, aliasAutoConnectDelay(""), "empty falls back to the default")
	assert.Equal(t, 150*time.Millisecond, aliasAutoConnectDelay("nonsense"), "invalid falls back to the default")
}

// reverseVideo is the escape sequence that swaps foreground and background,
// which is what keeps the chip label legible under any color scheme.
const reverseVideo = "\x1b[7m"

func TestAliasChip(t *testing.T) {
	m := newAliasModel(func(o *Options) { o.ShowIcons = true })
	chip := m.aliasChip("my-project", 0)
	assert.Contains(t, chip, "wp")
	assert.Contains(t, chip, chipLeftGlyph)
	assert.Contains(t, chip, chipRightGlyph)
	assert.Contains(t, chip, reverseVideo, "the chip label must be inverted to stay legible")
	assert.Equal(t, "", m.aliasChip("notes", 0), "sessions without an alias get no chip")

	plain := newAliasModel()
	assert.Contains(t, ansi.Strip(plain.aliasChip("my-project", 0)), "[wp]")
	assert.Contains(t, plain.aliasChip("my-project", 0), reverseVideo)
	assert.NotContains(t, plain.aliasChip("my-project", 0), chipLeftGlyph,
		"nerd font glyphs are only used when icons are enabled")
}

func TestAliasChip_HighlightsMatchedPrefix(t *testing.T) {
	m := newAliasModel()

	// The matched runes carry a color on top of reverse video, so they read as a
	// colored block; the rest of the label stays plain reverse video.
	assert.NotEqual(t, m.aliasChip("my-project", 0), m.aliasChip("my-project", 1),
		"a matched prefix must be styled differently from an unmatched label")
	assert.Equal(t, "[wp] ", ansi.Strip(m.aliasChip("my-project", 1)),
		"highlighting must not change the text of the chip")
	assert.Equal(t, m.aliasChip("my-project", 2), m.aliasChip("my-project", 99),
		"a match longer than the alias is clamped rather than panicking")
}

func TestView_AliasChip(t *testing.T) {
	m := newAliasModel()
	out := ansi.Strip(fmt.Sprintf("%v", m.View()))

	assert.Contains(t, out, "[wp] my-project")
	assert.Contains(t, out, "[dot] dotfiles")
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "notes") {
			assert.Equal(t, "  notes", strings.TrimRight(line, " \r"),
				"a session without an alias should render exactly as before")
		}
	}
}

func TestView_AliasChip_RowsFitContentWidth(t *testing.T) {
	dir := model.SeshSessionMap{
		"s1": {Name: "sesh", Src: "tmux", WindowNames: []string{
			"editor", "server", "logs", "database", "tests", "docs", "shell",
		}},
	}
	sessions := model.SeshSessions{OrderedIndex: []string{"s1"}, Directory: dir}
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) {
		o.ShowWindows = true
		o.Aliases = map[string]Alias{"sh": {Alias: "sh", Target: "sesh"}}
	}))
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)
	m.width = 60
	m.height = 24

	out := fmt.Sprintf("%v", m.View())
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "sesh") {
			assert.LessOrEqual(t, lipgloss.Width(line), m.contentWidth(),
				"the chip must be counted against the row's width budget")
			return
		}
	}
	t.Fatal("expected a row for the sesh session")
}

func TestApplyFilter_NonAliasQueriesAreUnaffected(t *testing.T) {
	plain := newTestModel()
	plain.filterInput.SetValue("o")
	plain.applyFilter()

	aliased := newAliasModel()
	aliased.filterInput.SetValue("o")
	aliased.applyFilter()

	assert.Equal(t, len(plain.filtered), len(aliased.filtered))
	for i := range plain.filtered {
		assert.Equal(t, plain.filtered[i].item.name, aliased.filtered[i].item.name,
			"anything short of an exact alias must rank exactly as before")
	}
}

func TestApplyFilter_ExactAliasIsTheOnlyResult(t *testing.T) {
	m := newAliasModel()
	m.filterInput.SetValue("wp")
	m.applyFilter()

	require.Len(t, m.filtered, 1, "an exact alias resolves to its session, whatever the ranking")
	assert.Equal(t, "my-project", m.filtered[0].item.name)
	assert.Empty(t, m.filtered[0].matchedIndexes,
		"the alias matched the session, not characters within its name")
	assert.Equal(t, "tmux", m.filtered[0].item.src,
		"the listed session is reused so its source and windows come along")
}

func TestApplyFilter_ExactAliasIsCaseInsensitive(t *testing.T) {
	m := newAliasModel()
	m.filterInput.SetValue("WP")
	m.applyFilter()

	require.Len(t, m.filtered, 1)
	assert.Equal(t, "my-project", m.filtered[0].item.name)
}

func TestApplyFilter_AliasTargetMissingFromList(t *testing.T) {
	m := newAliasModel(func(o *Options) {
		o.Aliases = map[string]Alias{"wp": {Alias: "wp", Target: "wallpaper"}}
	})
	m.filterInput.SetValue("wp")
	m.applyFilter()

	require.Len(t, m.filtered, 1,
		"an alias always names a [[session]], so it stays selectable when unlisted")
	assert.Equal(t, "wallpaper", m.filtered[0].item.name)
	assert.Equal(t, "config", m.filtered[0].item.src)
}

func TestApplyFilter_PartialAliasFallsBackToFuzzy(t *testing.T) {
	m := newAliasModel()
	m.filterInput.SetValue("do")
	m.applyFilter()

	require.Len(t, m.filtered, 1)
	assert.NotEmpty(t, m.filtered[0].matchedIndexes,
		"a partial alias is fuzzy-matched like any other query")
}

func TestApplyFilter_TypingPastAnAliasFallsBackToFuzzy(t *testing.T) {
	m := newAliasModel()
	m.filterInput.SetValue("dotf")
	m.applyFilter()

	require.Len(t, m.filtered, 1)
	assert.NotEmpty(t, m.filtered[0].matchedIndexes,
		"past the alias the query is fuzzy-matched again")
}

func TestApplyFilter_ExactAliasIgnoresSeparatorNormalization(t *testing.T) {
	m := newAliasModel(func(o *Options) {
		o.SeparatorAware = true
		o.Aliases = map[string]Alias{"w-p": {Alias: "w-p", Target: "my-project"}}
	})

	m.filterInput.SetValue("w p")
	m.applyFilter()
	assert.NotEqual(t, 1, len(m.filtered),
		"normalized input must not resolve an alias containing a separator")

	m.filterInput.SetValue("w-p")
	m.applyFilter()
	require.Len(t, m.filtered, 1)
	assert.Equal(t, "my-project", m.filtered[0].item.name)
}

func TestEnter_AfterTypingAliasSelectsTheTarget(t *testing.T) {
	m, _ := typeFilter(newAliasModel(), "dot")

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	assert.Equal(t, "dotfiles", result.(Model).Chosen(),
		"enter must land on the aliased session even without auto-connect")
}

func TestAliasAutoConnect_FiresOnExactAlias(t *testing.T) {
	m, cmd := typeFilter(newAliasModel(), "wp")
	tick := aliasTick(cmd)
	if assert.NotNil(t, tick, "typing an auto-connect alias should schedule a tick") {
		assert.Equal(t, m.aliasSeq, tick.seq)
	}

	result, quitCmd := m.Update(*tick)
	resultModel := result.(Model)

	assert.Equal(t, "my-project", resultModel.Chosen())
	assert.NotNil(t, quitCmd)
	assert.False(t, resultModel.Quit(), "auto-connect is a selection, not a cancel")
}

func TestAliasAutoConnect_MatchesCaseInsensitively(t *testing.T) {
	m, cmd := typeFilter(newAliasModel(), "WP")
	assert.NotNil(t, aliasTick(cmd))

	result, _ := m.Update(aliasAutoConnectMsg{seq: m.aliasSeq, alias: "wp"})
	assert.Equal(t, "my-project", result.(Model).Chosen())
}

func TestAliasAutoConnect_IgnoresStaleTick(t *testing.T) {
	m, _ := typeFilter(newAliasModel(), "wp")

	result, _ := m.Update(aliasAutoConnectMsg{seq: m.aliasSeq - 1, alias: "wp"})
	assert.Equal(t, "", result.(Model).Chosen(),
		"a tick from an earlier keystroke must not connect")
}

func TestAliasAutoConnect_IgnoresTickAfterFilterChanged(t *testing.T) {
	m, _ := typeFilter(newAliasModel(), "wp")
	seq := m.aliasSeq
	m.filterInput.SetValue("wpx")

	result, _ := m.Update(aliasAutoConnectMsg{seq: seq, alias: "wp"})
	assert.Equal(t, "", result.(Model).Chosen(),
		"keeping typing past an alias must not connect to it")
}

func TestAliasAutoConnect_SkipsAliasesWithoutAutoConnect(t *testing.T) {
	m, cmd := typeFilter(newAliasModel(), "dot")
	assert.Nil(t, aliasTick(cmd), "an alias without alias_auto_connect should not schedule a tick")

	result, _ := m.Update(aliasAutoConnectMsg{seq: m.aliasSeq, alias: "dot"})
	assert.Equal(t, "", result.(Model).Chosen(),
		"even a hand-delivered tick must not connect without alias_auto_connect")
}

func TestAliasAutoConnect_PartialAliasDoesNothing(t *testing.T) {
	_, cmd := typeFilter(newAliasModel(), "w")
	assert.Nil(t, aliasTick(cmd), "a partial alias behaves like a normal fuzzy query")
}

func TestAliasAutoConnect_TrailingWhitespaceDoesNotMatch(t *testing.T) {
	_, cmd := typeFilter(newAliasModel(), "wp ")
	assert.Nil(t, aliasTick(cmd), "the alias must match the raw input exactly")
}

func TestAliasAutoConnect_Suppressed(t *testing.T) {
	m, cmd := typeFilter(newAliasModel(func(o *Options) { o.DisableAliasAutoConnect = true }), "wp")
	assert.Nil(t, aliasTick(cmd), "--no-alias-auto should stop the tick from being scheduled")

	result, _ := m.Update(aliasAutoConnectMsg{seq: m.aliasSeq, alias: "wp"})
	assert.Equal(t, "", result.(Model).Chosen())
}

func TestAliasAutoConnect_FiresWhileLoading(t *testing.T) {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) {
		o.Aliases = testAliases()
		o.AliasAutoConnectDelay = time.Millisecond
	}))
	assert.True(t, m.loading)

	m, cmd := typeFilter(m, "wp")
	assert.NotNil(t, aliasTick(cmd),
		"the alias target comes from config, so loading must not block it")

	result, _ := m.Update(aliasAutoConnectMsg{seq: m.aliasSeq, alias: "wp"})
	assert.Equal(t, "my-project", result.(Model).Chosen())
}

func TestScheduleAliasAutoConnect_ZeroDelaySkipsTimer(t *testing.T) {
	m := newAliasModel(func(o *Options) { o.AliasAutoConnectDelay = 0 })
	m.filterInput.SetValue("wp")

	cmd := m.scheduleAliasAutoConnect()
	if assert.NotNil(t, cmd) {
		assert.Equal(t, aliasAutoConnectMsg{seq: m.aliasSeq, alias: "wp"}, cmd(),
			"a zero delay should not go through a timer")
	}
}

func TestScheduleAliasAutoConnect_BumpsSeqOnEveryChange(t *testing.T) {
	m := newAliasModel()
	m.filterInput.SetValue("notes")

	assert.Nil(t, m.scheduleAliasAutoConnect(), "a non-alias schedules nothing")
	assert.Equal(t, 1, m.aliasSeq,
		"the sequence still advances so any pending tick is invalidated")
}

func TestAliasAutoConnect_SeparatorAwareDoesNotMangleAliases(t *testing.T) {
	m := newAliasModel(func(o *Options) {
		o.SeparatorAware = true
		o.Aliases = map[string]Alias{"w-p": {Alias: "w-p", Target: "my-project", AutoConnect: true}}
	})

	_, cmd := typeFilter(m, "w p")
	assert.Nil(t, aliasTick(cmd), "normalized input must not match an alias containing a separator")

	_, cmd = typeFilter(m, "w-p")
	assert.NotNil(t, aliasTick(cmd), "the raw alias still matches with separator_aware on")
}

// filteredNames lists the session names currently shown, for asserting on both
// membership and order.
func filteredNames(m Model) []string {
	names := make([]string, 0, len(m.filtered))
	for _, item := range m.filtered {
		names = append(names, item.item.name)
	}
	return names
}

// filterAliasMode types a query into alias-filter mode, sigil included.
func filterAliasMode(m Model, query string) Model {
	m.filterInput.SetValue(m.aliasFilterPrefix + query)
	m.applyFilter()
	return m
}

func TestApplyFilter_AliasModeShowsEveryAlias(t *testing.T) {
	m := filterAliasMode(newAliasModel(), "")

	assert.Equal(t, []string{"my-project", "dotfiles"}, filteredNames(m),
		"the bare sigil narrows to aliased sessions, in list order")
	for _, item := range m.filtered {
		assert.Equal(t, 0, item.chipMatchLen, "nothing is typed yet, so nothing is highlighted")
	}
}

func TestApplyFilter_AliasModeMatchesAliasPrefix(t *testing.T) {
	m := filterAliasMode(newAliasModel(), "d")

	require.Equal(t, []string{"dotfiles"}, filteredNames(m))
	assert.Equal(t, 1, m.filtered[0].chipMatchLen, "the matched part of the chip is highlighted")

	// `xy` on a session whose name shares no letters with the query isolates the
	// alias tier from the session-name fallback.
	mid := newAliasModel(func(o *Options) {
		o.Aliases = map[string]Alias{"xy": {Alias: "xy", Target: "notes"}}
	})
	assert.Empty(t, filteredNames(filterAliasMode(mid, "y")),
		"aliases match by prefix, so a mid-alias substring finds nothing")
}

func TestApplyFilter_AliasModeIsCaseInsensitive(t *testing.T) {
	m := filterAliasMode(newAliasModel(), "DO")
	assert.Equal(t, []string{"dotfiles"}, filteredNames(m))
}

func TestApplyFilter_AliasModeFallsBackToSessionName(t *testing.T) {
	m := filterAliasMode(newAliasModel(), "project")

	require.Equal(t, []string{"my-project"}, filteredNames(m),
		"a session name is matched by substring, since names are multi-word")
	assert.Equal(t, 0, m.filtered[0].chipMatchLen,
		"the query matched the name, not the chip, so the chip is left plain")
}

func TestApplyFilter_AliasModeRanksAliasMatchesFirst(t *testing.T) {
	m := newAliasModel(func(o *Options) {
		o.Aliases = map[string]Alias{
			"do": {Alias: "do", Target: "my-project"},
			"z":  {Alias: "z", Target: "dotfiles"},
		}
	})
	m = filterAliasMode(m, "do")

	assert.Equal(t, []string{"my-project", "dotfiles"}, filteredNames(m),
		"an alias match outranks a session whose name merely contains the query")
}

func TestApplyFilter_AliasModeIncludesUnlistedAliases(t *testing.T) {
	m := newAliasModel(func(o *Options) {
		o.Aliases = map[string]Alias{
			"dot": {Alias: "dot", Target: "dotfiles"},
			"wp":  {Alias: "wp", Target: "wallpaper"},
			"aa":  {Alias: "aa", Target: "archive"},
		}
	})
	m = filterAliasMode(m, "")

	assert.Equal(t, []string{"dotfiles", "archive", "wallpaper"}, filteredNames(m),
		"unlisted aliases trail the listed ones, sorted so the order never shifts")
	for _, item := range m.filtered[1:] {
		assert.Equal(t, "config", item.item.src)
	}
}

func TestApplyFilter_AliasModeWithNoMatches(t *testing.T) {
	m := filterAliasMode(newAliasModel(), "zzz")
	assert.Empty(t, m.filtered)
}

func TestApplyFilter_AliasModeWithNoAliasesConfigured(t *testing.T) {
	m := newTestModel()
	m.width, m.height = 60, 24
	m = filterAliasMode(m, "")

	assert.Empty(t, m.filtered)
	assert.Contains(t, ansi.Strip(fmt.Sprintf("%v", m.View())), "No aliases configured",
		"the sigil should read as unconfigured rather than broken")
}

func TestApplyFilter_SigilPastTheStartIsANormalQuery(t *testing.T) {
	m := newAliasModel()
	m.filterInput.SetValue("code/app")
	m.applyFilter()

	require.Equal(t, []string{"~/code/app"}, filteredNames(m),
		"only a leading sigil enters alias mode, so path-like queries still work")
	assert.NotEmpty(t, m.filtered[0].matchedIndexes)
}

func TestApplyFilter_CustomAliasFilterPrefix(t *testing.T) {
	m := newAliasModel(func(o *Options) { o.AliasFilterPrefix = "@" })

	m = filterAliasMode(m, "")
	assert.Equal(t, []string{"my-project", "dotfiles"}, filteredNames(m))

	m.filterInput.SetValue("/")
	m.applyFilter()
	assert.NotEmpty(t, m.filtered, "the default sigil is inert once one is configured")
}

func TestApplyFilter_AliasFilterPrefixDisabled(t *testing.T) {
	m := newAliasModel(func(o *Options) { o.AliasFilterPrefix = "" })
	m.filterInput.SetValue("/")
	m.applyFilter()

	assert.Equal(t, []string{"~/code/app"}, filteredNames(m),
		"an empty prefix disables the mode, leaving the sigil a normal query")
}

func ptr[T any](v T) *T { return &v }

func TestAliasFilterPrefix(t *testing.T) {
	assert.Equal(t, "/", aliasFilterPrefix(nil), "an absent key falls back to the default")
	assert.Equal(t, "", aliasFilterPrefix(ptr("")), "an explicit empty string disables the mode")
	assert.Equal(t, "@", aliasFilterPrefix(ptr("@")))
}

func TestAliasAutoConnect_AliasModeFiresWithoutTheOptIn(t *testing.T) {
	m, cmd := typeFilter(newAliasModel(), "/dot")
	require.NotNil(t, aliasTick(cmd),
		"reaching for the sigil is itself the intent, so alias_auto_connect isn't required")

	result, quitCmd := m.Update(aliasAutoConnectMsg{seq: m.aliasSeq, alias: "dot"})
	assert.Equal(t, "dotfiles", result.(Model).Chosen())
	assert.NotNil(t, quitCmd)
}

func TestAliasAutoConnect_AliasModePartialAliasDoesNothing(t *testing.T) {
	_, cmd := typeFilter(newAliasModel(), "/d")
	assert.Nil(t, aliasTick(cmd), "only an exact alias auto-connects")
}

func TestAliasAutoConnect_AliasModeSessionNameMatchDoesNothing(t *testing.T) {
	_, cmd := typeFilter(newAliasModel(), "/project")
	assert.Nil(t, aliasTick(cmd), "matching a session name is not typing an alias")
}

func TestAliasAutoConnect_AliasModeLeavesRoomForALongerAlias(t *testing.T) {
	m := newAliasModel(func(o *Options) {
		o.Aliases = map[string]Alias{
			"w":  {Alias: "w", Target: "my-project"},
			"wp": {Alias: "wp", Target: "dotfiles"},
		}
	})

	m, cmd := typeFilter(m, "/w")
	tick := aliasTick(cmd)
	require.NotNil(t, tick)

	m, _ = typeFilter(m, "p")
	result, _ := m.Update(*tick)
	assert.Equal(t, "", result.(Model).Chosen(),
		"the tick armed by the shorter alias must not fire once the longer one is typed")
}

func TestAliasAutoConnect_AliasModeSuppressed(t *testing.T) {
	m, cmd := typeFilter(newAliasModel(func(o *Options) { o.DisableAliasAutoConnect = true }), "/wp")
	assert.Nil(t, aliasTick(cmd), "--no-alias-auto covers alias mode too")

	result, _ := m.Update(aliasAutoConnectMsg{seq: m.aliasSeq, alias: "wp"})
	assert.Equal(t, "", result.(Model).Chosen())
}

func TestEnter_InAliasModeSelectsTheTarget(t *testing.T) {
	m, _ := typeFilter(newAliasModel(), "/do")

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	assert.Equal(t, "dotfiles", result.(Model).Chosen(),
		"the sigil must never leak into the chosen session name")
}

func TestWindowsText(t *testing.T) {
	tests := []struct {
		name     string
		names    []string
		budget   int
		expected string
	}{
		{"no names", nil, 40, ""},
		{"zero budget", []string{"editor"}, 0, ""},
		{"negative budget", []string{"editor"}, -5, ""},
		{"all fit", []string{"editor", "server"}, 40, " editor server"},
		{"one fits, rest elided", []string{"editor", "server", "logs"}, 14, " editor +2"},
		{"none fit, count still shown", []string{"editorwindow"}, 6, " +1"},
		{"nothing fits at all", []string{"editorwindow"}, 2, ""},
		{"exact fit", []string{"editor"}, 7, " editor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windowsText(tt.names, tt.budget)
			assert.Equal(t, tt.expected, got)
			assert.LessOrEqual(t, lipgloss.Width(got), max(tt.budget, 0))
		})
	}
}

// --- preview pane -----------------------------------------------------------

// testPreviewFunc records every session it is asked about, so tests can assert
// which fetches actually happened.
func testPreviewFunc(asked *[]string) PreviewFunc {
	return func(name string) (string, error) {
		*asked = append(*asked, name)
		return "preview of " + name, nil
	}
}

// previewModel builds a loaded model with the preview pane on and a terminal
// wide enough to split.
func previewModel(tweak ...func(*Options)) (Model, *[]string) {
	asked := &[]string{}
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) {
		o.Preview = true
		o.PreviewFunc = testPreviewFunc(asked)
		for _, fn := range tweak {
			fn(o)
		}
	}))
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)
	m.width = 120
	m.height = 24
	return m, asked
}

// previewFetch digs the preview fetch message out of a command, returning nil
// when none was scheduled. Cursor keys return no other commands, but filter
// keystrokes come back batched with a cursor blink.
func previewFetch(cmd tea.Cmd) *previewFetchMsg {
	if cmd == nil {
		return nil
	}
	switch msg := cmd().(type) {
	case previewFetchMsg:
		return &msg
	case tea.BatchMsg:
		for _, batched := range msg {
			if found := previewFetch(batched); found != nil {
				return found
			}
		}
	}
	return nil
}

// loadPreview runs the fetch a command scheduled and feeds the result back in,
// standing in for the debounce tick and the goroutine.
func loadPreview(t *testing.T, m Model, cmd tea.Cmd) Model {
	t.Helper()
	fetch := previewFetch(cmd)
	require.NotNil(t, fetch, "expected a preview fetch to be scheduled")
	result, started := m.Update(*fetch)
	m = result.(Model)
	require.NotNil(t, started, "the fetch message should start the preview")
	result, _ = m.Update(started())
	return result.(Model)
}

func TestPreview_OffByDefault(t *testing.T) {
	m := newTestModel()
	m.width = 200
	m.height = 24

	assert.False(t, m.previewOn)
	assert.False(t, m.splitActive())
	assert.Equal(t, 0, m.previewCols())
	assert.Equal(t, 60, m.contentWidth(), "list width must be untouched when off")
}

func TestPreview_NoFetchWhenDisabled(t *testing.T) {
	asked := &[]string{}
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) {
		o.PreviewFunc = testPreviewFunc(asked)
	}))
	result, cmd := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)

	assert.Nil(t, previewFetch(cmd))
	result, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Nil(t, previewFetch(cmd))
	assert.Empty(t, *asked)
	assert.False(t, result.(Model).splitActive())
}

func TestPreview_FetchedForFirstSessionOnLoad(t *testing.T) {
	asked := &[]string{}
	sessions := testSessions()
	m := New(testFetchFunc(sessions), testOptionsWith(func(o *Options) {
		o.Preview = true
		o.PreviewFunc = testPreviewFunc(asked)
	}))
	result, cmd := m.Update(sessionsLoadedMsg{sessions: sessions})

	m = loadPreview(t, result.(Model), cmd)
	assert.Equal(t, []string{"my-project"}, *asked)
	assert.Equal(t, "my-project", m.previewName)
	assert.Equal(t, "preview of my-project", m.previewContent)
}

func TestPreview_CursorMoveFetchesNewSession(t *testing.T) {
	m, asked := previewModel()
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})

	m = loadPreview(t, result.(Model), cmd)
	assert.Equal(t, "dotfiles", m.previewName)
	assert.Contains(t, *asked, "dotfiles")
}

func TestPreview_AlreadyPreviewedSessionIsNotRefetched(t *testing.T) {
	m, asked := previewModel()
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = loadPreview(t, result.(Model), cmd)
	before := len(*asked)

	// Move away, then straight back before the new preview lands: the pane is
	// still showing dotfiles, so there is nothing to fetch.
	result, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = result.(Model)
	require.NotNil(t, previewFetch(cmd), "moving away schedules a fetch")

	result, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(Model)

	assert.Nil(t, previewFetch(cmd), "the displayed preview should not be refetched")
	assert.Len(t, *asked, before, "no preview command should have run")
}

func TestPreview_StaleResultIsDiscarded(t *testing.T) {
	m, _ := previewModel()
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = loadPreview(t, result.(Model), cmd)
	require.Equal(t, "dotfiles", m.previewName)

	// A result from a cursor position the user has since left behind.
	result, _ = m.Update(previewLoadedMsg{
		seq:     m.previewSeq - 1,
		name:    "my-project",
		content: "stale content",
	})
	m = result.(Model)

	assert.Equal(t, "dotfiles", m.previewName)
	assert.Equal(t, "preview of dotfiles", m.previewContent,
		"the pane should keep the preview that belongs to the highlighted row")
}

func TestPreview_StaleFetchTickIsDiscarded(t *testing.T) {
	m, asked := previewModel()
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(Model)
	fetch := previewFetch(cmd)
	require.NotNil(t, fetch)

	// The cursor moves again before the debounce elapses, so the earlier tick
	// arrives stale and must not shell out.
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(Model)
	before := len(*asked)
	result, started := m.Update(*fetch)

	assert.Nil(t, started, "a superseded tick should not start a fetch")
	assert.Len(t, *asked, before)
	assert.False(t, result.(Model).loading)
}

func TestPreview_KeepsPreviousContentWhileLoading(t *testing.T) {
	m, _ := previewModel()
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = loadPreview(t, result.(Model), cmd)

	// Move on, but don't deliver the new preview yet.
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(Model)

	assert.Equal(t, "preview of dotfiles", m.previewContent)
	assert.Contains(t, m.View().Content, "preview of dotfiles")
}

func TestPreview_ErrorIsShownWithoutQuitting(t *testing.T) {
	m, _ := previewModel()
	result, _ := m.Update(previewLoadedMsg{
		seq:  m.previewSeq,
		name: "my-project",
		err:  errors.New("preview_command exploded"),
	})
	m = result.(Model)

	assert.False(t, m.Quit())
	assert.NoError(t, m.LoadErr())
	assert.Contains(t, m.View().Content, "Preview unavailable")
}

func TestPreview_EmptyFilterResultClearsPane(t *testing.T) {
	m, _ := previewModel()
	m.previewName = "my-project"
	m.previewContent = "preview of my-project"

	m, _ = typeFilter(m, "zzzznomatch")

	require.Empty(t, m.filtered)
	assert.Empty(t, m.previewContent, "no highlighted row means nothing to preview")
	assert.Empty(t, m.previewName)
}

func TestPreview_CtrlOToggles(t *testing.T) {
	m, asked := previewModel()

	result, _ := m.Update(tea.KeyPressMsg{Code: 'o', Mod: tea.ModCtrl})
	m = result.(Model)
	assert.False(t, m.previewOn)
	assert.False(t, m.splitActive())
	assert.Equal(t, 60, m.contentWidth(), "the list reclaims the width")

	before := len(*asked)
	result, cmd := m.Update(tea.KeyPressMsg{Code: 'o', Mod: tea.ModCtrl})
	m = result.(Model)
	assert.True(t, m.previewOn)
	assert.True(t, m.splitActive())
	m = loadPreview(t, m, cmd)
	assert.Len(t, *asked, before+1, "turning the pane back on fetches a preview")
}

func TestPreview_CtrlOIsInertWithoutAPreviewer(t *testing.T) {
	m := newTestModel()
	m.width = 120
	m.height = 24

	result, cmd := m.Update(tea.KeyPressMsg{Code: 'o', Mod: tea.ModCtrl})
	m = result.(Model)

	assert.Nil(t, cmd)
	assert.False(t, m.previewOn)
}

func TestPreview_NarrowTerminalFallsBackToListOnly(t *testing.T) {
	m, _ := previewModel()
	m.width = 99 // one column under the default minimum

	assert.True(t, m.previewOn, "the setting stays on")
	assert.False(t, m.splitActive(), "but the pane isn't rendered")
	assert.Equal(t, 60, m.contentWidth())

	out := fmt.Sprintf("%v", m.View())
	assert.NotContains(t, out, "preview of", "no preview content on a narrow terminal")

	// Growing the window is enough to bring it back.
	m.width = 120
	assert.True(t, m.splitActive())
}

func TestPreview_SplitWidths(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		pct     int
		preview int
		list    int
	}{
		{"at the minimum the list keeps its floor", 100, 60, 60, 40},
		{"percent governs a middling terminal", 120, 60, 72, 48},
		{"preview absorbs width past the list cap", 200, 60, 140, 60},
		{"a small percent still leaves the cap alone", 200, 20, 140, 60},
		{"a large percent squeezes the list to its floor", 120, 90, 80, 40},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := previewModel(func(o *Options) { o.PreviewWidth = tt.pct })
			m.width = tt.width

			assert.Equal(t, tt.preview, m.previewCols(), "preview columns")
			assert.Equal(t, tt.list, m.contentWidth(), "list columns")
			assert.Equal(t, tt.width, m.previewCols()+m.contentWidth(),
				"the two panes should use the full terminal width")
		})
	}
}

func TestPreview_ViewFitsTerminalWidth(t *testing.T) {
	m, _ := previewModel()
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = loadPreview(t, result.(Model), cmd)

	out := m.View().Content
	assert.Contains(t, out, "preview of dotfiles")
	for i, line := range strings.Split(out, "\n") {
		assert.LessOrEqual(t, lipgloss.Width(line), m.width,
			"line %d may not exceed the terminal width", i)
	}
}

func TestPreview_LongContentIsClipped(t *testing.T) {
	m, _ := previewModel(func(o *Options) {
		o.PreviewFunc = func(string) (string, error) {
			return strings.Repeat(strings.Repeat("x", 400)+"\n", 200), nil
		}
	})
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = loadPreview(t, result.(Model), cmd)

	lines := strings.Split(m.View().Content, "\n")
	assert.Len(t, lines, headerLines+m.visibleCount(),
		"an oversized preview must not grow the picker")
	for i, line := range lines {
		assert.LessOrEqual(t, lipgloss.Width(line), m.width,
			"line %d may not exceed the terminal width", i)
	}
}

func TestClipLines(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		rows     int
		expected string
	}{
		{"fits", "a\nb", 10, 5, "a\nb"},
		{"clips rows", "a\nb\nc", 10, 2, "a\nb"},
		{"clips columns", "abcdef", 3, 1, "abc"},
		{"zero width", "abc", 0, 1, ""},
		{"zero rows", "abc", 10, 0, ""},
		{"carriage returns normalized", "a\r\nb", 10, 2, "a\nb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, clipLines(tt.text, tt.width, tt.rows))
		})
	}
}

func TestClipLines_PreservesColorAndResets(t *testing.T) {
	// tmux previews come from `capture-pane -e`, so the content is full of
	// escapes that must survive truncation and not leak into the next line.
	line := "\x1b[31mred text here\x1b[0m"
	got := clipLines(line, 6, 1)

	assert.Contains(t, got, "\x1b[31m", "color must survive truncation")
	assert.True(t, strings.HasSuffix(got, "\x1b[0m"), "the line must be reset")
	assert.Equal(t, 6, lipgloss.Width(got))
}

func TestPreviewWidth_Resolution(t *testing.T) {
	assert.Equal(t, model.DefaultPreviewWidth, previewWidth(0), "unset falls back")
	assert.Equal(t, model.DefaultPreviewWidth, previewWidth(-10))
	assert.Equal(t, model.MinPreviewWidth, previewWidth(1), "clamped, not rejected")
	assert.Equal(t, model.MaxPreviewWidth, previewWidth(300))
	assert.Equal(t, 45, previewWidth(45))
}

func TestPreviewMinWidth_Resolution(t *testing.T) {
	assert.Equal(t, model.DefaultPreviewMinWidth, previewMinWidth(0))
	assert.Equal(t, model.DefaultPreviewMinWidth, previewMinWidth(-1))
	assert.Equal(t, 80, previewMinWidth(80))
}

func TestPreviewBorder_Resolution(t *testing.T) {
	assert.Equal(t, model.DefaultPreviewBorder, previewBorder(""), "unset falls back")
	assert.Equal(t, model.DefaultPreviewBorder, previewBorder("squiggly"),
		"an unrecognized style falls back rather than failing the picker")
	assert.Equal(t, model.PreviewBorderNone, previewBorder("none"))
	assert.Equal(t, model.PreviewBorderThick, previewBorder("thick"))
	assert.Equal(t, model.PreviewBorderDouble, previewBorder("double"))
}

func TestPreview_BorderStyles(t *testing.T) {
	tests := []struct {
		name       string
		configured string
		divider    string
		chrome     int
	}{
		{"unset draws a line", "", "│", 2},
		{"line", model.PreviewBorderLine, "│", 2},
		{"thick", model.PreviewBorderThick, "┃", 2},
		{"double", model.PreviewBorderDouble, "║", 2},
		{"none draws nothing", model.PreviewBorderNone, "", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := previewModel(func(o *Options) { o.PreviewBorder = tt.configured })
			result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
			m = loadPreview(t, result.(Model), cmd)

			assert.Equal(t, tt.chrome, m.previewChrome(),
				"a hidden divider should give its column to the preview text")

			out := m.View().Content
			require.Contains(t, out, "preview of dotfiles")
			for _, glyph := range []string{"│", "┃", "║"} {
				if glyph == tt.divider {
					assert.Contains(t, out, glyph)
					continue
				}
				assert.NotContains(t, out, glyph, "only the configured divider may be drawn")
			}
			for i, line := range strings.Split(out, "\n") {
				assert.LessOrEqual(t, lipgloss.Width(line), m.width,
					"line %d may not exceed the terminal width", i)
			}
		})
	}
}

func TestView_FillsTerminalHeight(t *testing.T) {
	for _, height := range []int{10, 24, 50, 120} {
		t.Run(fmt.Sprintf("height %d", height), func(t *testing.T) {
			m, asked := previewModel()
			_ = asked
			m.height = height
			result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
			m = loadPreview(t, result.(Model), cmd)

			lines := strings.Split(m.View().Content, "\n")
			assert.Len(t, lines, height, "the frame should use every row")
			assert.True(t, m.View().AltScreen, "full height needs the alt screen")
		})
	}
}

func TestView_FillsTerminalHeight_ListOnly(t *testing.T) {
	m := newTestModel()
	m.width = 80
	m.height = 40

	lines := strings.Split(m.View().Content, "\n")
	assert.Len(t, lines, 40, "the list-only frame should use every row too")
}

func TestVisibleCount_FallsBackBeforeSizeIsKnown(t *testing.T) {
	m := newTestModel()
	assert.Equal(t, 0, m.height, "no WindowSizeMsg has arrived yet")
	assert.Equal(t, fallbackVisibleCount, m.visibleCount())
}
