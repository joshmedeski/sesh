package worktree

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/browser"
	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
)

// testIssueCache points the issue namespace at a scratch directory so tests
// never read or write the developer's real cache.
func testIssueCache(t *testing.T) *cache.Namespace[github.Issue] {
	t.Helper()
	return cache.NewNamespaceInDir[github.Issue](
		t.TempDir(), IssueCacheName, IssueCacheVersion, IssueCacheTTL,
	)
}

// fakeDirEntry is a minimal fs.DirEntry; only Name and IsDir are consulted.
type fakeDirEntry struct {
	name  string
	isDir bool
}

func (f fakeDirEntry) Name() string               { return f.name }
func (f fakeDirEntry) IsDir() bool                { return f.isDir }
func (f fakeDirEntry) Type() fs.FileMode          { return 0 }
func (f fakeDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func dirs(names ...string) []os.DirEntry {
	entries := make([]os.DirEntry, 0, len(names))
	for _, name := range names {
		entries = append(entries, fakeDirEntry{name: name, isDir: true})
	}
	return entries
}

// listFixture wires a RealWorktree for List with the repo root at /repo and
// worktrees under /repo/w.
type listFixture struct {
	worktree Worktree
	gh       *github.MockGithub
	os       *oswrap.MockOs
	issues   *cache.Namespace[github.Issue]
}

func newListFixture(t *testing.T, issues *cache.Namespace[github.Issue]) listFixture {
	t.Helper()
	mOs := oswrap.NewMockOs(t)
	mGh := github.NewMockGithub(t)

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	// None of the paths under test contain ~ or $VARs, so expansion is identity.
	mOs.EXPECT().ExpandEnv(mock.Anything).
		RunAndReturn(func(s string) string { return s }).Maybe()

	w := NewWorktree(
		nuConfig(), git.NewMockGit(t), mGh, connector.NewMockConnector(t),
		browser.NewMockBrowser(t), home.NewHome(mOs), mOs, pathwrap.NewPath(), issues,
	)
	return listFixture{worktree: w, gh: mGh, os: mOs, issues: issues}
}

func TestList_FetchesTitlesOnColdCache(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir("/repo/w").Return(dirs("409", "426"), nil)
	f.gh.EXPECT().Issues("nutiliti/nutiliti", []int{409, 426}).Return(map[int]github.Issue{
		409: {Number: 409, Title: "worktree support", State: "OPEN"},
		426: {Number: 426, Title: "cap fuzzy penalty", State: "CLOSED"},
	}, nil, nil).Once()

	entries, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, []model.WorktreeEntry{
		{Number: 409, Path: "/repo/w/409", Title: "worktree support", State: "OPEN"},
		{Number: 426, Path: "/repo/w/426", Title: "cap fuzzy penalty", State: "CLOSED"},
	}, entries)
}

func TestList_WarmCacheMakesNoRequest(t *testing.T) {
	issues := testIssueCache(t)

	// First run populates the cache.
	first := newListFixture(t, issues)
	first.os.EXPECT().ReadDir("/repo/w").Return(dirs("409"), nil)
	first.gh.EXPECT().Issues("nutiliti/nutiliti", []int{409}).Return(map[int]github.Issue{
		409: {Number: 409, Title: "worktree support", State: "OPEN"},
	}, nil, nil).Once()
	_, err := first.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})
	require.NoError(t, err)

	// Second run reuses it. No Issues expectation is registered, so any call
	// fails the test.
	second := newListFixture(t, issues)
	second.os.EXPECT().ReadDir("/repo/w").Return(dirs("409"), nil)

	entries, err := second.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, "worktree support", entries[0].Title)
}

func TestList_OnlyFetchesTheStaleNumbers(t *testing.T) {
	issues := testIssueCache(t)
	entries := cache.Entries[github.Issue]{}
	entries.Put(issueKey("nutiliti/nutiliti", 409), github.Issue{Number: 409, Title: "cached", State: "OPEN"})
	require.NoError(t, issues.Save(entries))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir("/repo/w").Return(dirs("409", "426"), nil)
	// Only 426 is requested — 409 is already fresh.
	f.gh.EXPECT().Issues("nutiliti/nutiliti", []int{426}).Return(map[int]github.Issue{
		426: {Number: 426, Title: "fetched", State: "OPEN"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, "cached", got[0].Title)
	assert.Equal(t, "fetched", got[1].Title)
}

func TestList_ExpiredEntriesAreRefetched(t *testing.T) {
	dir := t.TempDir()
	issues := cache.NewNamespaceInDir[github.Issue](dir, IssueCacheName, IssueCacheVersion, time.Hour)
	stored := cache.Entries[github.Issue]{
		issueKey("nutiliti/nutiliti", 409): {
			Value:     github.Issue{Number: 409, Title: "old title", State: "OPEN"},
			FetchedAt: time.Now().Add(-2 * time.Hour),
		},
	}
	require.NoError(t, issues.Save(stored))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir("/repo/w").Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("nutiliti/nutiliti", []int{409}).Return(map[int]github.Issue{
		409: {Number: 409, Title: "new title", State: "CLOSED"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, "new title", got[0].Title)
}

func TestList_ShowsStaleTitleWhenRefetchFails(t *testing.T) {
	dir := t.TempDir()
	issues := cache.NewNamespaceInDir[github.Issue](dir, IssueCacheName, IssueCacheVersion, time.Hour)
	stored := cache.Entries[github.Issue]{
		issueKey("nutiliti/nutiliti", 409): {
			Value:     github.Issue{Number: 409, Title: "stale but useful", State: "OPEN"},
			FetchedAt: time.Now().Add(-2 * time.Hour),
		},
	}
	require.NoError(t, issues.Save(stored))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir("/repo/w").Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("nutiliti/nutiliti", []int{409}).
		Return(nil, nil, assert.AnError).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	// An offline refresh must not fail the listing or blank out the title.
	require.NoError(t, err)
	assert.Equal(t, "stale but useful", got[0].Title)
}

func TestList_FetchFailureOnColdCacheStillLists(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir("/repo/w").Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("nutiliti/nutiliti", []int{409}).
		Return(nil, nil, assert.AnError).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, []model.WorktreeEntry{{Number: 409, Path: "/repo/w/409"}}, got)
}

func TestList_MissingIssuesAreNotRefetched(t *testing.T) {
	issues := testIssueCache(t)

	first := newListFixture(t, issues)
	first.os.EXPECT().ReadDir("/repo/w").Return(dirs("99999999"), nil)
	first.gh.EXPECT().Issues("nutiliti/nutiliti", []int{99999999}).
		Return(map[int]github.Issue{}, []int{99999999}, nil).Once()
	got, err := first.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})
	require.NoError(t, err)
	assert.Empty(t, got[0].Title)

	// The absence was cached, so the second run does not ask again.
	second := newListFixture(t, issues)
	second.os.EXPECT().ReadDir("/repo/w").Return(dirs("99999999"), nil)

	got, err = second.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, []model.WorktreeEntry{{Number: 99999999, Path: "/repo/w/99999999"}}, got)
}

func TestList_IgnoresNonNumericAndNonDirEntries(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir("/repo/w").Return([]os.DirEntry{
		fakeDirEntry{name: "409", isDir: true},
		fakeDirEntry{name: "scratch", isDir: true},
		fakeDirEntry{name: ".DS_Store", isDir: false},
		fakeDirEntry{name: "426", isDir: false}, // a file, not a worktree
	}, nil)
	f.gh.EXPECT().Issues("nutiliti/nutiliti", []int{409}).Return(map[int]github.Issue{
		409: {Number: 409, Title: "worktree support", State: "OPEN"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, 409, got[0].Number)
}

func TestList_SortsNumerically(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))
	// ReadDir returns lexical order, in which "1000" precedes "99".
	f.os.EXPECT().ReadDir("/repo/w").Return(dirs("1000", "409", "99"), nil)
	f.gh.EXPECT().Issues("nutiliti/nutiliti", []int{99, 409, 1000}).
		Return(map[int]github.Issue{}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, []int{99, 409, 1000}, []int{got[0].Number, got[1].Number, got[2].Number})
}

func TestList_AbsentRootIsEmptyNotAnError(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir("/repo/w").Return(nil, os.ErrNotExist)

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err, "a repo with no worktrees yet is not a failure")
	assert.Empty(t, got)
}

func TestList_ReadDirFailureIsAnError(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir("/repo/w").Return(nil, os.ErrPermission)

	_, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "/repo/w")
}

func TestList_EmptyRootMakesNoRequest(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))
	// No Issues expectation: an empty root must not trigger a fetch.
	f.os.EXPECT().ReadDir("/repo/w").Return(dirs(), nil)

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestList_SelectsConfigByPath(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir("/repo/w").Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("nutiliti/nutiliti", []int{409}).
		Return(map[int]github.Issue{409: {Number: 409, Title: "t", State: "OPEN"}}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Path: "/repo/w"})

	require.NoError(t, err)
	assert.Equal(t, "t", got[0].Title)
}

func TestList_UnknownPathErrorsWithKnownRoots(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ExpandEnv("/nope").Return("/nope").Maybe()

	_, err := f.worktree.List(model.WorktreeListOpts{Path: "/nope"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "/nope")
	assert.Contains(t, err.Error(), "/repo/w", "the error should name the roots that do exist")
}

func TestList_MissingRepoKeyErrors(t *testing.T) {
	// A [[worktree]] block using an older key name leaves Repo empty. That is a
	// config mistake, so it must surface rather than silently listing bare
	// numbers with no titles.
	config := model.Config{WorktreeConfigs: []model.WorktreeConfig{{
		Path: "/repo", WorktreeDir: "w",
	}}}
	mOs := oswrap.NewMockOs(t)
	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv(mock.Anything).
		RunAndReturn(func(s string) string { return s }).Maybe()
	w := NewWorktree(
		config, git.NewMockGit(t), github.NewMockGithub(t), connector.NewMockConnector(t),
		browser.NewMockBrowser(t), home.NewHome(mOs), mOs, pathwrap.NewPath(), testIssueCache(t),
	)

	_, err := w.List(model.WorktreeListOpts{Path: "/repo/w"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "repo")
}

func TestList_UnknownRepoErrors(t *testing.T) {
	f := newListFixture(t, testIssueCache(t))

	_, err := f.worktree.List(model.WorktreeListOpts{Repo: "other/repo"})

	require.Error(t, err)
}

func TestList_KeysAreScopedByRepo(t *testing.T) {
	// A cached issue 409 belonging to a different repo must not be reused.
	issues := testIssueCache(t)
	entries := cache.Entries[github.Issue]{}
	entries.Put(issueKey("other/repo", 409), github.Issue{Number: 409, Title: "wrong repo", State: "OPEN"})
	require.NoError(t, issues.Save(entries))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir("/repo/w").Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("nutiliti/nutiliti", []int{409}).Return(map[int]github.Issue{
		409: {Number: 409, Title: "right repo", State: "OPEN"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, "right repo", got[0].Title)
}
