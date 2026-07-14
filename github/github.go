package github

import (
	"encoding/json"
	"strconv"

	"github.com/joshmedeski/sesh/v2/shell"
)

// PullRequest is the subset of gh pr view data the worktree flow needs.
type PullRequest struct {
	Author        string
	ClosingIssues []int
	Title         string
	Body          string
}

type Github interface {
	// PrView returns PR metadata. found is false (with nil error) when the
	// number is not a pull request (e.g. it is a plain issue).
	PrView(repo string, number int) (PullRequest, bool, error)
	// CurrentUser returns the authenticated gh user's login.
	CurrentUser() (string, error)
	// PrCheckout runs `gh pr checkout` inside dir (a worktree).
	PrCheckout(dir string, repo string, number int) (string, error)
}

type RealGithub struct {
	shell shell.Shell
}

func NewGithub(shell shell.Shell) Github {
	return &RealGithub{shell}
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
