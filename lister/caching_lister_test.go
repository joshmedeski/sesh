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
	opts := lister.ListOptions{Tmux: true}

	sessions := fakeSessions()
	inner.On("List", opts).Return(sessions, nil).Once()

	cl := lister.NewCachingLister(inner, fc)
	got, err := cl.List(opts)
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
	opts := lister.ListOptions{Tmux: true}

	inner.On("List", opts).Return(model.SeshSessions{}, fmt.Errorf("tmux not running")).Once()

	cl := lister.NewCachingLister(inner, fc)
	_, err := cl.List(opts)
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
	// Inner should be called WITHOUT HideAttached (cache stores full list)
	inner.On("List", lister.ListOptions{Tmux: true, HideAttached: false}).Return(sessions, nil).Once()
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
