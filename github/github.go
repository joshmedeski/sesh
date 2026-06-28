package github

import (
	"encoding/json"
	"regexp"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
)

// Issue is the subset of GitHub issue data sesh renders in the status bar.
type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"` // "OPEN" | "CLOSED"
}

type Github interface {
	// Issue returns the GitHub issue for the branch checked out at path.
	// The bool is false (with a nil error) for every "nothing to show" case.
	Issue(path string) (Issue, bool, error)
}

type RealGithub struct {
	shell shell.Shell
	git   git.Git
}

func NewGithub(shell shell.Shell, git git.Git) Github {
	return &RealGithub{shell, git}
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

	number, ok := parseIssueNumber(branch)
	if !ok {
		return Issue{}, false, nil
	}

	out, err := g.shell.Cmd("gh", "issue", "view", number, "--json", "number,title,state")
	if err != nil || out == "" {
		return Issue{}, false, nil
	}

	var issue Issue
	if err := json.Unmarshal([]byte(out), &issue); err != nil {
		return Issue{}, false, nil
	}
	return issue, true, nil
}
