package git

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeShell stubs shell.Shell for RealGit tests.
type fakeShell struct {
	out string
	err error
}

func (f *fakeShell) Cmd(cmd string, args ...string) (string, error) {
	return f.out, f.err
}

func (f *fakeShell) CmdInDir(dir string, cmd string, args ...string) (string, error) {
	return f.out, f.err
}

func (f *fakeShell) CmdWithOutput(cmd string, args ...string) (string, error) {
	return f.out, f.err
}

func (f *fakeShell) CmdCapture(cmd string, args ...string) (string, error) {
	return f.out, f.err
}

func (f *fakeShell) ListCmd(cmd string, args ...string) ([]string, error) {
	return nil, nil
}

func (f *fakeShell) PrepareCmd(cmd string, replacements map[string]string) ([]string, error) {
	return nil, nil
}

func (f *fakeShell) ShellCmd(cmd string, replacements map[string]string) (string, error) {
	return f.out, f.err
}

func statusSummary(t *testing.T, porcelain string) StatusSummary {
	t.Helper()
	g := NewGit(&fakeShell{out: porcelain})
	s, err := g.StatusSummary("/some/path")
	assert.NoError(t, err)
	return s
}

func TestStatusSummary_CleanRepo(t *testing.T) {
	assert.Equal(t, StatusSummary{}, statusSummary(t, ""))
	assert.Equal(t, StatusSummary{}, statusSummary(t, "\n"))
}

// Regression: the first line of porcelain output must not be misclassified.
// TrimSpace previously stripped the leading space of " M file", turning an
// unstaged modification into a staged one.
func TestStatusSummary_FirstLineUnstaged(t *testing.T) {
	s := statusSummary(t, " M file.txt\n")
	assert.Equal(t, StatusSummary{Unstaged: 1}, s)
}

func TestStatusSummary_OnlyUnstaged(t *testing.T) {
	s := statusSummary(t, " M a.go\n M b.go\n M c.go\n")
	assert.Equal(t, StatusSummary{Unstaged: 3}, s)
}

func TestStatusSummary_OnlyStaged(t *testing.T) {
	s := statusSummary(t, "M  a.go\nA  b.go\n")
	assert.Equal(t, StatusSummary{Staged: 2}, s)
}

func TestStatusSummary_Mixed(t *testing.T) {
	s := statusSummary(t, " M a.go\nM  b.go\nMM c.go\n?? d.go\n D e.go\n")
	assert.Equal(t, StatusSummary{Staged: 2, Unstaged: 2, Deleted: 1, Untracked: 1}, s)
}

func TestStatusSummary_UntrackedFirstLine(t *testing.T) {
	s := statusSummary(t, "?? new/\n M a.go\n")
	assert.Equal(t, StatusSummary{Unstaged: 1, Untracked: 1}, s)
}

func TestStatusSummary_Deleted(t *testing.T) {
	s := statusSummary(t, " D a.go\nD  b.go\n")
	assert.Equal(t, StatusSummary{Staged: 1, Deleted: 2}, s)
}

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
