package github

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrViewFound(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().
		Cmd("gh", "pr", "view", "2345", "--repo", "nutiliti/nutiliti", "--json", "author,closingIssuesReferences,title,body").
		Return(`{"author":{"login":"octocat"},"closingIssuesReferences":[{"number":42}],"title":"Fix","body":"closes #7"}`, nil)
	g := NewGithub(s)
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
	g := NewGithub(s)
	_, found, err := g.PrView("nutiliti/nutiliti", 7)
	require.NoError(t, err)
	assert.False(t, found)
}

func TestCurrentUser(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().Cmd("gh", "api", "user", "--jq", ".login").Return("me", nil)
	g := NewGithub(s)
	user, err := g.CurrentUser()
	require.NoError(t, err)
	assert.Equal(t, "me", user)
}

func TestPrCheckout(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().
		CmdInDir("/repo/w/2345", "gh", "pr", "checkout", "2345", "--repo", "nutiliti/nutiliti").
		Return("", nil)
	g := NewGithub(s)
	_, err := g.PrCheckout("/repo/w/2345", "nutiliti/nutiliti", 2345)
	require.NoError(t, err)
}
