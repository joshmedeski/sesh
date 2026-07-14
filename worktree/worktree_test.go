package worktree

import (
	"testing"

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
			Name:           "nutiliti/nutiliti",
			Path:           "/repo",
			WorktreeDir:    "w",
			BranchTemplate: "jam/{number}-1",
			BaseBranch:     "origin/main",
			Fetch:          falsePtr(),
			StartupCommand: "nu_setup",
		}},
	}
}

func TestCreateIssueWithRepoOverride(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
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

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, h, mOs, p)
	_, err := w.Create(model.WorktreeCreateOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestCreateOwnPrWithClosingIssue(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
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

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, h, mOs, p)
	_, err := w.Create(model.WorktreeCreateOpts{Number: 100, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestCreateForeignPrDetachCheckout(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
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

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, h, mOs, p)
	_, err := w.Create(model.WorktreeCreateOpts{Number: 200, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestCreateNoConfigMatch(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, h, mOs, p)
	_, err := w.Create(model.WorktreeCreateOpts{Number: 1, Repo: "unknown/repo"})
	require.Error(t, err)
}
