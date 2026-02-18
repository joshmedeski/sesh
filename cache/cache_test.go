package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/model"
)

func testSessions() model.SeshSessions {
	return model.SeshSessions{
		OrderedIndex: []string{"tmux:main", "zoxide:~/code"},
		Directory: model.SeshSessionMap{
			"tmux:main": {
				Src:  "tmux",
				Name: "main",
				Path: "/home/user",
			},
			"zoxide:~/code": {
				Src:   "zoxide",
				Name:  "~/code",
				Path:  "/home/user/code",
				Score: 42.5,
			},
		},
	}
}

func TestFileCache_WriteAndRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sesh", "sessions.gob")
	c := NewFileCacheWithPath(path)

	sessions := testSessions()
	err := c.Write(sessions)
	require.NoError(t, err)

	_, err = os.Stat(path)
	require.NoError(t, err)

	got, err := c.Read()
	require.NoError(t, err)

	assert.Equal(t, sessions.OrderedIndex, got.Sessions.OrderedIndex)
	assert.Equal(t, sessions.Directory, got.Sessions.Directory)
	assert.WithinDuration(t, time.Now(), got.Timestamp, 2*time.Second)
}

func TestFileCache_ReadMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.gob")
	c := NewFileCacheWithPath(path)

	_, err := c.Read()
	assert.Error(t, err)
}

func TestFileCache_AtomicWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sesh", "sessions.gob")
	c := NewFileCacheWithPath(path)

	sessions1 := testSessions()
	require.NoError(t, c.Write(sessions1))

	sessions2 := model.SeshSessions{
		OrderedIndex: []string{"tmux:other"},
		Directory: model.SeshSessionMap{
			"tmux:other": {Src: "tmux", Name: "other", Path: "/tmp"},
		},
	}
	require.NoError(t, c.Write(sessions2))

	// No temp file left behind
	_, err := os.Stat(path + ".tmp")
	assert.True(t, os.IsNotExist(err))

	got, err := c.Read()
	require.NoError(t, err)
	assert.Equal(t, sessions2.OrderedIndex, got.Sessions.OrderedIndex)
	assert.Equal(t, sessions2.Directory, got.Sessions.Directory)
}

func TestNewFileCache_XDGCacheHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)

	c := NewFileCache()
	expected := filepath.Join(dir, "sesh", "sessions.gob")
	assert.Equal(t, expected, c.path)
}

func TestNewFileCache_FallbackToHomeCache(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("HOME", "/fakehome")

	c := NewFileCache()
	assert.Contains(t, c.path, filepath.Join(".cache", "sesh", "sessions.gob"))
}
