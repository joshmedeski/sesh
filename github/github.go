package github

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
)

// Issue is the subset of GitHub issue data sesh renders in the status bar.
type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"` // "OPEN" | "CLOSED"
}

// BranchRef identifies the repo, branch, and (optional) issue number for a path.
type BranchRef struct {
	RepoRoot  string
	Branch    string
	Number    int
	HasNumber bool
}

type Github interface {
	// Issue returns the GitHub issue for the branch checked out at path.
	// The bool is false (with a nil error) for every "nothing to show" case.
	Issue(path string) (Issue, bool, error)
	// Resolve returns repo root, branch, and the issue number parsed from the
	// branch (HasNumber=false if none). ok is false only when path is not a repo.
	Resolve(path string) (BranchRef, bool)
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

func (g *RealGithub) Resolve(path string) (BranchRef, bool) {
	ok, branch, err := g.git.CurrentBranch(path)
	if err != nil || !ok {
		return BranchRef{}, false
	}
	topOk, repoRoot, err := g.git.ShowTopLevel(path)
	if err != nil || !topOk {
		return BranchRef{}, false
	}
	ref := BranchRef{RepoRoot: repoRoot, Branch: branch}
	if numStr, has := parseIssueNumber(branch); has {
		if n, err := strconv.Atoi(numStr); err == nil {
			ref.Number = n
			ref.HasNumber = true
		}
	}
	return ref, true
}

func (g *RealGithub) Issue(path string) (Issue, bool, error) {
	ref, ok := g.Resolve(path)
	if !ok || !ref.HasNumber {
		return Issue{}, false, nil
	}

	out, err := g.shell.Cmd("gh", "issue", "view", strconv.Itoa(ref.Number), "--json", "number,title,state")
	if err != nil || out == "" {
		return Issue{}, false, nil
	}

	var issue Issue
	if err := json.Unmarshal([]byte(out), &issue); err != nil {
		return Issue{}, false, nil
	}
	return issue, true, nil
}
