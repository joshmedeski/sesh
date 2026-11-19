package worktree

import (
	"os"
	"slices"
	"strconv"
	"testing"

	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type browseFixture struct {
	w  *RealWorktree
	gh *forge.MockForge
	gl *forge.MockForge
}

func newBrowseFixture(t *testing.T, jjDirs ...string) browseFixture {
	mGh := forge.NewMockForge(t)
	mGl := forge.NewMockForge(t)
	mOs := oswrap.NewMockOs(t)
	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().ExpandEnv("/dotfiles").Return("/dotfiles").Maybe()
	mOs.EXPECT().ExpandEnv("/gitlab-repo").Return("/gitlab-repo").Maybe()
	mOs.EXPECT().ExpandEnv("/internal-repo").Return("/internal-repo").Maybe()
	mOs.EXPECT().ExpandEnv("/enterprise-repo").Return("/enterprise-repo").Maybe()
	mOs.EXPECT().Stat(mock.Anything).RunAndReturn(func(name string) (os.FileInfo, error) {
		if slices.Contains(jjDirs, name) {
			return dirInfo{}, nil
		}
		return nil, os.ErrNotExist
	}).Maybe()
	cfg := nuConfig()
	cfg.WorktreeConfigs = append([]model.WorktreeConfig{
		{Repo: "joshmedeski/dotfiles", Path: "/dotfiles"},
		{Repo: "enterprise/repo", Host: "github.example.com", Path: "/enterprise-repo"},
		{Repo: "inkscape/inkscape", Forge: "gitlab", Path: "/gitlab-repo"},
		{Repo: "internal/repo", Forge: "gitlab", Host: "gitlab.example.com", Path: "/internal-repo"},
	}, cfg.WorktreeConfigs...)

	fr := forge.NewResolver(
		func(c model.WorktreeConfig) (forge.Forge, error) {
			if c.Forge == model.ForgeGitlab {
				return mGl, nil
			}
			return mGh, nil
		},
		func(path string) (forge.Forge, error) {
			return mGh, nil
		},
	)

	w := &RealWorktree{config: cfg, forges: fr, home: home.NewHome(mOs), os: mOs, path: pathwrap.NewPath()}
	return browseFixture{w: w, gh: mGh, gl: mGl}
}

func TestBrowseURLOpensTheWorktreeIssue(t *testing.T) {
	f := newBrowseFixture(t)
	f.gh.EXPECT().WorktreeNumber("2345").Return(2345, false, nil).Twice()
	f.gh.EXPECT().IssueURL("github.com", "nutiliti/nutiliti", 2345, false).
		Return("https://github.com/nutiliti/nutiliti/issues/2345").Twice()

	for _, path := range []string{"/repo/w/2345", "/repo/w/2345/apps/web"} {
		url, err := f.w.BrowseURL(path, false)
		require.NoError(t, err)
		assert.Equal(t, "https://github.com/nutiliti/nutiliti/issues/2345", url)
	}
}

func TestBrowseURLDefaultWorktreeDir(t *testing.T) {
	f := newBrowseFixture(t)
	f.gh.EXPECT().WorktreeNumber("12").Return(12, false, nil)
	f.gh.EXPECT().IssueURL("github.com", "joshmedeski/dotfiles", 12, false).
		Return("https://github.com/joshmedeski/dotfiles/issues/12")

	url, err := f.w.BrowseURL("/dotfiles/.wk/12", false)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/joshmedeski/dotfiles/issues/12", url)
}

func TestBrowseURLGitLab(t *testing.T) {
	f := newBrowseFixture(t)
	f.gl.EXPECT().WorktreeNumber("42").Return(42, false, nil)
	f.gl.EXPECT().IssueURL("gitlab.com", "inkscape/inkscape", 42, false).
		Return("https://gitlab.com/inkscape/inkscape/-/work_items/42")

	glURL, err := f.w.BrowseURL("/gitlab-repo/.wk/42", false)
	require.NoError(t, err)
	assert.Equal(t, "https://gitlab.com/inkscape/inkscape/-/work_items/42", glURL)
}

func TestBrowseURLGitLabMergeRequest(t *testing.T) {
	f := newBrowseFixture(t)
	f.gl.EXPECT().WorktreeNumber("mr-42").Return(42, true, nil)
	f.gl.EXPECT().IssueURL("gitlab.com", "inkscape/inkscape", 42, true).
		Return("https://gitlab.com/inkscape/inkscape/-/merge_requests/42")

	mrURL, err := f.w.BrowseURL("/gitlab-repo/.wk/mr-42", false)
	require.NoError(t, err)
	assert.Equal(t, "https://gitlab.com/inkscape/inkscape/-/merge_requests/42", mrURL)
}

func TestBrowseURLGitLabCustomHost(t *testing.T) {
	f := newBrowseFixture(t)
	f.gl.EXPECT().WorktreeNumber("99").Return(99, false, nil)
	f.gl.EXPECT().IssueURL("gitlab.example.com", "internal/repo", 99, false).
		Return("https://gitlab.example.com/internal/repo/-/work_items/99")

	customURL, err := f.w.BrowseURL("/internal-repo/.wk/99", false)
	require.NoError(t, err)
	assert.Equal(t, "https://gitlab.example.com/internal/repo/-/work_items/99", customURL)
}

func TestBrowseURLGitHubCustomHost(t *testing.T) {
	f := newBrowseFixture(t)
	f.gh.EXPECT().WorktreeNumber("77").Return(77, false, nil)
	f.gh.EXPECT().IssueURL("github.example.com", "enterprise/repo", 77, false).
		Return("https://github.example.com/enterprise/repo/issues/77")

	customURL, err := f.w.BrowseURL("/enterprise-repo/.wk/77", false)
	require.NoError(t, err)
	assert.Equal(t, "https://github.example.com/enterprise/repo/issues/77", customURL)
}

func TestBrowseURLRejectsPathsOutsideWorktrees(t *testing.T) {
	f := newBrowseFixture(t)
	f.gh.EXPECT().WorktreeNumber("notes").Return(0, false, strconv.ErrSyntax)

	for _, path := range []string{"/repo", "/repo/w", "/repo/w/notes", "/repo/wk/2345", "/elsewhere"} {
		_, err := f.w.BrowseURL(path, false)
		assert.ErrorContains(t, err, "is not a worktree", path)
	}
}

func TestBrowseURLRejectsEmptySessionPath(t *testing.T) {
	f := newBrowseFixture(t)
	_, err := f.w.BrowseURL("", false)
	assert.ErrorContains(t, err, "no focused tmux session")
}

func TestBrowseURLPrUsesTheBranchPullRequest(t *testing.T) {
	p := pathwrap.NewPath()

	f := newBrowseFixture(t)
	f.gh.EXPECT().WorktreeNumber("2345").Return(2345, false, nil)
	f.gh.EXPECT().PrURL(p.FromSlash("/repo/w/2345")).Return("https://github.com/nutiliti/nutiliti/pull/2400", true, nil)
	url, err := f.w.BrowseURL("/repo/w/2345/apps", true)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/nutiliti/nutiliti/pull/2400", url)
}

func TestBrowseURLPrWithoutPullRequest(t *testing.T) {
	p := pathwrap.NewPath()

	f := newBrowseFixture(t)
	f.gh.EXPECT().WorktreeNumber("2345").Return(2345, false, nil)
	f.gh.EXPECT().PrURL(p.FromSlash("/repo/w/2345")).Return("", false, nil)
	_, err := f.w.BrowseURL("/repo/w/2345", true)
	assert.EqualError(t, err, "no pull request found for worktree 2345")
}

func TestBrowseURLPrJjUsesTheBookmarkPullRequest(t *testing.T) {
	p := pathwrap.NewPath()

	f := newBrowseFixture(t, p.FromSlash("/repo/w/2345/.jj"))
	f.gh.EXPECT().WorktreeNumber("2345").Return(2345, false, nil)
	f.gh.EXPECT().PrURLForHead("github.com", "nutiliti/nutiliti", "jam/2345-1").Return("https://github.com/nutiliti/nutiliti/pull/2400", true, nil)
	url, err := f.w.BrowseURL("/repo/w/2345/apps", true)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/nutiliti/nutiliti/pull/2400", url)
}

func TestBrowseURLPrJjFallsBackToThePrNumber(t *testing.T) {
	p := pathwrap.NewPath()

	f := newBrowseFixture(t, p.FromSlash("/repo/w/200/.jj"))
	f.gh.EXPECT().WorktreeNumber("200").Return(200, false, nil)
	f.gh.EXPECT().PrURLForHead("github.com", "nutiliti/nutiliti", "jam/200-1").Return("", false, nil)
	f.gh.EXPECT().PrView("github.com", "nutiliti/nutiliti", 200).Return(forge.PullRequest{URL: "https://github.com/nutiliti/nutiliti/pull/200"}, true, nil)
	url, err := f.w.BrowseURL("/repo/w/200", true)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/nutiliti/nutiliti/pull/200", url)
}
