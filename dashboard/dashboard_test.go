package dashboard

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/model"
)

// stubSection is a minimal Section implementation for testing the Model.
type stubSection struct {
	name     string
	width    float64
	chosen   string
	items    int
	lastView struct {
		width, height int
		focused       bool
	}
	updateCount int
}

func (s *stubSection) Name() string    { return s.name }
func (s *stubSection) Init() tea.Cmd   { return nil }
func (s *stubSection) Chosen() string  { return s.chosen }
func (s *stubSection) TotalItems() int { return s.items }
func (s *stubSection) Width() float64  { return s.width }
func (s *stubSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	s.updateCount++
	return s, nil
}
func (s *stubSection) ViewBorderless(width, height int, focused bool) (string, string) {
	s.lastView = struct {
		width, height int
		focused       bool
	}{width, height, focused}
	return s.name, s.name
}

// testModel builds a Model with stub widgets and default dimensions.
func testModel(widgets ...Section) Model {
	m := Model{
		config:     model.DashboardConfig{},
		sessions:   &SessionsSection{},
		configured: &ConfiguredSection{},
		widgets:    widgets,
		page:       pageOpen,
		focus:      0,
		width:      80,
		height:     24,
	}
	return m.withLayout()
}

// updateModel runs m.Update and type-asserts the result back to Model.
func updateModel(m Model, msg tea.Msg) Model {
	result, _ := m.Update(msg)
	return result.(Model)
}

// pressKey constructs a tea.KeyPressMsg whose String() matches the given key name.
func pressKey(key string) tea.KeyPressMsg {
	switch key {
	case "tab":
		return tea.KeyPressMsg{Code: tea.KeyTab}
	case "shift+tab":
		return tea.KeyPressMsg{Mod: 1, Code: tea.KeyTab} // ModShift = 1
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEsc}
	case "backspace":
		return tea.KeyPressMsg{Code: tea.KeyBackspace}
	case "ctrl+c":
		return tea.KeyPressMsg{Mod: 4, Code: 'c'} // ModCtrl = 1 << 2 = 4
	case "ctrl+d":
		return tea.KeyPressMsg{Mod: 4, Code: 'd'}
	case "ctrl+h":
		return tea.KeyPressMsg{Mod: 4, Code: 'h'}
	case "ctrl+l":
		return tea.KeyPressMsg{Mod: 4, Code: 'l'}
	case "ctrl+j":
		return tea.KeyPressMsg{Mod: 4, Code: 'j'}
	case "ctrl+k":
		return tea.KeyPressMsg{Mod: 4, Code: 'k'}
	default:
		return tea.KeyPressMsg{Text: key, Code: rune(key[0])}
	}
}

// --- BuildSections tests ---

func TestBuildSections_AlwaysBuildsLists(t *testing.T) {
	built := BuildSections(model.DashboardConfig{}, SectionDeps{})
	require.NotNil(t, built.Sessions)
	require.NotNil(t, built.Configured)
	assert.Empty(t, built.Widgets)
}

func TestBuildSections_WidgetsFromConfig(t *testing.T) {
	cfg := model.DashboardConfig{
		Sections: []model.DashboardSectionConfig{
			{Type: "ssh", Title: "SSH"},
			{Type: "docker", Title: "Docker"},
		},
	}
	built := BuildSections(cfg, SectionDeps{})
	require.NotNil(t, built.Sessions)
	require.NotNil(t, built.Configured)
	require.Len(t, built.Widgets, 2)
	assert.Equal(t, "SSH", built.Widgets[0].Name())
	assert.Equal(t, "Docker", built.Widgets[1].Name())
}

func TestBuildSections_SkipsSessionsAndUnknown(t *testing.T) {
	cfg := model.DashboardConfig{
		Sections: []model.DashboardSectionConfig{
			{Type: "sessions", Title: "Sesh"},
			{Type: "bogus"},
			{Type: "details", Title: "Details"},
		},
	}
	built := BuildSections(cfg, SectionDeps{})
	// "sessions" is implicit, "bogus" is unknown, and "details" is no longer
	// in the registry, so no widgets are built.
	assert.Empty(t, built.Widgets)
}

func TestBuildSections_WorkmuxWidget(t *testing.T) {
	cfg := model.DashboardConfig{
		Sections: []model.DashboardSectionConfig{{Type: "workmux", Title: "Agents"}},
	}
	built := BuildSections(cfg, SectionDeps{})
	require.Len(t, built.Widgets, 1)
	assert.Equal(t, "Agents", built.Widgets[0].Name())
}

func TestBuildSections_AIAgentSkipped(t *testing.T) {
	// "aiagent" was removed from the registry; it now hits the unknown-type
	// skip path (and is therefore absent from the built widgets).
	cfg := model.DashboardConfig{
		Sections: []model.DashboardSectionConfig{
			{Type: "aiagent", Title: "AI"},
			{Type: "workmux", Title: "Agents"},
		},
	}
	built := BuildSections(cfg, SectionDeps{})
	require.Len(t, built.Widgets, 1)
	assert.Equal(t, "Agents", built.Widgets[0].Name())
}

func TestBuildSections_SessionsEntryCarriesTitleNotGroups(t *testing.T) {
	cfg := model.DashboardConfig{
		Sections: []model.DashboardSectionConfig{
			{
				Type:   "sessions",
				Title:  "Open Sessions",
				Groups: []model.DashboardGroup{{Name: "dev", Patterns: []string{"~/dev/*"}}},
			},
		},
	}
	built := BuildSections(cfg, SectionDeps{})
	assert.Empty(t, built.Widgets)
	assert.Equal(t, "Open Sessions", built.Sessions.Name())
	// Groups are parsed for backward compatibility but no longer applied.
	assert.Empty(t, built.Sessions.config.Groups)
}

// --- New ---

func TestNewBuildsDefaultModel(t *testing.T) {
	m := New(model.DashboardConfig{}, nil, nil, nil, nil, nil, "/home/user")
	require.NotNil(t, m.sessions)
	require.NotNil(t, m.configured)
	assert.Equal(t, pageOpen, m.page)
	assert.Equal(t, 0, m.focus)
	assert.Equal(t, 21, m.contentHeight) // default height 24 - 3
}

// --- Model: tab switching ---

func TestTabSwitchesPage(t *testing.T) {
	m := testModel()
	m = updateModel(m, pressKey("tab"))
	assert.Equal(t, pageConfigured, m.page)
	m = updateModel(m, pressKey("tab"))
	assert.Equal(t, pageOpen, m.page)
}

func TestShiftTabSwitchesPageBackward(t *testing.T) {
	m := testModel()
	m = updateModel(m, pressKey("shift+tab"))
	assert.Equal(t, pageConfigured, m.page)
}

func TestTabResetsFocus(t *testing.T) {
	m := testModel(&stubSection{name: "a"}, &stubSection{name: "b"})
	m.focus = 2
	m = updateModel(m, pressKey("tab"))
	assert.Equal(t, pageConfigured, m.page)
	assert.Equal(t, 0, m.focus)
}

// --- Model: pane focus navigation ---

func TestCtrlLFocusNext(t *testing.T) {
	m := testModel(&stubSection{name: "a"}, &stubSection{name: "b"})
	m = updateModel(m, pressKey("ctrl+l"))
	assert.Equal(t, 1, m.focus)
	m = updateModel(m, pressKey("ctrl+l"))
	assert.Equal(t, 2, m.focus)
	m = updateModel(m, pressKey("ctrl+l"))
	assert.Equal(t, 0, m.focus)
}

func TestCtrlHFocusPrev(t *testing.T) {
	m := testModel(&stubSection{name: "a"}, &stubSection{name: "b"})
	m = updateModel(m, pressKey("ctrl+h"))
	assert.Equal(t, 2, m.focus)
	m = updateModel(m, pressKey("ctrl+h"))
	assert.Equal(t, 1, m.focus)
	m = updateModel(m, pressKey("ctrl+h"))
	assert.Equal(t, 0, m.focus)
}

func TestBackspaceAliasForCtrlH(t *testing.T) {
	m := testModel(&stubSection{name: "a"}, &stubSection{name: "b"})
	m = updateModel(m, pressKey("backspace"))
	assert.Equal(t, 2, m.focus)
}

func TestCtrlJKMoveBetweenRows(t *testing.T) {
	// Row 1 = [sessions, details], row 2 = [a, b].
	// Flat: [sessions(0), details(1), a(2), b(3)].
	details := NewDetailsSection(model.DashboardSectionConfig{Type: "details", Title: "Details"}, SectionDeps{})
	m := testModel(details, &stubSection{name: "a"}, &stubSection{name: "b"})

	m.focus = 0
	m = updateModel(m, pressKey("ctrl+j")) // down: sessions → row2 col 0 (a)
	assert.Equal(t, 2, m.focus)
	m = updateModel(m, pressKey("ctrl+k")) // up: a → row1 col 0 (sessions)
	assert.Equal(t, 0, m.focus)
}

func TestCtrlJKDownClampsColumn(t *testing.T) {
	// Row 1 = [sessions, details], row 2 = [a] (only one column).
	details := NewDetailsSection(model.DashboardSectionConfig{Type: "details", Title: "Details"}, SectionDeps{})
	m := testModel(details, &stubSection{name: "a"})
	m.focus = 1 // details (row1 col 1)
	m = updateModel(m, pressKey("ctrl+j"))
	// down → row2 col min(1, 0) = 0 (a)
	assert.Equal(t, 2, m.focus)
}

func TestCtrlJKNoopWithoutSecondRow(t *testing.T) {
	m := testModel() // no widgets → only row 1
	m = updateModel(m, pressKey("ctrl+j"))
	assert.Equal(t, 0, m.focus)
	m = updateModel(m, pressKey("ctrl+k"))
	assert.Equal(t, 0, m.focus)
}

func TestCtrlJKNoopOnConfiguredPage(t *testing.T) {
	m := testModel(&stubSection{name: "a"})
	m.page = pageConfigured
	m = updateModel(m, pressKey("ctrl+j"))
	assert.Equal(t, 0, m.focus)
	m = updateModel(m, pressKey("ctrl+k"))
	assert.Equal(t, 0, m.focus)
}

func TestPaneNavigationNoopOnConfiguredPage(t *testing.T) {
	m := testModel(&stubSection{name: "a"})
	m.page = pageConfigured
	m = updateModel(m, pressKey("ctrl+l"))
	assert.Equal(t, 0, m.focus)
}

// decodeKey decodes a single raw C0 control byte the way bubbletea's terminal
// reader does (via ultraviolet's EventDecoder) and returns the resulting
// tea.KeyPressMsg. This exercises the real decode path rather than hand-built
// Key values.
func decodeKey(t *testing.T, b byte) tea.KeyPressMsg {
	t.Helper()
	var d uv.EventDecoder
	_, ev := d.Decode([]byte{b})
	ke, ok := ev.(uv.KeyPressEvent)
	require.Truef(t, ok, "byte 0x%02X should decode to a key press, got %T", b, ev)
	return tea.KeyPressMsg(ke)
}

func TestRealDecodeCtrlNavKeys(t *testing.T) {
	// Regression guard: feed the actual decoded ctrl+h/j/k/l and backspace
	// bytes through handleKey and assert they reach the nav handlers. This
	// fails if a future decoder change makes these bytes map to a different
	// string/code.
	details := NewDetailsSection(model.DashboardSectionConfig{Type: "details", Title: "Details"}, SectionDeps{})
	m := testModel(details, &stubSection{name: "a"}, &stubSection{name: "b"})
	// flat panes: [sessions(0), details(1), a(2), b(3)]

	step := func(b byte) {
		rm, _ := m.handleKey(decodeKey(t, b))
		m = rm.(Model)
	}

	// ctrl+l (0x0C) → next, wrapping.
	step(0x0C)
	assert.Equal(t, 1, m.focus)
	step(0x0C)
	assert.Equal(t, 2, m.focus)
	step(0x0C)
	assert.Equal(t, 3, m.focus)
	step(0x0C)
	assert.Equal(t, 0, m.focus) // wrap

	// ctrl+h (0x08) → prev, wrapping.
	step(0x08)
	assert.Equal(t, 3, m.focus)

	// backspace (0x7F) → alias for ctrl+h.
	step(0x7F)
	assert.Equal(t, 2, m.focus)

	// ctrl+j (0x0A) → down a row: row2 col0, already bottom → no-op.
	step(0x0A)
	assert.Equal(t, 2, m.focus)

	// ctrl+k (0x0B) → up a row: row2 col0 → row1 col0 (sessions).
	step(0x0B)
	assert.Equal(t, 0, m.focus)
}

// --- Model: pane order / row separation ---

func TestRowSeparation_WithDetails(t *testing.T) {
	details := NewDetailsSection(model.DashboardSectionConfig{Type: "details", Title: "Details"}, SectionDeps{})
	ssh := NewSSHSection(model.DashboardSectionConfig{Type: "ssh", Title: "SSH"}, SectionDeps{})
	git := NewGitSection(model.DashboardSectionConfig{Type: "git", Title: "Git"}, SectionDeps{})
	m := testModel(ssh, git, details) // config order: ssh, git, details

	row1 := m.row1Panes()
	row2 := m.row2Panes()
	require.Len(t, row1, 2)
	assert.Same(t, m.sessions, row1[0])
	assert.Same(t, details, row1[1]) // details pulled into row 1
	require.Len(t, row2, 2)
	assert.Same(t, ssh, row2[0])
	assert.Same(t, git, row2[1])
}

func TestRowSeparation_NoDetails(t *testing.T) {
	ssh := NewSSHSection(model.DashboardSectionConfig{Type: "ssh", Title: "SSH"}, SectionDeps{})
	m := testModel(ssh)

	row1 := m.row1Panes()
	row2 := m.row2Panes()
	require.Len(t, row1, 1)
	assert.Same(t, m.sessions, row1[0])
	require.Len(t, row2, 1)
	assert.Same(t, ssh, row2[0])
}

func TestRowSeparation_NoWidgets(t *testing.T) {
	m := testModel()
	assert.Len(t, m.row1Panes(), 1)
	assert.Empty(t, m.row2Panes())
}

func TestSessionsPaneIsFlexByDefault(t *testing.T) {
	m := testModel(&stubSection{name: "a", width: 0.5})
	assert.Equal(t, float64(0), m.sessions.Width())
}

// --- Model: quit ---

func TestQuitKeys(t *testing.T) {
	for _, k := range []string{"q", "esc", "ctrl+c"} {
		t.Run(k, func(t *testing.T) {
			m := testModel()
			m = updateModel(m, pressKey(k))
			assert.True(t, m.quit)
			assert.True(t, m.Quit())
		})
	}
}

// --- Model: select ---

func TestEnterOnWidgetSetsChosenAndQuits(t *testing.T) {
	w := &stubSection{name: "w", chosen: "my-session"}
	m := testModel(w)
	m.focus = 1 // first widget
	result, cmd := m.Update(pressKey("enter"))
	rm := result.(Model)
	assert.Equal(t, "my-session", rm.Chosen())
	assert.NotNil(t, cmd)
}

func TestEnterWithoutChosenDoesNotQuit(t *testing.T) {
	w := &stubSection{name: "w", chosen: ""}
	m := testModel(w)
	m.focus = 1
	m = updateModel(m, pressKey("enter"))
	assert.False(t, m.quit)
	assert.Equal(t, "", m.Chosen())
}

// --- Model: key dispatch ---

func TestKeysForwardedToFocusedWidget(t *testing.T) {
	a := &stubSection{name: "a"}
	b := &stubSection{name: "b"}
	m := testModel(a, b)
	m.focus = 2 // second widget (b)
	m = updateModel(m, pressKey("j"))
	assert.Equal(t, 0, a.updateCount)
	assert.Equal(t, 1, b.updateCount)
}

func TestKeysForwardedToSessionsList(t *testing.T) {
	m := testModel(&stubSection{name: "w"})
	m.focus = 0
	m.sessions = &SessionsSection{
		sessions: []model.SeshSession{{Name: "a"}, {Name: "b"}},
	}
	m = updateModel(m, pressKey("j"))
	assert.Equal(t, 1, m.sessions.cursor)
}

// --- Model: configured page ---

func TestConfiguredNavigationAndSelect(t *testing.T) {
	m := testModel()
	m.configured = &ConfiguredSection{
		sessions: []model.SeshSession{{Name: "a"}, {Name: "b"}},
		running:  map[string]bool{},
	}
	m.page = pageConfigured
	m = updateModel(m, pressKey("j"))
	assert.Equal(t, 1, m.configured.cursor)
	m = updateModel(m, pressKey("enter"))
	assert.Equal(t, "b", m.Chosen())
}

func TestSyncHoveredSessionSkippedOnConfiguredPage(t *testing.T) {
	ds := NewDetailsSection(model.DashboardSectionConfig{Title: "Details"}, SectionDeps{})
	m := testModel(ds)
	m.page = pageConfigured
	_, cmd := m.syncHoveredSession()
	assert.Nil(t, cmd)
}

// --- Model: View ---

func TestViewQuitReturnsEmpty(t *testing.T) {
	m := testModel()
	m.quit = true
	v := m.View()
	assert.Equal(t, "", v.Content)
}

func TestViewTooSmallReturnsMessage(t *testing.T) {
	m := testModel()
	m.tooSmall = true
	v := m.View()
	assert.Contains(t, v.Content, "Terminal too small")
}

func TestViewRendersTabs(t *testing.T) {
	m := testModel()
	v := m.View()
	assert.Contains(t, v.Content, "Open")
	assert.Contains(t, v.Content, "Configured")
}

// --- WindowSizeMsg ---

func TestWindowSizeUpdatesDimensions(t *testing.T) {
	m := testModel()
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	m = updateModel(m, msg)
	assert.Equal(t, 120, m.width)
	assert.Equal(t, 40, m.height)
	assert.False(t, m.tooSmall)
	assert.Equal(t, 37, m.contentHeight)
}

func TestWindowSizeTooSmall(t *testing.T) {
	m := testModel()
	msg := tea.WindowSizeMsg{Width: 10, Height: 3}
	m = updateModel(m, msg)
	assert.True(t, m.tooSmall)
}

// --- Width allocation (per row) ---

func TestComputePaneWidths_SinglePane(t *testing.T) {
	m := Model{width: 100}
	// Single pane: available = 100 - 0 junctions - 2 corners = 98.
	assert.Equal(t, []int{98}, m.computePaneWidths([]Section{&stubSection{}}))
}

func TestComputePaneWidths_FlexSplit(t *testing.T) {
	m := Model{width: 100}
	// 3 flex panes, available = 100 - 2 junctions - 2 corners = 96 → 32 each.
	pw := m.computePaneWidths([]Section{&stubSection{}, &stubSection{}, &stubSection{}})
	assert.Equal(t, []int{32, 32, 32}, pw)
}

func TestComputePaneWidths_FixedAndFlex(t *testing.T) {
	m := Model{width: 100}
	// a(0.3 fixed) + two flex; available = 96 → a=28, remainder 68 → 34 each.
	pw := m.computePaneWidths([]Section{&stubSection{width: 0.3}, &stubSection{}, &stubSection{}})
	assert.Equal(t, []int{28, 34, 34}, pw)
}

func TestComputePaneWidths_AllFixedRemainderToLast(t *testing.T) {
	m := Model{width: 100}
	// all fixed 0.3 + 0.3: available = 97 → 29 each, remainder 39 → last pane.
	pw := m.computePaneWidths([]Section{&stubSection{width: 0.3}, &stubSection{width: 0.3}})
	assert.Equal(t, []int{29, 68}, pw)
}

func TestComputePaneWidths_TotalFractionScale(t *testing.T) {
	m := Model{width: 100}
	// fixed 0.7 + 0.7 sum > 1 → scale 1/1.4: int(97*0.7/1.4)=48 each,
	// remainder 1 → last pane (all fixed → remainder to last).
	pw := m.computePaneWidths([]Section{&stubSection{width: 0.7}, &stubSection{width: 0.7}})
	assert.Equal(t, []int{48, 49}, pw)
}

// --- Row height split ---

func TestRowHeights_NoWidgets(t *testing.T) {
	m := testModel()
	m.height = 30
	m = m.withLayout()
	assert.Equal(t, 27, m.contentHeight)
	assert.Equal(t, 27, m.row1Height) // row 1 takes full height
	assert.Equal(t, 0, m.row2Height)
}

func TestRowHeights_WithWidgets(t *testing.T) {
	m := testModel(&stubSection{name: "a"})
	m.height = 30
	m = m.withLayout()
	assert.Equal(t, 27, m.contentHeight)
	assert.Equal(t, 16, m.row1Height) // 27 * 3/5 = 16
	assert.Equal(t, 11, m.row2Height) // 27 - 16
}

func TestRowHeights_DegenerateTiny(t *testing.T) {
	m := testModel(&stubSection{name: "a"})
	m.height = 10 // contentHeight = 7
	m = m.withLayout()
	// row1 = 4, row2 = 3 (< 4) → clamp: row1 = 7-4 = 3, row2 = 4
	assert.Equal(t, 3, m.row1Height)
	assert.Equal(t, 4, m.row2Height)
}

// --- Header / footer rendering ---

func TestRenderHeaderContainsTabsAndCount(t *testing.T) {
	h := renderHeader(0, 3, 80)
	assert.Contains(t, h, "Open")
	assert.Contains(t, h, "Configured")
	assert.Contains(t, h, "3 active")
}

func TestRenderHeaderSmallDropsCount(t *testing.T) {
	h := renderHeader(0, 3, 40)
	assert.NotContains(t, h, "active")
}

func TestRenderFooterPage0(t *testing.T) {
	f := renderFooter(0, 120, "name", false, "")
	assert.Contains(t, f, "ctrl+d")
	assert.Contains(t, f, "ctrl+h/j/k/l")
	assert.Contains(t, f, "panes")
	assert.Contains(t, f, "sort:name")
	assert.NotContains(t, f, "widgets")
}

func TestRenderFooterPage0DropsLabelsWhenOverflow(t *testing.T) {
	// The labeled tab-1 footer (after dropping the `t group` bind) is 94 cols,
	// so at 90 cols labels are dropped (keys only) rather than wrapping.
	f := renderFooter(0, 90, "name", false, "")
	assert.Contains(t, f, "ctrl+h/j/k/l")
	assert.NotContains(t, f, "panes")
}

func TestRenderFooterPage1(t *testing.T) {
	f := renderFooter(1, 100, "name", false, "")
	assert.NotContains(t, f, "ctrl+d")
	assert.Contains(t, f, "refresh")
}

func TestRenderFooterNarrowKeysOnly(t *testing.T) {
	f := renderFooter(0, 50, "name", false, "")
	assert.Contains(t, f, "ctrl+d")
	assert.NotContains(t, f, "kill")
}

func TestRenderFooterTiny(t *testing.T) {
	f := renderFooter(0, 20, "name", false, "")
	assert.Contains(t, f, "tab")
	assert.Contains(t, f, "j/k")
	assert.Contains(t, f, "enter")
	assert.NotContains(t, f, "ctrl+d")
}

func TestRenderFooterFiltering(t *testing.T) {
	f := renderFooter(0, 120, "name", true, "foo")
	assert.Contains(t, f, "filter:")
	assert.Contains(t, f, "foo")
	assert.Contains(t, f, "esc")
	assert.Contains(t, f, "enter")
	// The filter line replaces all binds, so no quit/help/panes binds remain.
	assert.NotContains(t, f, "quit")
	assert.NotContains(t, f, "help")
	assert.NotContains(t, f, "ctrl+d")
}

// --- Frame rendering ---

func TestRenderFrame_ContainsTitlesAndJunctions(t *testing.T) {
	panes := []framePane{
		{title: "Sessions", content: "row1\nrow2", width: 16, focused: true},
		{title: "Details", content: "x", width: 16, focused: false},
	}
	out := renderFrame(panes, 6)
	assert.Contains(t, out, "Sessions")
	assert.Contains(t, out, "Details")
	for _, j := range []string{"┬", "┴", "┌", "┐", "└", "┘"} {
		assert.Contains(t, out, j)
	}
}

func TestRenderFrame_TopBorderTitlesInOrder(t *testing.T) {
	panes := []framePane{
		{title: "Sessions", content: "", width: 16, focused: false},
		{title: "Details", content: "", width: 16, focused: false},
		{title: "Git", content: "", width: 16, focused: false},
	}
	out := renderFrame(panes, 4)
	top := strings.Split(out, "\n")[0]
	iS := strings.Index(top, "Sessions")
	iD := strings.Index(top, "Details")
	iG := strings.Index(top, "Git")
	assert.True(t, iS >= 0 && iD > iS && iG > iD)
}

func TestRenderFrame_FocusedTitleAccent(t *testing.T) {
	panes := []framePane{
		{title: "Sessions", content: "", width: 16, focused: true},
		{title: "Details", content: "", width: 16, focused: false},
	}
	out := renderFrame(panes, 4)
	assert.Contains(t, out, accentStyle().Render(" Sessions "))
}

func TestRenderFrame_ContentHeight(t *testing.T) {
	panes := []framePane{{title: "A", content: "line", width: 6, focused: false}}
	out := renderFrame(panes, 5)
	lines := strings.Split(out, "\n")
	// 1 top + 3 content + 1 bottom = 5 lines
	assert.Len(t, lines, 5)
}

// --- Configured page frame ---

func TestViewConfiguredPageSinglePaneFrame(t *testing.T) {
	m := testModel()
	m.configured = &ConfiguredSection{config: model.DashboardSectionConfig{Title: "Configured"}}
	m.page = pageConfigured
	m.width = 50
	m.height = 10
	m = m.withLayout()
	out := m.viewConfiguredPage()
	assert.Contains(t, out, "Configured")
	assert.NotContains(t, out, "┬")
	assert.NotContains(t, out, "┴")
}

func TestViewFrameWidthMatchesModel(t *testing.T) {
	// Regression: the shared frame's chrome (n-1 junctions + 2 corners) must be
	// subtracted from pane widths so the frame is exactly m.width wide. If not,
	// the outer .Width() word-wraps the trailing chars onto the next row.
	m := testModel(&stubSection{name: "a"}, &stubSection{name: "b"})
	m.width = 120
	m.height = 30
	m = m.withLayout()
	v := m.View()
	lines := strings.Split(v.Content, "\n")
	assert.Len(t, lines, 30) // header 2 + content 27 + footer 1
	for i, l := range lines {
		assert.Equalf(t, 120, lipgloss.Width(l), "line %d width", i)
	}
}

func TestViewConfiguredPageWidthMatchesModel(t *testing.T) {
	m := testModel()
	m.page = pageConfigured
	m.width = 120
	m.height = 30
	m = m.withLayout()
	v := m.View()
	lines := strings.Split(v.Content, "\n")
	assert.Len(t, lines, 30)
	for i, l := range lines {
		assert.Equalf(t, 120, lipgloss.Width(l), "line %d width", i)
	}
}

// --- Row rendering ---

func TestRenderOpenRow_FullColumns(t *testing.T) {
	row := renderOpenRow(100, true, false, "mysession", 0, 3, "~/code/proj", "main", "+1 ~2", nil, nil)
	assert.Contains(t, row, "mysession")
	assert.Contains(t, row, "(main)")
	assert.Contains(t, row, "+1 ~2")
}

func TestRenderOpenRow_DropsStatusUnder90(t *testing.T) {
	row := renderOpenRow(80, false, false, "s", 0, 1, "~/d", "main", "+1", nil, nil)
	assert.Contains(t, row, "(main)")
	assert.NotContains(t, row, "+1")
}

func TestRenderOpenRow_DropsBranchUnder70(t *testing.T) {
	row := renderOpenRow(60, false, false, "foo", 0, 1, "~/d", "main", "+1", nil, nil)
	assert.NotContains(t, row, "(main)")
}

func TestRenderOpenRow_DropsWindowsUnder50(t *testing.T) {
	row := renderOpenRow(40, false, false, "foo", 0, 1, "~/d", "main", "+1", nil, nil)
	assert.NotContains(t, row, "1w")
}

func TestRenderOpenRow_AttachedIndicator(t *testing.T) {
	row := renderOpenRow(100, false, false, "s", 1, 1, "~/d", "", "", nil, nil)
	assert.Contains(t, row, "●")
	row = renderOpenRow(100, false, false, "s", 0, 1, "~/d", "", "", nil, nil)
	assert.NotContains(t, row, "●")
}

func TestRenderOpenRow_Age(t *testing.T) {
	twoH := time.Now().Add(-2 * time.Hour)
	row := renderOpenRow(100, false, false, "s", 0, 1, "~/d", "", "", &twoH, nil)
	assert.Contains(t, row, "2h")

	threeD := time.Now().Add(-3 * 24 * time.Hour)
	row = renderOpenRow(100, false, false, "s", 0, 1, "~/d", "", "", &threeD, nil)
	assert.Contains(t, row, "3d")

	fourMo := time.Now().Add(-4 * 30 * 24 * time.Hour)
	row = renderOpenRow(100, false, false, "s", 0, 1, "~/d", "", "", &fourMo, nil)
	assert.Contains(t, row, "4mo")

	// Nil/zero created → blank age.
	assert.Equal(t, "", formatAge(nil))
	zero := time.Time{}
	assert.Equal(t, "", formatAge(&zero))
}

func TestRenderOpenRow_Alerts(t *testing.T) {
	row := renderOpenRow(100, false, false, "s", 0, 1, "~/d", "", "", nil, []string{"bell"})
	assert.Contains(t, row, "!")
	row = renderOpenRow(100, false, false, "s", 0, 1, "~/d", "", "", nil, nil)
	assert.NotContains(t, row, "!")
}

func TestRenderOpenRow_CurrentHighlight(t *testing.T) {
	row := renderOpenRow(100, false, true, "mysession", 0, 1, "~/d", "", "", nil, nil)
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

func TestFormatGitStatusColored(t *testing.T) {
	got := formatGitStatus(git.StatusSummary{Staged: 1, Unstaged: 2, Deleted: 3, Untracked: 4})
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
	got := formatGitStatus(git.StatusSummary{Staged: 2, Untracked: 1})
	assert.NotContains(t, got, "~")
	assert.NotContains(t, got, "-")
	assert.Contains(t, got, "\x1b[38;5;10m")
	assert.Contains(t, got, "\x1b[38;5;5m")
}

func TestRenderOpenRow_StyledStatusKeepsWidth(t *testing.T) {
	styled := formatGitStatus(git.StatusSummary{Staged: 1, Unstaged: 2})
	row := renderOpenRow(100, false, false, "mysession", 0, 3, "~/code/proj", "main", styled, nil, nil)
	// The status cell has no outer foreground, so the embedded colours remain.
	assert.Contains(t, row, "\x1b[38;5;10m")
	assert.Contains(t, row, "\x1b[38;5;11m")
	// Embedded escapes do not inflate the visible row width.
	plain := renderOpenRow(100, false, false, "mysession", 0, 3, "~/code/proj", "main", "+1 ~2", nil, nil)
	assert.Equal(t, lipgloss.Width(plain), lipgloss.Width(row))
}

func TestRenderConfiguredRow_StyledStatusKeepsWidth(t *testing.T) {
	styled := formatGitStatus(git.StatusSummary{Deleted: 3, Untracked: 4})
	row := renderConfiguredRow(100, false, "proj", "", false, "~/code/proj", "main", styled)
	assert.Contains(t, row, "\x1b[38;5;9m")
	assert.Contains(t, row, "\x1b[38;5;5m")
	assert.Equal(t, 100, lipgloss.Width(row))
}

// --- Sessions section (flat list) ---

func TestFlattenSessionsSortedAlphabetically(t *testing.T) {
	sessions := model.SeshSessions{
		OrderedIndex: []string{"z", "a", "m"},
		Directory: model.SeshSessionMap{
			"z": {Name: "z"},
			"a": {Name: "a"},
			"m": {Name: "m"},
		},
	}
	flat := flattenSessions(sessions)
	require.Len(t, flat, 3)
	assert.Equal(t, []string{"a", "m", "z"}, []string{flat[0].Name, flat[1].Name, flat[2].Name})
}

func TestSessionsSectionTKeyIsNoop(t *testing.T) {
	s := &SessionsSection{sessions: []model.SeshSession{{Name: "a"}, {Name: "b"}}}
	updated, cmd := s.handleKey(pressKey("t"))
	assert.Nil(t, cmd)
	// The list is flat; `t` no longer collapses/expands anything, so the
	// cursor and list are untouched.
	assert.Equal(t, 0, updated.cursor)
	assert.Len(t, updated.sessions, 2)
}

// --- Configured section ---

func TestConfiguredSectionBasic(t *testing.T) {
	cs := NewConfiguredSection(model.DashboardSectionConfig{Type: "configured", Title: "Configured"}, SectionDeps{})
	assert.Equal(t, "Configured", cs.Name())
	assert.Equal(t, "", cs.Chosen())
	assert.Equal(t, 0, cs.TotalItems())
}

func TestConfiguredSectionViewLoading(t *testing.T) {
	cs := NewConfiguredSection(model.DashboardSectionConfig{Type: "configured", Title: "Configured"}, SectionDeps{})
	title, content := cs.ViewBorderless(80, 10, true)
	assert.Equal(t, "Configured", title)
	assert.Contains(t, content, "Loading")
}

func TestConfiguredSectionSelectItem(t *testing.T) {
	cs := &ConfiguredSection{sessions: []model.SeshSession{{Name: "a"}, {Name: "b"}}}
	cs.cursor = 1
	cs.selectItem()
	assert.Equal(t, "b", cs.chosen)
}

// --- Sort modes ---

func sessionNames(ss []model.SeshSession) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = s.Name
	}
	return out
}

func TestSessionsSectionSortModeCycle(t *testing.T) {
	now := time.Now()
	t1 := now.Add(-1 * time.Hour)
	t2 := now.Add(-2 * time.Hour)
	t3 := now.Add(-3 * time.Hour)
	c1 := now.Add(-10 * time.Hour)
	c2 := now.Add(-20 * time.Hour)
	c3 := now.Add(-30 * time.Hour)

	s := &SessionsSection{
		sessions: []model.SeshSession{
			{Name: "b", LastAttached: &t1, Created: &c2},
			{Name: "a", LastAttached: &t2, Created: &c1},
			{Name: "c", LastAttached: &t3, Created: &c3},
		},
		sortMode: "name",
	}
	s.applySort()
	assert.Equal(t, []string{"a", "b", "c"}, sessionNames(s.sessions))

	s.cycleSortMode() // name → recent
	assert.Equal(t, "recent", s.SortLabel())
	assert.Equal(t, []string{"b", "a", "c"}, sessionNames(s.sessions))

	s.cycleSortMode() // recent → created
	assert.Equal(t, "created", s.SortLabel())
	assert.Equal(t, []string{"a", "b", "c"}, sessionNames(s.sessions))

	s.cycleSortMode() // created → name
	assert.Equal(t, "name", s.SortLabel())
	assert.Equal(t, []string{"a", "b", "c"}, sessionNames(s.sessions))
}

func TestSessionsSectionSKeyCyclesSort(t *testing.T) {
	s := &SessionsSection{
		sessions: []model.SeshSession{{Name: "b"}, {Name: "a"}},
		sortMode: "name",
	}
	updated, _ := s.handleKey(pressKey("s"))
	assert.Equal(t, "recent", updated.sortMode)
	updated, _ = updated.handleKey(pressKey("s"))
	assert.Equal(t, "created", updated.sortMode)
	updated, _ = updated.handleKey(pressKey("s"))
	assert.Equal(t, "name", updated.sortMode)
}

// --- Type-to-filter ---

func TestFilterMatchesCaseInsensitive(t *testing.T) {
	s := &SessionsSection{
		sessions:  []model.SeshSession{{Name: "Alpha"}, {Name: "BETA"}, {Name: "alpine"}},
		filtering: true,
	}
	s.filterQuery = "ALP"
	s.applyFilter()
	require.Len(t, s.filtered, 2)
	assert.Equal(t, "Alpha", s.filtered[0].Name)
	assert.Equal(t, "alpine", s.filtered[1].Name)
}

func TestConfiguredFilterMatchesCaseInsensitive(t *testing.T) {
	s := &ConfiguredSection{
		sessions:  []model.SeshSession{{Name: "Alpha"}, {Name: "beta"}},
		filtering: true,
	}
	s.filterQuery = "ALP"
	s.applyFilter()
	require.Len(t, s.filtered, 1)
	assert.Equal(t, "Alpha", s.filtered[0].Name)
}

func TestFilterRoutingWhileTyping(t *testing.T) {
	m := testModel()
	m.sessions = &SessionsSection{
		sessions:  []model.SeshSession{{Name: "a"}, {Name: "b"}},
		filtering: true,
	}

	// Printable keys append to the query (q doesn't quit).
	m = updateModel(m, pressKey("q"))
	assert.False(t, m.quit)
	assert.Equal(t, "q", m.sessions.filterQuery)

	// tab does not switch page while filtering.
	m = updateModel(m, pressKey("tab"))
	assert.Equal(t, pageOpen, m.page)
	assert.True(t, m.sessions.filtering)

	// ctrl+h does not move focus while filtering (deletes last rune instead).
	m = updateModel(m, pressKey("ctrl+h"))
	assert.Equal(t, 0, m.focus)
	assert.Equal(t, "", m.sessions.filterQuery)

	// esc exits and clears.
	m = updateModel(m, pressKey("esc"))
	assert.False(t, m.sessions.filtering)
	assert.Equal(t, "", m.sessions.filterQuery)

	// ctrl+c always quits even while filtering.
	m.sessions.filtering = true
	m = updateModel(m, pressKey("ctrl+c"))
	assert.True(t, m.quit)
}

func TestFilterEnterExitsWithoutSelecting(t *testing.T) {
	m := testModel()
	m.sessions = &SessionsSection{
		sessions:    []model.SeshSession{{Name: "a"}, {Name: "b"}},
		filtering:   true,
		filterQuery: "a",
	}
	result, _ := m.Update(pressKey("enter"))
	rm := result.(Model)
	assert.False(t, rm.sessions.filtering)
	assert.Equal(t, "", rm.sessions.filterQuery)
	assert.Equal(t, "", rm.Chosen()) // enter exits filter, does not select
}

func TestFilterSlashToggles(t *testing.T) {
	m := testModel()
	m.sessions = &SessionsSection{sessions: []model.SeshSession{{Name: "a"}}}
	m = updateModel(m, pressKey("/"))
	assert.True(t, m.sessions.filtering)
}

// --- Current-session highlight ---

func TestSessionsSectionRenderItemCurrentHighlight(t *testing.T) {
	s := &SessionsSection{
		sessions:    []model.SeshSession{{Name: "active", Path: "/home/u/active"}},
		cursor:      1, // row 0 is not selected, so only the name accent shows
		currentName: "active",
	}
	row := s.renderItem(0, 100)
	assert.Contains(t, row, "\x1b[1;38;5;14m") // bold cyan accent on the name
}

func TestSessionsSectionRenderItemNonCurrentNotAccent(t *testing.T) {
	s := &SessionsSection{
		sessions:    []model.SeshSession{{Name: "active", Path: "/home/u/active"}},
		cursor:      1, // row 0 not selected → no marker accent either
		currentName: "other",
	}
	row := s.renderItem(0, 100)
	assert.NotContains(t, row, "\x1b[38;5;14m")
}

// --- Live preview (DetailsSection) ---

func TestDetailsSectionHoverKicksCapture(t *testing.T) {
	ds := NewDetailsSection(model.DashboardSectionConfig{Title: "Details"}, SectionDeps{}).(*DetailsSection)
	_, cmd := ds.Update(hoveredSessionMsg{Name: "sesh", Path: "/x", Windows: 1})
	assert.NotNil(t, cmd)
}

func TestDetailsSectionPreviewLoadedUpdatesState(t *testing.T) {
	ds := &DetailsSection{hoveredName: "sesh"}
	updated, _ := ds.Update(previewLoadedMsg{name: "sesh", output: "hello\nworld"})
	assert.Equal(t, "hello\nworld", updated.(*DetailsSection).previewOutput)

	// Stale capture for a previous hover is ignored.
	updated2, _ := updated.(*DetailsSection).Update(previewLoadedMsg{name: "other", output: "stale"})
	assert.Equal(t, "hello\nworld", updated2.(*DetailsSection).previewOutput)
}

func TestDetailsSectionPreviewTickContinues(t *testing.T) {
	ds := &DetailsSection{hoveredName: "sesh"}
	_, cmd := ds.Update(previewTickMsg{name: "sesh"})
	assert.NotNil(t, cmd)

	// Stale tick (hover moved) stops the ticker.
	_, cmd = ds.Update(previewTickMsg{name: "other"})
	assert.Nil(t, cmd)

	// Empty hover stops the ticker.
	_, cmd = ds.Update(previewTickMsg{name: ""})
	assert.Nil(t, cmd)
}

func TestDetailsSectionHoverClearStopsPreview(t *testing.T) {
	ds := &DetailsSection{hoveredName: "sesh", previewOutput: "x"}
	updated, cmd := ds.Update(hoveredSessionMsg{Name: ""})
	assert.Nil(t, cmd)
	assert.Equal(t, "", updated.(*DetailsSection).hoveredName)
	assert.Equal(t, "", updated.(*DetailsSection).previewOutput)
}

func TestHitTestPageOpen(t *testing.T) {
	m := testModel(
		&WorkmuxSection{agents: []wmAgent{{}, {}}},
		&SSHSection{hosts: []SSHHost{{Host: "a"}, {Host: "b"}}},
	)
	m.width = 100
	m.height = 30
	m = m.withLayout()

	// Sessions pane: row 1, x 1..98, content y 2..2+row1Height-1.
	idx, sec, row, ok := m.hitTest(10, 3)
	require.True(t, ok)
	assert.Equal(t, 0, idx)
	assert.Equal(t, m.sessions, sec)
	assert.Equal(t, 0, row)

	// Row 2: workmux is pane 0, click its second row.
	row2Top := 2 + m.row1Height
	idx, sec, row, ok = m.hitTest(10, row2Top+2)
	require.True(t, ok)
	assert.Equal(t, 1, idx)
	assert.IsType(t, &WorkmuxSection{}, sec)
	assert.Equal(t, 1, row)

	// Header click → miss.
	_, _, _, ok = m.hitTest(10, 0)
	assert.False(t, ok)

	// Right corner chrome → miss.
	_, _, _, ok = m.hitTest(99, 3)
	assert.False(t, ok)

	// Below content (footer) → miss.
	_, _, _, ok = m.hitTest(10, 2+m.contentHeight)
	assert.False(t, ok)
}

func TestHitTestPageConfigured(t *testing.T) {
	m := testModel()
	m.page = pageConfigured
	m.width = 100
	m.height = 30
	m = m.withLayout()

	idx, sec, row, ok := m.hitTest(50, 3)
	require.True(t, ok)
	assert.Equal(t, 0, idx)
	assert.Equal(t, m.configured, sec)
	assert.Equal(t, 0, row)

	// Left corner chrome → miss.
	_, _, _, ok = m.hitTest(0, 3)
	assert.False(t, ok)
}

func TestMouseClickFocusesPaneAndSelectsRow(t *testing.T) {
	wm := &WorkmuxSection{agents: []wmAgent{{}, {}, {}}, viewHeight: 10}
	m := testModel(wm)
	m.width = 100
	m.height = 30
	m = m.withLayout()

	row2Top := 2 + m.row1Height
	updated := updateModel(m, tea.MouseClickMsg{X: 10, Y: row2Top + 2, Button: tea.MouseLeft})
	assert.Equal(t, 1, updated.focus)
	assert.Equal(t, 1, wm.cursor) // clicked view row 1 → absolute row 1

	// Non-left clicks are ignored.
	updated = updateModel(m, tea.MouseClickMsg{X: 10, Y: row2Top + 2, Button: tea.MouseRight})
	assert.Equal(t, 0, updated.focus)
	assert.Equal(t, 1, wm.cursor)
}

func TestMouseClickScrollsSectionIntoView(t *testing.T) {
	wm := &WorkmuxSection{agents: make([]wmAgent, 30), viewHeight: 10, offset: 20, cursor: 20}
	m := testModel(wm)
	m.width = 100
	m.height = 30
	m = m.withLayout()

	row2Top := 2 + m.row1Height
	updateModel(m, tea.MouseClickMsg{X: 10, Y: row2Top + 5, Button: tea.MouseLeft})
	assert.Equal(t, 24, wm.cursor)
	assert.Equal(t, 20, wm.offset)
}
