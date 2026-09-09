package picker

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/zoxide"
)

// ctrlX is the removal binding as the terminal delivers it.
var ctrlX = tea.KeyPressMsg{Code: 'x', Mod: tea.ModCtrl}

// newRemovableModel loads the test sessions with a removal func that records
// the paths it was asked to remove, and returns that record alongside the
// model.
func newRemovableModel(err error) (Model, *[]string) {
	removed := &[]string{}
	sessions := testSessions()
	opts := testOptionsWith(func(o *Options) {
		o.Remove = func(path string) error {
			*removed = append(*removed, path)
			return err
		}
	})
	m := New(testFetchFunc(sessions), opts)
	result, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	return result.(Model), removed
}

// cursorOnZoxide moves the cursor to the zoxide row of testSessions.
func cursorOnZoxide(t *testing.T, m Model) Model {
	t.Helper()
	for i, item := range m.filtered {
		if item.item.src == zoxideSrc {
			m.cursor = i
			return m
		}
	}
	t.Fatal("no zoxide row in the test sessions")
	return m
}

func press(m Model, msg tea.Msg) Model {
	result, _ := m.Update(msg)
	return result.(Model)
}

func TestCtrlX_OpensDialogOnZoxideRow(t *testing.T) {
	m, _ := newRemovableModel(nil)
	m = cursorOnZoxide(t, m)

	m = press(m, ctrlX)

	require.NotNil(t, m.confirm, "ctrl+x on a zoxide row should open the dialog")
	assert.Equal(t, "~/code/app", m.confirm.name)
	assert.Equal(t, "/home/user/code/app", m.confirm.path)
	assert.True(t, m.confirm.yes, "the dialog should open with Yes focused")
}

func TestCtrlX_InertOnNonZoxideRows(t *testing.T) {
	for _, src := range []string{"tmux", "config", "tmuxinator"} {
		t.Run(src, func(t *testing.T) {
			m, _ := newRemovableModel(nil)
			for i, item := range m.filtered {
				if item.item.src == src {
					m.cursor = i
					break
				}
			}
			m = press(m, ctrlX)

			assert.Nil(t, m.confirm, "ctrl+x should not open the dialog on a %s row", src)
			assert.Contains(t, m.status, "Only zoxide entries")
		})
	}
}

func TestCtrlX_InertWhileLoading(t *testing.T) {
	opts := testOptionsWith(func(o *Options) {
		o.Remove = func(string) error { return nil }
	})
	m := New(testFetchFunc(testSessions()), opts)
	require.True(t, m.loading)

	m = press(m, ctrlX)

	assert.Nil(t, m.confirm)
	assert.Empty(t, m.status)
}

func TestCtrlX_InertWithoutRemoveFunc(t *testing.T) {
	m := newTestModel()
	m = cursorOnZoxide(t, m)

	m = press(m, ctrlX)

	assert.Nil(t, m.confirm)
}

func TestCtrlX_DoesNotTypeIntoTheFilter(t *testing.T) {
	m, _ := newRemovableModel(nil)
	m = cursorOnZoxide(t, m)

	m = press(m, ctrlX)

	assert.Equal(t, "", m.filterInput.Value())
}

func TestConfirm_SwallowsKeysWhileOpen(t *testing.T) {
	m, _ := newRemovableModel(nil)
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	m = press(m, tea.KeyPressMsg{Code: 'a'})

	assert.Equal(t, "", m.filterInput.Value(), "the dialog should swallow typing")
	assert.NotNil(t, m.confirm, "an unrelated key should leave the dialog open")
}

func TestConfirm_YesRemovesTheEntry(t *testing.T) {
	m, removed := newRemovableModel(nil)
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	result, cmd := m.Update(tea.KeyPressMsg{Code: 'y'})
	m = result.(Model)
	assert.Nil(t, m.confirm, "confirming should close the dialog")
	require.NotNil(t, cmd, "confirming should start the removal")

	msg := cmd()
	assert.Equal(t, []string{"/home/user/code/app"}, *removed)

	m = press(m, msg)
	assert.Len(t, m.allItems, 4)
	for _, item := range m.allItems {
		assert.NotEqual(t, "~/code/app", item.name)
	}
}

func TestConfirm_EnterConfirmsByDefault(t *testing.T) {
	m, removed := newRemovableModel(nil)
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	require.NotNil(t, cmd)
	cmd()

	assert.Equal(t, []string{"/home/user/code/app"}, *removed)
}

func TestConfirm_EnterOnNoCancels(t *testing.T) {
	m, removed := newRemovableModel(nil)
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	m = press(m, tea.KeyPressMsg{Code: tea.KeyRight})
	require.NotNil(t, m.confirm)
	assert.False(t, m.confirm.yes, "right should move focus to No")

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(Model)
	assert.Nil(t, m.confirm)
	assert.Nil(t, cmd)
	assert.Empty(t, *removed)
	assert.Len(t, m.allItems, 5)
}

func TestConfirm_CancelKeys(t *testing.T) {
	cancels := map[string]tea.KeyPressMsg{
		"n":   {Code: 'n'},
		"q":   {Code: 'q'},
		"esc": {Code: tea.KeyEscape},
	}
	for name, key := range cancels {
		t.Run(name, func(t *testing.T) {
			m, removed := newRemovableModel(nil)
			m = cursorOnZoxide(t, m)
			m = press(m, ctrlX)

			result, cmd := m.Update(key)
			m = result.(Model)

			assert.Nil(t, m.confirm, "%s should close the dialog", name)
			assert.Nil(t, cmd, "%s should not start a removal", name)
			assert.Empty(t, *removed)
			assert.Len(t, m.allItems, 5)
		})
	}
}

func TestConfirm_EscDoesNotQuitThePicker(t *testing.T) {
	m, _ := newRemovableModel(nil)
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	m = press(m, tea.KeyPressMsg{Code: tea.KeyEscape})

	assert.False(t, m.quit, "esc should close the dialog, not the picker")
}

func TestConfirm_TabMovesBetweenButtons(t *testing.T) {
	m, _ := newRemovableModel(nil)
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	require.NotNil(t, m.confirm)
	assert.False(t, m.confirm.yes)

	m = press(m, tea.KeyPressMsg{Code: tea.KeyTab})
	require.NotNil(t, m.confirm)
	assert.True(t, m.confirm.yes)
}

func TestRemoval_FailureKeepsTheRow(t *testing.T) {
	m, _ := newRemovableModel(errors.New("exit status 1"))
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	result, cmd := m.Update(tea.KeyPressMsg{Code: 'y'})
	m = press(result.(Model), cmd())

	assert.Len(t, m.allItems, 5, "a failed removal must not drop the row")
	assert.Contains(t, m.status, "Couldn't remove entry")
	assert.Contains(t, m.status, "exit status 1")
}

func TestRemoval_PreservesTheFilter(t *testing.T) {
	m, _ := newRemovableModel(nil)
	m.filterInput.SetValue("code/app")
	m.applyFilter()
	require.Len(t, m.filtered, 1)

	m = press(m, ctrlX)
	require.NotNil(t, m.confirm)

	result, cmd := m.Update(tea.KeyPressMsg{Code: 'y'})
	m = press(result.(Model), cmd())

	assert.Equal(t, "code/app", m.filterInput.Value())
	assert.Empty(t, m.filtered)
	assert.Equal(t, 0, m.cursor)
}

func TestRemoval_PullsTheCursorBackFromTheLastRow(t *testing.T) {
	sessions := model.SeshSessions{
		OrderedIndex: []string{"s1", "s2"},
		Directory: model.SeshSessionMap{
			"s1": {Name: "one", Src: "tmux", Path: "/one"},
			"s2": {Name: "~/two", Src: zoxideSrc, Path: "/two"},
		},
	}
	opts := testOptionsWith(func(o *Options) {
		o.Remove = func(string) error { return nil }
	})
	result, _ := New(testFetchFunc(sessions), opts).Update(sessionsLoadedMsg{sessions: sessions})
	m := result.(Model)
	m.cursor = 1

	m = press(m, ctrlX)
	confirmed, cmd := m.Update(tea.KeyPressMsg{Code: 'y'})
	m = press(confirmed.(Model), cmd())

	assert.Len(t, m.filtered, 1)
	assert.Equal(t, 0, m.cursor, "the cursor should be pulled back onto the new last row")
}

func TestStatus_ClearedByTheNextKeypress(t *testing.T) {
	m, _ := newRemovableModel(nil)
	m.cursor = 0 // a tmux row
	m = press(m, ctrlX)
	require.NotEmpty(t, m.status)

	m = press(m, tea.KeyPressMsg{Code: tea.KeyDown})

	assert.Empty(t, m.status)
}

func TestConfirmView_ShowsThePromptAndButtons(t *testing.T) {
	m, _ := newRemovableModel(nil)
	m = press(m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	out := ansi.Strip(m.View().Content)

	assert.Contains(t, out, "Do you want to remove this directory from zoxide?")
	assert.Contains(t, out, "~/code/app")
	assert.Contains(t, out, "Yes")
	assert.Contains(t, out, "No")
}

func TestConfirmView_KeepsTheListVisibleBehindIt(t *testing.T) {
	m, _ := newRemovableModel(nil)
	m = press(m, tea.WindowSizeMsg{Width: 120, Height: 24})
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	out := ansi.Strip(m.View().Content)

	assert.Contains(t, out, "my-project", "the dialog overlays the list rather than replacing it")
}

func TestTruncate(t *testing.T) {
	assert.Equal(t, "short", truncate("short", 10))
	assert.Equal(t, "abcd…", truncate("abcdefgh", 5))
}

func TestConfirm_CtrlCQuitsThePicker(t *testing.T) {
	m, removed := newRemovableModel(nil)
	m = cursorOnZoxide(t, m)
	m = press(m, ctrlX)

	m = press(m, tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})

	assert.True(t, m.quit, "ctrl+c should quit the picker, dialog or not")
	assert.Nil(t, m.confirm)
	assert.Empty(t, *removed)
}

func TestRemoval_DropsTheZoxideRowNotItsNamesake(t *testing.T) {
	sessions := model.SeshSessions{
		OrderedIndex: []string{"t", "z"},
		Directory: model.SeshSessionMap{
			"t": {Name: "app", Src: "tmux", Path: "/live/app"},
			"z": {Name: "app", Src: zoxideSrc, Path: "/stale/app"},
		},
	}
	opts := testOptionsWith(func(o *Options) {
		o.Remove = func(string) error { return nil }
	})
	result, _ := New(testFetchFunc(sessions), opts).Update(sessionsLoadedMsg{sessions: sessions})
	m := cursorOnZoxide(t, result.(Model))

	m = press(m, ctrlX)
	confirmed, cmd := m.Update(tea.KeyPressMsg{Code: 'y'})
	m = press(confirmed.(Model), cmd())

	require.Len(t, m.allItems, 1)
	assert.Equal(t, "tmux", m.allItems[0].src, "the live session must survive its zoxide namesake")
}

func TestRealPicker_RemoveEntry(t *testing.T) {
	t.Run("refreshes the cache after a successful removal", func(t *testing.T) {
		mockZoxide := new(zoxide.MockZoxide)
		mockZoxide.EXPECT().Remove("/stale/app").Return(nil)
		refreshed := 0
		p := &RealPicker{zoxide: mockZoxide, refreshCache: func() { refreshed++ }}

		assert.NoError(t, p.removeEntry("/stale/app"))
		assert.Equal(t, 1, refreshed)
	})

	t.Run("leaves the cache alone when the removal failed", func(t *testing.T) {
		mockZoxide := new(zoxide.MockZoxide)
		mockZoxide.EXPECT().Remove("/stale/app").Return(errors.New("exit status 1"))
		refreshed := 0
		p := &RealPicker{zoxide: mockZoxide, refreshCache: func() { refreshed++ }}

		assert.Error(t, p.removeEntry("/stale/app"))
		assert.Zero(t, refreshed)
	})

	t.Run("removes without a cache to refresh", func(t *testing.T) {
		mockZoxide := new(zoxide.MockZoxide)
		mockZoxide.EXPECT().Remove("/stale/app").Return(nil)
		p := &RealPicker{zoxide: mockZoxide}

		assert.NoError(t, p.removeEntry("/stale/app"))
	})
}
