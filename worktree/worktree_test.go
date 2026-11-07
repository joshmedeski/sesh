package worktree

import (
	"os"
	"testing"

	"github.com/joshmedeski/sesh/v2/browser"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/jj"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
			// Deliberately different values: every Connect test asserts which of
			// the two the connector was handed, which a shared value would hide.
			StartupCommand: "nu_setup",
			CreateCommand:  "nu_install",
		}},
	}
}

func TestConnectIssueWithRepoOverride(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/2345")

	// path expansion: "/repo" has no ~, but home.ExpandPath calls os.UserHomeDir
	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	// Not a PR => issue path
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	// target does not exist yet
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)

	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)

	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectOwnPrWithClosingIssue(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/42")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	// It is a PR authored by me, closing issue 42
	mGh.EXPECT().PrView("nutiliti/nutiliti", 100).
		Return(github.PullRequest{Author: "me", ClosingIssues: []int{42}}, true, nil)
	mGh.EXPECT().CurrentUser().Return("me", nil)

	// Keyed on the issue (42), normal branch create
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/42-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 100, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectForeignPrDetachCheckout(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/200")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	// PR authored by someone else
	mGh.EXPECT().PrView("nutiliti/nutiliti", 200).
		Return(github.PullRequest{Author: "octocat"}, true, nil)
	mGh.EXPECT().CurrentUser().Return("me", nil)

	// Keyed on PR number, detached worktree + gh pr checkout
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAddDetached("/repo", wPath, "origin/main").Return("", nil)
	mGh.EXPECT().PrCheckout(wPath, "nutiliti/nutiliti", 200).Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 200, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectResolvesConfigFromCwd(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
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

	// No --repo override: resolveConfig detects the repo from cwd via
	// `git worktree list --porcelain` on the main worktree.
	mOs.EXPECT().Getwd().Return(wPath, nil)
	mJj.EXPECT().Roots(wPath).Return(false, "", "")
	mGit.EXPECT().WorktreeList(wPath).Return(true,
		"worktree /repo\nHEAD abc\nbranch refs/heads/main\n\n"+
			"worktree /repo/w/2345\nHEAD def\nbranch refs/heads/jam/2345-1\n",
		nil)

	// Not a PR => issue path
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Switch: true})
	require.NoError(t, err)
}

func TestConnectOwnPrFirstIssueRefFallback(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/7")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	// PR authored by me, no ClosingIssues, but body references "#7"
	mGh.EXPECT().PrView("nutiliti/nutiliti", 100).
		Return(github.PullRequest{Author: "me", Body: "fixes #7"}, true, nil)
	mGh.EXPECT().CurrentUser().Return("me", nil)

	// Keyed on issue 7 scanned from title+body, normal branch create (not detached/PrCheckout)
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/7-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 100, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectExistingWorktreeRunsStartupCommand(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wPath := p.FromSlash("/repo/w/2345")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	// The worktree is already there, so nothing is created — and the command is
	// the startup one, not the setup that already ran when it was created.
	mOs.EXPECT().Stat(wPath).Return(nil, nil)
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_setup"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectExistingPrWorktreeRunsStartupCommand(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wPath := p.FromSlash("/repo/w/200")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	// Foreign PR whose worktree already exists: the checkout is refreshed, which
	// is not the same thing as creating it.
	mGh.EXPECT().PrView("nutiliti/nutiliti", 200).
		Return(github.PullRequest{Author: "someone-else"}, true, nil)
	mGh.EXPECT().CurrentUser().Return("me", nil)
	mOs.EXPECT().Stat(wPath).Return(nil, nil)
	mGh.EXPECT().PrCheckout(wPath, "nutiliti/nutiliti", 200).Return("", nil)
	mGit.EXPECT().Pull(wPath).Return("", nil)
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_setup"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 200, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectCreateWithoutCreateCommandRunsNothing(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/2345")

	cfg := nuConfig()
	cfg.WorktreeConfigs[0].CreateCommand = ""

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)

	// startup_command does not stand in for an unset create_command: creating is
	// the one time it is deliberately not run.
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: ""}).
		Return("", nil)

	w := NewWorktree(cfg, mGit, mJj, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectNoConfigMatch(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	w := NewWorktree(nuConfig(), mGit, mJj, mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 1, Repo: "unknown/repo"})
	require.Error(t, err)
}

type dirInfo struct{ os.FileInfo }

func (dirInfo) IsDir() bool { return true }

type jjFixture struct {
	jj   *jj.MockJj
	gh   *github.MockGithub
	conn *connector.MockConnector
	os   *oswrap.MockOs
	w    Worktree
}

func newJjFixture(t *testing.T, targetExists bool) jjFixture {
	p := pathwrap.NewPath()
	mJj := jj.NewMockJj(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mOs := oswrap.NewMockOs(t)
	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(p.FromSlash("/repo/.jj")).Return(dirInfo{}, nil)
	mOs.EXPECT().MkdirAll(p.FromSlash("/repo/w"), mock.Anything).Return(nil).Maybe()
	mOs.EXPECT().Stat(mock.Anything).RunAndReturn(func(string) (os.FileInfo, error) {
		if targetExists {
			return dirInfo{}, nil
		}
		return nil, os.ErrNotExist
	})
	w := NewWorktree(nuConfig(), git.NewMockGit(t), mJj, mGh, mConn, browser.NewMockBrowser(t),
		home.NewHome(mOs), mOs, p, testIssueCache(t))
	return jjFixture{mJj, mGh, mConn, mOs, w}
}

func TestConnectJjIssueAddsWorkspaceAndBookmark(t *testing.T) {
	wPath := pathwrap.NewPath().FromSlash("/repo/w/2345")
	f := newJjFixture(t, false)
	f.gh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)
	f.jj.EXPECT().WorkspaceAdd("/repo", "2345", `"main"@"origin"`, wPath).Return("", nil)
	f.jj.EXPECT().BookmarkCreate(wPath, "jam/2345-1", "@").Return("", nil)
	f.conn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	_, err := f.w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectJjForkPrAddsRemoteAndFetchesBranch(t *testing.T) {
	wPath := pathwrap.NewPath().FromSlash("/repo/w/200")
	f := newJjFixture(t, false)
	f.gh.EXPECT().PrView("nutiliti/nutiliti", 200).Return(github.PullRequest{
		Author: "octocat", HeadRefName: "fix-a", HeadOwner: "octocat", HeadRepo: "nutiliti", IsFork: true,
	}, true, nil)
	f.gh.EXPECT().CurrentUser().Return("me", nil)
	f.jj.EXPECT().RemoteAdd("/repo", "octocat", "https://github.com/octocat/nutiliti.git").Return("", nil)
	f.jj.EXPECT().Fetch("/repo", "--remote", "octocat", "--branch", "fix-a").Return("", nil)
	f.jj.EXPECT().WorkspaceAdd("/repo", "200", `"fix-a"@"octocat"`, wPath).Return("", nil)
	f.conn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	_, err := f.w.Connect(model.WorktreeConnectOpts{Number: 200, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectJjExistingPrMovesToUpdatedHead(t *testing.T) {
	wPath := pathwrap.NewPath().FromSlash("/repo/w/200")
	f := newJjFixture(t, true)
	f.gh.EXPECT().PrView("nutiliti/nutiliti", 200).Return(github.PullRequest{
		Author: "octocat", HeadRefName: "fix-a", HeadOwner: "nutiliti", HeadRepo: "nutiliti",
	}, true, nil)
	f.gh.EXPECT().CurrentUser().Return("me", nil)
	f.jj.EXPECT().Fetch("/repo", "--remote", "origin", "--branch", "fix-a").Return("", nil)
	f.jj.EXPECT().New(wPath, `"fix-a"@"origin"`).Return("", nil)
	f.conn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_setup"}).Return("", nil)

	_, err := f.w.Connect(model.WorktreeConnectOpts{Number: 200, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestJjRevset(t *testing.T) {
	assert.Equal(t, `"main"@"origin"`, jjRevset("origin/main"))
	assert.Equal(t, `"feat/x"@"upstream"`, jjRevset("upstream/feat/x"))
	assert.Equal(t, "main@origin", jjRevset("main@origin"))
	assert.Equal(t, "trunk()", jjRevset("trunk()"))
	assert.Equal(t, "main", jjRevset("main"))
}

func TestConnectWithCommandOpensNewWindow(t *testing.T) {
	for _, background := range []bool{false, true} {
		mGit := git.NewMockGit(t)
		mGh := github.NewMockGithub(t)
		mConn := connector.NewMockConnector(t)
		mBrowser := browser.NewMockBrowser(t)
		mOs := oswrap.NewMockOs(t)
		h := home.NewHome(mOs)
		p := pathwrap.NewPath()

		wPath := p.FromSlash("/repo/w/2345")

		mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
		mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
		mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)
		mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()
		mOs.EXPECT().Stat(wPath).Return(nil, nil)

		mConn.EXPECT().
			Connect(wPath, model.ConnectOpts{Command: "nu_setup", Background: true}).
			Return("nutiliti/2345", nil)
		mConn.EXPECT().
			ConnectWindow(model.WindowConnectOpts{
				Name:       "claude",
				New:        true,
				Session:    "nutiliti/2345",
				Command:    `/usr/bin/claude "Implement #2345"`,
				Background: background,
			}).
			Return("nutiliti/2345:3", nil)

		w := NewWorktree(nuConfig(), mGit, jj.NewMockJj(t), mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
		_, err := w.Connect(model.WorktreeConnectOpts{
			Number:     2345,
			Repo:       "nutiliti/nutiliti",
			Command:    `/usr/bin/claude "Implement #2345"`,
			Background: background,
		})
		require.NoError(t, err)
	}
}

func TestConnectBackgroundWithoutCommand(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wRoot := p.FromSlash("/repo/w")
	wPath := p.FromSlash("/repo/w/2345")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Command: "nu_install", Background: true}).
		Return("nutiliti/2345", nil)

	w := NewWorktree(nuConfig(), mGit, jj.NewMockJj(t), mGh, mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Background: true})
	require.NoError(t, err)
}
