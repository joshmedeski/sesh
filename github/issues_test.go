package github

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
)

// expectGraphql sets up the one CmdCapture call Issues makes for numbers,
// returning out and runErr.
func expectGraphql(t *testing.T, s *shell.MockShell, numbers []int, out string, runErr error) {
	t.Helper()
	query, _ := buildIssueQuery(numbers)
	s.EXPECT().
		CmdCapture("gh", "api", "graphql",
			"-f", "query="+query,
			"-f", "owner=joshmedeski",
			"-f", "name=sesh",
		).
		Return(out, runErr).
		Once()
}

func TestIssues_BatchesIntoOneCall(t *testing.T) {
	s := shell.NewMockShell(t)
	expectGraphql(t, s, []int{409, 426}, `{"data":{"repository":{
		"i409":{"number":409,"title":"worktree support","state":"OPEN"},
		"i426":{"number":426,"title":"cap fuzzy penalty","state":"CLOSED"}
	}}}`, nil)

	g := NewGithub(s, git.NewMockGit(t))
	found, missing, err := g.Issues("joshmedeski/sesh", []int{409, 426})

	require.NoError(t, err)
	assert.Empty(t, missing)
	assert.Equal(t, map[int]Issue{
		409: {Number: 409, Title: "worktree support", State: "OPEN"},
		426: {Number: 426, Title: "cap fuzzy penalty", State: "CLOSED"},
	}, found)
}

// This is a verbatim response from `gh api graphql`, which exits 1 for the
// unresolvable number while still returning the sibling that resolved.
func TestIssues_KeepsSiblingsWhenOneIsNotFound(t *testing.T) {
	s := shell.NewMockShell(t)
	expectGraphql(t, s, []int{409, 99999999},
		`{"data":{"repository":{"i409":{"number":409,"title":"Add git worktree support","state":"OPEN"},"i99999999":null}},"errors":[{"type":"NOT_FOUND","path":["repository","i99999999"],"locations":[{"line":1,"column":117}],"message":"Could not resolve to an Issue with the number of 99999999."}]}`,
		fmt.Errorf("gh: Could not resolve to an Issue with the number of 99999999"),
	)

	g := NewGithub(s, git.NewMockGit(t))
	found, missing, err := g.Issues("joshmedeski/sesh", []int{409, 99999999})

	require.NoError(t, err, "a NOT_FOUND for one number must not fail the batch")
	assert.Equal(t, map[int]Issue{
		409: {Number: 409, Title: "Add git worktree support", State: "OPEN"},
	}, found)
	assert.Equal(t, []int{99999999}, missing, "confirmed-absent numbers are reported for negative caching")
}

func TestIssues_NullWithoutNotFoundIsNotMissing(t *testing.T) {
	s := shell.NewMockShell(t)
	expectGraphql(t, s, []int{409},
		`{"data":{"repository":{"i409":null}},"errors":[{"type":"RATE_LIMITED","path":["repository","i409"],"message":"API rate limit exceeded"}]}`,
		fmt.Errorf("gh: API rate limit exceeded"),
	)

	g := NewGithub(s, git.NewMockGit(t))
	found, missing, err := g.Issues("joshmedeski/sesh", []int{409})

	require.NoError(t, err)
	assert.Empty(t, found)
	// A transient failure must not be cached as "this issue does not exist".
	assert.Empty(t, missing)
}

func TestIssues_ErrorsWhenNoBodyAtAll(t *testing.T) {
	s := shell.NewMockShell(t)
	expectGraphql(t, s, []int{409}, "", fmt.Errorf("gh: not authenticated"))

	g := NewGithub(s, git.NewMockGit(t))
	_, _, err := g.Issues("joshmedeski/sesh", []int{409})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not authenticated")
}

func TestIssues_ErrorsOnEmptySuccessfulResponse(t *testing.T) {
	s := shell.NewMockShell(t)
	expectGraphql(t, s, []int{409}, "", nil)

	g := NewGithub(s, git.NewMockGit(t))
	_, _, err := g.Issues("joshmedeski/sesh", []int{409})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

func TestIssues_MalformedBodyReportsGhFailure(t *testing.T) {
	s := shell.NewMockShell(t)
	expectGraphql(t, s, []int{409}, "not json", fmt.Errorf("gh: could not connect"))

	g := NewGithub(s, git.NewMockGit(t))
	_, _, err := g.Issues("joshmedeski/sesh", []int{409})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not connect")
}

func TestIssues_MalformedBodyWithoutGhError(t *testing.T) {
	s := shell.NewMockShell(t)
	expectGraphql(t, s, []int{409}, "not json", nil)

	g := NewGithub(s, git.NewMockGit(t))
	_, _, err := g.Issues("joshmedeski/sesh", []int{409})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "decoding gh graphql response")
}

func TestIssues_DedupesAndSorts(t *testing.T) {
	s := shell.NewMockShell(t)
	// Callers pass worktree dirs in arbitrary order, possibly repeated; the
	// query is built from a deduped, sorted set.
	expectGraphql(t, s, []int{409, 426}, `{"data":{"repository":{
		"i409":{"number":409,"title":"a","state":"OPEN"},
		"i426":{"number":426,"title":"b","state":"OPEN"}
	}}}`, nil)

	g := NewGithub(s, git.NewMockGit(t))
	found, _, err := g.Issues("joshmedeski/sesh", []int{426, 409, 426, 409})

	require.NoError(t, err)
	assert.Len(t, found, 2)
}

func TestIssues_NoNumbersMakesNoCall(t *testing.T) {
	// No CmdCapture is registered, so any gh call fails the test.
	g := NewGithub(shell.NewMockShell(t), git.NewMockGit(t))

	found, missing, err := g.Issues("joshmedeski/sesh", nil)

	require.NoError(t, err)
	assert.Empty(t, found)
	assert.Empty(t, missing)
}

func TestIssues_ChunksLargeRequests(t *testing.T) {
	numbers := make([]int, issueBatchSize+5)
	for i := range numbers {
		numbers[i] = i + 1
	}

	s := shell.NewMockShell(t)
	expectGraphql(t, s, numbers[:issueBatchSize],
		`{"data":{"repository":{"i1":{"number":1,"title":"first","state":"OPEN"}}}}`, nil)
	expectGraphql(t, s, numbers[issueBatchSize:],
		fmt.Sprintf(`{"data":{"repository":{"i%d":{"number":%d,"title":"last","state":"OPEN"}}}}`,
			issueBatchSize+5, issueBatchSize+5), nil)

	g := NewGithub(s, git.NewMockGit(t))
	found, _, err := g.Issues("joshmedeski/sesh", numbers)

	require.NoError(t, err)
	assert.Equal(t, "first", found[1].Title)
	assert.Equal(t, "last", found[issueBatchSize+5].Title)
}

func TestIssues_RejectsMalformedRepo(t *testing.T) {
	g := NewGithub(shell.NewMockShell(t), git.NewMockGit(t))

	for _, repo := range []string{"sesh", "", "/sesh", "joshmedeski/", "a/b/c"} {
		_, _, err := g.Issues(repo, []int{409})
		assert.Error(t, err, repo)
	}
}

func TestBuildIssueQuery(t *testing.T) {
	query, aliases := buildIssueQuery([]int{409, 426})

	assert.Equal(t,
		"query($owner: String!, $name: String!) { repository(owner: $owner, name: $name) {"+
			" i409: issueOrPullRequest(number: 409) "+issueSelection+
			" i426: issueOrPullRequest(number: 426) "+issueSelection+" } }",
		query)
	assert.Equal(t, map[string]int{"i409": 409, "i426": 426}, aliases)
}

// Worktrees are keyed on PR numbers as well as issue numbers, so the query must
// resolve both rather than reporting a PR as a nonexistent issue.
func TestBuildIssueQuery_ResolvesPullRequestsToo(t *testing.T) {
	query, _ := buildIssueQuery([]int{411})

	assert.Contains(t, query, "issueOrPullRequest(number: 411)")
	assert.Contains(t, query, "... on PullRequest { number title state }")
	assert.NotContains(t, query, "i411: issue(number:")
}

func TestSplitRepo(t *testing.T) {
	cases := []struct {
		repo        string
		owner, name string
		ok          bool
	}{
		{"joshmedeski/sesh", "joshmedeski", "sesh", true},
		{"a/b/c", "", "", false},
		{"sesh", "", "", false},
		{"/sesh", "", "", false},
		{"joshmedeski/", "", "", false},
		{"", "", "", false},
	}
	for _, c := range cases {
		owner, name, ok := splitRepo(c.repo)
		assert.Equal(t, c.ok, ok, c.repo)
		assert.Equal(t, c.owner, owner, c.repo)
		assert.Equal(t, c.name, name, c.repo)
	}
}

func TestBatchNumbers(t *testing.T) {
	assert.Nil(t, batchNumbers(nil, 2))
	assert.Equal(t, [][]int{{1, 2}, {3}}, batchNumbers([]int{3, 1, 2, 1}, 2))
	assert.Equal(t, [][]int{{1, 2, 3}}, batchNumbers([]int{1, 2, 3}, 10))
}

func TestNotFoundAliases(t *testing.T) {
	errs := []graphqlError{
		{Type: "NOT_FOUND", Path: []any{"repository", "i1"}},
		{Type: "RATE_LIMITED", Path: []any{"repository", "i2"}},
		{Type: "NOT_FOUND", Path: []any{}},
		{Type: "NOT_FOUND", Path: []any{"repository", float64(3)}},
	}

	assert.Equal(t, map[string]bool{"i1": true}, notFoundAliases(errs))
}
