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

// cachedSources are the sources whose results the cache keeps: the ones that
// start a slow process. Every other source is fetched live on every list:
// tmux sessions and panes because one call over a local socket takes
// milliseconds and a cached copy can show an out-of-date or empty list, and
// config sessions because they are read from the loaded configuration.
var cachedSources = []string{"tmuxinator", "zoxide"}

// CachingLister wraps a SourceLister with stale-while-revalidate caching of the
// slow sources. It implements Lister: List serves the slow sources from the
// cache and lists everything else live, and all other methods delegate to the
// inner lister.
type CachingLister struct {
	inner    SourceLister
	cache    cache.Cache
	wg       sync.WaitGroup
	skipRead bool
}

// NewCachingLister creates a CachingLister that decorates inner with file-based caching.
func NewCachingLister(inner SourceLister, c cache.Cache) *CachingLister {
	return &CachingLister{inner: inner, cache: c}
}

// List implements Lister. The cached sources come from the cache when it can be
// read, and a cache older than the soft TTL triggers a background refresh. The
// inner lister assembles the list, so source order, blacklist, HideAttached,
// HideDuplicates and formatting are exactly those of an uncached list.
func (cl *CachingLister) List(opts ListOptions) (model.SeshSessions, error) {
	// A list that shows no cached source has nothing to read from the cache,
	// and its result must not depend on whether the slow sources work.
	if !requestsCachedSource(opts) {
		if cl.skipRead {
			cl.RefreshCache()
		}
		return cl.inner.ListWith(opts, nil)
	}

	if cl.skipRead {
		slog.Debug("cache: read skipped")
		return cl.listLive(opts)
	}

	cached, err := cl.cache.Read()
	if err == nil {
		age := time.Since(cached.Timestamp)
		if age < softTTL {
			slog.Debug("cache: hit (fresh)", "age", age)
		} else {
			// Stale -- use it now but revalidate in background
			slog.Debug("cache: hit (stale, revalidating)", "age", age)
			cl.RefreshCache()
		}
		return cl.inner.ListWith(opts, splitBySource(cached.Sessions))
	}

	// Cold start -- fetch synchronously
	slog.Debug("cache: miss (cold start)")
	return cl.listLive(opts)
}

// listLive fetches the cached sources, writes them to the cache and lists with them.
func (cl *CachingLister) listLive(opts ListOptions) (model.SeshSessions, error) {
	fetched, err := cl.inner.ListSources(cachedSources)
	if err != nil {
		if requestsEveryCachedSource(opts) {
			return model.SeshSessions{}, err
		}
		// The failing source may not be one this list shows. The uncached list
		// decides: it fails only when a source it needs fails. A cached source
		// the list does show is fetched a second time there, which only happens
		// while a cached source fails.
		slog.Debug("cache: fetching the cached sources failed, listing uncached", "error", err)
		return cl.inner.ListWith(opts, nil)
	}
	if writeErr := cl.cache.Write(flatten(fetched)); writeErr != nil {
		slog.Warn("cache: write failed", "error", writeErr)
	}
	return cl.inner.ListWith(opts, fetched)
}

// SkipRead makes List always fetch live data, still writing it to the cache.
func (cl *CachingLister) SkipRead() {
	cl.skipRead = true
}

// requestsCachedSource reports whether opts ask for any cached source.
func requestsCachedSource(opts ListOptions) bool {
	return slices.ContainsFunc(srcs(opts), func(src string) bool {
		return slices.Contains(cachedSources, src)
	})
}

// requestsEveryCachedSource reports whether opts ask for all of the cached sources.
func requestsEveryCachedSource(opts ListOptions) bool {
	requested := srcs(opts)
	return !slices.ContainsFunc(cachedSources, func(src string) bool {
		return !slices.Contains(requested, src)
	})
}

// flatten joins the results of the cached sources into the one list the cache
// stores, in cachedSources order.
func flatten(fetched SourceResults) model.SeshSessions {
	flat := model.SeshSessions{
		OrderedIndex: make([]string, 0),
		Directory:    make(model.SeshSessionMap),
	}
	for _, src := range cachedSources {
		for _, key := range fetched[src].OrderedIndex {
			flat.OrderedIndex = append(flat.OrderedIndex, key)
			flat.Directory[key] = fetched[src].Directory[key]
		}
	}
	return flat
}

// splitBySource is the inverse of flatten. Every cached source gets a result,
// empty when the cache holds no row for it, so that an empty source counts as
// fetched. Rows of any other source are dropped.
func splitBySource(stored model.SeshSessions) SourceResults {
	results := make(SourceResults, len(cachedSources))
	for _, src := range cachedSources {
		results[src] = model.SeshSessions{
			OrderedIndex: make([]string, 0),
			Directory:    make(model.SeshSessionMap),
		}
	}
	for _, key := range stored.OrderedIndex {
		session, ok := stored.Directory[key]
		if !ok {
			continue
		}
		result, cachedSource := results[session.Src]
		if !cachedSource {
			continue
		}
		result.OrderedIndex = append(result.OrderedIndex, key)
		result.Directory[key] = session
		results[session.Src] = result
	}
	return results
}

// revalidate fetches the cached sources and writes them to the cache.
func (cl *CachingLister) revalidate() {
	fetched, err := cl.inner.ListSources(cachedSources)
	if err != nil {
		slog.Warn("cache: background revalidation fetch failed", "error", err)
		return
	}
	if err := cl.cache.Write(flatten(fetched)); err != nil {
		slog.Warn("cache: background revalidation write failed", "error", err)
	}
}

// RefreshCache fetches the cached sources and writes them to the cache in the
// background. It never reads the cache. Call Wait before process exit.
func (cl *CachingLister) RefreshCache() {
	cl.wg.Add(1)
	go func() {
		defer cl.wg.Done()
		cl.revalidate()
	}()
}

// Wait blocks until all background refresh goroutines have completed.
// Call this before process exit to avoid truncated cache writes.
func (cl *CachingLister) Wait() {
	cl.wg.Wait()
}

// --- Delegate all other Lister methods to inner ---

func (cl *CachingLister) Format(sessions model.SeshSessions, opts ListOptions) (model.SeshSessions, error) {
	return cl.inner.Format(sessions, opts)
}

func (cl *CachingLister) ListTmuxPanes() (model.SeshSessions, error) {
	return cl.inner.ListTmuxPanes()
}

func (cl *CachingLister) FindTmuxSession(name string) (model.SeshSession, bool) {
	return cl.inner.FindTmuxSession(name)
}

func (cl *CachingLister) FindTmuxSessionByBase(base string) (model.SeshSession, bool) {
	return cl.inner.FindTmuxSessionByBase(base)
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
