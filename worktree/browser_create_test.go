package worktree

import (
	"os"
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
)

func TestCreateFromBrowserIssue(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()

	// Browser resolves the issue URL for the configured repo.
	mBrowser.EXPECT().ActiveTabURL().
		Return("https://github.com/nutiliti/nutiliti/issues/2345", true, nil)

	// Not a PR => issue path.
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	mOs.EXPECT().Stat("/repo/w/2345").Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll("/repo/w", mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", "/repo/w/2345", "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().
		Connect("/repo/w/2345", model.ConnectOpts{Switch: true, Command: "nu_setup"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, mBrowser, h, mOs, p)
	_, err := w.Create(model.WorktreeCreateOpts{FromBrowser: true, Switch: true})
	require.NoError(t, err)
}

func TestCreateFromBrowserUnparseableURL(t *testing.T) {
	mBrowser := browser.NewMockBrowser(t)
	mBrowser.EXPECT().ActiveTabURL().Return("https://example.com/", true, nil)

	w := NewWorktree(
		nuConfig(),
		git.NewMockGit(t), github.NewMockGithub(t), connector.NewMockConnector(t),
		mBrowser, home.NewHome(oswrap.NewMockOs(t)), oswrap.NewMockOs(t), pathwrap.NewPath(),
	)
	_, err := w.Create(model.WorktreeCreateOpts{FromBrowser: true})
	require.Error(t, err)
}

func TestCreateFromBrowserUnavailable(t *testing.T) {
	mBrowser := browser.NewMockBrowser(t)
	// Non-macOS / no application configured => skipped.
	mBrowser.EXPECT().ActiveTabURL().Return("", false, nil)

	w := NewWorktree(
		nuConfig(),
		git.NewMockGit(t), github.NewMockGithub(t), connector.NewMockConnector(t),
		mBrowser, home.NewHome(oswrap.NewMockOs(t)), oswrap.NewMockOs(t), pathwrap.NewPath(),
	)
	_, err := w.Create(model.WorktreeCreateOpts{FromBrowser: true})
	require.Error(t, err)
}
