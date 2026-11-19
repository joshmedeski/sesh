package gitlab

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestIssues_EmptySlice(t *testing.T) {
	mockShell := shell.NewMockShell(t)
	g := NewGitlab(mockShell, git.NewMockGit(t))

	found, missing, err := g.Issues("gitlab.com", "group/repo", []int{})
	require.NoError(t, err)
	assert.Empty(t, found)
	assert.Empty(t, missing)
}

func TestIssues_SubprocessError(t *testing.T) {
	mockShell := shell.NewMockShell(t)
	g := NewGitlab(mockShell, git.NewMockGit(t))
	mockShell.EXPECT().
		Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
			"--raw-field", mock.Anything,
			"--raw-field", "fullPath=group/repo").
		Return("", errors.New("network error"))

	_, _, err := g.Issues("gitlab.com", "group/repo", []int{42})
	require.Error(t, err)
}

func TestIssues_UnmarshalError(t *testing.T) {
	mockShell := shell.NewMockShell(t)
	g := NewGitlab(mockShell, git.NewMockGit(t))
	mockShell.EXPECT().
		Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
			"--raw-field", mock.Anything,
			"--raw-field", "fullPath=group/repo").
		Return("not json", nil)

	_, _, err := g.Issues("gitlab.com", "group/repo", []int{42})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decoding glab graphql response")
}

func TestIssues_ReportsMissingNumbers(t *testing.T) {
	mockShell := shell.NewMockShell(t)
	g := NewGitlab(mockShell, git.NewMockGit(t))
	mockShell.EXPECT().
		Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
			"--raw-field", queryField(`issue(iid: "42")`),
			"--raw-field", "fullPath=group/repo").
		Return(`{"data":{"project":{"i42":null}}}`, nil)

	found, missing, err := g.Issues("gitlab.com", "group/repo", []int{42})
	require.NoError(t, err)
	assert.Empty(t, found)
	assert.Equal(t, []int{42}, missing)
}

func TestIssues_BatchesOverOneHundred(t *testing.T) {
	mockShell := shell.NewMockShell(t)
	g := NewGitlab(mockShell, git.NewMockGit(t))

	// Create 101 distinct numbers to trigger chunking across issueBatchSize (100)
	numbers := make([]int, 101)
	for i := 0; i < 101; i++ {
		numbers[i] = i + 1
	}

	// First batch: 100 items (1..100)
	mockShell.EXPECT().
		Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
			"--raw-field", queryField(`i100:`),
			"--raw-field", "fullPath=group/repo").
		Return(`{"data":{"project":{"i1":{"iid":"1","title":"One","state":"opened"}}}}`, nil).
		Once()

	// Second batch: 1 item (101)
	mockShell.EXPECT().
		Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
			"--raw-field", queryField(`i101:`),
			"--raw-field", "fullPath=group/repo").
		Return(`{"data":{"project":{"i101":{"iid":"101","title":"One-O-One","state":"closed"}}}}`, nil).
		Once()

	found, missing, err := g.Issues("gitlab.com", "group/repo", numbers)
	require.NoError(t, err)
	// Issues 2..100 were null in batch 1, so they're reported missing
	assert.Len(t, found, 2)
	assert.Equal(t, "One", found[1].Title)
	assert.Equal(t, "One-O-One", found[101].Title)
	assert.Len(t, missing, 99)
}
