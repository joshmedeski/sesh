package github

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
)

// githubIssue is the GitHub-specific internal wire DTO for CLI and GraphQL issue payloads.
type githubIssue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"` // "OPEN" | "CLOSED" | "MERGED" (for PRs)
}

// githubIssueStates maps GitHub’s state strings to forge-neutral enum values.
var githubIssueStates = map[string]forge.IssueState{
	"OPEN":   forge.IssueStateOpen,
	"CLOSED": forge.IssueStateClosed,
	"MERGED": forge.IssueStateMerged,
}

// toIssue converts a GitHub-specific issue DTO to a forge-neutral Issue.
func (g githubIssue) toIssue() forge.Issue {
	upper := strings.ToUpper(g.State)
	state, ok := githubIssueStates[upper]
	if !ok {
		state = forge.IssueState(upper)
	}
	return forge.Issue{
		Number: g.Number,
		Title:  g.Title,
		State:  state,
	}
}

type RealGithub struct {
	shell shell.Shell
	git   git.Git
}

func NewGithub(shell shell.Shell, git git.Git) forge.Forge {
	return &RealGithub{shell, git}
}

// PrsAsIssues fetches many pull requests from one repo in as few requests as possible,
// returning them as forge.Issue entries. This method exists only to satisfy the forge.Forge interface;
// it simply calls Issues().
func (g *RealGithub) PrsAsIssues(host string, repo string, numbers []int) (map[int]forge.Issue, []int, error) {
	return g.Issues(host, repo, numbers)
}

// WorktreeName generates the directory name (and default branch name) for a worktree
// from a given issue or pull request number. GitHub does not distinguish between
// issue and pull request numbers, so this implementation ignores the isPr parameter.
func (g *RealGithub) WorktreeName(number int, _ bool) string {
	return strconv.Itoa(number)
}

// WorktreeNumber derives the issue number and whether it explicitly refers to a pull request
// from the given worktreeName. GitHub does not distinguish between issue and
// pull request numbers, so this implementation always returns false for the isPr value.
func (g *RealGithub) WorktreeNumber(worktreeName string) (int, bool, error) {
	number, err := strconv.Atoi(worktreeName)
	return number, false, err
}

// IssueKey generates a cache/lookup key for an issue or pull request in a repository.
// GitHub does not distinguish between issue and pull request numbers, so this
// implementation ignores the isPr parameter. The host is part of the key: two hosts can
// carry the same owner/name.
func (g *RealGithub) IssueKey(host string, repo string, number int, _ bool) string {
	return fmt.Sprintf("%s/%s#%d", host, repo, number)
}

// IssueURL returns the issue or pull request URL for the given host, repo, and number.
func (g *RealGithub) IssueURL(host string, repo string, number int, isPr bool) string {
	if isPr {
		return fmt.Sprintf("https://%s/%s/pull/%d", host, repo, number)
	}
	return fmt.Sprintf("https://%s/%s/issues/%d", host, repo, number)
}

func (g *RealGithub) PrView(host string, repo string, number int) (forge.PullRequest, bool, error) {
	out, err := g.shell.Cmd(
		"gh", "pr", "view", strconv.Itoa(number),
		"--repo", ghRepoArg(host, repo),
		"--json", "author,closingIssuesReferences,title,body,url,headRefName,headRepository,headRepositoryOwner,isCrossRepository",
	)
	if err != nil || out == "" {
		// Not a PR (or gh reported no such PR); treat as "not found", not an error.
		return forge.PullRequest{}, false, nil
	}

	var raw struct {
		Author struct {
			Login string `json:"login"`
		} `json:"author"`
		ClosingIssuesReferences []struct {
			Number int `json:"number"`
		} `json:"closingIssuesReferences"`
		Title          string `json:"title"`
		Body           string `json:"body"`
		URL            string `json:"url"`
		HeadRefName    string `json:"headRefName"`
		HeadRepository struct {
			Name string `json:"name"`
		} `json:"headRepository"`
		HeadRepositoryOwner struct {
			Login string `json:"login"`
		} `json:"headRepositoryOwner"`
		IsCrossRepository bool `json:"isCrossRepository"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return forge.PullRequest{}, false, err
	}

	pr := forge.PullRequest{
		Author:      raw.Author.Login,
		Title:       raw.Title,
		Body:        raw.Body,
		URL:         raw.URL,
		HeadRefName: raw.HeadRefName,
		HeadOwner:   raw.HeadRepositoryOwner.Login,
		HeadRepo:    raw.HeadRepository.Name,
		IsFork:      raw.IsCrossRepository,
	}
	for _, ref := range raw.ClosingIssuesReferences {
		pr.ClosingIssues = append(pr.ClosingIssues, ref.Number)
	}
	return pr, true, nil
}

func (g *RealGithub) CurrentUser(host string) (string, error) {
	return g.shell.Cmd("gh", "api", "user", "--hostname", host, "--jq", ".login")
}

func (g *RealGithub) PrCheckout(dir string, host string, repo string, number int) (string, error) {
	return g.shell.CmdInDir(dir, "gh", "pr", "checkout", strconv.Itoa(number), "--repo", ghRepoArg(host, repo))
}

func (g *RealGithub) PrURL(dir string) (string, bool, error) {
	out, err := g.shell.CmdInDir(dir, "gh", "pr", "view", "--json", "url", "--jq", ".url")
	if err != nil || out == "" {
		return "", false, nil
	}
	return out, true, nil
}

func (g *RealGithub) PrURLForHead(host string, repo string, head string) (string, bool, error) {
	out, err := g.shell.Cmd(
		"gh", "pr", "list", "--repo", ghRepoArg(host, repo), "--head", head,
		"--state", "all", "--limit", "1", "--json", "url", "--jq", ".[0].url",
	)
	if err != nil || out == "" {
		return "", false, nil
	}
	return out, true, nil
}

// ghRepoArg returns the [host/]repo form gh accepts for its --repo flag.
// github.com is gh's own default, so it is omitted to keep the common case short.
func ghRepoArg(host string, repo string) string {
	if host == "" || host == "github.com" {
		return repo
	}
	return fmt.Sprintf("%s/%s", host, repo)
}

// Issue returns the issue for the branch checked out at the given path, whose name is
// expected to contain the issue number (e.g., "1234-some-feature"). found is false with
// a nil error whenever there is nothing to report.
func (g *RealGithub) Issue(path string) (forge.Issue, bool, error) {
	ok, branch, err := g.git.CurrentBranch(path)
	if err != nil || !ok {
		return forge.Issue{}, false, nil
	}
	numStr, has := forge.ParseIssueNumber(branch)
	if !has {
		return forge.Issue{}, false, nil
	}

	out, err := g.shell.CmdInDir(path, "gh", "issue", "view", numStr, "--json", "number,title,state")
	if err != nil || out == "" {
		return forge.Issue{}, false, nil
	}

	var dto githubIssue
	if err := json.Unmarshal([]byte(out), &dto); err != nil {
		slog.Debug("github: could not decode issue view response", "path", path, "error", err)
		return forge.Issue{}, false, nil
	}
	return dto.toIssue(), true, nil
}
