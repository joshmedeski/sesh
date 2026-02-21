package lister

import (
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/model"
)

const softTTL = 5 * time.Second

// CachingLister wraps a Lister with stale-while-revalidate caching for List().
// All other Lister methods delegate directly to the inner lister.
type CachingLister struct {
	inner Lister
	cache cache.Cache
	wg    sync.WaitGroup
}

// NewCachingLister creates a CachingLister that decorates inner with file-based caching.
func NewCachingLister(inner Lister, c cache.Cache) *CachingLister {
	return &CachingLister{inner: inner, cache: c}
}

// List implements Lister. It returns cached data when available, triggering a
// background refresh when the cache is older than the soft TTL.
// The cache always stores the full unfiltered list; source filters,
// HideDuplicates, and HideAttached are applied after reading from cache.
func (cl *CachingLister) List(opts ListOptions) (model.SeshSessions, error) {
	// Always fetch/store the full list (all sources, no dedup, no hide);
	// we apply all view-level filters ourselves after cache read.
	innerOpts := ListOptions{}

	cached, err := cl.cache.Read()
	if err == nil {
		age := time.Since(cached.Timestamp)
		if age < softTTL {
			slog.Debug("cache: hit (fresh)", "age", age)
			return cl.applyFilters(cached.Sessions, opts), nil
		}
		// Stale -- return immediately but revalidate in background
		slog.Debug("cache: hit (stale, revalidating)", "age", age)
		cl.wg.Add(1)
		go func() {
			defer cl.wg.Done()
			cl.revalidate(innerOpts)
		}()
		return cl.applyFilters(cached.Sessions, opts), nil
	}

	// Cold start -- fetch synchronously
	slog.Debug("cache: miss (cold start)")
	sessions, err := cl.inner.List(innerOpts)
	if err != nil {
		return sessions, err
	}
	if writeErr := cl.cache.Write(sessions); writeErr != nil {
		slog.Warn("cache: write failed on cold start", "error", writeErr)
	}
	return cl.applyFilters(sessions, opts), nil
}

// applyFilters applies view-level filters that should not affect what gets
// stored in the cache: source filtering, HideDuplicates, and HideAttached.
func (cl *CachingLister) applyFilters(sessions model.SeshSessions, opts ListOptions) model.SeshSessions {
	filtered := sessions.OrderedIndex

	// 1. Source filtering: if any source flag is set, keep only matching sessions.
	if allowed := sourceSet(opts); allowed != nil {
		result := make([]string, 0, len(filtered))
		for _, index := range filtered {
			if allowed[sessions.Directory[index].Src] {
				result = append(result, index)
			}
		}
		filtered = result
	}

	// 2. HideDuplicates: deduplicate by name, then by path.
	if opts.HideDuplicates {
		nameHash := make(map[string]bool)
		pathHash := make(map[string]bool)
		destIndex := 0
		for _, index := range filtered {
			session := sessions.Directory[index]
			nameIsDuplicate := nameHash[session.Name]
			pathIsDuplicate := session.Path != "" && pathHash[session.Path]
			if !nameIsDuplicate && !pathIsDuplicate {
				filtered[destIndex] = index
				nameHash[session.Name] = true
				pathHash[session.Path] = true
				destIndex++
			}
		}
		filtered = filtered[:destIndex]
	}

	// 3. HideAttached: remove the currently attached tmux session.
	if opts.HideAttached {
		attached, ok := cl.inner.GetAttachedTmuxSession()
		if ok {
			for i, index := range filtered {
				if sessions.Directory[index].Name == attached.Name {
					filtered = slices.Delete(slices.Clone(filtered), i, i+1)
					break
				}
			}
		}
	}

	if len(filtered) == len(sessions.OrderedIndex) {
		return sessions
	}
	return model.SeshSessions{
		OrderedIndex: filtered,
		Directory:    sessions.Directory,
	}
}

// sourceSet returns a set of allowed source names based on opts, or nil if
// no source flags are set (meaning all sources are allowed).
func sourceSet(opts ListOptions) map[string]bool {
	if !opts.Tmux && !opts.Config && !opts.Zoxide && !opts.Tmuxinator {
		return nil
	}
	m := make(map[string]bool)
	if opts.Tmux {
		m["tmux"] = true
	}
	if opts.Config {
		m["config"] = true
	}
	if opts.Zoxide {
		m["zoxide"] = true
	}
	if opts.Tmuxinator {
		m["tmuxinator"] = true
	}
	return m
}

func (cl *CachingLister) revalidate(opts ListOptions) {
	sessions, err := cl.inner.List(opts)
	if err != nil {
		slog.Warn("cache: background revalidation fetch failed", "error", err)
		return
	}
	if err := cl.cache.Write(sessions); err != nil {
		slog.Warn("cache: background revalidation write failed", "error", err)
	}
}

// RefreshCache fetches live data from the inner lister and writes it to the cache.
// This bypasses the cache read entirely, intended for use after sesh connect.
func (cl *CachingLister) RefreshCache(opts ListOptions) {
	cl.wg.Add(1)
	go func() {
		defer cl.wg.Done()
		cl.revalidate(opts)
	}()
}

// Wait blocks until all background refresh goroutines have completed.
// Call this before process exit to avoid truncated cache writes.
func (cl *CachingLister) Wait() {
	cl.wg.Wait()
}

// --- Delegate all other Lister methods to inner ---

func (cl *CachingLister) FindTmuxSession(name string) (model.SeshSession, bool) {
	return cl.inner.FindTmuxSession(name)
}

func (cl *CachingLister) GetAttachedTmuxSession() (model.SeshSession, bool) {
	return cl.inner.GetAttachedTmuxSession()
}

func (cl *CachingLister) GetLastTmuxSession() (model.SeshSession, bool) {
	return cl.inner.GetLastTmuxSession()
}

func (cl *CachingLister) FindConfigSession(name string) (model.SeshSession, bool) {
	return cl.inner.FindConfigSession(name)
}

func (cl *CachingLister) FindConfigWildcard(path string) (model.WildcardConfig, bool) {
	return cl.inner.FindConfigWildcard(path)
}

func (cl *CachingLister) FindZoxideSession(name string) (model.SeshSession, bool) {
	return cl.inner.FindZoxideSession(name)
}

func (cl *CachingLister) FindTmuxinatorConfig(name string) (model.SeshSession, bool) {
	return cl.inner.FindTmuxinatorConfig(name)
}
