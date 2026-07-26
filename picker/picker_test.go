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
	m.height = 12

	visible := m.visibleCount()
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
	chip := m.aliasChip("my-project")
	assert.Contains(t, chip, "wp")
	assert.Contains(t, chip, chipLeftGlyph)
	assert.Contains(t, chip, chipRightGlyph)
	assert.Contains(t, chip, reverseVideo, "the chip label must be inverted to stay legible")
	assert.Equal(t, "", m.aliasChip("notes"), "sessions without an alias get no chip")

	plain := newAliasModel()
	assert.Contains(t, plain.aliasChip("my-project"), "[wp]")
	assert.Contains(t, plain.aliasChip("my-project"), reverseVideo)
	assert.NotContains(t, plain.aliasChip("my-project"), chipLeftGlyph,
		"nerd font glyphs are only used when icons are enabled")
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
