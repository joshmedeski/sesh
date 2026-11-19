package worktree

import (
	"io/fs"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/browser"
	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/jj"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
)

// testIssueCache points the issue namespace at a scratch directory so tests
// never read or write the developer's real cache.
func testIssueCache(t *testing.T) *cache.Namespace[forge.Issue] {
	t.Helper()
	return cache.NewNamespaceInDir[forge.Issue](
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
	gh       *forge.MockForge
	gl       *forge.MockForge
	os       *oswrap.MockOs
	issues   *cache.Namespace[forge.Issue]
	created  map[string]time.Time
}

func (f *listFixture) expectWorktree(id int, isPr bool) {
	name := strconv.Itoa(id)
	if isPr {
		name = "mr-" + name
	}
	f.gh.EXPECT().WorktreeNumber(name).Return(id, isPr, nil).Maybe()
	f.gh.EXPECT().IssueKey("github.com", "nutiliti/nutiliti", id, isPr).
		Return("github.com/nutiliti/nutiliti#" + strconv.Itoa(id)).Maybe()
	f.gh.EXPECT().WorktreeName(id, isPr).Return(name).Maybe()
}

func (f *listFixture) expectWorktreeCustomHost(host, repo string, id int, isPr bool) {
	name := strconv.Itoa(id)
	if isPr {
		name = "mr-" + name
	}
	f.gh.EXPECT().WorktreeNumber(name).Return(id, isPr, nil).Maybe()
	f.gh.EXPECT().IssueKey(host, repo, id, isPr).
		Return(host + "/" + repo + "#" + strconv.Itoa(id)).Maybe()
	f.gh.EXPECT().WorktreeName(id, isPr).Return(name).Maybe()
}

func (f *listFixture) expectWorktrees(ids ...int) {
	for _, id := range ids {
		f.expectWorktree(id, false)
	}
}

func (f *listFixture) expectGitLabWorktreeCustomHost(host, repo string, id int, isPr bool) {
	name := strconv.Itoa(id)
	if isPr {
		name = "mr-" + name
	}
	key := host + "/" + repo + "#" + strconv.Itoa(id)
	if isPr {
		key = host + "/" + repo + "!" + strconv.Itoa(id)
	}
	f.gl.EXPECT().WorktreeNumber(name).Return(id, isPr, nil).Maybe()
	f.gl.EXPECT().IssueKey(host, repo, id, isPr).Return(key).Maybe()
	f.gl.EXPECT().WorktreeName(id, isPr).Return(name).Maybe()
}

func (f *listFixture) expectGitLabWorktree(repo string, id int, isPr bool) {
	f.expectGitLabWorktreeCustomHost("gitlab.com", repo, id, isPr)
}

type fakeFileInfo struct {
	os.FileInfo
	modTime time.Time
}

func (f fakeFileInfo) ModTime() time.Time { return f.modTime }

func newListFixture(t *testing.T, issues *cache.Namespace[forge.Issue]) listFixture {
	t.Helper()
	mOs := oswrap.NewMockOs(t)
	mGh := forge.NewMockForge(t)
	mGl := forge.NewMockForge(t)

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	// None of the paths under test contain ~ or $VARs, so expansion is identity.
	mOs.EXPECT().ExpandEnv(mock.Anything).
		RunAndReturn(func(s string) string { return s }).Maybe()
	mGh.EXPECT().PrsAsIssues("github.com", "nutiliti/nutiliti", []int{}).Return(map[int]forge.Issue{}, nil, nil).Maybe()
	created := map[string]time.Time{}
	mOs.EXPECT().Stat(mock.Anything).RunAndReturn(func(name string) (os.FileInfo, error) {
		if t, ok := created[name]; ok {
			return fakeFileInfo{modTime: t}, nil
		}
		return nil, os.ErrNotExist
	}).Maybe()

	cfg := nuConfig()
	cfg.WorktreeConfigs = append(cfg.WorktreeConfigs,
		model.WorktreeConfig{
			Repo:        "enterprise/repo",
			Host:        "github.example.com",
			Path:        "/enterprise-repo",
			WorktreeDir: ".wk",
		},
		model.WorktreeConfig{
			Repo:        "inkscape/inkscape",
			Forge:       "gitlab",
			Path:        "/gitlab-repo",
			WorktreeDir: ".wk",
		},
		model.WorktreeConfig{
			Repo:        "internal/repo",
			Forge:       "gitlab",
			Host:        "gitlab.example.com",
			Path:        "/internal-repo",
			WorktreeDir: ".wk",
		},
	)

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

	w := NewWorktree(
		cfg, git.NewMockGit(t), jj.NewMockJj(t), fr, connector.NewMockConnector(t),
		browser.NewMockBrowser(t), home.NewHome(mOs), mOs, pathwrap.NewPath(), issues,
	)
	f := listFixture{worktree: w, gh: mGh, gl: mGl, os: mOs, issues: issues, created: created}
	f.expectWorktrees(409, 426)
	return f
}

func TestList_FetchesTitlesOnColdCache(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("409", "426"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409, 426}).Return(map[int]forge.Issue{
		409: {Number: 409, Title: "worktree support", State: "OPEN"},
		426: {Number: 426, Title: "cap fuzzy penalty", State: "CLOSED"},
	}, nil, nil).Once()

	entries, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, []model.WorktreeEntry{
		{Number: 409, Name: "409", Path: p.FromSlash("/repo/w/409"), Title: "worktree support", State: "OPEN", RefType: model.RefTypeIssue},
		{Number: 426, Name: "426", Path: p.FromSlash("/repo/w/426"), Title: "cap fuzzy penalty", State: "CLOSED", RefType: model.RefTypeIssue},
	}, entries)
}

func TestList_WarmCacheMakesNoRequest(t *testing.T) {
	p := pathwrap.NewPath()
	wRoot := p.FromSlash("/repo/w")

	issues := testIssueCache(t)

	// First run populates the cache.
	first := newListFixture(t, issues)
	first.os.EXPECT().ReadDir(wRoot).Return(dirs("409"), nil)
	first.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).Return(map[int]forge.Issue{
		409: {Number: 409, Title: "worktree support", State: "OPEN"},
	}, nil, nil).Once()
	_, err := first.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})
	require.NoError(t, err)

	// Second run reuses it. No Issues expectation is registered, so any call
	// fails the test.
	second := newListFixture(t, issues)
	second.os.EXPECT().ReadDir(wRoot).Return(dirs("409"), nil)

	entries, err := second.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, "worktree support", entries[0].Title)
}

func TestList_OnlyFetchesTheStaleNumbers(t *testing.T) {
	p := pathwrap.NewPath()

	issues := testIssueCache(t)
	entries := cache.Entries[forge.Issue]{}
	entries.Put("github.com/nutiliti/nutiliti#409", forge.Issue{Number: 409, Title: "cached", State: "OPEN"})
	require.NoError(t, issues.Save(entries))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("409", "426"), nil)
	// Only 426 is requested — 409 is already fresh.
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{426}).Return(map[int]forge.Issue{
		426: {Number: 426, Title: "fetched", State: "OPEN"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, "cached", got[0].Title)
	assert.Equal(t, "fetched", got[1].Title)
}

func TestList_RefreshRefetchesEvenFreshEntries(t *testing.T) {
	p := pathwrap.NewPath()
	wRoot := p.FromSlash("/repo/w")

	issues := testIssueCache(t)
	cached := cache.Entries[forge.Issue]{}
	cached.Put("github.com/nutiliti/nutiliti#409", forge.Issue{Number: 409, Title: "stale title", State: "OPEN"})
	require.NoError(t, issues.Save(cached))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir(wRoot).Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).Return(map[int]forge.Issue{
		409: {Number: 409, Title: "renamed since", State: "CLOSED"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti", Refresh: true})

	require.NoError(t, err)
	assert.Equal(t, "renamed since", got[0].Title)
	assert.Equal(t, "CLOSED", got[0].State)

	// The refetched values are written back, so the next listing without
	// --refresh sees them too.
	next := newListFixture(t, issues)
	next.os.EXPECT().ReadDir(wRoot).Return(dirs("409"), nil)
	got, err = next.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})
	require.NoError(t, err)
	assert.Equal(t, "renamed since", got[0].Title)
}

func TestList_RefreshRefetchesNegativelyCachedNumbers(t *testing.T) {
	p := pathwrap.NewPath()

	issues := testIssueCache(t)
	cached := cache.Entries[forge.Issue]{}
	cached.PutMissing("github.com/nutiliti/nutiliti#409")
	require.NoError(t, issues.Save(cached))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).Return(map[int]forge.Issue{
		409: {Number: 409, Title: "created since", State: "OPEN"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti", Refresh: true})

	require.NoError(t, err)
	assert.Equal(t, "created since", got[0].Title)
}

func TestList_RefreshOnAnEmptyRootMakesNoRequest(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(nil, nil)

	entries, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti", Refresh: true})

	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestList_ExpiredEntriesAreRefetched(t *testing.T) {
	p := pathwrap.NewPath()

	dir := t.TempDir()
	issues := cache.NewNamespaceInDir[forge.Issue](dir, IssueCacheName, IssueCacheVersion, time.Hour)
	stored := cache.Entries[forge.Issue]{
		"github.com/nutiliti/nutiliti#409": {
			Value:     forge.Issue{Number: 409, Title: "old title", State: "OPEN"},
			FetchedAt: time.Now().Add(-2 * time.Hour),
		},
	}
	require.NoError(t, issues.Save(stored))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).Return(map[int]forge.Issue{
		409: {Number: 409, Title: "new title", State: "CLOSED"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, "new title", got[0].Title)
}

func TestList_ShowsStaleTitleWhenRefetchFails(t *testing.T) {
	p := pathwrap.NewPath()

	dir := t.TempDir()
	issues := cache.NewNamespaceInDir[forge.Issue](dir, IssueCacheName, IssueCacheVersion, time.Hour)
	stored := cache.Entries[forge.Issue]{
		"github.com/nutiliti/nutiliti#409": {
			Value:     forge.Issue{Number: 409, Title: "stale but useful", State: "OPEN"},
			FetchedAt: time.Now().Add(-2 * time.Hour),
		},
	}
	require.NoError(t, issues.Save(stored))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).
		Return(nil, nil, assert.AnError).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	// An offline refresh must not fail the listing or blank out the title.
	require.NoError(t, err)
	assert.Equal(t, "stale but useful", got[0].Title)
}

func TestList_FetchFailureOnColdCacheStillLists(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).
		Return(nil, nil, assert.AnError).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, []model.WorktreeEntry{{Number: 409, Name: "409", Path: p.FromSlash("/repo/w/409"), RefType: model.RefTypeIssue}}, got)
}

func TestList_MissingIssuesAreNotRefetched(t *testing.T) {
	p := pathwrap.NewPath()
	wRoot := p.FromSlash("/repo/w")

	issues := testIssueCache(t)

	first := newListFixture(t, issues)
	first.expectWorktrees(99999999)
	first.os.EXPECT().ReadDir(wRoot).Return(dirs("99999999"), nil)
	first.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{99999999}).
		Return(map[int]forge.Issue{}, []int{99999999}, nil).Once()
	got, err := first.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})
	require.NoError(t, err)
	assert.Empty(t, got[0].Title)

	// The absence was cached, so the second run does not ask again.
	second := newListFixture(t, issues)
	second.expectWorktrees(99999999)
	second.os.EXPECT().ReadDir(wRoot).Return(dirs("99999999"), nil)

	got, err = second.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, []model.WorktreeEntry{{Number: 99999999, Name: "99999999", Path: p.FromSlash("/repo/w/99999999"), RefType: model.RefTypeIssue}}, got)
}

func TestList_IgnoresNonNumericAndNonDirEntries(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.gh.EXPECT().WorktreeNumber("scratch").Return(0, false, strconv.ErrSyntax)
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return([]os.DirEntry{
		fakeDirEntry{name: "409", isDir: true},
		fakeDirEntry{name: "scratch", isDir: true},
		fakeDirEntry{name: ".DS_Store", isDir: false},
		fakeDirEntry{name: "426", isDir: false}, // a file, not a worktree
	}, nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).Return(map[int]forge.Issue{
		409: {Number: 409, Title: "worktree support", State: "OPEN"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, 409, got[0].Number)
}

func TestList_SortsNumerically(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.expectWorktrees(99, 1000)
	// ReadDir returns lexical order, in which "1000" precedes "99".
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("1000", "409", "99"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{99, 409, 1000}).
		Return(map[int]forge.Issue{}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, []int{99, 409, 1000}, []int{got[0].Number, got[1].Number, got[2].Number})
}

func TestList_AbsentRootIsEmptyNotAnError(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(nil, os.ErrNotExist)

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err, "a repo with no worktrees yet is not a failure")
	assert.Empty(t, got)
}

func TestList_ReadDirFailureIsAnError(t *testing.T) {
	p := pathwrap.NewPath()
	wRoot := p.FromSlash("/repo/w")

	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir(wRoot).Return(nil, os.ErrPermission)

	_, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), strconv.Quote(wRoot))
}

func TestList_EmptyRootMakesNoRequest(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	// No Issues expectation: an empty root must not trigger a fetch.
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs(), nil)

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestList_SelectsConfigByPath(t *testing.T) {
	p := pathwrap.NewPath()
	wRoot := p.FromSlash("/repo/w")

	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ReadDir(wRoot).Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).
		Return(map[int]forge.Issue{409: {Number: 409, Title: "t", State: "OPEN"}}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Path: wRoot})

	require.NoError(t, err)
	assert.Equal(t, "t", got[0].Title)
}

func TestList_UnknownPathErrorsWithKnownRoots(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.os.EXPECT().ExpandEnv("/nope").Return("/nope").Maybe()

	_, err := f.worktree.List(model.WorktreeListOpts{Path: "/nope"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "/nope")
	assert.Contains(t, err.Error(), p.FromSlash("/repo/w"), "the error should name the roots that do exist")
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
	created := map[string]time.Time{}
	mOs.EXPECT().Stat(mock.Anything).RunAndReturn(func(name string) (os.FileInfo, error) {
		if t, ok := created[name]; ok {
			return fakeFileInfo{modTime: t}, nil
		}
		return nil, os.ErrNotExist
	}).Maybe()
	w := NewWorktree(
		config, git.NewMockGit(t), jj.NewMockJj(t), forgeResolver(forge.NewMockForge(t)), connector.NewMockConnector(t),
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
	p := pathwrap.NewPath()

	// A cached issue 409 belonging to a different repo must not be reused.
	issues := testIssueCache(t)
	entries := cache.Entries[forge.Issue]{}
	entries.Put("github.com/other/repo#409", forge.Issue{Number: 409, Title: "wrong repo", State: "OPEN"})
	require.NoError(t, issues.Save(entries))

	f := newListFixture(t, issues)
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).Return(map[int]forge.Issue{
		409: {Number: 409, Title: "right repo", State: "OPEN"},
	}, nil, nil).Once()

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})

	require.NoError(t, err)
	assert.Equal(t, "right repo", got[0].Title)
}

func TestList_CreatedFromDotGitOrDotJjModTime(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	added := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	f.created[p.FromSlash("/repo/w/409/.git")] = added
	f.created[p.FromSlash("/repo/w/426/.jj")] = added.Add(time.Hour)
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("409", "426"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409, 426}).Return(map[int]forge.Issue{}, nil, nil)

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, added, got[0].Created)
	assert.Equal(t, added.Add(time.Hour), got[1].Created)
}

func TestList_CarriesConfiguredIcon(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.worktree.(*RealWorktree).config.WorktreeConfigs[0].Icon = "🏠"
	f.os.EXPECT().ReadDir(p.FromSlash("/repo/w")).Return(dirs("409"), nil)
	f.gh.EXPECT().Issues("github.com", "nutiliti/nutiliti", []int{409}).Return(map[int]forge.Issue{}, nil, nil)

	got, err := f.worktree.List(model.WorktreeListOpts{Repo: "nutiliti/nutiliti"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "🏠", got[0].Icon)
}

func TestList_GitLabIssuesAndMergeRequests(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.expectGitLabWorktree("inkscape/inkscape", 42, false)
	f.expectGitLabWorktree("inkscape/inkscape", 101, true)

	f.os.EXPECT().ReadDir(p.FromSlash("/gitlab-repo/.wk")).Return(dirs("42", "mr-101"), nil)
	f.gl.EXPECT().Issues("gitlab.com", "inkscape/inkscape", []int{42}).Return(map[int]forge.Issue{
		42: {Number: 42, Title: "GitLab Issue", State: "OPEN"},
	}, nil, nil).Once()
	f.gl.EXPECT().PrsAsIssues("gitlab.com", "inkscape/inkscape", []int{101}).Return(map[int]forge.Issue{
		101: {Number: 101, Title: "GitLab MR", State: "OPEN"},
	}, nil, nil).Once()

	entries, err := f.worktree.List(model.WorktreeListOpts{Repo: "inkscape/inkscape"})

	require.NoError(t, err)
	assert.Equal(t, []model.WorktreeEntry{
		{Number: 42, Name: "42", Path: p.FromSlash("/gitlab-repo/.wk/42"), Title: "GitLab Issue", State: "OPEN", RefType: model.RefTypeIssue},
		{Number: 101, Name: "mr-101", Path: p.FromSlash("/gitlab-repo/.wk/mr-101"), Title: "GitLab MR", State: "OPEN", RefType: model.RefTypePr},
	}, entries)
}

func TestList_GitLabCustomHost(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.expectGitLabWorktreeCustomHost("gitlab.example.com", "internal/repo", 99, false)
	f.expectGitLabWorktreeCustomHost("gitlab.example.com", "internal/repo", 202, true)

	f.os.EXPECT().ReadDir(p.FromSlash("/internal-repo/.wk")).Return(dirs("99", "mr-202"), nil)
	f.gl.EXPECT().Issues("gitlab.example.com", "internal/repo", []int{99}).Return(map[int]forge.Issue{
		99: {Number: 99, Title: "Custom Host Issue", State: "OPEN"},
	}, nil, nil).Once()
	f.gl.EXPECT().PrsAsIssues("gitlab.example.com", "internal/repo", []int{202}).Return(map[int]forge.Issue{
		202: {Number: 202, Title: "Custom Host MR", State: "OPEN"},
	}, nil, nil).Once()

	entries, err := f.worktree.List(model.WorktreeListOpts{Repo: "internal/repo"})

	require.NoError(t, err)
	assert.Equal(t, []model.WorktreeEntry{
		{Number: 99, Name: "99", Path: p.FromSlash("/internal-repo/.wk/99"), Title: "Custom Host Issue", State: "OPEN", RefType: model.RefTypeIssue},
		{Number: 202, Name: "mr-202", Path: p.FromSlash("/internal-repo/.wk/mr-202"), Title: "Custom Host MR", State: "OPEN", RefType: model.RefTypePr},
	}, entries)
}

func TestList_GitHubCustomHost(t *testing.T) {
	p := pathwrap.NewPath()

	f := newListFixture(t, testIssueCache(t))
	f.expectWorktreeCustomHost("github.example.com", "enterprise/repo", 55, false)

	f.os.EXPECT().ReadDir(p.FromSlash("/enterprise-repo/.wk")).Return(dirs("55"), nil)
	f.gh.EXPECT().Issues("github.example.com", "enterprise/repo", []int{55}).Return(map[int]forge.Issue{
		55: {Number: 55, Title: "Enterprise Issue", State: "OPEN"},
	}, nil, nil).Once()
	f.gh.EXPECT().PrsAsIssues("github.example.com", "enterprise/repo", []int{}).Return(map[int]forge.Issue{}, nil, nil).Once()

	entries, err := f.worktree.List(model.WorktreeListOpts{Repo: "enterprise/repo"})

	require.NoError(t, err)
	assert.Equal(t, []model.WorktreeEntry{
		{Number: 55, Name: "55", Path: p.FromSlash("/enterprise-repo/.wk/55"), Title: "Enterprise Issue", State: "OPEN", RefType: model.RefTypeIssue},
	}, entries)
}
