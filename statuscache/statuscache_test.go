package statuscache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCache(t *testing.T) *FileStatusCache {
	t.Helper()
	return NewFileStatusCacheWithDir(filepath.Join(t.TempDir(), "status"))
}

func TestKeyStableAndDistinct(t *testing.T) {
	a := Key("/repo/one", "400")
	assert.Equal(t, a, Key("/repo/one", "400"), "same inputs => same key")
	assert.NotEqual(t, a, Key("/repo/two", "400"), "different repo => different key")
	assert.NotEqual(t, a, Key("/repo/one", "401"), "different branch => different key")
}

func TestWriteReadRoundTrip(t *testing.T) {
	c := newTestCache(t)
	entry := Entry{Issue: &Ref{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"}}
	require.NoError(t, c.Write("k1", entry))

	got, found, err := c.Read("k1")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Nil(t, got.PR)
	assert.Equal(t, &Ref{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"}, got.Issue)
}

func TestNegativeEntryRoundTrips(t *testing.T) {
	c := newTestCache(t)
	require.NoError(t, c.Write("k2", Entry{}))

	got, found, err := c.Read("k2")
	require.NoError(t, err)
	assert.True(t, found, "negative entry is still a hit")
	assert.Nil(t, got.PR)
	assert.Nil(t, got.Issue)
}

func TestMissingFileIsMiss(t *testing.T) {
	c := newTestCache(t)
	_, found, err := c.Read("does-not-exist")
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestCorruptFileIsMiss(t *testing.T) {
	c := newTestCache(t)
	require.NoError(t, c.Write("k3", Entry{Issue: &Ref{Number: 1}}))
	// Corrupt the file on disk.
	require.NoError(t, writeRaw(c, "k3", []byte("not gob data")))

	_, found, err := c.Read("k3")
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestPreferred(t *testing.T) {
	pr := &Ref{Number: 401, Title: "pr", State: "OPEN"}
	iss := &Ref{Number: 400, Title: "iss", State: "OPEN"}

	got, ok := Entry{PR: pr, Issue: iss}.Preferred()
	assert.True(t, ok)
	assert.Equal(t, pr, got, "PR preferred over issue")

	got, ok = Entry{Issue: iss}.Preferred()
	assert.True(t, ok)
	assert.Equal(t, iss, got)

	_, ok = Entry{}.Preferred()
	assert.False(t, ok)
}

func writeRaw(c *FileStatusCache, key string, data []byte) error {
	return os.WriteFile(c.path(key), data, 0o644)
}
