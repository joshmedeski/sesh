package gitlab

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
)

// iid is a GitLab issue or merge request number as it appears on the wire. The two
// transports this package reads disagree on its JSON type: the REST API, which backs
// `glab issue view --output json`, sends a number, while GraphQL, which backs
// `glab api graphql`, sends a string (its ID scalar is always one). Both shapes reach
// the same DTOs, so decoding is the type's job.
//
// Zero means "no issue or merge request": PrView's not-found check reads it that way.
type iid int

// UnmarshalJSON accepts either wire form of an iid, 1234 or "1234".
func (i *iid) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*i = 0
		return nil
	}
	// Number form, as sent by the REST API.
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*i = iid(n)
		return nil
	}
	// String form, as sent by GraphQL.
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("gitlab iid must be a number or a numeric string, got %s", data)
	}
	number, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("gitlab iid %q is not a number: %w", s, err)
	}
	*i = iid(number)
	return nil
}

// gitlabIssue is the GitLab-specific internal wire DTO for CLI and GraphQL issue payloads.
type gitlabIssue struct {
	Iid   iid    `json:"iid"`
	Title string `json:"title"`
	State string `json:"state"` // "opened" | "closed"
}

type gitlabMergeRequest struct {
	Iid         iid    `json:"iid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"` // "opened" | "closed" | "merged"
	Author      struct {
		Username string `json:"username"`
	} `json:"author"`
}

// gitlabIssueStates maps GitLab’s state strings to forge-neutral enum values.
// "merged" exists here only so legacy GitHub-centric workflows that treat PRs as issues
// can convert GitLab MRs to forge.Issue entries.
var gitlabIssueStates = map[string]forge.IssueState{
	"opened": forge.IssueStateOpen,
	"closed": forge.IssueStateClosed,
	"merged": forge.IssueStateMerged,
}

// toIssue converts a GitLab-specific issue DTO to a forge-neutral Issue.
func (g gitlabIssue) toIssue() forge.Issue {
	lower := strings.ToLower(g.State)
	state, ok := gitlabIssueStates[lower]
	if !ok {
		state = forge.IssueState(g.State)
	}
	return forge.Issue{
		Number: int(g.Iid),
		Title:  g.Title,
		State:  state,
	}
}

// toPullRequest converts a GitLab-specific merge request DTO to a forge-neutral PullRequest.
func (g gitlabMergeRequest) toPullRequest() forge.PullRequest {
	return forge.PullRequest{
		Author: g.Author.Username,
		Title:  g.Title,
		Body:   g.Description,
	}
}

// toIssue converts a GitLab-specific merge request DTO to a forge-neutral Issue.
// This is a non-idiomatic conversion from the GitLab perspective, but it's necessary
// for supporting legacy GitHub-centric workflows that treat PRs as issues.
func (g gitlabMergeRequest) toIssue() forge.Issue {
	lower := strings.ToLower(g.State)
	state, ok := gitlabIssueStates[lower]
	if !ok {
		state = forge.IssueState(g.State)
	}
	return forge.Issue{
		Number: int(g.Iid),
		Title:  g.Title,
		State:  state,
	}
}

type RealGitlab struct {
	shell shell.Shell
	git   git.Git
}

func NewGitlab(shell shell.Shell, git git.Git) forge.Forge {
	return &RealGitlab{shell, git}
}

// WorktreeName generates the directory name for a worktree from a given issue or merge
// request number. A merge request takes an "mr-" prefix, so its directory cannot collide
// with the issue of the same number.
func (g *RealGitlab) WorktreeName(number int, isPr bool) string {
	if isPr {
		return fmt.Sprintf("mr-%d", number)
	}
	return strconv.Itoa(number)
}

// WorktreeNumber derives the issue number and whether it explicitly refers to a merge request
// from the given worktreeName. The worktree name is expected to be either
// a number (e.g., "1234") or a merge request reference (e.g., "mr-1234").
func (g *RealGitlab) WorktreeNumber(worktreeName string) (int, bool, error) {
	if mrKey, found := strings.CutPrefix(worktreeName, "mr-"); found {
		number, err := strconv.Atoi(mrKey)
		return number, true, err
	}
	number, err := strconv.Atoi(worktreeName)
	return number, false, err
}

// IssueKey generates a cache/lookup key for an issue or merge request in a repository.
func (g *RealGitlab) IssueKey(host string, repo string, number int, isPr bool) string {
	if isPr {
		return fmt.Sprintf("%s/%s!%d", host, repo, number)
	}
	return fmt.Sprintf("%s/%s#%d", host, repo, number)
}

// IssueURL returns the URL of the issue for the given host, repo, and issue number.
func (g *RealGitlab) IssueURL(host string, repo string, number int, isPr bool) string {
	if isPr {
		return fmt.Sprintf("https://%s/%s/-/merge_requests/%d", host, repo, number)
	}
	return fmt.Sprintf("https://%s/%s/-/work_items/%d", host, repo, number)
}

// PrView returns merge request metadata for a given host and repo and MR number.
// The bool return value is false (with a nil error) when the MR does not exist.
func (g *RealGitlab) PrView(host string, repo string, number int) (forge.PullRequest, bool, error) {
	// A GraphQL query is faster than `glab mr view`, which has no filter options.
	out, err := g.shell.Cmd(
		"glab", "api", "graphql", "--hostname", host,
		"--raw-field", fmt.Sprintf(`query={
			project(fullPath: "%s") {
				mergeRequest(iid: "%d") {
					iid
					author { username }
					title
					description
					state
				}
			}
		}`, repo, number),
	)
	if err != nil {
		return forge.PullRequest{}, false, err
	}
	var raw struct {
		Data struct {
			Project struct {
				MergeRequest gitlabMergeRequest `json:"mergeRequest"`
			} `json:"project"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return forge.PullRequest{}, false, err
	}
	rawMr := raw.Data.Project.MergeRequest
	if rawMr.Iid == 0 {
		// Indicates that MR wasn't found (query returned mergeRequest: null); not an error.
		return forge.PullRequest{}, false, nil
	}
	pr := rawMr.toPullRequest()

	// Unfortunately, GitLab's GraphQL API does not expose the closing issues of MRs,
	// so we additionally use the most targeted API endpoint available for this.
	issuesOut, err := g.shell.Cmd(
		"glab", "api", "--hostname", host,
		fmt.Sprintf("/projects/%s/merge_requests/%d/closes_issues", url.PathEscape(repo), number),
	)
	if err != nil {
		return forge.PullRequest{}, false, err
	}
	var closedIssues []struct {
		Iid int `json:"iid"`
	}
	if err := json.Unmarshal([]byte(issuesOut), &closedIssues); err != nil {
		return forge.PullRequest{}, false, err
	}
	for _, issue := range closedIssues {
		if issue.Iid > 0 {
			pr.ClosingIssues = append(pr.ClosingIssues, issue.Iid)
		}
	}
	return pr, true, nil
}

// CurrentUser returns the authenticated user’s username for host.
func (g *RealGitlab) CurrentUser(host string) (string, error) {
	// A GraphQL query is the only option here.
	out, err := g.shell.Cmd(
		"glab", "api", "graphql", "--hostname", host,
		"--raw-field", "query=query { currentUser { username } }",
	)
	if err != nil {
		return "", err
	}
	var raw struct {
		Data struct {
			CurrentUser struct {
				Username string `json:"username"`
			} `json:"currentUser"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return "", err
	}
	return raw.Data.CurrentUser.Username, nil
}

// PrCheckout checks out the merge request with the given number.
func (g *RealGitlab) PrCheckout(dir string, host string, repo string, number int) (string, error) {
	// `glab mr checkout` has no --hostname flag, but its --repo flag accepts https://host/repo.
	return g.shell.CmdInDir(
		dir, "glab", "mr", "checkout", strconv.Itoa(number),
		"--repo", glabRepoArg(host, repo),
	)
}

// PrURL returns the URL of the merge request for the branch checked out in dir.
// The bool return value is false (with a nil error) when the branch has no MR.
func (g *RealGitlab) PrURL(dir string) (string, bool, error) {
	// A glab api graphql query would be faster, but would require repo "group/name" and MR number,
	// which are not directly available here.
	out, err := g.shell.CmdInDir(dir, "glab", "mr", "view", "--output", "json", "--jq", ".web_url")
	if err != nil || out == "" {
		return "", false, nil
	}
	return out, true, nil
}

// PrURLForHead returns the URL of the newest merge request in repo whose
// head branch is head, in any state.
func (g *RealGitlab) PrURLForHead(host string, repo string, head string) (string, bool, error) {
	return "", false, nil
}

// glabRepoArg returns the [https://host/]repo form glab accepts for its --repo flag.
// gitlab.com is glab's own default, so it is omitted to keep the common case short.
func glabRepoArg(host string, repo string) string {
	if host == "" || host == "gitlab.com" {
		return repo
	}
	return fmt.Sprintf("https://%s/%s", host, repo)
}

// Issue returns the issue for the branch checked out at the given path, whose name is
// expected to contain the issue number (e.g., "1234-some-feature"). found is false with
// a nil error whenever there is nothing to report.
func (g *RealGitlab) Issue(path string) (forge.Issue, bool, error) {
	ok, branch, err := g.git.CurrentBranch(path)
	if err != nil || !ok {
		return forge.Issue{}, false, nil
	}
	numStr, has := forge.ParseIssueNumber(branch)
	if !has {
		return forge.Issue{}, false, nil
	}
	out, err := g.shell.CmdInDir(path, "glab", "issue", "view", numStr, "--output", "json")
	if err != nil || out == "" {
		return forge.Issue{}, false, nil
	}
	var dto gitlabIssue
	if err := json.Unmarshal([]byte(out), &dto); err != nil {
		slog.Debug("gitlab: could not decode issue view response", "path", path, "error", err)
		return forge.Issue{}, false, nil
	}
	return dto.toIssue(), true, nil
}
