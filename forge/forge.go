package forge

import (
	"fmt"
	"regexp"

	"github.com/joshmedeski/sesh/v2/model"
)

// IssueState represents the forge-neutral lifecycle status of an issue for worktree workflows.
type IssueState string

const (
	IssueStateOpen    IssueState = "OPEN"
	IssueStateClosed  IssueState = "CLOSED"
	IssueStateMerged  IssueState = "MERGED" // when representing a pull request, not a plain issue – see Issue struct
	IssueStateUnknown IssueState = ""
)

// PullRequest models common metadata for GitHub PRs, GitLab MRs, etc., as needed by worktree workflows.
type PullRequest struct {
	Author        string
	ClosingIssues []int
	Title         string
	Body          string
	URL           string
	HeadRefName   string
	HeadOwner     string
	HeadRepo      string
	IsFork        bool
}

// Issue is the forge-neutral domain entity for issues in worktree workflows.
// In forges that have a common number/reference scheme for issues and pull requests (e.g., GitHub),
// it can also represent a pull request.
type Issue struct {
	Number int
	Title  string
	State  IssueState
}

// Forge defines the common contract for interacting with SCM hosting platforms (GitHub, GitLab, etc.).
type Forge interface {
	// WorktreeName generates the directory name (and default branch name) for a worktree
	// from a given issue or pull/merge request number.
	WorktreeName(number int, isPr bool) string

	// WorktreeNumber derives the issue number and whether it explicitly refers to a pull request
	// from the given worktreeName.
	WorktreeNumber(worktreeName string) (int, bool, error)

	// IssueKey generates a cache/lookup key for an issue or pull/merge request in a repository.
	IssueKey(host string, repo string, number int, isPr bool) string

	// IssueURL returns the issue or pull/merge request URL for the given host, repo, and number.
	IssueURL(host string, repo string, number int, isPr bool) string

	// PrView returns pull/merge request metadata. found is false (with nil error) when the
	// number is not a pull/merge request (e.g. it is a plain issue).
	PrView(host string, repo string, number int) (PullRequest, bool, error)

	// CurrentUser returns the authenticated user's login or username for host.
	CurrentUser(host string) (string, error)

	// PrURL returns the URL of the pull/merge request for the branch checked out in dir.
	// found is false (with nil error) when the branch has no PR/MR.
	PrURL(dir string) (url string, found bool, err error)

	// PrURLForHead returns the URL of the newest pull/merge request in repo whose
	// head branch is head, in any state.
	PrURLForHead(host string, repo string, head string) (url string, found bool, err error)

	// PrCheckout checks out the pull/merge request inside dir (a worktree).
	PrCheckout(dir string, host string, repo string, number int) (string, error)

	// Issue returns the issue for the branch checked out at path. found is false with a
	// nil error whenever there is nothing to report: path is not a git repository, its
	// branch carries no issue number, or the lookup failed, printed nothing, or printed
	// something that cannot be decoded.
	Issue(path string) (Issue, bool, error)

	// Issues fetches many issues from host/repo in as few requests as possible.
	// found holds the issues that resolved; missing holds the numbers confirmed not to exist,
	// so callers can cache those negatively instead of retrying them forever. Numbers that
	// neither resolved nor were confirmed absent appear in neither result and should simply
	// be retried later.
	Issues(host string, repo string, numbers []int) (found map[int]Issue, missing []int, err error)

	// PrsAsIssues fetches many pull/merge requests from host/repo in as few requests as possible,
	// returning them as forge.Issue entries. This is intended for legacy GitHub-centric workflows
	// that treat PRs as issues. See Issues for the meaning of the return values.
	PrsAsIssues(host string, repo string, numbers []int) (found map[int]Issue, missing []int, err error)
}

// Host returns the host a worktree's forge calls should target: the configured host when one
// is set, otherwise the forge's public default.
func Host(cfg model.WorktreeConfig) (string, error) {
	if cfg.Host != "" {
		return cfg.Host, nil
	}
	normalized, err := cfg.Forge.Normalize()
	if err != nil {
		return "", err
	}
	switch normalized {
	case model.ForgeGithub:
		return "github.com", nil

	case model.ForgeGitlab:
		return "gitlab.com", nil
	}
	return "", fmt.Errorf("no default host for forge %q", cfg.Forge)
}

// Resolver provides methods for determining forge drivers for repositories.
type Resolver interface {
	// FromConfig maps a worktree config to the forge driver that should handle it.
	FromConfig(cfg model.WorktreeConfig) (Forge, error)

	// FromPath maps a repository path to the forge driver that should handle it.
	FromPath(path string) (Forge, error)
}

// RealResolver is a functional Resolver implementation.
type RealResolver struct {
	FromConfigFunc func(cfg model.WorktreeConfig) (Forge, error)
	FromPathFunc   func(path string) (Forge, error)
}

func (r RealResolver) FromConfig(cfg model.WorktreeConfig) (Forge, error) {
	if r.FromConfigFunc == nil {
		return nil, fmt.Errorf("FromConfigFunc is not implemented for %v", r)
	}
	return r.FromConfigFunc(cfg)
}

func (r RealResolver) FromPath(path string) (Forge, error) {
	if r.FromPathFunc == nil {
		return nil, fmt.Errorf("FromPathFunc is not implemented for %v", r)
	}
	return r.FromPathFunc(path)
}

func NewResolver(
	fromConfig func(cfg model.WorktreeConfig) (Forge, error),
	fromPath func(path string) (Forge, error),
) RealResolver {
	return RealResolver{
		FromConfigFunc: fromConfig,
		FromPathFunc:   fromPath,
	}
}

var issueNumberRe = regexp.MustCompile(`\d+`)

// ParseIssueNumber returns the first run of digits in a branch name. Both forges
// encode the issue number there, e.g. "1234-some-feature".
func ParseIssueNumber(branch string) (string, bool) {
	match := issueNumberRe.FindString(branch)
	if match == "" {
		return "", false
	}
	return match, true
}
