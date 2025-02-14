package namer

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func TestFromPath(t *testing.T) {
	t.Run("when path does not contain a symlink", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		n := NewNamer(mockPathwrap, mockGit, mockHome)

		t.Run("name for git repo", func(t *testing.T) {
			mockPathwrap.On("EvalSymlinks", "/Users/josh/config/dotfiles/.config/neovim").Return("/Users/josh/config/dotfiles/.config/neovim", nil)
			mockGit.On("ShowTopLevel", "/Users/josh/config/dotfiles/.config/neovim").Return(true, "/Users/josh/config/dotfiles", nil)
			mockGit.On("GitCommonDir", "/Users/josh/config/dotfiles/.config/neovim").Return(true, "", nil)
			mockPathwrap.On("Base", "/Users/josh/config/dotfiles").Return("dotfiles")
			name, _ := n.Name("/Users/josh/config/dotfiles/.config/neovim")
			assert.Equal(t, "dotfiles/_config/neovim", name)
		})

		t.Run("name for git worktree", func(t *testing.T) {
			mockPathwrap.On("EvalSymlinks", "/Users/josh/config/sesh/main").Return("/Users/josh/config/sesh/main", nil)
			mockGit.On("ShowTopLevel", "/Users/josh/config/sesh/main").Return(true, "/Users/josh/config/sesh/main", nil)
			mockGit.On("GitCommonDir", "/Users/josh/config/sesh/main").Return(true, "/Users/josh/config/sesh/.bare", nil)
			mockPathwrap.On("Base", "/Users/josh/config/sesh").Return("sesh")
			name, _ := n.Name("/Users/josh/config/sesh/main")
			assert.Equal(t, "sesh/main", name)
		})

		t.Run("returns base on non-git dir", func(t *testing.T) {
			mockPathwrap.On("EvalSymlinks", "/Users/josh/.config/neovim").Return("/Users/josh/.config/neovim", nil)
			mockGit.On("ShowTopLevel", "/Users/josh/.config/neovim").Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
			mockGit.On("GitCommonDir", "/Users/josh/.config/neovim").Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
			mockPathwrap.On("Base", "/Users/josh/.config/neovim").Return("neovim")
			name, _ := n.Name("/Users/josh/.config/neovim")
			assert.Equal(t, "neovim", name)
		})
	})

	t.Run("when path contains a symlink", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		n := NewNamer(mockPathwrap, mockGit, mockHome)

		t.Run("name for symlinked file in symlinked git repo", func(t *testing.T) {
			mockPathwrap.On("EvalSymlinks", "/Users/josh/d/.c/neovim").Return("/Users/josh/dotfiles/.config/neovim", nil)
			mockGit.On("ShowTopLevel", "/Users/josh/dotfiles/.config/neovim").Return(true, "/Users/josh/dotfiles", nil)
			mockGit.On("GitCommonDir", "/Users/josh/dotfiles/.config/neovim").Return(true, "", nil)
			mockPathwrap.On("Base", "/Users/josh/dotfiles").Return("dotfiles")
			name, _ := n.Name("/Users/josh/d/.c/neovim")
			assert.Equal(t, "dotfiles/_config/neovim", name)
		})

		t.Run("name for git worktree", func(t *testing.T) {
			mockPathwrap.On("EvalSymlinks", "/Users/josh/p/sesh/main").Return("/Users/josh/projects/sesh/main", nil)
			mockGit.On("ShowTopLevel", "/Users/josh/projects/sesh/main").Return(true, "/Users/josh/projects/sesh/main", nil)
			mockGit.On("GitCommonDir", "/Users/josh/projects/sesh/main").Return(true, "/Users/josh/projects/sesh/.bare", nil)
			mockPathwrap.On("Base", "/Users/josh/projects/sesh").Return("sesh")
			name, _ := n.Name("/Users/josh/p/sesh/main")
			assert.Equal(t, "sesh/main", name)
		})

		t.Run("returns base on non-git dir", func(t *testing.T) {
			mockPathwrap.On("EvalSymlinks", "/Users/josh/c/neovim").Return("/Users/josh/.config/neovim", nil)
			mockGit.On("ShowTopLevel", "/Users/josh/.config/neovim").Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
			mockGit.On("GitCommonDir", "/Users/josh/.config/neovim").Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
			mockPathwrap.On("Base", "/Users/josh/.config/neovim").Return("neovim")
			name, _ := n.Name("/Users/josh/c/neovim")
			assert.Equal(t, "neovim", name)
		})
	})
}
