package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// fakeShell stubs shell.Shell for RealGit tests.
type fakeShell struct {
	out string
	err error
}

func (f *fakeShell) Cmd(cmd string, args ...string) (string, error) {
	return f.out, f.err
}

func (f *fakeShell) CmdWithOutput(cmd string, args ...string) (string, error) {
	return f.out, f.err
}

func (f *fakeShell) ListCmd(cmd string, args ...string) ([]string, error) {
	return nil, nil
}

func (f *fakeShell) PrepareCmd(cmd string, replacements map[string]string) ([]string, error) {
	return nil, nil
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
