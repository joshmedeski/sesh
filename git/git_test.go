package git

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestGitRoot(t *testing.T) {
	t.Run("run should find worktree root", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockShell.On("Cmd", "git", "-C", "~/code/project/sesh/main", "worktree", "list").Return(`
/Users/hansolo/code/project/sesh             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`, nil)
		git := &RealGit{shell: mockShell}
		isGit, out, err := git.GitRoot("~/code/project/sesh/main")
		assert.True(t, isGit)
		assert.Nil(t, err)
		assert.Equal(t, "/Users/hansolo/code/project/sesh", out)
	})

	t.Run("run should find non-worktree root", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockShell.On("Cmd", "git", "-C", "~/.dotfiles/nvim", "worktree", "list").Return(`
/Users/hansolo/.dotfiles        ba04ca494 [5.x]
`, nil)
		git := &RealGit{shell: mockShell}
		isGit, out, err := git.GitRoot("~/.dotfiles/nvim")
		assert.True(t, isGit)
		assert.Nil(t, err)
		assert.Equal(t, "/Users/hansolo/.dotfiles", out)
	})

	t.Run("run should fail when not in git repo", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockShell.On("Cmd", "git", "-C", "~/not-a-repo", "worktree", "list").Return("", errors.New(`
fatal: not a git repository (or any of the parent directories): .git
`))
		git := &RealGit{shell: mockShell}
		isGit, out, _ := git.GitRoot("~/not-a-repo")
		assert.False(t, isGit)
		assert.Equal(t, "", out)
	})
}
