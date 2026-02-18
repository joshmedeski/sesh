package picker

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"

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

// newTestModel creates a model and simulates the async load completing.
func newTestModel() Model {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), false)
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
	m := New(testFetchFunc(sessions), false)
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
	m := New(testFetchFunc(sessions), false)
	assert.True(t, m.loading)

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	resultModel := result.(Model)

	assert.Equal(t, "", resultModel.Chosen(), "enter while loading should not select anything")
	assert.Nil(t, cmd, "enter while loading should not quit")
	assert.True(t, resultModel.loading, "should still be loading")
}

func TestUpdate_Escape_WhileLoading(t *testing.T) {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), false)

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	resultModel := result.(Model)

	assert.True(t, resultModel.Quit(), "escape while loading should quit")
}

func TestUpdate_SessionsLoaded(t *testing.T) {
	sessions := testSessions()
	m := New(testFetchFunc(sessions), false)
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
	m := New(testFetchFunc(sessions), false)

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
	}, false)

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
	m := New(testFetchFunc(sessions), false)
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
	m := New(testFetchFunc(sessions), false)
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	m = result.(Model)
	m.height = 20

	half := m.visibleCount() / 2
	result, _ = m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	resultModel := result.(Model)
	assert.Equal(t, half, resultModel.cursor)
}
