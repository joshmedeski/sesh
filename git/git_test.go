package git

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/shell"
	"github.com/stretchr/testify/assert"
)

func TestGitMainWorktree(t *testing.T) {
	t.Run("run should find worktree root", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockShell.On("Cmd", "git", "-C", "~/.dotfiles/nvim", "worktree", "list").Return(`
/Users/hansolo/.dotfiles             (bare)
/Users/hansolo/.dotfiles/nvim        ba04ca494 [5.x]
`, nil)
		git := &RealGit{shell: mockShell}
		isWorktree, out, err := git.GitMainWorktree("~/.dotfiles/nvim")
		assert.True(t, isWorktree)
		assert.Nil(t, err)
		assert.Equal(t, "/Users/hansolo/.dotfiles", out)
	})

	t.Run("run should find non-worktree root", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockShell.On("Cmd", "git", "-C", "~/.dotfiles/nvim", "worktree", "list").Return(`
/Users/hansolo/.dotfiles/nvim        ba04ca494 [5.x]
`, nil)
		git := &RealGit{shell: mockShell}
		isWorktree, out, err := git.GitMainWorktree("~/.dotfiles/nvim")
		assert.True(t, isWorktree)
		assert.Nil(t, err)
		assert.Equal(t, "/Users/hansolo/.dotfiles/nvim", out)
	})

	t.Run("run should fail when not in git repo", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockShell.On("Cmd", "git", "-C", "~/.dotfiles/nvim", "worktree", "list").Return("", errors.New(`
fatal: not a git repository (or any of the parent directories): .git
`))
		git := &RealGit{shell: mockShell}
		isWorktree, out, _ := git.GitMainWorktree("~/.dotfiles/nvim")
		assert.False(t, isWorktree)
		assert.Equal(t, "", out)
	})
}
