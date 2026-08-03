package github

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
)

// PullRequest is the subset of gh pr view data the worktree flow needs.
type PullRequest struct {
	Author        string
	ClosingIssues []int
	Title         string
	Body          string
}

// Issue is the subset of GitHub issue data sesh renders in the status bar.
type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"` // "OPEN" | "CLOSED"
}

type Github interface {
	// PrView returns PR metadata. found is false (with nil error) when the
	// number is not a pull request (e.g. it is a plain issue).
	PrView(repo string, number int) (PullRequest, bool, error)
	// CurrentUser returns the authenticated gh user's login.
	CurrentUser() (string, error)
	// PrCheckout runs `gh pr checkout` inside dir (a worktree).
	PrCheckout(dir string, repo string, number int) (string, error)
	// Issue returns the GitHub issue for the branch checked out at path.
	// The bool is false (with a nil error) for every "nothing to show" case.
	Issue(path string) (Issue, bool, error)
	// Issues fetches many issues from one repo in as few requests as
	// possible. found holds the issues that resolved; missing holds the
	// numbers GitHub confirmed do not exist, so callers can cache those
	// negatively instead of retrying them forever. Numbers that neither
	// resolved nor were confirmed absent (a rate limit, say) appear in
	// neither result and should simply be retried later.
	Issues(repo string, numbers []int) (found map[int]Issue, missing []int, err error)
}

type RealGithub struct {
	shell shell.Shell
	git   git.Git
}

func NewGithub(shell shell.Shell, git git.Git) Github {
	return &RealGithub{shell, git}
}

func (g *RealGithub) PrView(repo string, number int) (PullRequest, bool, error) {
	out, err := g.shell.Cmd(
		"gh", "pr", "view", strconv.Itoa(number),
		"--repo", repo,
		"--json", "author,closingIssuesReferences,title,body",
	)
	if err != nil || out == "" {
		// Not a PR (or gh reported no such PR); treat as "not found", not an error.
		return PullRequest{}, false, nil
	}

	var raw struct {
		Author struct {
			Login string `json:"login"`
		} `json:"author"`
		ClosingIssuesReferences []struct {
			Number int `json:"number"`
		} `json:"closingIssuesReferences"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return PullRequest{}, false, err
	}

	pr := PullRequest{Author: raw.Author.Login, Title: raw.Title, Body: raw.Body}
	for _, ref := range raw.ClosingIssuesReferences {
		pr.ClosingIssues = append(pr.ClosingIssues, ref.Number)
	}
	return pr, true, nil
}

func (g *RealGithub) CurrentUser() (string, error) {
	return g.shell.Cmd("gh", "api", "user", "--jq", ".login")
}

func (g *RealGithub) PrCheckout(dir string, repo string, number int) (string, error) {
	return g.shell.CmdInDir(dir, "gh", "pr", "checkout", strconv.Itoa(number), "--repo", repo)
}

var issueNumberRe = regexp.MustCompile(`\d+`)

// parseIssueNumber returns the first run of digits in a branch name.
func parseIssueNumber(branch string) (string, bool) {
	match := issueNumberRe.FindString(branch)
	if match == "" {
		return "", false
	}
	return match, true
}

func (g *RealGithub) Issue(path string) (Issue, bool, error) {
	ok, branch, err := g.git.CurrentBranch(path)
	if err != nil || !ok {
		return Issue{}, false, nil
	}
	numStr, has := parseIssueNumber(branch)
	if !has {
		return Issue{}, false, nil
	}

	out, err := g.shell.Cmd("gh", "issue", "view", numStr, "--json", "number,title,state")
	if err != nil || out == "" {
		return Issue{}, false, nil
	}

	var issue Issue
	if err := json.Unmarshal([]byte(out), &issue); err != nil {
		return Issue{}, false, nil
	}
	return issue, true, nil
}
