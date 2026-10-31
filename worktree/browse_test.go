package worktree

import (
	"os"
	"slices"
	"testing"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newBrowseWorktree(t *testing.T, jjDirs ...string) (*RealWorktree, *github.MockGithub) {
	mGh := github.NewMockGithub(t)
	mOs := oswrap.NewMockOs(t)
	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()
	mOs.EXPECT().ExpandEnv("/dotfiles").Return("/dotfiles").Maybe()
	mOs.EXPECT().Stat(mock.Anything).RunAndReturn(func(name string) (os.FileInfo, error) {
		if slices.Contains(jjDirs, name) {
			return dirInfo{}, nil
		}
		return nil, os.ErrNotExist
	}).Maybe()
	cfg := nuConfig()
	cfg.WorktreeConfigs = append([]model.WorktreeConfig{{Repo: "joshmedeski/dotfiles", Path: "/dotfiles"}}, cfg.WorktreeConfigs...)
	w := &RealWorktree{config: cfg, github: mGh, home: home.NewHome(mOs), os: mOs, path: pathwrap.NewPath()}
	return w, mGh
}

func TestBrowseURLOpensTheWorktreeIssue(t *testing.T) {
	w, _ := newBrowseWorktree(t)
	for _, path := range []string{"/repo/w/2345", "/repo/w/2345/apps/web"} {
		url, err := w.BrowseURL(path, false)
		require.NoError(t, err)
		assert.Equal(t, "https://github.com/nutiliti/nutiliti/issues/2345", url)
	}
}

func TestBrowseURLDefaultWorktreeDir(t *testing.T) {
	w, _ := newBrowseWorktree(t)
	url, err := w.BrowseURL("/dotfiles/.wk/12", false)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/joshmedeski/dotfiles/issues/12", url)
}

func TestBrowseURLRejectsPathsOutsideWorktrees(t *testing.T) {
	w, _ := newBrowseWorktree(t)
	for _, path := range []string{"/repo", "/repo/w", "/repo/w/notes", "/repo/wk/2345", "/elsewhere"} {
		_, err := w.BrowseURL(path, false)
		assert.ErrorContains(t, err, "is not a worktree", path)
	}
}

func TestBrowseURLRejectsEmptySessionPath(t *testing.T) {
	w, _ := newBrowseWorktree(t)
	_, err := w.BrowseURL("", false)
	assert.ErrorContains(t, err, "no focused tmux session")
}

func TestBrowseURLPrUsesTheBranchPullRequest(t *testing.T) {
	p := pathwrap.NewPath()

	w, mGh := newBrowseWorktree(t)
	mGh.EXPECT().PrURL(p.FromSlash("/repo/w/2345")).Return("https://github.com/nutiliti/nutiliti/pull/2400", true, nil)
	url, err := w.BrowseURL("/repo/w/2345/apps", true)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/nutiliti/nutiliti/pull/2400", url)
}

func TestBrowseURLPrWithoutPullRequest(t *testing.T) {
	p := pathwrap.NewPath()

	w, mGh := newBrowseWorktree(t)
	mGh.EXPECT().PrURL(p.FromSlash("/repo/w/2345")).Return("", false, nil)
	_, err := w.BrowseURL("/repo/w/2345", true)
	assert.EqualError(t, err, "no pull request found for worktree 2345")
}

func TestBrowseURLPrJjUsesTheBookmarkPullRequest(t *testing.T) {
	p := pathwrap.NewPath()

	w, mGh := newBrowseWorktree(t, p.FromSlash("/repo/w/2345/.jj"))
	mGh.EXPECT().PrURLForHead("nutiliti/nutiliti", "jam/2345-1").Return("https://github.com/nutiliti/nutiliti/pull/2400", true, nil)
	url, err := w.BrowseURL("/repo/w/2345/apps", true)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/nutiliti/nutiliti/pull/2400", url)
}

func TestBrowseURLPrJjFallsBackToThePrNumber(t *testing.T) {
	p := pathwrap.NewPath()

	w, mGh := newBrowseWorktree(t, p.FromSlash("/repo/w/200/.jj"))
	mGh.EXPECT().PrURLForHead("nutiliti/nutiliti", "jam/200-1").Return("", false, nil)
	mGh.EXPECT().PrView("nutiliti/nutiliti", 200).Return(github.PullRequest{URL: "https://github.com/nutiliti/nutiliti/pull/200"}, true, nil)
	url, err := w.BrowseURL("/repo/w/200", true)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/nutiliti/nutiliti/pull/200", url)
}
