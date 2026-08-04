package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testIssue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

func testNamespace(t *testing.T, ttl time.Duration) *Namespace[testIssue] {
	t.Helper()
	return NewNamespaceInDir[testIssue](t.TempDir(), "issues", 1, ttl)
}

func TestNamespace_SaveAndLoad(t *testing.T) {
	n := testNamespace(t, time.Hour)

	entries := Entries[testIssue]{}
	entries.Put("sesh#409", testIssue{Number: 409, Title: "worktree support"})
	require.NoError(t, n.Save(entries))

	got := n.Load()
	require.Len(t, got, 1)
	assert.Equal(t, testIssue{Number: 409, Title: "worktree support"}, got["sesh#409"].Value)
	assert.WithinDuration(t, time.Now(), got["sesh#409"].FetchedAt, 2*time.Second)
}

func TestNamespace_LoadMissingFileIsEmpty(t *testing.T) {
	n := testNamespace(t, time.Hour)

	got := n.Load()
	assert.Empty(t, got)
	assert.NotNil(t, got, "callers should be able to write into the result")
}

func TestNamespace_LoadCorruptFileIsEmpty(t *testing.T) {
	n := testNamespace(t, time.Hour)
	require.NoError(t, os.MkdirAll(filepath.Dir(n.Path()), 0o755))
	require.NoError(t, os.WriteFile(n.Path(), []byte("{not json"), 0o644))

	got := n.Load()
	assert.Empty(t, got)
}

func TestNamespace_VersionIsInFilename(t *testing.T) {
	dir := t.TempDir()
	v1 := NewNamespaceInDir[testIssue](dir, "issues", 1, time.Hour)
	v2 := NewNamespaceInDir[testIssue](dir, "issues", 2, time.Hour)

	assert.Equal(t, filepath.Join(dir, "issues.v1.json"), v1.Path())
	assert.Equal(t, filepath.Join(dir, "issues.v2.json"), v2.Path())

	entries := Entries[testIssue]{}
	entries.Put("sesh#409", testIssue{Number: 409, Title: "worktree support"})
	require.NoError(t, v1.Save(entries))

	// Bumping the version starts cold rather than reading v1's data.
	assert.Empty(t, v2.Load())
}

func TestNamespace_Prune(t *testing.T) {
	dir := t.TempDir()
	v1 := NewNamespaceInDir[testIssue](dir, "issues", 1, time.Hour)
	v2 := NewNamespaceInDir[testIssue](dir, "issues", 2, time.Hour)
	other := NewNamespaceInDir[testIssue](dir, "pr-state", 1, time.Hour)

	require.NoError(t, v1.Save(Entries[testIssue]{}))
	require.NoError(t, v2.Save(Entries[testIssue]{}))
	require.NoError(t, other.Save(Entries[testIssue]{}))

	v2.Prune()

	_, err := os.Stat(v1.Path())
	assert.True(t, os.IsNotExist(err), "old version should be removed")

	_, err = os.Stat(v2.Path())
	assert.NoError(t, err, "current version should survive")

	_, err = os.Stat(other.Path())
	assert.NoError(t, err, "other namespaces should be untouched")
}

func TestNamespace_LookupFreshAndStale(t *testing.T) {
	n := testNamespace(t, time.Hour)
	entries := Entries[testIssue]{
		"fresh": {Value: testIssue{Number: 1, Title: "fresh"}, FetchedAt: time.Now()},
		"stale": {Value: testIssue{Number: 2, Title: "stale"}, FetchedAt: time.Now().Add(-2 * time.Hour)},
	}

	val, found, fresh := n.Lookup(entries, "fresh")
	assert.True(t, found)
	assert.True(t, fresh)
	assert.Equal(t, "fresh", val.Title)

	// Stale entries are still returned so a picker can render them immediately.
	val, found, fresh = n.Lookup(entries, "stale")
	assert.True(t, found)
	assert.False(t, fresh)
	assert.Equal(t, "stale", val.Title)

	val, found, fresh = n.Lookup(entries, "absent")
	assert.False(t, found)
	assert.False(t, fresh)
	assert.Zero(t, val)
}

func TestNamespace_ZeroTTLNeverExpires(t *testing.T) {
	n := testNamespace(t, 0)
	entries := Entries[testIssue]{
		"old": {Value: testIssue{Number: 1}, FetchedAt: time.Now().Add(-100 * time.Hour)},
	}

	_, _, fresh := n.Lookup(entries, "old")
	assert.True(t, fresh)
	assert.Empty(t, n.Stale(entries, []string{"old"}))
}

func TestNamespace_MissingIsCachedNegatively(t *testing.T) {
	n := testNamespace(t, time.Hour)

	entries := Entries[testIssue]{}
	entries.PutMissing("sesh#99999")
	require.NoError(t, n.Save(entries))

	got := n.Load()
	require.True(t, got["sesh#99999"].Missing)

	// A negatively cached key has no value, but is not stale — so it does not
	// get refetched on every render.
	_, found, fresh := n.Lookup(got, "sesh#99999")
	assert.False(t, found)
	assert.True(t, fresh)
	assert.Empty(t, n.Stale(got, []string{"sesh#99999"}))
}

func TestNamespace_MissingTTLIsSeparate(t *testing.T) {
	n := testNamespace(t, time.Hour).WithMissingTTL(time.Minute)
	entries := Entries[testIssue]{
		"value":   {Value: testIssue{Number: 1}, FetchedAt: time.Now().Add(-30 * time.Minute)},
		"missing": {FetchedAt: time.Now().Add(-30 * time.Minute), Missing: true},
	}

	assert.Equal(t, []string{"missing"}, n.Stale(entries, []string{"value", "missing"}))
}

func TestNamespace_StaleReturnsAbsentAndExpired(t *testing.T) {
	n := testNamespace(t, time.Hour)
	entries := Entries[testIssue]{
		"fresh": {Value: testIssue{Number: 1}, FetchedAt: time.Now()},
		"old":   {Value: testIssue{Number: 2}, FetchedAt: time.Now().Add(-2 * time.Hour)},
	}

	stale := n.Stale(entries, []string{"fresh", "old", "absent"})
	assert.Equal(t, []string{"old", "absent"}, stale)
}

func TestNamespace_SaveIsAtomic(t *testing.T) {
	n := testNamespace(t, time.Hour)

	first := Entries[testIssue]{}
	first.Put("a", testIssue{Number: 1, Title: "one"})
	require.NoError(t, n.Save(first))

	second := Entries[testIssue]{}
	second.Put("b", testIssue{Number: 2, Title: "two"})
	require.NoError(t, n.Save(second))

	matches, err := filepath.Glob(filepath.Join(filepath.Dir(n.Path()), "*.tmp"))
	require.NoError(t, err)
	assert.Empty(t, matches, "no temp files left behind")

	got := n.Load()
	require.Len(t, got, 1)
	assert.Equal(t, "two", got["b"].Value.Title)
}

func TestNamespace_FileIsHumanReadable(t *testing.T) {
	n := testNamespace(t, time.Hour)

	entries := Entries[testIssue]{}
	entries.Put("sesh#409", testIssue{Number: 409, Title: "worktree support"})
	require.NoError(t, n.Save(entries))

	data, err := os.ReadFile(n.Path())
	require.NoError(t, err)

	var raw map[string]struct {
		Value     testIssue `json:"value"`
		FetchedAt time.Time `json:"fetched_at"`
	}
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.Equal(t, "worktree support", raw["sesh#409"].Value.Title)
}

func TestNewNamespace_UsesCacheDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)

	n := NewNamespace[testIssue]("github-issues", 1, time.Hour)
	assert.Equal(t, filepath.Join(dir, "sesh", "github-issues.v1.json"), n.Path())
}

func TestDir_XDGCacheHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)

	assert.Equal(t, filepath.Join(dir, "sesh"), Dir())
}

func TestDir_FallbackToHomeCache(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("HOME", "/fakehome")

	assert.Contains(t, Dir(), filepath.Join(".cache", "sesh"))
}
