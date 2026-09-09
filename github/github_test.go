package github

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrViewFound(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().
		Cmd("gh", "pr", "view", "2345", "--repo", "nutiliti/nutiliti", "--json", "author,closingIssuesReferences,title,body").
		Return(`{"author":{"login":"octocat"},"closingIssuesReferences":[{"number":42}],"title":"Fix","body":"closes #7"}`, nil)
	g := NewGithub(s, git.NewMockGit(t))
	pr, found, err := g.PrView("nutiliti/nutiliti", 2345)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "octocat", pr.Author)
	assert.Equal(t, []int{42}, pr.ClosingIssues)
}

func TestPrViewNotAPr(t *testing.T) {
	s := shell.NewMockShell(t)
	// gh returns empty (swallowed error) when the number is an issue, not a PR
	s.EXPECT().
		Cmd("gh", "pr", "view", "7", "--repo", "nutiliti/nutiliti", "--json", "author,closingIssuesReferences,title,body").
		Return("", nil)
	g := NewGithub(s, git.NewMockGit(t))
	_, found, err := g.PrView("nutiliti/nutiliti", 7)
	require.NoError(t, err)
	assert.False(t, found)
}

func TestCurrentUser(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().Cmd("gh", "api", "user", "--jq", ".login").Return("me", nil)
	g := NewGithub(s, git.NewMockGit(t))
	user, err := g.CurrentUser()
	require.NoError(t, err)
	assert.Equal(t, "me", user)
}

func TestPrCheckout(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().
		CmdInDir("/repo/w/2345", "gh", "pr", "checkout", "2345", "--repo", "nutiliti/nutiliti").
		Return("", nil)
	g := NewGithub(s, git.NewMockGit(t))
	_, err := g.PrCheckout("/repo/w/2345", "nutiliti/nutiliti", 2345)
	require.NoError(t, err)
}

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

func TestIssue(t *testing.T) {
	t.Run("returns the issue on success", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "feat/400-status-bar", nil)
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
		mockShell.On("Cmd", "gh", "issue", "view", "400", "--json", "number,title,state").
			Return("not json", nil)

		_, found, err := gh.Issue(path)

		assert.False(t, found)
		assert.NoError(t, err)
	})
}
