package statuscache

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"
)

// Ref is one GitHub entity (issue or PR) as rendered in the status bar.
type Ref struct {
	Number int
	Title  string
	State  string
}

// Entry is the cached status for one branch. Both pointers may be nil (a
// "negative" entry: the branch has nothing to show), which still counts as a
// cache hit so the reader respects the TTL instead of refreshing every tick.
type Entry struct {
	PR        *Ref
	Issue     *Ref
	Timestamp time.Time
}

// Preferred returns the entity to render: the PR if present, otherwise the
// issue. ok is false for a negative entry.
func (e Entry) Preferred() (*Ref, bool) {
	if e.PR != nil {
		return e.PR, true
	}
	if e.Issue != nil {
		return e.Issue, true
	}
	return nil, false
}

// StatusCache reads and writes per-branch status entries.
type StatusCache interface {
	Read(key string) (Entry, bool, error) // bool=found; false (nil err) on miss/corrupt
	Write(key string, entry Entry) error
}

// Key builds the cache key (and filename stem) for a repo root + branch.
func Key(repoRoot, branch string) string {
	sum := sha256.Sum256([]byte(repoRoot + "\x00" + branch))
	return hex.EncodeToString(sum[:])
}

// FileStatusCache stores one gob file per key under a directory.
type FileStatusCache struct {
	dir string
}

// NewFileStatusCache stores entries under $XDG_CACHE_HOME/sesh/status
// (falling back to ~/.cache/sesh/status).
func NewFileStatusCache() *FileStatusCache {
	dir := os.Getenv("XDG_CACHE_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		dir = filepath.Join(home, ".cache")
	}
	return &FileStatusCache{dir: filepath.Join(dir, "sesh", "status")}
}

// NewFileStatusCacheWithDir stores entries under an explicit directory (tests).
func NewFileStatusCacheWithDir(dir string) *FileStatusCache {
	return &FileStatusCache{dir: dir}
}

func (c *FileStatusCache) path(key string) string {
	return filepath.Join(c.dir, key+".gob")
}

func (c *FileStatusCache) Read(key string) (Entry, bool, error) {
	data, err := os.ReadFile(c.path(key))
	if err != nil {
		return Entry{}, false, nil // missing → miss
	}
	var entry Entry
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&entry); err != nil {
		return Entry{}, false, nil // corrupt → miss
	}
	return entry, true, nil
}

func (c *FileStatusCache) Write(key string, entry Entry) error {
	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(entry); err != nil {
		return err
	}
	tmp := c.path(key) + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, c.path(key))
}
