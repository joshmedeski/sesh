package lister_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

func TestCachingLister_RefreshCacheWritesTheSlowSourcesWithoutTouchingTmux(t *testing.T) {
	f := newFixture(t, nil)
	cl, fc := f.cached()

	cl.RefreshCache()
	cl.Wait()

	stored, err := fc.Read()
	require.NoError(t, err)
	assert.Equal(t, []string{"tmuxinator:work", "zoxide:~/code", "zoxide:~/dev", "zoxide:~/project"}, stored.Sessions.OrderedIndex)
	assert.Zero(t, f.tmux.sessionCalls.Load())
}

func TestCachingLister_RefreshCacheKeepsTheCacheWhenAFetchFails(t *testing.T) {
	f := newFixture(t, nil)
	stored := storedZoxide("/home/user/old")
	writeCacheFile(t, f.cachePath, stored, time.Now())
	f.zoxide.fail(errors.New("zoxide is broken"))
	cl, fc := f.cached()

	cl.RefreshCache()
	cl.Wait()

	got, err := fc.Read()
	require.NoError(t, err)
	assert.Equal(t, stored.OrderedIndex, got.Sessions.OrderedIndex)
}

func TestCachingLister_RefreshCacheThenListUsesTheRefreshedSources(t *testing.T) {
	f := newFixture(t, nil)
	f.config.SessionConfigs = nil
	f.tmuxinator.configs = nil
	f.zoxide.set(&model.ZoxideResult{Path: "/home/user/fresh", Score: 1})
	cl, _ := f.cached()

	cl.RefreshCache()
	cl.Wait()

	assert.Equal(t, []string{"main", "dev", "~/fresh"}, listNames(t, cl, lister.ListOptions{}))
	assert.EqualValues(t, 1, f.zoxide.calls.Load(), "the list read the refreshed cache")
}

func TestCachingLister_WaitReturnsOnlyOnceTheRefreshHasWritten(t *testing.T) {
	f := newFixture(t, nil)
	started, release := f.zoxide.hold()
	t.Cleanup(release)
	cl, fc := f.cached()

	cl.RefreshCache()
	<-started
	go release()
	cl.Wait()

	_, err := fc.Read()
	assert.NoError(t, err)
}

func TestCachingLister_RefreshCacheSurvivesAnUnwritableCache(t *testing.T) {
	f := newFixture(t, nil)
	blocker := filepath.Join(t.TempDir(), "file")
	require.NoError(t, os.WriteFile(blocker, nil, 0o644))
	f.cachePath = filepath.Join(blocker, "sessions.v2.gob")
	cl, _ := f.cached()

	cl.RefreshCache()
	cl.Wait()

	assert.EqualValues(t, 1, f.zoxide.calls.Load())
}
