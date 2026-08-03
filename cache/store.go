package cache

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// Dir returns sesh's cache directory: $XDG_CACHE_HOME/sesh, falling back to
// ~/.cache/sesh.
func Dir() string {
	dir := os.Getenv("XDG_CACHE_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		dir = filepath.Join(home, ".cache")
	}
	return filepath.Join(dir, "sesh")
}

// writeAtomic writes data to path via a uniquely named temp file in the same
// directory, then renames it into place. The unique name means two concurrent
// sesh processes writing the same cache resolve to last-rename-wins rather than
// corrupting a shared temp file.
func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("cache mkdir: %w", err)
	}

	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("cache create tmp: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath) // no-op once the rename succeeds

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("cache write tmp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("cache close tmp: %w", err)
	}
	if err := os.Chmod(tmpPath, 0o644); err != nil {
		return fmt.Errorf("cache chmod tmp: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("cache rename: %w", err)
	}
	return nil
}

// Entry is a cached value along with the time it was fetched.
type Entry[T any] struct {
	Value     T         `json:"value"`
	FetchedAt time.Time `json:"fetched_at"`
	// Missing records that the source confirmed there is no value for this
	// key, so lookups can be cached negatively instead of refetched forever.
	Missing bool `json:"missing,omitempty"`
}

// Entries maps cache keys to their entries.
type Entries[T any] map[string]Entry[T]

// Put stores val under key, marking it as fetched now.
func (e Entries[T]) Put(key string, val T) {
	e[key] = Entry[T]{Value: val, FetchedAt: time.Now()}
}

// PutMissing records that key has no value at the source.
func (e Entries[T]) PutMissing(key string) {
	e[key] = Entry[T]{FetchedAt: time.Now(), Missing: true}
}

// Namespace is a keyed, TTL'd cache of T backed by a single JSON file.
//
// The version is part of the filename, so changing the shape of T is handled by
// bumping the version rather than migrating: the new version reads a file that
// does not exist yet, starts empty, and refetches. Cached values are always
// refetchable, so discarding them is never lossy.
type Namespace[T any] struct {
	name       string
	path       string
	ttl        time.Duration
	missingTTL time.Duration
}

// NewNamespace creates a namespace stored at <cache dir>/<name>.v<version>.json.
// A zero ttl means entries never expire.
func NewNamespace[T any](name string, version int, ttl time.Duration) *Namespace[T] {
	return NewNamespaceInDir[T](Dir(), name, version, ttl)
}

// NewNamespaceInDir creates a namespace rooted at a specific directory (useful
// for testing).
func NewNamespaceInDir[T any](dir, name string, version int, ttl time.Duration) *Namespace[T] {
	return &Namespace[T]{
		name:       name,
		path:       filepath.Join(dir, fmt.Sprintf("%s.v%d.json", name, version)),
		ttl:        ttl,
		missingTTL: ttl,
	}
}

// WithMissingTTL sets a separate, usually shorter, TTL for negative entries.
func (n *Namespace[T]) WithMissingTTL(ttl time.Duration) *Namespace[T] {
	n.missingTTL = ttl
	return n
}

// Path returns the file backing this namespace.
func (n *Namespace[T]) Path() string {
	return n.path
}

// Load reads every entry in the namespace. A namespace that has never been
// written yields an empty set and no error. Unreadable or corrupt files also
// yield an empty set, so a damaged cache degrades to a cold one instead of
// failing the caller.
func (n *Namespace[T]) Load() Entries[T] {
	data, err := os.ReadFile(n.path)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Debug("cache: could not read namespace", "path", n.path, "error", err)
		}
		return Entries[T]{}
	}

	entries := Entries[T]{}
	if err := json.Unmarshal(data, &entries); err != nil {
		slog.Debug("cache: discarding corrupt namespace", "path", n.path, "error", err)
		return Entries[T]{}
	}
	return entries
}

// Save writes entries to disk atomically.
func (n *Namespace[T]) Save(entries Entries[T]) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("cache encode %s: %w", n.name, err)
	}
	return writeAtomic(n.path, data)
}

// Fresh reports whether an entry is still within its TTL.
func (n *Namespace[T]) Fresh(entry Entry[T]) bool {
	ttl := n.ttl
	if entry.Missing {
		ttl = n.missingTTL
	}
	if ttl <= 0 {
		return true
	}
	return time.Since(entry.FetchedAt) < ttl
}

// Lookup returns the cached value for key. found is true for expired entries
// too, so callers can render stale data immediately and refresh behind it;
// check fresh to decide whether a refetch is needed. A negatively cached key
// reports found as false.
func (n *Namespace[T]) Lookup(entries Entries[T], key string) (val T, found bool, fresh bool) {
	entry, ok := entries[key]
	if !ok {
		var zero T
		return zero, false, false
	}
	if entry.Missing {
		var zero T
		return zero, false, n.Fresh(entry)
	}
	return entry.Value, true, n.Fresh(entry)
}

// Stale returns the subset of keys that are absent or past their TTL — the set
// a caller should refetch, ideally in one batch.
func (n *Namespace[T]) Stale(entries Entries[T], keys []string) []string {
	var stale []string
	for _, key := range keys {
		entry, ok := entries[key]
		if !ok || !n.Fresh(entry) {
			stale = append(stale, key)
		}
	}
	return stale
}

// Prune deletes files for this namespace at other versions. Best-effort: a
// leftover file is harmless, so failures are logged rather than returned.
func (n *Namespace[T]) Prune() {
	pattern := filepath.Join(filepath.Dir(n.path), n.name+".v*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		slog.Debug("cache: could not glob old versions", "pattern", pattern, "error", err)
		return
	}
	for _, match := range matches {
		if match == n.path {
			continue
		}
		if err := os.Remove(match); err != nil {
			slog.Debug("cache: could not remove old version", "path", match, "error", err)
		}
	}
}
