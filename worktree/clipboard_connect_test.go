package worktree

import (
	"errors"
	"os"
	"testing"

	"github.com/joshmedeski/sesh/v2/browser"
	"github.com/joshmedeski/sesh/v2/clipboard"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/jj"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConnectFromClipboardIssueURL(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mClipboard := clipboard.NewMockClipboard(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/2345")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	mClipboard.EXPECT().Read().Return("https://github.com/nutiliti/nutiliti/issues/2345", nil)

	// Not a PR => issue path.
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, mClipboard, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{FromClipboard: true, Switch: true})
	require.NoError(t, err)
}

func TestConnectFromClipboardPRURL(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mClipboard := clipboard.NewMockClipboard(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/678")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	mClipboard.EXPECT().Read().Return("https://github.com/nutiliti/nutiliti/pull/678", nil)

	// PR authored by someone else => foreign PR, detached checkout keyed on the PR number.
	mGh.EXPECT().PrView("nutiliti/nutiliti", 678).Return(github.PullRequest{Author: "octocat"}, true, nil)
	mGh.EXPECT().CurrentUser().Return("me", nil)

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAddDetached("/repo", wPath, "origin/main").Return("", nil)
	mGh.EXPECT().PrCheckout(wPath, "nutiliti/nutiliti", 678).Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, mClipboard, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{FromClipboard: true, Switch: true})
	require.NoError(t, err)
}

func TestConnectFromClipboardBareNumber(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mClipboard := clipboard.NewMockClipboard(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/2345")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	mClipboard.EXPECT().Read().Return("2345", nil)

	// Not a PR => issue path.
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, mClipboard, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{
		FromClipboard: true,
		Repo:          "nutiliti/nutiliti",
		Switch:        true,
	})
	require.NoError(t, err)
}

// A clipboard URL naming a different repo than --repo: the override wins and
// only the number comes from the clipboard.
func TestConnectFromClipboardRepoOverrideWins(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mClipboard := clipboard.NewMockClipboard(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/2345")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	// The clipboard URL names a DIFFERENT repo than the caller's --repo. The
	// override must win; only the number (2345) comes from the clipboard.
	mClipboard.EXPECT().Read().Return("https://github.com/someone/other/issues/2345", nil)

	// Not a PR => issue path.
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, mClipboard, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{
		FromClipboard: true,
		Repo:          "nutiliti/nutiliti",
		Switch:        true,
	})
	require.NoError(t, err)
}

// A bare number with no --repo falls back to cwd detection, the same route a
// number on the command line takes.
func TestConnectFromClipboardBareNumberWithoutRepo(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mClipboard := clipboard.NewMockClipboard(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/2345")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	mClipboard.EXPECT().Read().Return("2345", nil)

	// No --repo, so the config is detected from cwd.
	mOs.EXPECT().Getwd().Return(wPath, nil)
	mJj.EXPECT().Roots(wPath).Return(false, "", "")
	mGit.EXPECT().WorktreeList(wPath).Return(
		true,
		"worktree /repo\nHEAD abc\nbranch refs/heads/main\n\nworktree /repo/w/2345\nHEAD def\nbranch refs/heads/jam/2345-1\n",
		nil)

	// Not a PR => issue path.
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, mClipboard, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{FromClipboard: true, Switch: true})
	require.NoError(t, err)
}

func TestConnectFromClipboardUnparseableURL(t *testing.T) {
	mClipboard := clipboard.NewMockClipboard(t)
	mClipboard.EXPECT().Read().Return("https://example.com/", nil)

	w := NewWorktree(
		nuConfig(),
		git.NewMockGit(t), jj.NewMockJj(t), github.NewMockGithub(t), connector.NewMockConnector(t),
		browser.NewMockBrowser(t), mClipboard, home.NewHome(oswrap.NewMockOs(t)), oswrap.NewMockOs(t),
		pathwrap.NewPath(), testIssueCache(t),
	)
	_, err := w.Connect(model.WorktreeConnectOpts{FromClipboard: true})
	require.EqualError(t, err, "could not parse a GitHub issue/PR from \"https://example.com/\"")
}

// A clipboard that cannot be read (e.g., no X display, no Wayland data-control)
// surfaces its error rather than being treated as an empty reference.
func TestConnectFromClipboardReadError(t *testing.T) {
	mClipboard := clipboard.NewMockClipboard(t)
	mClipboard.EXPECT().Read().Return("", errors.New("clipboard unavailable"))

	w := NewWorktree(
		nuConfig(),
		git.NewMockGit(t), jj.NewMockJj(t), github.NewMockGithub(t), connector.NewMockConnector(t),
		browser.NewMockBrowser(t), mClipboard, home.NewHome(oswrap.NewMockOs(t)), oswrap.NewMockOs(t),
		pathwrap.NewPath(), testIssueCache(t),
	)
	_, err := w.Connect(model.WorktreeConnectOpts{FromClipboard: true})
	require.EqualError(t, err, "clipboard unavailable")
}
