package dashboard

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
func (s *stubSection) View(width, height int, focused bool) string {
	s.lastView = struct {
		width, height int
		focused       bool
	}{width, height, focused}
	return s.name
}

func stubs(names ...string) []Section {
	sections := make([]Section, len(names))
	for i, n := range names {
		sections[i] = &stubSection{name: n}
	}
	return sections
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
	case "ctrl+c":
		return tea.KeyPressMsg{Mod: 4, Code: 'c'} // ModCtrl = 1 << 2 = 4
	default:
		return tea.KeyPressMsg{Text: key, Code: rune(key[0])}
	}
}

// --- BuildSections tests ---

func TestBuildSections_DefaultsToSessions(t *testing.T) {
	cfg := model.DashboardConfig{}
	sections := BuildSections(cfg, SectionDeps{})
	require.Len(t, sections, 1)
	assert.Equal(t, "Sessions", sections[0].Name())
}

func TestBuildSections_SkipsUnknownTypes(t *testing.T) {
	cfg := model.DashboardConfig{
		Sections: []model.DashboardSectionConfig{
			{Type: "bogus"},
			{Type: "sessions", Title: "Test"},
		},
	}
	sections := BuildSections(cfg, SectionDeps{})
	require.Len(t, sections, 1)
	assert.Equal(t, "Test", sections[0].Name())
}

func TestBuildSections_MultipleTypes(t *testing.T) {
	cfg := model.DashboardConfig{
		Sections: []model.DashboardSectionConfig{
			{Type: "ssh", Title: "SSH"},
			{Type: "docker", Title: "Docker"},
		},
	}
	sections := BuildSections(cfg, SectionDeps{})
	require.Len(t, sections, 2)
	assert.Equal(t, "SSH", sections[0].Name())
	assert.Equal(t, "Docker", sections[1].Name())
}

func TestBuildSections_EmptySectionsList(t *testing.T) {
	cfg := model.DashboardConfig{
		Sections: []model.DashboardSectionConfig{},
	}
	sections := BuildSections(cfg, SectionDeps{})
	require.Len(t, sections, 1)
	assert.Equal(t, "Sessions", sections[0].Name())
}

// --- Model: focus switching ---

func modelWithStubs(names ...string) Model {
	return Model{
		sections: stubs(names...),
		focused:  0,
		width:    80,
		height:   24,
	}
}

func TestTabMovesFocusForward(t *testing.T) {
	m := modelWithStubs("a", "b", "c")
	m = updateModel(m, pressKey("tab"))
	assert.Equal(t, 1, m.focused)
	m = updateModel(m, pressKey("tab"))
	assert.Equal(t, 2, m.focused)
}

func TestTabWrapsToFirst(t *testing.T) {
	m := modelWithStubs("a", "b", "c")
	m.focused = 2
	m = updateModel(m, pressKey("tab"))
	assert.Equal(t, 0, m.focused)
}

func TestShiftTabMovesFocusBackward(t *testing.T) {
	m := modelWithStubs("a", "b", "c")
	m.focused = 2
	m = updateModel(m, pressKey("shift+tab"))
	assert.Equal(t, 1, m.focused)
}

func TestShiftTabWrapsToLast(t *testing.T) {
	m := modelWithStubs("a", "b", "c")
	m = updateModel(m, pressKey("shift+tab"))
	assert.Equal(t, 2, m.focused)
}

func TestTabWithSingleSection(t *testing.T) {
	m := modelWithStubs("a")
	m = updateModel(m, pressKey("tab"))
	assert.Equal(t, 0, m.focused)
}

func TestTabWithNoSections(t *testing.T) {
	m := Model{width: 80, height: 24}
	m = updateModel(m, pressKey("tab"))
	assert.Equal(t, 0, m.focused)
}

// --- Model: quit ---

func TestQuitKeys(t *testing.T) {
	for _, k := range []string{"q", "esc", "ctrl+c"} {
		t.Run(k, func(t *testing.T) {
			m := modelWithStubs("a")
			m = updateModel(m, pressKey(k))
			assert.True(t, m.quit)
			assert.True(t, m.Quit())
		})
	}
}

// --- Model: select ---

func TestEnterSelectsAndQuits(t *testing.T) {
	s := &stubSection{name: "a", chosen: "my-session"}
	m := Model{
		sections: []Section{s},
		focused:  0,
		width:    80,
		height:   24,
	}
	result, cmd := m.Update(pressKey("enter"))
	rm := result.(Model)
	assert.Equal(t, "my-session", rm.Chosen())
	// The quit is signaled via tea.Quit command, not m.quit
	assert.NotNil(t, cmd)
}

func TestEnterWithoutChosenDoesNotQuit(t *testing.T) {
	s := &stubSection{name: "a", chosen: ""}
	m := Model{
		sections: []Section{s},
		focused:  0,
		width:    80,
		height:   24,
	}
	m = updateModel(m, pressKey("enter"))
	assert.False(t, m.quit)
}

// --- Model: key dispatch ---

func TestKeysForwardedToFocusedSection(t *testing.T) {
	a := &stubSection{name: "a"}
	b := &stubSection{name: "b"}
	m := Model{
		sections: []Section{a, b},
		focused:  1,
		width:    80,
		height:   24,
	}
	m = updateModel(m, pressKey("j"))
	assert.Equal(t, 0, a.updateCount)
	assert.Equal(t, 1, b.updateCount)
}

// --- Model: View ---

func TestViewPassesFocusedToSections(t *testing.T) {
	a := &stubSection{name: "a"}
	b := &stubSection{name: "b"}
	m := Model{
		sections: []Section{a, b},
		focused:  1,
		width:    80,
		height:   24,
		rows:     [][]int{{0, 1}},
	}
	m.View()
	assert.False(t, a.lastView.focused)
	assert.True(t, b.lastView.focused)
}

func TestViewQuitReturnsEmpty(t *testing.T) {
	m := modelWithStubs("a")
	m.quit = true
	v := m.View()
	assert.Equal(t, "", v.Content)
}

func TestViewTooSmallReturnsMessage(t *testing.T) {
	m := modelWithStubs("a")
	m.tooSmall = true
	v := m.View()
	assert.Contains(t, v.Content, "Terminal too small")
}

// --- WindowSizeMsg ---

func TestWindowSizeUpdatesDimensions(t *testing.T) {
	m := modelWithStubs("a", "b")
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	m = updateModel(m, msg)
	assert.Equal(t, 120, m.width)
	assert.Equal(t, 40, m.height)
	assert.False(t, m.tooSmall)
}

func TestWindowSizeTooSmall(t *testing.T) {
	m := modelWithStubs("a")
	msg := tea.WindowSizeMsg{Width: 10, Height: 3}
	m = updateModel(m, msg)
	assert.True(t, m.tooSmall)
}

// --- Width allocation ---

func TestWidthAllocation_FlexSections(t *testing.T) {
	a := &stubSection{name: "a", width: 0}
	b := &stubSection{name: "b", width: 0}
	m := Model{
		sections: []Section{a, b},
		rows:     [][]int{{0, 1}},
		width:    100,
		height:   24,
	}
	msg := tea.WindowSizeMsg{Width: 100, Height: 24}
	m = updateModel(m, msg)
	// 100 cols - 1 sep = 99, split evenly = 49, 50
	assert.Equal(t, 49, m.sectionWidths[0])
	assert.Equal(t, 50, m.sectionWidths[1])
}

func TestWidthAllocation_FixedAndFlex(t *testing.T) {
	a := &stubSection{name: "a", width: 0.3}
	b := &stubSection{name: "b", width: 0}
	m := Model{
		sections: []Section{a, b},
		rows:     [][]int{{0, 1}},
		width:    100,
		height:   24,
	}
	msg := tea.WindowSizeMsg{Width: 100, Height: 24}
	m = updateModel(m, msg)
	// fixed: 99 * 0.3 = 29, flex gets rest = 70
	assert.Equal(t, 29, m.sectionWidths[0])
	assert.Equal(t, 70, m.sectionWidths[1])
}

func TestWidthAllocation_AllFixed(t *testing.T) {
	a := &stubSection{name: "a", width: 0.5}
	b := &stubSection{name: "b", width: 0.5}
	m := Model{
		sections: []Section{a, b},
		rows:     [][]int{{0, 1}},
		width:    100,
		height:   24,
	}
	msg := tea.WindowSizeMsg{Width: 100, Height: 24}
	m = updateModel(m, msg)
	assert.Equal(t, 49, m.sectionWidths[0])
	assert.Equal(t, 50, m.sectionWidths[1])
}

// --- Row layout ---

func TestRowLayout_MultipleRows(t *testing.T) {
	a := &stubSection{name: "a", width: 0.5}
	b := &stubSection{name: "b", width: 0.5}
	c := &stubSection{name: "c"}
	m := Model{
		sections: []Section{a, b, c},
		rows:     [][]int{{0, 1}, {2}},
		width:    100,
		height:   24,
	}
	msg := tea.WindowSizeMsg{Width: 100, Height: 24}
	m = updateModel(m, msg)
	// Row 0: 99 * 0.5 = 49 each, remainder 1 goes to last → 49, 50
	assert.Equal(t, 49, m.sectionWidths[0])
	assert.Equal(t, 50, m.sectionWidths[1])
	// Row 1: full width
	assert.Equal(t, 100, m.sectionWidths[2])
}
