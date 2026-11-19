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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func falsePtr() *bool { b := false; return &b }

func forgeResolver(f forge.Forge) forge.Resolver {
	return forge.NewResolver(
		func(cfg model.WorktreeConfig) (forge.Forge, error) {
			return f, nil
		},
		func(path string) (forge.Forge, error) {
			return f, nil
		},
	)
}

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
	mGh := forge.NewMockForge(t)
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
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 2345).Return(forge.PullRequest{}, false, nil)
	mGh.EXPECT().WorktreeName(2345, false).Return("2345")

	// target does not exist yet
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)

	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)

	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectOwnPrWithClosingIssue(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
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
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 100).
		Return(forge.PullRequest{Author: "me", ClosingIssues: []int{42}}, true, nil)
	mGh.EXPECT().CurrentUser("github.com").Return("me", nil)
	mGh.EXPECT().WorktreeName(42, false).Return("42")

	// Keyed on the issue (42), normal branch create
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/42-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 100, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectForeignPrDetachCheckout(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
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
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 200).
		Return(forge.PullRequest{Author: "octocat"}, true, nil)
	mGh.EXPECT().CurrentUser("github.com").Return("me", nil)
	mGh.EXPECT().WorktreeName(200, true).Return("200")

	// Keyed on PR number, detached worktree + gh pr checkout
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAddDetached("/repo", wPath, "origin/main").Return("", nil)
	mGh.EXPECT().PrCheckout(wPath, "github.com", "nutiliti/nutiliti", 200).Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 200, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectResolvesConfigFromCwd(t *testing.T) {
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

	// No --repo override: resolveConfig detects the repo from cwd via
	// `git worktree list --porcelain` on the main worktree.
	mOs.EXPECT().Getwd().Return(wPath, nil)
	mJj.EXPECT().Roots(wPath).Return(false, "", "")
	mGit.EXPECT().WorktreeList(wPath).Return(true,
		"worktree /repo\nHEAD abc\nbranch refs/heads/main\n\n"+
			"worktree /repo/w/2345\nHEAD def\nbranch refs/heads/jam/2345-1\n",
		nil)

	// Not a PR => issue path
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 2345).Return(forge.PullRequest{}, false, nil)
	mGh.EXPECT().WorktreeName(2345, false).Return("2345")

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Switch: true})
	require.NoError(t, err)
}

func TestConnectOwnPrFirstIssueRefFallback(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
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
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 100).
		Return(forge.PullRequest{Author: "me", Body: "fixes #7"}, true, nil)
	mGh.EXPECT().CurrentUser("github.com").Return("me", nil)
	mGh.EXPECT().WorktreeName(7, false).Return("7")

	// Keyed on issue 7 scanned from title+body, normal branch create (not detached/PrCheckout)
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/7-1", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 100, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectExistingWorktreeRunsStartupCommand(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	wPath := p.FromSlash("/repo/w/2345")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 2345).Return(forge.PullRequest{}, false, nil)
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()
	mGh.EXPECT().WorktreeName(2345, false).Return("2345")

	// The worktree is already there, so nothing is created — and the command is
	// the startup one, not the setup that already ran when it was created.
	mOs.EXPECT().Stat(wPath).Return(nil, nil)
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_setup"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectExistingPrWorktreeRunsStartupCommand(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
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
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 200).
		Return(forge.PullRequest{Author: "someone-else"}, true, nil)
	mGh.EXPECT().CurrentUser("github.com").Return("me", nil)
	mGh.EXPECT().WorktreeName(200, true).Return("200")
	mOs.EXPECT().Stat(wPath).Return(nil, nil)
	mGh.EXPECT().PrCheckout(wPath, "github.com", "nutiliti/nutiliti", 200).Return("", nil)
	mGit.EXPECT().Pull(wPath).Return("", nil)
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_setup"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 200, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectCreateWithoutCreateCommandRunsNothing(t *testing.T) {
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

	cfg := nuConfig()
	cfg.WorktreeConfigs[0].CreateCommand = ""

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 2345).Return(forge.PullRequest{}, false, nil)
	mGh.EXPECT().WorktreeName(2345, false).Return("2345")
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)

	// startup_command does not stand in for an unset create_command: creating is
	// the one time it is deliberately not run.
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Switch: true, Command: ""}).
		Return("", nil)

	w := NewWorktree(cfg, mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectNoConfigMatch(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	w := NewWorktree(nuConfig(), mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 1, Repo: "unknown/repo"})
	require.Error(t, err)
}

type dirInfo struct{ os.FileInfo }

func (dirInfo) IsDir() bool { return true }

type jjFixture struct {
	jj   *jj.MockJj
	gh   *forge.MockForge
	conn *connector.MockConnector
	os   *oswrap.MockOs
	w    Worktree
}

func newJjFixture(t *testing.T, targetExists bool) jjFixture {
	p := pathwrap.NewPath()
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
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
	w := NewWorktree(nuConfig(), git.NewMockGit(t), mJj, forgeResolver(mGh), mConn, browser.NewMockBrowser(t),
		home.NewHome(mOs), mOs, p, testIssueCache(t))
	return jjFixture{mJj, mGh, mConn, mOs, w}
}

func TestConnectJjIssueAddsWorkspaceAndBookmark(t *testing.T) {
	wPath := pathwrap.NewPath().FromSlash("/repo/w/2345")
	f := newJjFixture(t, false)
	f.gh.EXPECT().WorktreeName(2345, false).Return("2345")
	f.gh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 2345).Return(forge.PullRequest{}, false, nil)
	f.jj.EXPECT().WorkspaceAdd("/repo", "2345", `"main"@"origin"`, wPath).Return("", nil)
	f.jj.EXPECT().BookmarkCreate(wPath, "jam/2345-1", "@").Return("", nil)
	f.conn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	_, err := f.w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectJjForkPrAddsRemoteAndFetchesBranch(t *testing.T) {
	wPath := pathwrap.NewPath().FromSlash("/repo/w/200")
	f := newJjFixture(t, false)
	f.gh.EXPECT().WorktreeName(200, true).Return("200")
	f.gh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 200).Return(forge.PullRequest{
		Author: "octocat", HeadRefName: "fix-a", HeadOwner: "octocat", HeadRepo: "nutiliti", IsFork: true,
	}, true, nil)
	f.gh.EXPECT().CurrentUser("github.com").Return("me", nil)
	f.jj.EXPECT().RemoteAdd("/repo", "octocat", "https://github.com/octocat/nutiliti.git").Return("", nil)
	f.jj.EXPECT().Fetch("/repo", "--remote", "octocat", "--branch", "fix-a").Return("", nil)
	f.jj.EXPECT().WorkspaceAdd("/repo", "200", `"fix-a"@"octocat"`, wPath).Return("", nil)
	f.conn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "nu_install"}).Return("", nil)

	_, err := f.w.Connect(model.WorktreeConnectOpts{Number: 200, Repo: "nutiliti/nutiliti", Switch: true})
	require.NoError(t, err)
}

func TestConnectJjForkPrUsesTheConfiguredHost(t *testing.T) {
	p := pathwrap.NewPath()
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
	mConn := connector.NewMockConnector(t)
	mOs := oswrap.NewMockOs(t)

	cfg := model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{
			Repo:           "enterprise/repo",
			Host:           "github.example.com",
			Path:           "/enterprise-repo",
			WorktreeDir:    ".wk",
			BranchTemplate: "feature/{number}",
			BaseBranch:     "origin/main",
			Fetch:          falsePtr(),
			CreateCommand:  "init_script",
		}},
	}

	wPath := p.FromSlash("/enterprise-repo/.wk/200")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/enterprise-repo").Return("/enterprise-repo").Maybe()
	mOs.EXPECT().Stat(p.FromSlash("/enterprise-repo/.jj")).Return(dirInfo{}, nil)
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(p.FromSlash("/enterprise-repo/.wk"), mock.Anything).Return(nil)

	mGh.EXPECT().WorktreeName(200, true).Return("200")
	mGh.EXPECT().PrView("github.example.com", "enterprise/repo", 200).Return(forge.PullRequest{
		Author: "octocat", HeadRefName: "fix-a", HeadOwner: "octocat", HeadRepo: "enterprise", IsFork: true,
	}, true, nil)
	mGh.EXPECT().CurrentUser("github.example.com").Return("me", nil)

	mJj.EXPECT().RemoteAdd("/enterprise-repo", "octocat", "https://github.example.com/octocat/enterprise.git").Return("", nil)
	mJj.EXPECT().Fetch("/enterprise-repo", "--remote", "octocat", "--branch", "fix-a").Return("", nil)
	mJj.EXPECT().WorkspaceAdd("/enterprise-repo", "200", `"fix-a"@"octocat"`, wPath).Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "init_script"}).Return("", nil)

	w := NewWorktree(cfg, git.NewMockGit(t), mJj, forgeResolver(mGh), mConn, browser.NewMockBrowser(t),
		home.NewHome(mOs), mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 200, Repo: "enterprise/repo", Switch: true})
	require.NoError(t, err)
}

func TestConnectJjExistingPrMovesToUpdatedHead(t *testing.T) {
	wPath := pathwrap.NewPath().FromSlash("/repo/w/200")
	f := newJjFixture(t, true)
	f.gh.EXPECT().WorktreeName(200, true).Return("200")
	f.gh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 200).Return(forge.PullRequest{
		Author: "octocat", HeadRefName: "fix-a", HeadOwner: "nutiliti", HeadRepo: "nutiliti",
	}, true, nil)
	f.gh.EXPECT().CurrentUser("github.com").Return("me", nil)
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
		mGh := forge.NewMockForge(t)
		mConn := connector.NewMockConnector(t)
		mBrowser := browser.NewMockBrowser(t)
		mOs := oswrap.NewMockOs(t)
		h := home.NewHome(mOs)
		p := pathwrap.NewPath()

		wPath := p.FromSlash("/repo/w/2345")

		mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
		mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
		mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 2345).Return(forge.PullRequest{}, false, nil)
		mGh.EXPECT().WorktreeName(2345, false).Return("2345")
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

		w := NewWorktree(nuConfig(), mGit, jj.NewMockJj(t), forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
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
	mGh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 2345).Return(forge.PullRequest{}, false, nil)
	mGh.EXPECT().WorktreeName(2345, false).Return("2345")
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/repo/.jj")).Return(nil, os.ErrNotExist).Maybe()
	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", wPath, "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().
		Connect(wPath, model.ConnectOpts{Command: "nu_install", Background: true}).
		Return("nutiliti/2345", nil)

	w := NewWorktree(nuConfig(), mGit, jj.NewMockJj(t), forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 2345, Repo: "nutiliti/nutiliti", Background: true})
	require.NoError(t, err)
}

func TestConnectGitHubCustomHost(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGh := forge.NewMockForge(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	cfg := model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{
			Repo:           "enterprise/repo",
			Host:           "github.example.com",
			Path:           "/enterprise-repo",
			WorktreeDir:    ".wk",
			BranchTemplate: "feature/{number}",
			BaseBranch:     "origin/main",
			Fetch:          falsePtr(),
			CreateCommand:  "init_script",
		}},
	}

	wRoot := p.FromSlash("/enterprise-repo/.wk")
	wPath := p.FromSlash("/enterprise-repo/.wk/55")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/enterprise-repo").Return("/enterprise-repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/enterprise-repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	mGh.EXPECT().PrView("github.example.com", "enterprise/repo", 55).Return(forge.PullRequest{}, false, nil)
	mGh.EXPECT().WorktreeName(55, false).Return("55")

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/enterprise-repo", wPath, "feature/55", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "init_script"}).Return("", nil)

	w := NewWorktree(cfg, mGit, mJj, forgeResolver(mGh), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 55, Repo: "enterprise/repo", Switch: true})
	require.NoError(t, err)
}

func TestConnectGitLabCustomHost(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGl := forge.NewMockForge(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	cfg := model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{
			Repo:           "internal/repo",
			Forge:          "gitlab",
			Host:           "gitlab.example.com",
			Path:           "/internal-repo",
			WorktreeDir:    ".wk",
			BranchTemplate: "feature/{number}",
			BaseBranch:     "origin/main",
			Fetch:          falsePtr(),
			CreateCommand:  "make setup",
		}},
	}

	wRoot := p.FromSlash("/internal-repo/.wk")
	wPath := p.FromSlash("/internal-repo/.wk/99")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/internal-repo").Return("/internal-repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/internal-repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	mGl.EXPECT().PrView("gitlab.example.com", "internal/repo", 99).Return(forge.PullRequest{}, false, nil)
	mGl.EXPECT().WorktreeName(99, false).Return("99")

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/internal-repo", wPath, "feature/99", "origin/main").Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "make setup"}).Return("", nil)

	w := NewWorktree(cfg, mGit, mJj, forgeResolver(mGl), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 99, Repo: "internal/repo", Switch: true})
	require.NoError(t, err)
}

func TestConnectGitLabMergeRequest(t *testing.T) {
	mGit := git.NewMockGit(t)
	mJj := jj.NewMockJj(t)
	mGl := forge.NewMockForge(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	cfg := model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{
			Repo:          "inkscape/inkscape",
			Forge:         "gitlab",
			Path:          "/gitlab-repo",
			WorktreeDir:   ".wk",
			BaseBranch:    "origin/master",
			Fetch:         falsePtr(),
			CreateCommand: "gl_setup",
		}},
	}

	wRoot := p.FromSlash("/gitlab-repo/.wk")
	wPath := p.FromSlash("/gitlab-repo/.wk/mr-101")

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/gitlab-repo").Return("/gitlab-repo").Maybe()
	mOs.EXPECT().Stat(pathwrap.NewPath().FromSlash("/gitlab-repo/.jj")).Return(nil, os.ErrNotExist).Maybe()

	mGl.EXPECT().PrView("gitlab.com", "inkscape/inkscape", 101).
		Return(forge.PullRequest{Author: "someone-else"}, true, nil)
	mGl.EXPECT().CurrentUser("gitlab.com").Return("me", nil)
	mGl.EXPECT().WorktreeName(101, true).Return("mr-101")

	mOs.EXPECT().Stat(wPath).Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll(wRoot, mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAddDetached("/gitlab-repo", wPath, "origin/master").Return("", nil)
	mGl.EXPECT().PrCheckout(wPath, "gitlab.com", "inkscape/inkscape", 101).Return("", nil)
	mConn.EXPECT().Connect(wPath, model.ConnectOpts{Switch: true, Command: "gl_setup"}).Return("", nil)

	w := NewWorktree(cfg, mGit, mJj, forgeResolver(mGl), mConn, mBrowser, h, mOs, p, testIssueCache(t))
	_, err := w.Connect(model.WorktreeConnectOpts{Number: 101, Repo: "inkscape/inkscape", Switch: true})
	require.NoError(t, err)
}
