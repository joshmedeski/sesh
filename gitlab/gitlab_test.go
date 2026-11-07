package gitlab

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// The iid arrives in two different JSON shapes depending on the transport, and the
// same DTOs have to decode both.
func TestIidDecodesBothWireForms(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    iid
		wantErr bool
	}{
		{"number, as sent by the REST API", `1234`, 1234, false},
		{"string, as sent by GraphQL", `"1234"`, 1234, false},
		{"null", `null`, 0, false},
		{"empty string", `""`, 0, true},
		{"non-numeric string", `"abc"`, 0, true},
		{"fractional number", `12.5`, 0, true},
		{"boolean", `true`, 0, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got struct {
				Iid iid `json:"iid"`
			}
			err := json.Unmarshal([]byte(`{"iid":`+c.in+`}`), &got)
			if c.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, c.want, got.Iid)
		})
	}
}

func TestIssueDecodesRESTPayload(t *testing.T) {
	// `glab issue view --output json` is the REST API, which sends iid as a number.
	// This is the shape whose rejection broke `sesh rename --enrich`.
	mockShell := shell.NewMockShell(t)
	mockGit := git.NewMockGit(t)
	g := NewGitlab(mockShell, mockGit)
	path := "/repo/w/1234"
	mockGit.EXPECT().CurrentBranch(path).Return(true, "1234-some-feature", nil)
	mockShell.EXPECT().
		CmdInDir(path, "glab", "issue", "view", "1234", "--output", "json").
		Return(`{"iid":1234,"title":"Warm the cache","state":"opened"}`, nil)

	issue, found, err := g.Issue(path)

	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, forge.Issue{
		Number: 1234,
		Title:  "Warm the cache",
		State:  forge.IssueStateOpen,
	}, issue)
}

func TestIssueNotFoundsOnDecodeFailure(t *testing.T) {
	// An undecodable payload is "no issue here", matching the GitHub driver: a missing
	// title must not become a failure for rename --enrich.
	mockShell := shell.NewMockShell(t)
	mockGit := git.NewMockGit(t)
	g := NewGitlab(mockShell, mockGit)
	path := "/repo/w/1234"
	mockGit.EXPECT().CurrentBranch(path).Return(true, "1234-some-feature", nil)
	mockShell.EXPECT().
		CmdInDir(path, "glab", "issue", "view", "1234", "--output", "json").
		Return("not json", nil)

	issue, found, err := g.Issue(path)

	require.NoError(t, err)
	assert.False(t, found)
	assert.Equal(t, forge.Issue{}, issue)
}

// queryField matches the --raw-field argument carrying a GraphQL query without
// pinning the whole generated string, which is multi-line and whitespace-sensitive.
func queryField(fragment string) any {
	return mock.MatchedBy(func(q string) bool {
		return strings.HasPrefix(q, "query=") && strings.Contains(q, fragment)
	})
}

func TestIssuesDecodesGraphQLPayload(t *testing.T) {
	// `glab api graphql` sends iid as a string, because GraphQL's ID scalar is one.
	mockShell := shell.NewMockShell(t)
	g := NewGitlab(mockShell, git.NewMockGit(t))
	mockShell.EXPECT().
		Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
			"--raw-field", queryField(`issue(iid: "42")`),
			"--raw-field", "fullPath=group/repo").
		Return(`{"data":{"project":{"i42":{"iid":"42","title":"Warm the cache","state":"opened"}}}}`, nil)

	found, missing, err := g.Issues("gitlab.com", "group/repo", []int{42})

	require.NoError(t, err)
	assert.Empty(t, missing)
	assert.Equal(t, map[int]forge.Issue{
		42: {Number: 42, Title: "Warm the cache", State: forge.IssueStateOpen},
	}, found)
}

func TestPrsAsIssuesDecodesGraphQLPayload(t *testing.T) {
	mockShell := shell.NewMockShell(t)
	g := NewGitlab(mockShell, git.NewMockGit(t))
	mockShell.EXPECT().
		Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
			"--raw-field", queryField(`mergeRequest(iid: "101")`),
			"--raw-field", "fullPath=group/repo").
		Return(`{"data":{"project":{"mr101":{"iid":"101","author":{"username":"me"},"title":"Ship it","description":"body","state":"merged"}}}}`, nil)

	found, missing, err := g.PrsAsIssues("gitlab.com", "group/repo", []int{101})

	require.NoError(t, err)
	assert.Empty(t, missing)
	assert.Equal(t, map[int]forge.Issue{
		101: {Number: 101, Title: "Ship it", State: forge.IssueStateMerged},
	}, found)
}

func TestPrViewDecodesGraphQLPayload(t *testing.T) {
	t.Run("with closing issues", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		mockShell.EXPECT().
			Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
				"--raw-field", queryField(`mergeRequest(iid: "7")`)).
			Return(`{"data":{"project":{"mergeRequest":{"iid":"7","author":{"username":"me"},"title":"Ship it","description":"body","state":"opened"}}}}`, nil)
		mockShell.EXPECT().
			Cmd("glab", "api", "--hostname", "gitlab.com", "/projects/group%2Frepo/merge_requests/7/closes_issues").
			Return(`[{"iid": 42}, {"iid": 43}]`, nil)

		pr, found, err := g.PrView("gitlab.com", "group/repo", 7)

		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "me", pr.Author)
		assert.Equal(t, "Ship it", pr.Title)
		assert.Equal(t, "body", pr.Body)
		assert.Equal(t, []int{42, 43}, pr.ClosingIssues)
	})

	t.Run("without closing issues", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		mockShell.EXPECT().
			Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
				"--raw-field", queryField(`mergeRequest(iid: "7")`)).
			Return(`{"data":{"project":{"mergeRequest":{"iid":"7","author":{"username":"me"},"title":"Ship it","description":"body","state":"opened"}}}}`, nil)
		mockShell.EXPECT().
			Cmd("glab", "api", "--hostname", "gitlab.com", "/projects/group%2Frepo/merge_requests/7/closes_issues").
			Return(`[]`, nil)

		pr, found, err := g.PrView("gitlab.com", "group/repo", 7)

		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "me", pr.Author)
		assert.Equal(t, "Ship it", pr.Title)
		assert.Equal(t, "body", pr.Body)
		assert.Empty(t, pr.ClosingIssues)
	})
}

func TestPrViewReportsNoMergeRequest(t *testing.T) {
	// A null mergeRequest means "no such MR", which the DTO reports as iid 0.
	mockShell := shell.NewMockShell(t)
	g := NewGitlab(mockShell, git.NewMockGit(t))
	mockShell.EXPECT().
		Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
			"--raw-field", queryField(`mergeRequest(iid: "7")`)).
		Return(`{"data":{"project":{"mergeRequest":null}}}`, nil)

	_, found, err := g.PrView("gitlab.com", "group/repo", 7)

	require.NoError(t, err)
	assert.False(t, found)
}

func TestWorktreeName(t *testing.T) {
	g := NewGitlab(nil, nil)
	assert.Equal(t, "42", g.WorktreeName(42, false))
	assert.Equal(t, "mr-42", g.WorktreeName(42, true))
}

func TestWorktreeNumber(t *testing.T) {
	g := NewGitlab(nil, nil)

	num, isPr, err := g.WorktreeNumber("42")
	require.NoError(t, err)
	assert.Equal(t, 42, num)
	assert.False(t, isPr)

	num, isPr, err = g.WorktreeNumber("mr-42")
	require.NoError(t, err)
	assert.Equal(t, 42, num)
	assert.True(t, isPr)

	_, _, err = g.WorktreeNumber("not-a-number")
	require.Error(t, err)

	_, isPr, err = g.WorktreeNumber("mr-not-a-number")
	require.Error(t, err)
	assert.True(t, isPr)
}

func TestIssueKey(t *testing.T) {
	g := NewGitlab(nil, nil)
	assert.Equal(t, "gitlab.com/group/repo#42", g.IssueKey("gitlab.com", "group/repo", 42, false))
	assert.Equal(t, "gitlab.com/group/repo!42", g.IssueKey("gitlab.com", "group/repo", 42, true))
	assert.Equal(t, "gitlab.example.com/group/repo#100", g.IssueKey("gitlab.example.com", "group/repo", 100, false))
	assert.Equal(t, "gitlab.example.com/group/repo!100", g.IssueKey("gitlab.example.com", "group/repo", 100, true))
}

func TestIssueURL(t *testing.T) {
	g := NewGitlab(nil, nil)
	assert.Equal(t, "https://gitlab.com/group/repo/-/work_items/42", g.IssueURL("gitlab.com", "group/repo", 42, false))
	assert.Equal(t, "https://gitlab.com/group/repo/-/merge_requests/42", g.IssueURL("gitlab.com", "group/repo", 42, true))
	assert.Equal(t, "https://gitlab.example.com/group/repo/-/work_items/42", g.IssueURL("gitlab.example.com", "group/repo", 42, false))
	assert.Equal(t, "https://gitlab.example.com/group/repo/-/merge_requests/42", g.IssueURL("gitlab.example.com", "group/repo", 42, true))
}

func TestGlabRepoArg(t *testing.T) {
	assert.Equal(t, "group/repo", glabRepoArg("", "group/repo"))
	assert.Equal(t, "group/repo", glabRepoArg("gitlab.com", "group/repo"))
	assert.Equal(t, "https://gitlab.example.com/group/repo", glabRepoArg("gitlab.example.com", "group/repo"))
}

func TestToIssueStateMapping(t *testing.T) {
	cases := []struct {
		state string
		want  forge.IssueState
	}{
		{"opened", forge.IssueStateOpen},
		{"closed", forge.IssueStateClosed},
		{"merged", forge.IssueStateMerged},
		{"custom", forge.IssueState("custom")},
	}
	for _, c := range cases {
		issue := gitlabIssue{Iid: 1, Title: "Test", State: c.state}
		assert.Equal(t, c.want, issue.toIssue().State)

		mr := gitlabMergeRequest{Iid: 1, Title: "Test", State: c.state}
		assert.Equal(t, c.want, mr.toIssue().State)
	}
}

func TestCurrentUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		mockShell.EXPECT().
			Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
				"--raw-field", "query=query { currentUser { username } }").
			Return(`{"data":{"currentUser":{"username":"josh"}}}`, nil)

		user, err := g.CurrentUser("gitlab.com")
		require.NoError(t, err)
		assert.Equal(t, "josh", user)
	})

	t.Run("subprocess error", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		mockShell.EXPECT().
			Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
				"--raw-field", "query=query { currentUser { username } }").
			Return("", errors.New("auth failed"))

		_, err := g.CurrentUser("gitlab.com")
		require.Error(t, err)
	})

	t.Run("unmarshal error", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		mockShell.EXPECT().
			Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
				"--raw-field", "query=query { currentUser { username } }").
			Return("invalid json", nil)

		_, err := g.CurrentUser("gitlab.com")
		require.Error(t, err)
	})
}

func TestPrCheckout(t *testing.T) {
	t.Run("default host", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		dir := "/repo/w/42"
		mockShell.EXPECT().
			CmdInDir(dir, "glab", "mr", "checkout", "42", "--repo", "group/repo").
			Return("switched to branch", nil)

		out, err := g.PrCheckout(dir, "gitlab.com", "group/repo", 42)
		require.NoError(t, err)
		assert.Equal(t, "switched to branch", out)
	})

	t.Run("custom host", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		dir := "/repo/w/42"
		mockShell.EXPECT().
			CmdInDir(dir, "glab", "mr", "checkout", "42", "--repo", "https://gitlab.example.com/group/repo").
			Return("switched to branch", nil)

		out, err := g.PrCheckout(dir, "gitlab.example.com", "group/repo", 42)
		require.NoError(t, err)
		assert.Equal(t, "switched to branch", out)
	})
}

func TestPrURL(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		dir := "/repo/w/42"
		mockShell.EXPECT().
			CmdInDir(dir, "glab", "mr", "view", "--output", "json", "--jq", ".web_url").
			Return("https://gitlab.com/group/repo/-/merge_requests/42", nil)

		url, found, err := g.PrURL(dir)
		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "https://gitlab.com/group/repo/-/merge_requests/42", url)
	})

	t.Run("empty output", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		dir := "/repo/w/42"
		mockShell.EXPECT().
			CmdInDir(dir, "glab", "mr", "view", "--output", "json", "--jq", ".web_url").
			Return("", nil)

		url, found, err := g.PrURL(dir)
		require.NoError(t, err)
		assert.False(t, found)
		assert.Empty(t, url)
	})

	t.Run("command error returns not found", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		dir := "/repo/w/42"
		mockShell.EXPECT().
			CmdInDir(dir, "glab", "mr", "view", "--output", "json", "--jq", ".web_url").
			Return("", errors.New("no mr found"))

		url, found, err := g.PrURL(dir)
		require.NoError(t, err)
		assert.False(t, found)
		assert.Empty(t, url)
	})
}

func TestPrViewErrors(t *testing.T) {
	t.Run("graphql subprocess error", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		mockShell.EXPECT().
			Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
				"--raw-field", queryField(`mergeRequest(iid: "7")`)).
			Return("", errors.New("network error"))

		_, found, err := g.PrView("gitlab.com", "group/repo", 7)
		require.Error(t, err)
		assert.False(t, found)
	})

	t.Run("graphql unmarshal error", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		mockShell.EXPECT().
			Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
				"--raw-field", queryField(`mergeRequest(iid: "7")`)).
			Return("not json", nil)

		_, found, err := g.PrView("gitlab.com", "group/repo", 7)
		require.Error(t, err)
		assert.False(t, found)
	})

	t.Run("closes_issues subprocess error", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		mockShell.EXPECT().
			Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
				"--raw-field", queryField(`mergeRequest(iid: "7")`)).
			Return(`{"data":{"project":{"mergeRequest":{"iid":"7","author":{"username":"me"},"title":"Ship it","description":"body","state":"opened"}}}}`, nil)
		mockShell.EXPECT().
			Cmd("glab", "api", "--hostname", "gitlab.com", "/projects/group%2Frepo/merge_requests/7/closes_issues").
			Return("", errors.New("api rate limit or auth error"))

		_, found, err := g.PrView("gitlab.com", "group/repo", 7)
		require.Error(t, err)
		assert.False(t, found)
	})

	t.Run("closes_issues unmarshal error", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		g := NewGitlab(mockShell, git.NewMockGit(t))
		mockShell.EXPECT().
			Cmd("glab", "api", "graphql", "--hostname", "gitlab.com",
				"--raw-field", queryField(`mergeRequest(iid: "7")`)).
			Return(`{"data":{"project":{"mergeRequest":{"iid":"7","author":{"username":"me"},"title":"Ship it","description":"body","state":"opened"}}}}`, nil)
		mockShell.EXPECT().
			Cmd("glab", "api", "--hostname", "gitlab.com", "/projects/group%2Frepo/merge_requests/7/closes_issues").
			Return("invalid json", nil)

		_, found, err := g.PrView("gitlab.com", "group/repo", 7)
		require.Error(t, err)
		assert.False(t, found)
	})
}

func TestIssueEdgeCases(t *testing.T) {
	t.Run("git CurrentBranch error", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		mockGit := git.NewMockGit(t)
		g := NewGitlab(mockShell, mockGit)
		path := "/not/a/repo"
		mockGit.EXPECT().CurrentBranch(path).Return(false, "", errors.New("not a git repo"))

		_, found, err := g.Issue(path)
		require.NoError(t, err)
		assert.False(t, found)
	})

	t.Run("branch has no issue number", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		mockGit := git.NewMockGit(t)
		g := NewGitlab(mockShell, mockGit)
		path := "/repo/main"
		mockGit.EXPECT().CurrentBranch(path).Return(true, "main", nil)

		_, found, err := g.Issue(path)
		require.NoError(t, err)
		assert.False(t, found)
	})

	t.Run("glab issue view error", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		mockGit := git.NewMockGit(t)
		g := NewGitlab(mockShell, mockGit)
		path := "/repo/w/1234"
		mockGit.EXPECT().CurrentBranch(path).Return(true, "1234-feature", nil)
		mockShell.EXPECT().
			CmdInDir(path, "glab", "issue", "view", "1234", "--output", "json").
			Return("", errors.New("issue not found"))

		_, found, err := g.Issue(path)
		require.NoError(t, err)
		assert.False(t, found)
	})

	t.Run("glab issue view empty output", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		mockGit := git.NewMockGit(t)
		g := NewGitlab(mockShell, mockGit)
		path := "/repo/w/1234"
		mockGit.EXPECT().CurrentBranch(path).Return(true, "1234-feature", nil)
		mockShell.EXPECT().
			CmdInDir(path, "glab", "issue", "view", "1234", "--output", "json").
			Return("", nil)

		_, found, err := g.Issue(path)
		require.NoError(t, err)
		assert.False(t, found)
	})
}
