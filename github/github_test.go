package github

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestParseIssueNumber(t *testing.T) {
	cases := []struct {
		branch string
		want   string
		ok     bool
	}{
		{"400", "400", true},
		{"feat/400-status-bar", "400", true},
		{"bugfix/issue-12", "12", true},
		{"feat/400-then-401", "400", true},
		{"main", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		got, ok := parseIssueNumber(c.branch)
		assert.Equal(t, c.ok, ok, c.branch)
		assert.Equal(t, c.want, got, c.branch)
	}
}

func TestResolve(t *testing.T) {
	t.Run("numeric branch resolves repo, branch, and number", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "feat/400-status-bar", nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/c/sesh", nil)

		ref, ok := gh.Resolve(path)

		assert.True(t, ok)
		assert.Equal(t, BranchRef{RepoRoot: "/Users/josh/c/sesh", Branch: "feat/400-status-bar", Number: 400, HasNumber: true}, ref)
	})

	t.Run("non-numeric branch resolves with HasNumber false", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "main", nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/c/sesh", nil)

		ref, ok := gh.Resolve(path)

		assert.True(t, ok)
		assert.Equal(t, "main", ref.Branch)
		assert.False(t, ref.HasNumber)
		assert.Equal(t, 0, ref.Number)
	})

	t.Run("not a git repo => ok false", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/tmp/x"
		mockGit.On("CurrentBranch", path).Return(false, "", fmt.Errorf("not a git repo"))

		_, ok := gh.Resolve(path)
		assert.False(t, ok)
	})
}

func TestIssue(t *testing.T) {
	t.Run("returns the issue on success", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "feat/400-status-bar", nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/c/sesh", nil)
		mockShell.On("Cmd", "gh", "issue", "view", "400", "--json", "number,title,state").
			Return(`{"number":400,"state":"OPEN","title":"Dynamic tmux status bar"}`, nil)

		issue, found, err := gh.Issue(path)

		assert.True(t, found)
		assert.NoError(t, err)
		assert.Equal(t, Issue{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"}, issue)
	})

	t.Run("not found when branch has no number", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "main", nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/c/sesh", nil)

		issue, found, err := gh.Issue(path)

		assert.False(t, found)
		assert.NoError(t, err)
		assert.Equal(t, Issue{}, issue)
		// No gh call is set up on mockShell; testify panics on an
		// unexpected call, so reaching gh would fail this test.
	})

	t.Run("not found when not a git repo", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/tmp/x"
		mockGit.On("CurrentBranch", path).Return(false, "", fmt.Errorf("not a git repo"))

		_, found, err := gh.Issue(path)

		assert.False(t, found)
		assert.NoError(t, err)
	})

	t.Run("not found when gh errors", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "400", nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/c/sesh", nil)
		mockShell.On("Cmd", "gh", "issue", "view", "400", "--json", "number,title,state").
			Return("", fmt.Errorf("gh: not found"))

		_, found, err := gh.Issue(path)

		assert.False(t, found)
		assert.NoError(t, err)
	})

	t.Run("not found when gh returns malformed json", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "400", nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/c/sesh", nil)
		mockShell.On("Cmd", "gh", "issue", "view", "400", "--json", "number,title,state").
			Return("not json", nil)

		_, found, err := gh.Issue(path)

		assert.False(t, found)
		assert.NoError(t, err)
	})
}
