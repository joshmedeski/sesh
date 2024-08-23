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

	t.Run("name for git repo", func(t *testing.T) {
		mockGit.On("ShowTopLevel", "/Users/josh/c/dotfiles/.config/neovim").Return(true, "/Users/josh/c/dotfiles", nil)
		mockGit.On("GitCommonDir", "/Users/josh/c/dotfiles/.config/neovim").Return(true, "", nil)
		mockPathwrap.On("Base", "/Users/josh/c/dotfiles").Return("dotfiles")
		name, _ := n.Name("/Users/josh/c/dotfiles/.config/neovim")
		assert.Equal(t, "dotfiles/_config/neovim", name)
	})

	t.Run("name for git worktree", func(t *testing.T) {
		mockGit.On("ShowTopLevel", "/Users/josh/c/sesh/main").Return(true, "/Users/josh/c/sesh/main", nil)
		mockGit.On("GitCommonDir", "/Users/josh/c/sesh/main").Return(true, "/Users/josh/c/sesh/.bare", nil)
		mockPathwrap.On("Base", "/Users/josh/c/sesh").Return("sesh")
		name, _ := n.Name("/Users/josh/c/sesh/main")
		assert.Equal(t, "sesh/main", name)
	})

	t.Run("returns base on non-git dir", func(t *testing.T) {
		mockGit.On("ShowTopLevel", "/Users/josh/.config/neovim").Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
		mockGit.On("GitCommonDir", "/Users/josh/.config/neovim").Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
		mockPathwrap.On("Base", "/Users/josh/.config/neovim").Return("neovim")
		name, _ := n.Name("/Users/josh/.config/neovim")
		assert.Equal(t, "neovim", name)
	})
}
