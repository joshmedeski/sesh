package worktree

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/browser"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
)

func falsePtr() *bool { b := false; return &b }

func nuConfig() model.Config {
	return model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{
			Repo:           "nutiliti/nutiliti",
			Path:           "/repo",
			WorktreeDir:    "w",
			BranchTemplate: "jam/{number}-1",
			BaseBranch:     "origin/main",
			Fetch:          falsePtr(),
			StartupCommand: "nu_setup",
		}},
	}
}

func TestConnectIssueWithRepoOverride(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	// path expansion: "/repo" has no ~, but home.ExpandPath calls os.UserHomeDir
	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()

	// Not a PR => issue path
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	// target does not exist yet
	mOs.EXPECT().Stat("/repo/w/2345").Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll("/repo/w", mock.Anything).Return(nil)

	mGit.EXPECT().WorktreeAdd("/repo", "/repo/w/2345", "jam/2345-1", "origin/main").Return("", nil)

	mConn.EXPECT().
		Connect("/repo/w/2345", model.ConnectOpts{Switch: true, Command: "nu_setup"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectOwnPrWithClosingIssue(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()
	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()

	// It is a PR authored by me, closing issue 42
	mGh.EXPECT().PrView("nutiliti/nutiliti", 100).
		Return(github.PullRequest{Author: "me", ClosingIssues: []int{42}}, true, nil)
	mGh.EXPECT().CurrentUser().Return("me", nil)

	// Keyed on the issue (42), normal branch create
	mOs.EXPECT().Stat("/repo/w/42").Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll("/repo/w", mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", "/repo/w/42", "jam/42-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect("/repo/w/42", model.ConnectOpts{Switch: true, Command: "nu_setup"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 100, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectForeignPrDetachCheckout(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()
	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()

	// PR authored by someone else
	mGh.EXPECT().PrView("nutiliti/nutiliti", 200).
		Return(github.PullRequest{Author: "octocat"}, true, nil)
	mGh.EXPECT().CurrentUser().Return("me", nil)

	// Keyed on PR number, detached worktree + gh pr checkout
	mOs.EXPECT().Stat("/repo/w/200").Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll("/repo/w", mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAddDetached("/repo", "/repo/w/200", "origin/main").Return("", nil)
	mGh.EXPECT().PrCheckout("/repo/w/200", "nutiliti/nutiliti", 200).Return("", nil)
	mConn.EXPECT().Connect("/repo/w/200", model.ConnectOpts{Switch: true, Command: "nu_setup"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 200, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectResolvesConfigFromCwd(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()

	// No --repo override: resolveConfig detects the repo from cwd via
	// `git worktree list --porcelain` on the main worktree.
	mOs.EXPECT().Getwd().Return("/repo/w/2345", nil)
	mGit.EXPECT().WorktreeList("/repo/w/2345").Return(true,
		"worktree /repo\nHEAD abc\nbranch refs/heads/main\n\n"+
			"worktree /repo/w/2345\nHEAD def\nbranch refs/heads/jam/2345-1\n",
		nil)

	// Not a PR => issue path
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	mOs.EXPECT().Stat("/repo/w/2345").Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll("/repo/w", mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", "/repo/w/2345", "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect("/repo/w/2345", model.ConnectOpts{Switch: true, Command: "nu_setup"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Switch: true})
	require.NoError(t, err)
}

func TestConnectOwnPrFirstIssueRefFallback(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()
	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()

	// PR authored by me, no ClosingIssues, but body references "#7"
	mGh.EXPECT().PrView("nutiliti/nutiliti", 100).
		Return(github.PullRequest{Author: "me", Body: "fixes #7"}, true, nil)
	mGh.EXPECT().CurrentUser().Return("me", nil)

	// Keyed on issue 7 scanned from title+body, normal branch create (not detached/PrCheckout)
	mOs.EXPECT().Stat("/repo/w/7").Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll("/repo/w", mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", "/repo/w/7", "jam/7-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect("/repo/w/7", model.ConnectOpts{Switch: true, Command: "nu_setup"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 100, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectNoConfigMatch(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 1, Repo: "unknown/repo"})
	require.Error(t, err)
}
