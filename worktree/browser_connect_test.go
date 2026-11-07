package worktree

import (
	"os"
	"testing"

	"github.com/joshmedeski/sesh/v2/browser"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/jj"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConnectFromBrowserIssue(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/2345")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	// Browser resolves the issue URL for the configured repo.
	mBrowser.EXPECT().ActiveTabURL().
		Return("https://github.com/nutiliti/nutiliti/issues/2345", true, nil)

	// Not a PR => issue path.
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 2345).Return(forge.PullRequest{}, false, nil)
	mGh.EXPECT().WorktreeName(2345, false).Return("2345")

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{FromBrowser: true, Switch: true})
	require.NoError(t, err)
}

func TestConnectFromBrowserRepoOverrideWins(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/2345")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	// Browser tab points at a DIFFERENT repo than the caller's --repo.
	mBrowser.EXPECT().ActiveTabURL().
		Return("https://github.com/someone/other/issues/2345", true, nil)

	// The caller's --repo ("nutiliti/nutiliti") must win over the URL's repo
	// ("someone/other"); only the number (2345) comes from the URL.
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 2345).Return(forge.PullRequest{}, false, nil)
	mGh.EXPECT().WorktreeName(2345, false).Return("2345")

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{FromBrowser: true, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectFromBrowserPR(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/678")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	// Browser resolves a PR URL for the configured repo.
	mBrowser.EXPECT().ActiveTabURL().
		Return("https://github.com/nutiliti/nutiliti/pull/678", true, nil)

	// PR authored by someone else => foreign PR, detached checkout keyed on
	// the PR number.
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 678).
		Return(forge.PullRequest{Author: "octocat"}, true, nil)
	mGh.EXPECT().CurrentUser("github.com").Return("me", nil)
	mGh.EXPECT().WorktreeName(678, true).Return("678")

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAddDetached("/repo", wPath, "origin/main").Return("", nil)
	mGh.EXPECT().PrCheckout(wPath, "github.com", "nutiliti/nutiliti", 678).Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{FromBrowser: true, Switch: true})
	require.NoError(t, err)
}

func TestConnectFromBrowserUnparseableURL(t *testing.T) {
	mBrowser := browser.NewMockBrowser(t)
	mBrowser.EXPECT().ActiveTabURL().Return("https://example.com/", true, nil)

	w := NewWorktree(
		nuConfig(),
		git.NewMockGit(t), jj.NewMockJj(t), forgeResolver(forge.NewMockForge(t)), connector.NewMockConnector(t),
		mBrowser, home.NewHome(oswrap.NewMockOs(t)), oswrap.NewMockOs(t), pathwrap.NewPath(), testIssueCache(t),
	)
	_, err := w.Connect(model.WorktreeConnectOpts{FromBrowser: true})
	require.Error(t, err)
}

func TestConnectFromBrowserUnavailable(t *testing.T) {
	mBrowser := browser.NewMockBrowser(t)
	// Non-macOS / no application configured => skipped.
	mBrowser.EXPECT().ActiveTabURL().Return("", false, nil)

	w := NewWorktree(
		nuConfig(),
		git.NewMockGit(t), jj.NewMockJj(t), forgeResolver(forge.NewMockForge(t)), connector.NewMockConnector(t),
		mBrowser, home.NewHome(oswrap.NewMockOs(t)), oswrap.NewMockOs(t), pathwrap.NewPath(), testIssueCache(t),
	)
	_, err := w.Connect(model.WorktreeConnectOpts{FromBrowser: true})
	require.Error(t, err)
}
