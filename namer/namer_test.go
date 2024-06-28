package namer

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/pathwrap"
	"github.com/stretchr/testify/assert"
)

func TestFromPath(t *testing.T) {
	mockPathwrap := new(pathwrap.MockPath)
	mockGit := new(git.MockGit)
	n := NewNamer(mockPathwrap, mockGit)

	t.Run("returns base on git dir", func(t *testing.T) {
		mockGit.On("ShowTopLevel", "/Users/josh/c/dotfiles/.config/neovim").Return(true, "/Users/josh/c/dotfiles", nil)
		mockPathwrap.On("Base", "/Users/josh/c/dotfiles").Return("dotfiles")

		name, _ := n.FromPath("/Users/josh/c/dotfiles/.config/neovim")
		assert.Equal(t, "dotfiles/.config/neovim", name)
	})

	t.Run("returns base on git worktree dir", func(t *testing.T) {
		mockGit.On("ShowTopLevel", "/Users/josh/c/sesh/main/namer").Return(true, "/Users/josh/c/sesh", nil)
		mockPathwrap.On("Base", "/Users/josh/c/sesh").Return("sesh")

		name, _ := n.FromPath("/Users/josh/c/sesh/main/namer")
		assert.Equal(t, "sesh/main/namer", name)
	})

	t.Run("returns base on non-git dir", func(t *testing.T) {
		mockGit.On("ShowTopLevel", "/Users/josh/.config/neovim").Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
		mockPathwrap.On("Base", "/Users/josh/.config/neovim").Return("neovim")
		name, _ := n.FromPath("/Users/josh/.config/neovim")
		assert.Equal(t, "neovim", name)
	})
}
