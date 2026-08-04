package git

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorktreeAdd(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().
		CmdWithOutput("git", "-C", "/repo", "worktree", "add", "/repo/w/2345", "-b", "jam/2345-1", "--no-track", "origin/main").
		Return("", nil)
	g := NewGit(s)
	_, err := g.WorktreeAdd("/repo", "/repo/w/2345", "jam/2345-1", "origin/main")
	require.NoError(t, err)
}

func TestWorktreeAddDetached(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().
		CmdWithOutput("git", "-C", "/repo", "worktree", "add", "--detach", "/repo/w/2345", "origin/main").
		Return("", nil)
	g := NewGit(s)
	_, err := g.WorktreeAddDetached("/repo", "/repo/w/2345", "origin/main")
	require.NoError(t, err)
}

func TestFetchAndPull(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().CmdWithOutput("git", "-C", "/repo", "fetch").Return("", nil)
	s.EXPECT().CmdWithOutput("git", "-C", "/repo/w/2345", "pull", "--ff-only").Return("", nil)
	g := NewGit(s)
	_, err := g.Fetch("/repo")
	require.NoError(t, err)
	_, err = g.Pull("/repo/w/2345")
	assert.NoError(t, err)
}

func TestCurrentBranch(t *testing.T) {
	t.Run("returns the branch name on success", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		g := NewGit(mockShell)
		path := "/Users/josh/c/sesh"
		mockShell.On("Cmd", "git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD").
			Return("400", nil)

		ok, branch, err := g.CurrentBranch(path)

		assert.True(t, ok)
		assert.Equal(t, "400", branch)
		assert.NoError(t, err)
	})

	t.Run("returns false when not a git repo", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		g := NewGit(mockShell)
		path := "/tmp/not-a-repo"
		mockShell.On("Cmd", "git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD").
			Return("", fmt.Errorf("fatal: not a git repository"))

		ok, branch, err := g.CurrentBranch(path)

		assert.False(t, ok)
		assert.Equal(t, "", branch)
		assert.Error(t, err)
	})
}
