package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListWindows(t *testing.T) {
	t.Run("returns parsed windows", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell}
		mockShell.EXPECT().ListCmd("tmux", "list-windows", "-F", mock.Anything).Return(
			[]string{"0::editor::/Users/josh/c/sesh::0", "1::server::/Users/josh/c/sesh::1"},
			nil,
		)
		windows, err := tmux.ListWindows("")
		assert.Nil(t, err)
		assert.Len(t, windows, 2)
		assert.Equal(t, "editor", windows[0].Name)
		assert.Equal(t, "/Users/josh/c/sesh", windows[0].Path)
		assert.Equal(t, 0, windows[0].Index)
		assert.False(t, windows[0].Active)
		assert.Equal(t, "server", windows[1].Name)
		assert.True(t, windows[1].Active)
	})

	t.Run("target session flag is passed when non-empty", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell}
		mockShell.EXPECT().ListCmd("tmux", "list-windows", "-t", "work", "-F", mock.Anything).Return(
			[]string{"0::main::/home/user::0"},
			nil,
		)
		windows, err := tmux.ListWindows("work")
		assert.Nil(t, err)
		assert.Len(t, windows, 1)
	})

	t.Run("parseTmuxWindowsOutput", func(t *testing.T) {
		raw := []string{"0::editor::/Users/josh/c/sesh::0", "1::server::/Users/josh/c/sesh::1"}
		windows, err := parseTmuxWindowsOutput(raw)
		assert.Nil(t, err)
		assert.Len(t, windows, 2)
		assert.Equal(t, "editor", windows[0].Name)
		assert.Equal(t, 0, windows[0].Index)
		assert.False(t, windows[0].Active)
		assert.Equal(t, "server", windows[1].Name)
		assert.Equal(t, 1, windows[1].Index)
		assert.True(t, windows[1].Active)
	})
}

func TestListAllWindowNames(t *testing.T) {
	t.Run("groups window names by session", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().ListCmd("tmux", "list-windows", "-a", "-F", mock.Anything).Return(
			[]string{"sesh::editor", "sesh::server", "dotfiles::nvim"},
			nil,
		)
		windowNames, err := tmux.ListAllWindowNames()
		assert.Nil(t, err)
		assert.Equal(t, map[string][]string{
			"sesh":     {"editor", "server"},
			"dotfiles": {"nvim"},
		}, windowNames)
	})

	t.Run("returns error from shell", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().ListCmd("tmux", "list-windows", "-a", "-F", mock.Anything).Return(
			nil, assert.AnError,
		)
		windowNames, err := tmux.ListAllWindowNames()
		assert.Error(t, err)
		assert.Nil(t, windowNames)
	})

	t.Run("parseAllWindowNamesOutput skips malformed lines", func(t *testing.T) {
		raw := []string{"sesh::editor", "no-separator", "::orphan", "sesh::", "work::make::build"}
		windowNames := parseAllWindowNamesOutput(raw)
		assert.Equal(t, map[string][]string{
			"sesh": {"editor"},
			"work": {"make::build"},
		}, windowNames)
	})
}
