package lister_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

var longAgo = time.Now().Add(-time.Hour)

// The tests below put names in a fixed order: tmux, config, tmuxinator, zoxide,
// the default source order.

func TestCachingLister_TmuxSessionsAreLiveNotTheCachedSnapshot(t *testing.T) {
	f := newFixture(t, nil)
	f.config.SessionConfigs = nil
	f.zoxide.set()
	f.tmuxinator.configs = nil
	cl, _ := f.cached()

	f.tmux.set(&model.TmuxSession{Name: "a"}, &model.TmuxSession{Name: "b"})
	require.Equal(t, []string{"a", "b"}, listNames(t, cl, lister.ListOptions{}))

	f.tmux.set(&model.TmuxSession{Name: "a"}, &model.TmuxSession{Name: "c"})
	assert.Equal(t, []string{"a", "c"}, listNames(t, cl, lister.ListOptions{}))
}

// tmux.ListSessions turns a failed `tmux list-sessions` into an empty list
// with no error.
func TestCachingLister_EmptyTmuxListIsNotServedLater(t *testing.T) {
	f := newFixture(t, nil)
	f.config.SessionConfigs = nil
	f.zoxide.set()
	f.tmuxinator.configs = nil
	f.tmux.set()
	cl, _ := f.cached()
	require.Empty(t, listNames(t, cl, lister.ListOptions{}))

	f.tmux.set(&model.TmuxSession{Name: "a"}, &model.TmuxSession{Name: "b"})
	assert.Equal(t, []string{"a", "b"}, listNames(t, cl, lister.ListOptions{}))
}

func TestCachingLister_TmuxErrorSurfacesOnAWarmCache(t *testing.T) {
	f := newFixture(t, nil)
	cl, _ := f.cached()
	listNames(t, cl, lister.ListOptions{})

	f.tmux.fail(errors.New("tmux exploded"))
	_, err := cl.List(lister.ListOptions{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "couldn't list tmux sessions")
}

func TestCachingLister_CacheKeepsOnlyTheSlowSources(t *testing.T) {
	f := newFixture(t, nil)
	cl, fc := f.cached()

	listNames(t, cl, lister.ListOptions{})

	stored, err := fc.Read()
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"tmuxinator", "zoxide"}, sources(stored.Sessions))
}

func TestCachingLister_FreshHitRunsNoSlowSource(t *testing.T) {
	f := newFixture(t, nil)
	cl, _ := f.cached()
	listNames(t, cl, lister.ListOptions{})
	require.EqualValues(t, 1, f.zoxide.calls.Load())
	require.EqualValues(t, 1, f.tmuxinator.calls.Load())

	listNames(t, cl, lister.ListOptions{})
	listNames(t, cl, lister.ListOptions{})

	assert.EqualValues(t, 1, f.zoxide.calls.Load())
	assert.EqualValues(t, 1, f.tmuxinator.calls.Load())
	assert.EqualValues(t, 3, f.tmux.sessionCalls.Load(), "tmux is listed on every list")
}

func TestCachingLister_CachedSourceWithNoRowsIsNotFetchedAgain(t *testing.T) {
	f := newFixture(t, nil)
	f.tmuxinator.configs = nil
	cl, _ := f.cached()
	listNames(t, cl, lister.ListOptions{})

	listNames(t, cl, lister.ListOptions{})

	assert.EqualValues(t, 1, f.tmuxinator.calls.Load())
}

func TestCachingLister_StaleHitServesStoredSourcesAndRevalidates(t *testing.T) {
	f := newFixture(t, nil)
	f.config.SessionConfigs = nil
	f.tmuxinator.configs = nil
	writeCacheFile(t, f.cachePath, storedZoxide("/home/user/old"), longAgo)
	f.zoxide.set(&model.ZoxideResult{Path: "/home/user/new", Score: 1})
	started, release := f.zoxide.hold()
	t.Cleanup(release)
	cl, fc := f.cached()

	got, err := listWithin(t, cl, lister.ListOptions{})

	// The background fetch is still held open: the list did not wait for it.
	require.NoError(t, err)
	assert.Equal(t, []string{"main", "dev", "~/old"}, names(got), "tmux is live, zoxide is the stored snapshot")
	<-started
	release()
	cl.Wait()
	stored, err := fc.Read()
	require.NoError(t, err)
	assert.Equal(t, []string{"zoxide:~/new"}, stored.Sessions.OrderedIndex)
}

func TestCachingLister_StaleHitSurvivesAFailingRefresh(t *testing.T) {
	f := newFixture(t, nil)
	f.config.SessionConfigs = nil
	f.tmuxinator.configs = nil
	writeCacheFile(t, f.cachePath, storedZoxide("/home/user/old"), longAgo)
	before, err := os.ReadFile(f.cachePath)
	require.NoError(t, err)
	f.zoxide.fail(errors.New("zoxide is broken"))
	cl, _ := f.cached()

	got, err := listWithin(t, cl, lister.ListOptions{})
	cl.Wait()

	require.NoError(t, err)
	assert.Equal(t, []string{"main", "dev", "~/old"}, names(got))
	after, err := os.ReadFile(f.cachePath)
	require.NoError(t, err)
	assert.Equal(t, before, after, "a failed refresh leaves the stored sources in place")
}

func TestCachingLister_SkipReadFetchesTheSlowSourcesAndWritesThem(t *testing.T) {
	f := newFixture(t, nil)
	f.config.SessionConfigs = nil
	f.tmuxinator.configs = nil
	writeCacheFile(t, f.cachePath, storedZoxide("/home/user/old"), time.Now())
	f.zoxide.set(&model.ZoxideResult{Path: "/home/user/new", Score: 1})
	cl, fc := f.cached()
	cl.SkipRead()

	got := listNames(t, cl, lister.ListOptions{})

	assert.Equal(t, []string{"main", "dev", "~/new"}, got)
	stored, err := fc.Read()
	require.NoError(t, err)
	assert.Equal(t, []string{"zoxide:~/new"}, stored.Sessions.OrderedIndex)
}

func TestCachingLister_ListWithoutACachedSourceLeavesTheCacheAlone(t *testing.T) {
	f := newFixture(t, nil)
	writeCacheFile(t, f.cachePath, storedZoxide("/home/user/old"), longAgo)
	before, err := os.ReadFile(f.cachePath)
	require.NoError(t, err)
	cl, _ := f.cached()

	got := listNames(t, cl, lister.ListOptions{Tmux: true})
	cl.Wait()

	assert.Equal(t, []string{"main", "dev"}, got)
	assert.Zero(t, f.zoxide.calls.Load())
	assert.Zero(t, f.tmuxinator.calls.Load())
	after, err := os.ReadFile(f.cachePath)
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestCachingLister_ColdStartWithoutACachedSourceSkipsZoxide(t *testing.T) {
	f := newFixture(t, nil)
	f.zoxide.fail(errors.New("zoxide is broken"))
	cl, _ := f.cached()

	got, err := cl.List(lister.ListOptions{Tmux: true})

	require.NoError(t, err)
	assert.Equal(t, []string{"main", "dev"}, names(got))
	assert.Zero(t, f.zoxide.calls.Load())
	assert.NoFileExists(t, f.cachePath)
}

func TestCachingLister_SkipReadWithoutACachedSourceRefreshesInTheBackground(t *testing.T) {
	f := newFixture(t, nil)
	f.zoxide.fail(errors.New("zoxide is broken"))
	started, release := f.zoxide.hold()
	t.Cleanup(release)
	cl, _ := f.cached()
	cl.SkipRead()

	got, err := listWithin(t, cl, lister.ListOptions{Tmux: true})

	// The held zoxide call is the refresh: the list returned without it.
	require.NoError(t, err)
	assert.Equal(t, []string{"main", "dev"}, names(got))
	<-started
	release()
	cl.Wait()
	assert.EqualValues(t, 1, f.zoxide.calls.Load())
}

func TestCachingLister_ColdStartWithOnlyTmuxinatorRequestedSurvivesAFailingZoxide(t *testing.T) {
	f := newFixture(t, nil)
	f.zoxide.fail(errors.New("zoxide is broken"))
	cl, _ := f.cached()

	got, err := cl.List(lister.ListOptions{Tmuxinator: true})

	require.NoError(t, err)
	assert.Equal(t, []string{"work"}, names(got))
	assert.EqualValues(t, 2, f.tmuxinator.calls.Load(), "fetched with the cached sources, then again by the uncached list")
	assert.NoFileExists(t, f.cachePath)
}

func TestCachingLister_ColdStartWithZoxideRequestedFailsWhenZoxideFails(t *testing.T) {
	f := newFixture(t, nil)
	f.zoxide.fail(errors.New("zoxide is broken"))
	cl, _ := f.cached()

	_, err := cl.List(lister.ListOptions{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "zoxide is broken")
	assert.NoFileExists(t, f.cachePath)
}

func TestCachingLister_ListSucceedsWhenTheCacheCannotBeWritten(t *testing.T) {
	f := newFixture(t, nil)
	blocker := filepath.Join(t.TempDir(), "file")
	require.NoError(t, os.WriteFile(blocker, nil, 0o644))
	f.cachePath = filepath.Join(blocker, "sessions.v2.gob")
	cl, _ := f.cached()

	got := listNames(t, cl, lister.ListOptions{Tmux: true, Zoxide: true})

	assert.Equal(t, []string{"main", "dev", "~/code", "~/dev", "~/project"}, got)
}

// A cache file an older sesh wrote holds the whole list, tmux and config rows
// included. None of those rows may reach the output.
func TestCachingLister_StoredLiveSourcesNeverReachTheOutput(t *testing.T) {
	f := newFixture(t, nil)
	stored := storedZoxide("/home/user/stored")
	stored.OrderedIndex = append([]string{"tmux:gone", "config:project"}, stored.OrderedIndex...)
	stored.Directory["tmux:gone"] = model.SeshSession{Src: "tmux", Name: "gone", Path: "/gone"}
	stored.Directory["config:project"] = model.SeshSession{Src: "config", Name: "project", Path: "/old/project"}
	writeCacheFile(t, f.cachePath, stored, time.Now())
	cl, _ := f.cached()

	got, err := cl.List(lister.ListOptions{})

	require.NoError(t, err)
	assert.Equal(t, []string{"main", "dev", "project", "dev", "~/stored"}, names(got))
	assert.Equal(t, "/home/user/project", got.Directory["config:project"].Path)
	cl.Wait()
}

// A cached row the stored directory does not know is skipped instead of
// listed as an empty session.
func TestCachingLister_StoredRowWithoutAnEntryIsSkipped(t *testing.T) {
	f := newFixture(t, nil)
	stored := storedZoxide("/home/user/stored")
	stored.OrderedIndex = append(stored.OrderedIndex, "zoxide:~/dangling")
	writeCacheFile(t, f.cachePath, stored, time.Now())
	cl, _ := f.cached()

	assert.Equal(t, []string{"~/stored"}, listNames(t, cl, lister.ListOptions{Zoxide: true}))
}

func TestCachingLister_FileNamedByEarlierVersionIsNeverRead(t *testing.T) {
	cacheHome := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cacheHome)
	f := newFixture(t, nil)
	oldFile := filepath.Join(cacheHome, "sesh", "sessions.gob")
	stored := storedZoxide("/home/user/stored")
	stored.OrderedIndex = append([]string{"tmux:gone"}, stored.OrderedIndex...)
	stored.Directory["tmux:gone"] = model.SeshSession{Src: "tmux", Name: "gone", Path: "/gone"}
	writeCacheFile(t, oldFile, stored, time.Now())
	before, err := os.ReadFile(oldFile)
	require.NoError(t, err)
	cl := lister.NewCachingLister(f.real(), cache.NewFileCache())

	got := listNames(t, cl, lister.ListOptions{Tmux: true, Zoxide: true})
	cl.Wait()

	assert.Equal(t, []string{"main", "dev", "~/code", "~/dev", "~/project"}, got)
	after, err := os.ReadFile(oldFile)
	require.NoError(t, err)
	assert.Equal(t, before, after)
	assert.FileExists(t, filepath.Join(cacheHome, "sesh", "sessions.v2.gob"))
}

func TestCachingLister_DelegatesTheOtherListerMethods(t *testing.T) {
	inner := lister.NewMockSourceLister(t)
	session := model.SeshSession{Src: "tmux", Name: "test"}
	inner.On("FindTmuxSession", "test").Return(session, true)
	inner.On("FindTmuxSessionByBase", "test").Return(session, true)
	inner.On("GetAttachedTmuxSession").Return(session, true)
	inner.On("GetLastTmuxSession").Return(session, true)
	inner.On("FindConfigSession", "test").Return(session, true)
	inner.On("FindConfigWildcard", "/path").Return(model.WildcardConfig{}, false)
	inner.On("FindZoxideSession", "test").Return(session, true)
	inner.On("FindTmuxinatorConfig", "test").Return(session, true)
	inner.On("ListTmuxPanes").Return(model.SeshSessions{}, nil)
	inner.On("Format", model.SeshSessions{}, lister.ListOptions{}).Return(model.SeshSessions{}, nil)
	cl := lister.NewCachingLister(inner, cache.NewFileCacheWithPath(filepath.Join(t.TempDir(), "sessions.v2.gob")))

	s, ok := cl.FindTmuxSession("test")
	assert.True(t, ok)
	assert.Equal(t, "test", s.Name)
	_, ok = cl.FindTmuxSessionByBase("test")
	assert.True(t, ok)
	_, ok = cl.GetAttachedTmuxSession()
	assert.True(t, ok)
	_, ok = cl.GetLastTmuxSession()
	assert.True(t, ok)
	_, ok = cl.FindConfigSession("test")
	assert.True(t, ok)
	_, ok = cl.FindConfigWildcard("/path")
	assert.False(t, ok)
	_, ok = cl.FindZoxideSession("test")
	assert.True(t, ok)
	_, ok = cl.FindTmuxinatorConfig("test")
	assert.True(t, ok)
	panes, err := cl.ListTmuxPanes()
	assert.NoError(t, err)
	assert.Empty(t, panes.OrderedIndex)
	_, err = cl.Format(model.SeshSessions{}, lister.ListOptions{})
	assert.NoError(t, err)
}
