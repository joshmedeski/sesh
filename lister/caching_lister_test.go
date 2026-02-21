package lister_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

func fakeSessions() model.SeshSessions {
	return model.SeshSessions{
		OrderedIndex: []string{"tmux:main"},
		Directory: model.SeshSessionMap{
			"tmux:main": {Src: "tmux", Name: "main", Path: "/home/user"},
		},
	}
}

func updatedSessions() model.SeshSessions {
	return model.SeshSessions{
		OrderedIndex: []string{"tmux:main", "tmux:dev"},
		Directory: model.SeshSessionMap{
			"tmux:main": {Src: "tmux", Name: "main", Path: "/home/user"},
			"tmux:dev":  {Src: "tmux", Name: "dev", Path: "/home/user/dev"},
		},
	}
}

func TestCachingLister_ColdStart(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	sessions := fakeSessions()
	// Inner is always called with empty opts (full list)
	inner.On("List", lister.ListOptions{}).Return(sessions, nil).Once()

	cl := lister.NewCachingLister(inner, fc)
	got, err := cl.List(lister.ListOptions{Tmux: true})
	require.NoError(t, err)
	assert.Equal(t, sessions.OrderedIndex, got.OrderedIndex)

	cached, err := fc.Read()
	require.NoError(t, err)
	assert.Equal(t, sessions.OrderedIndex, cached.Sessions.OrderedIndex)

	cl.Wait()
}

func TestCachingLister_WarmHit(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)
	opts := lister.ListOptions{Tmux: true}

	sessions := fakeSessions()
	require.NoError(t, fc.Write(sessions))

	cl := lister.NewCachingLister(inner, fc)
	got, err := cl.List(opts)
	require.NoError(t, err)
	assert.Equal(t, sessions.OrderedIndex, got.OrderedIndex)

	cl.Wait()
	inner.AssertNotCalled(t, "List")
}

func TestCachingLister_RefreshThenWarmRead(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)
	opts := lister.ListOptions{Tmux: true}

	fresh := updatedSessions()
	inner.On("List", opts).Return(fresh, nil).Once()

	cl := lister.NewCachingLister(inner, fc)
	cl.RefreshCache(opts)
	cl.Wait()

	got, err := cl.List(opts)
	require.NoError(t, err)
	assert.Equal(t, fresh.OrderedIndex, got.OrderedIndex)
	assert.Equal(t, fresh.Directory, got.Directory)
}

func TestCachingLister_ColdStartError(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	// Inner is always called with empty opts
	inner.On("List", lister.ListOptions{}).Return(model.SeshSessions{}, fmt.Errorf("tmux not running")).Once()

	cl := lister.NewCachingLister(inner, fc)
	_, err := cl.List(lister.ListOptions{Tmux: true})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tmux not running")

	cl.Wait()
}

func TestCachingLister_RefreshCache(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)
	opts := lister.ListOptions{Tmux: true}

	freshSessions := updatedSessions()
	inner.On("List", opts).Return(freshSessions, nil).Once()

	cl := lister.NewCachingLister(inner, fc)
	cl.RefreshCache(opts)
	cl.Wait()

	cached, err := fc.Read()
	require.NoError(t, err)
	assert.Equal(t, freshSessions.OrderedIndex, cached.Sessions.OrderedIndex)
}

func TestCachingLister_DelegatesMethods(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	session := model.SeshSession{Src: "tmux", Name: "test"}
	inner.On("FindTmuxSession", "test").Return(session, true)
	inner.On("GetAttachedTmuxSession").Return(session, true)
	inner.On("GetLastTmuxSession").Return(session, true)
	inner.On("FindConfigSession", "test").Return(session, true)
	inner.On("FindConfigWildcard", "/path").Return(model.WildcardConfig{}, false)
	inner.On("FindZoxideSession", "test").Return(session, true)
	inner.On("FindTmuxinatorConfig", "test").Return(session, true)

	cl := lister.NewCachingLister(inner, fc)

	s, ok := cl.FindTmuxSession("test")
	assert.True(t, ok)
	assert.Equal(t, "test", s.Name)

	s, ok = cl.GetAttachedTmuxSession()
	assert.True(t, ok)

	s, ok = cl.GetLastTmuxSession()
	assert.True(t, ok)

	s, ok = cl.FindConfigSession("test")
	assert.True(t, ok)

	_, ok = cl.FindConfigWildcard("/path")
	assert.False(t, ok)

	s, ok = cl.FindZoxideSession("test")
	assert.True(t, ok)

	s, ok = cl.FindTmuxinatorConfig("test")
	assert.True(t, ok)
}

func sessionsWithAttached() model.SeshSessions {
	return model.SeshSessions{
		OrderedIndex: []string{"tmux:main", "tmux:dev", "tmux:work"},
		Directory: model.SeshSessionMap{
			"tmux:main": {Src: "tmux", Name: "main", Path: "/home/user", Attached: 1},
			"tmux:dev":  {Src: "tmux", Name: "dev", Path: "/home/user/dev"},
			"tmux:work": {Src: "tmux", Name: "work", Path: "/home/user/work"},
		},
	}
}

func TestCachingLister_HideAttached_ColdStart(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	sessions := sessionsWithAttached()
	// Inner is always called with empty opts (cache stores full list)
	inner.On("List", lister.ListOptions{}).Return(sessions, nil).Once()
	inner.On("GetAttachedTmuxSession").Return(sessions.Directory["tmux:main"], true)

	cl := lister.NewCachingLister(inner, fc)
	got, err := cl.List(lister.ListOptions{Tmux: true, HideAttached: true})
	require.NoError(t, err)
	// "main" should be filtered out
	assert.Equal(t, []string{"tmux:dev", "tmux:work"}, got.OrderedIndex)

	// Cache should still contain the full unfiltered list
	cached, err := fc.Read()
	require.NoError(t, err)
	assert.Equal(t, []string{"tmux:main", "tmux:dev", "tmux:work"}, cached.Sessions.OrderedIndex)

	cl.Wait()
}

func TestCachingLister_HideAttached_WarmHit(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	sessions := sessionsWithAttached()
	require.NoError(t, fc.Write(sessions))

	inner.On("GetAttachedTmuxSession").Return(sessions.Directory["tmux:main"], true)

	cl := lister.NewCachingLister(inner, fc)
	got, err := cl.List(lister.ListOptions{Tmux: true, HideAttached: true})
	require.NoError(t, err)
	assert.Equal(t, []string{"tmux:dev", "tmux:work"}, got.OrderedIndex)

	cl.Wait()
	inner.AssertNotCalled(t, "List")
}

func TestCachingLister_HideAttached_NoAttachedSession(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	sessions := fakeSessions()
	require.NoError(t, fc.Write(sessions))

	inner.On("GetAttachedTmuxSession").Return(model.SeshSession{}, false)

	cl := lister.NewCachingLister(inner, fc)
	got, err := cl.List(lister.ListOptions{Tmux: true, HideAttached: true})
	require.NoError(t, err)
	// Nothing filtered since no session is attached
	assert.Equal(t, sessions.OrderedIndex, got.OrderedIndex)

	cl.Wait()
}

func TestCachingLister_WaitBlocksUntilDone(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)
	opts := lister.ListOptions{}

	sessions := fakeSessions()
	inner.On("List", opts).Return(sessions, nil).Run(func(_ mock.Arguments) {
		time.Sleep(50 * time.Millisecond)
	})

	cl := lister.NewCachingLister(inner, fc)
	cl.RefreshCache(opts)

	cl.Wait()

	cached, err := fc.Read()
	require.NoError(t, err)
	assert.Equal(t, sessions.OrderedIndex, cached.Sessions.OrderedIndex)
}

// --- Source filtering and dedup tests ---

func mixedSessions() model.SeshSessions {
	return model.SeshSessions{
		OrderedIndex: []string{"tmux:main", "config:project", "zoxide:code"},
		Directory: model.SeshSessionMap{
			"tmux:main":      {Src: "tmux", Name: "main", Path: "/home/user"},
			"config:project": {Src: "config", Name: "project", Path: "/home/user/project"},
			"zoxide:code":    {Src: "zoxide", Name: "code", Path: "/home/user/code"},
		},
	}
}

func TestCachingLister_SourceFilter_FromCache(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	// Pre-populate cache with mixed-source sessions
	sessions := mixedSessions()
	require.NoError(t, fc.Write(sessions))

	cl := lister.NewCachingLister(inner, fc)

	// Request only tmux — should filter from full cache
	got, err := cl.List(lister.ListOptions{Tmux: true})
	require.NoError(t, err)
	assert.Equal(t, []string{"tmux:main"}, got.OrderedIndex)

	cl.Wait()
	inner.AssertNotCalled(t, "List")
}

func TestCachingLister_SourceFilter_DifferentFilters(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	// Pre-populate cache with mixed-source sessions
	sessions := mixedSessions()
	require.NoError(t, fc.Write(sessions))

	cl := lister.NewCachingLister(inner, fc)

	// First call: tmux only
	got, err := cl.List(lister.ListOptions{Tmux: true})
	require.NoError(t, err)
	assert.Equal(t, []string{"tmux:main"}, got.OrderedIndex)

	// Second call: zoxide only — must NOT return tmux sessions
	got, err = cl.List(lister.ListOptions{Zoxide: true})
	require.NoError(t, err)
	assert.Equal(t, []string{"zoxide:code"}, got.OrderedIndex)

	// Third call: no source flags — returns all
	got, err = cl.List(lister.ListOptions{})
	require.NoError(t, err)
	assert.Equal(t, sessions.OrderedIndex, got.OrderedIndex)

	cl.Wait()
	inner.AssertNotCalled(t, "List")
}

func TestCachingLister_SourceFilter_ColdStart(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	sessions := mixedSessions()
	// Inner is called with empty opts (fetches all sources)
	inner.On("List", lister.ListOptions{}).Return(sessions, nil).Once()

	cl := lister.NewCachingLister(inner, fc)

	// Request only zoxide — inner fetches all, cache stores all, filter returns zoxide
	got, err := cl.List(lister.ListOptions{Zoxide: true})
	require.NoError(t, err)
	assert.Equal(t, []string{"zoxide:code"}, got.OrderedIndex)

	// Cache should contain all sources
	cached, err := fc.Read()
	require.NoError(t, err)
	assert.Equal(t, sessions.OrderedIndex, cached.Sessions.OrderedIndex)

	cl.Wait()
}

func TestCachingLister_HideDuplicates_PostCache(t *testing.T) {
	dir := t.TempDir()
	fc := cache.NewFileCacheWithPath(filepath.Join(dir, "sessions.gob"))
	inner := lister.NewMockLister(t)

	// Sessions with duplicate name and path across sources
	sessions := model.SeshSessions{
		OrderedIndex: []string{"tmux:project", "config:project", "zoxide:other"},
		Directory: model.SeshSessionMap{
			"tmux:project":   {Src: "tmux", Name: "project", Path: "/home/user/project"},
			"config:project": {Src: "config", Name: "project", Path: "/home/user/project"},
			"zoxide:other":   {Src: "zoxide", Name: "other", Path: "/home/user/other"},
		},
	}
	require.NoError(t, fc.Write(sessions))

	cl := lister.NewCachingLister(inner, fc)

	// With HideDuplicates, "config:project" is a duplicate of "tmux:project" (same name + path)
	got, err := cl.List(lister.ListOptions{HideDuplicates: true})
	require.NoError(t, err)
	assert.Equal(t, []string{"tmux:project", "zoxide:other"}, got.OrderedIndex)

	// Without HideDuplicates, all returned
	got, err = cl.List(lister.ListOptions{})
	require.NoError(t, err)
	assert.Equal(t, sessions.OrderedIndex, got.OrderedIndex)

	cl.Wait()
	inner.AssertNotCalled(t, "List")
}
