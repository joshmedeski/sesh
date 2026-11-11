package lister_test

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/formatter"
	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
)

const fakeHomeDir = "/home/user"

// fakeHome shortens and expands paths against a fixed home directory.
type fakeHome struct{}

func (fakeHome) ShortenHome(path string) (string, error) {
	return strings.Replace(path, fakeHomeDir, "~", 1), nil
}

func (fakeHome) ExpandPath(path string) (string, error) {
	return strings.Replace(path, "~", fakeHomeDir, 1), nil
}

// fakeTmux answers the calls the list pipeline makes with state a test can
// change between lists. The embedded mock covers the rest of the interface, so
// a call outside the list pipeline panics.
type fakeTmux struct {
	*tmux.MockTmux

	mu          sync.Mutex
	sessions    []*model.TmuxSession
	err         error
	panes       []*model.TmuxPane
	windowNames map[string][]string

	sessionCalls atomic.Int32
}

func (t *fakeTmux) set(sessions ...*model.TmuxSession) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sessions = sessions
	t.err = nil
}

func (t *fakeTmux) fail(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.err = err
}

func (t *fakeTmux) ListSessions() ([]*model.TmuxSession, error) {
	t.sessionCalls.Add(1)
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.sessions, t.err
}

func (t *fakeTmux) ListTmuxPanes() ([]*model.TmuxPane, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.panes, nil
}

func (t *fakeTmux) ListAllWindowNames(string) (map[string][]string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.windowNames, nil
}

func (t *fakeTmux) ListAllWindows() (map[string][]model.TmuxWindow, error) {
	return nil, nil
}

// fakeZoxide counts its calls. A test can hold a call open on gate, and make
// it fail with err.
type fakeZoxide struct {
	mu      sync.Mutex
	results []*model.ZoxideResult
	err     error
	gate    chan struct{}
	started chan struct{}
	open    func()

	calls atomic.Int32
}

func (z *fakeZoxide) set(results ...*model.ZoxideResult) {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.results = results
}

func (z *fakeZoxide) fail(err error) {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.err = err
}

// hold makes every following call announce itself on started and then wait
// for the returned release. release is safe to call more than once.
func (z *fakeZoxide) hold() (started chan struct{}, release func()) {
	gate := make(chan struct{})
	started = make(chan struct{}, 8)
	var once sync.Once
	release = func() { once.Do(func() { close(gate) }) }
	z.mu.Lock()
	z.gate, z.started, z.open = gate, started, release
	z.mu.Unlock()
	return started, release
}

// release lets every held call go, and does nothing when nothing is held.
func (z *fakeZoxide) release() {
	z.mu.Lock()
	open := z.open
	z.mu.Unlock()
	if open != nil {
		open()
	}
}

func (z *fakeZoxide) ListResults() ([]*model.ZoxideResult, error) {
	z.calls.Add(1)
	z.mu.Lock()
	gate, started := z.gate, z.started
	z.mu.Unlock()
	if gate != nil {
		started <- struct{}{}
		<-gate
	}
	z.mu.Lock()
	defer z.mu.Unlock()
	return z.results, z.err
}

func (z *fakeZoxide) Add(string) error                          { return nil }
func (z *fakeZoxide) Query(string) (*model.ZoxideResult, error) { return nil, nil }
func (z *fakeZoxide) Remove(string) error                       { return nil }

type fakeTmuxinator struct {
	mu      sync.Mutex
	configs []*model.TmuxinatorConfig

	calls atomic.Int32
}

func (t *fakeTmuxinator) List() ([]*model.TmuxinatorConfig, error) {
	t.calls.Add(1)
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.configs, nil
}

func (t *fakeTmuxinator) Start(string) (string, error) { return "", nil }

// fixture is one machine: a tmux server, zoxide, tmuxinator and a config, plus
// the cache file a CachingLister over them reads and writes.
type fixture struct {
	tmux       *fakeTmux
	zoxide     *fakeZoxide
	tmuxinator *fakeTmuxinator
	config     model.Config
	cachePath  string
	t          *testing.T
}

// newFixture builds a machine with two tmux sessions (main is attached), two
// configured sessions, three zoxide entries and one tmuxinator project. Some
// entries share a name or a path across sources so dedup has work to do.
func newFixture(t *testing.T, tweak func(*model.Config)) *fixture {
	t.Helper()
	config := model.Config{
		SessionConfigs: []model.SessionConfig{
			{Name: "project", Path: "~/project"},
			{Name: "dev", Path: "~/dev"},
		},
	}
	if tweak != nil {
		tweak(&config)
	}
	f := &fixture{
		tmux: &fakeTmux{
			MockTmux: new(tmux.MockTmux),
			panes: []*model.TmuxPane{
				{WindowName: "editor", PaneID: "%1", PaneCommand: "vim", PanePath: "/home/user/main"},
			},
			windowNames: map[string][]string{"main": {"editor", "shell"}, "dev": {"server"}},
		},
		zoxide:     &fakeZoxide{},
		tmuxinator: &fakeTmuxinator{configs: []*model.TmuxinatorConfig{{Name: "work"}}},
		config:     config,
		cachePath:  filepath.Join(t.TempDir(), "sessions.v2.gob"),
		t:          t,
	}
	f.tmux.set(
		&model.TmuxSession{Name: "main", Path: "/home/user/main", Attached: 1, Windows: 2},
		&model.TmuxSession{Name: "dev", Path: "/home/user/dev", Windows: 1},
	)
	f.zoxide.set(
		&model.ZoxideResult{Path: "/home/user/code", Score: 10},
		&model.ZoxideResult{Path: "/home/user/dev", Score: 5},
		&model.ZoxideResult{Path: "/home/user/project", Score: 2},
	)
	return f
}

// real is the uncached lister over the machine.
func (f *fixture) real() lister.SourceLister {
	return lister.NewLister(f.config, fakeHome{}, f.tmux, f.zoxide, f.tmuxinator,
		formatter.NewFormatter(icon.NewIcon(f.config)), nil)
}

// cached is a CachingLister over the machine. When the test ends, any held
// zoxide call is released and background refreshes finish, so none outlives it
// and a failing gate test cannot deadlock the cleanup.
func (f *fixture) cached() (*lister.CachingLister, *cache.FileCache) {
	fc := cache.NewFileCacheWithPath(f.cachePath)
	cl := lister.NewCachingLister(f.real(), fc)
	f.t.Cleanup(func() {
		f.zoxide.release()
		cl.Wait()
	})
	return cl, fc
}

// listNames lists with opts and returns the session names in order.
func listNames(t *testing.T, l lister.Lister, opts lister.ListOptions) []string {
	t.Helper()
	got, err := l.List(opts)
	require.NoError(t, err)
	return names(got)
}

func names(sessions model.SeshSessions) []string {
	out := make([]string, 0, len(sessions.OrderedIndex))
	for _, key := range sessions.OrderedIndex {
		out = append(out, sessions.Directory[key].Name)
	}
	return out
}

// sources returns the distinct sources of a stored session list.
func sources(sessions model.SeshSessions) []string {
	seen := map[string]bool{}
	var out []string
	for _, key := range sessions.OrderedIndex {
		src := sessions.Directory[key].Src
		if !seen[src] {
			seen[src] = true
			out = append(out, src)
		}
	}
	return out
}

// writeCacheFile writes a cache file as FileCache would, but with a timestamp
// the test picks, so a test can make the cache stale without waiting.
func writeCacheFile(t *testing.T, path string, sessions model.SeshSessions, written time.Time) {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, gob.NewEncoder(&buf).Encode(cache.CachedData{Sessions: sessions, Timestamp: written}))
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, buf.Bytes(), 0o644))
}

func zoxideRow(path string, score float64) (string, model.SeshSession) {
	name := strings.Replace(path, fakeHomeDir, "~", 1)
	return "zoxide:" + name, model.SeshSession{Src: "zoxide", Name: name, Path: path, Score: score}
}

// storedZoxide is a cached snapshot holding one zoxide row per path.
func storedZoxide(paths ...string) model.SeshSessions {
	stored := model.SeshSessions{Directory: model.SeshSessionMap{}}
	for _, path := range paths {
		key, session := zoxideRow(path, 1)
		stored.OrderedIndex = append(stored.OrderedIndex, key)
		stored.Directory[key] = session
	}
	return stored
}

// listWithin lists in a goroutine and fails the test when the list does not
// return in time, so a list that wrongly waits on a held backend fails instead
// of hanging the suite.
func listWithin(t *testing.T, l lister.Lister, opts lister.ListOptions) (model.SeshSessions, error) {
	t.Helper()
	type result struct {
		sessions model.SeshSessions
		err      error
	}
	done := make(chan result, 1)
	go func() {
		sessions, err := l.List(opts)
		done <- result{sessions, err}
	}()
	select {
	case r := <-done:
		return r.sessions, r.err
	case <-time.After(2 * time.Second):
		t.Fatal("List did not return while a backend was held open")
		return model.SeshSessions{}, nil
	}
}
