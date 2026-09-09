package lister

import (
	"fmt"
	"strings"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
)

// benchSizes are the session counts every benchmark in this file runs at. 1000
// is not absurd: a zoxide history that size is ordinary, and it is the size at
// which an accidental per-session shell-out or list copy shows up.
var benchSizes = []int{10, 100, 1000}

const benchHomeDir = "/home/bench/"

// benchHome is a hand-written Home rather than the generated mock: the mocks
// take a lock and walk their expectations per call, which would be most of what
// a 1000-session benchmark measured.
type benchHome struct{}

func (benchHome) ShortenHome(path string) (string, error) {
	return strings.TrimPrefix(path, benchHomeDir), nil
}

func (benchHome) ExpandPath(path string) (string, error) {
	if after, ok := strings.CutPrefix(path, "~/"); ok {
		return benchHomeDir + after, nil
	}
	return path, nil
}

// benchTmux answers the two calls the list pipeline makes. The embedded mock
// covers the rest of the interface; a benchmark that reaches one of those is
// exercising a path it did not mean to and panics on the nil mock.
type benchTmux struct {
	*tmux.MockTmux
	sessions    []*model.TmuxSession
	windowNames map[string][]string
}

func (t *benchTmux) ListSessions() ([]*model.TmuxSession, error) {
	return t.sessions, nil
}

func (t *benchTmux) ListAllWindowNames(string) (map[string][]string, error) {
	return t.windowNames, nil
}

type benchZoxide struct {
	results []*model.ZoxideResult
}

func (z *benchZoxide) ListResults() ([]*model.ZoxideResult, error) { return z.results, nil }
func (z *benchZoxide) Add(string) error                            { return nil }
func (z *benchZoxide) Query(string) (*model.ZoxideResult, error)   { return nil, nil }
func (z *benchZoxide) Remove(string) error                         { return nil }

type benchTmuxinator struct {
	configs []*model.TmuxinatorConfig
}

func (t *benchTmuxinator) List() ([]*model.TmuxinatorConfig, error) { return t.configs, nil }
func (t *benchTmuxinator) Start(string) (string, error)             { return "", nil }

// benchSources splits n sessions across the sources in roughly the proportions
// a real machine has them: a handful of live tmux sessions, a handful of
// configured ones, and a long zoxide tail.
//
// Every fifth zoxide entry reuses a tmux session's path and every fifth config
// entry reuses a tmux session's name, so dedup has real matches to drop rather
// than walking a list that can never collide.
func benchSources(n int) (*benchTmux, *benchZoxide, *benchTmuxinator, []model.SessionConfig) {
	tmuxCount := max(n/10, 1)
	tmuxinatorCount := max(n/20, 1)
	configCount := max(n/10, 1)
	zoxideCount := max(n-tmuxCount-tmuxinatorCount-configCount, 0)

	tmuxSessions := make([]*model.TmuxSession, tmuxCount)
	windowNames := make(map[string][]string, tmuxCount)
	for i := range tmuxSessions {
		name := fmt.Sprintf("tmux-session-%d", i)
		tmuxSessions[i] = &model.TmuxSession{
			Name:    name,
			Path:    fmt.Sprintf("%sprojects/tmux-session-%d", benchHomeDir, i),
			Windows: 3,
		}
		windowNames[name] = []string{"editor", "server", "shell"}
	}

	configSessions := make([]model.SessionConfig, configCount)
	for i := range configSessions {
		name := fmt.Sprintf("config-session-%d", i)
		if i%5 == 0 {
			// Same name as a live tmux session: config loses the dedup.
			name = fmt.Sprintf("tmux-session-%d", i%tmuxCount)
		}
		configSessions[i] = model.SessionConfig{
			Name: name,
			Path: fmt.Sprintf("~/projects/config-session-%d", i),
		}
	}

	tmuxinatorConfigs := make([]*model.TmuxinatorConfig, tmuxinatorCount)
	for i := range tmuxinatorConfigs {
		tmuxinatorConfigs[i] = &model.TmuxinatorConfig{Name: fmt.Sprintf("tmuxinator-%d", i)}
	}

	zoxideResults := make([]*model.ZoxideResult, zoxideCount)
	for i := range zoxideResults {
		path := fmt.Sprintf("%scode/repo-%d/packages/app-%d", benchHomeDir, i, i%7)
		if i%5 == 0 {
			// Same path as a live tmux session: zoxide loses the dedup.
			path = fmt.Sprintf("%sprojects/tmux-session-%d", benchHomeDir, i%tmuxCount)
		}
		zoxideResults[i] = &model.ZoxideResult{Path: path, Score: float64(zoxideCount - i)}
	}

	return &benchTmux{sessions: tmuxSessions, windowNames: windowNames},
		&benchZoxide{results: zoxideResults},
		&benchTmuxinator{configs: tmuxinatorConfigs},
		configSessions
}

// benchLister wires a RealLister over the fixtures, with config overrides
// applied before the lister is built (the lister copies its config).
func benchLister(n int, tweak func(*model.Config)) *RealLister {
	mockTmux, mockZoxide, mockTmuxinator, configSessions := benchSources(n)
	config := model.Config{SessionConfigs: configSessions}
	if tweak != nil {
		tweak(&config)
	}
	return NewLister(config, benchHome{}, mockTmux, mockZoxide, mockTmuxinator).(*RealLister)
}

// benchSessions is the merged, unfiltered list the post-merge stages operate
// on: what applyDedup and CachingLister.applyFilters get handed.
func benchSessions(n int) model.SeshSessions {
	l := benchLister(n, nil)
	sessions, err := l.List(ListOptions{})
	if err != nil {
		panic(err)
	}
	return sessions
}

func BenchmarkList(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			l := benchLister(n, nil)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := l.List(ListOptions{}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkListHideDuplicates(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			l := benchLister(n, nil)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := l.List(ListOptions{HideDuplicates: true}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkListBlacklist(b *testing.B) {
	patterns := []string{"^scratch$", "^temp", "-wip$", "node_modules"}
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			l := benchLister(n, func(c *model.Config) { c.Blacklist = patterns })
			b.ReportAllocs()
			for b.Loop() {
				if _, err := l.List(ListOptions{}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkListShowWindows covers attachWindowNames, the extra full pass over
// the list that `tui.show_windows` turns on.
func BenchmarkListShowWindows(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			l := benchLister(n, func(c *model.Config) { c.TUI.ShowWindows = true })
			b.ReportAllocs()
			for b.Loop() {
				if _, err := l.List(ListOptions{}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkApplyDedup isolates the tier walk from the source fan-out: it is
// currently one pass over N per tier, so it is the number to watch when a tier
// or a rule is added.
func BenchmarkApplyDedup(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			sessions := benchSessions(n)
			b.ReportAllocs()
			for b.Loop() {
				applyDedup(sessions)
			}
		})
	}
}

// BenchmarkBlacklistFilter isolates the pattern matching from the rest of the
// pipeline: N sessions times W patterns.
func BenchmarkBlacklistFilter(b *testing.B) {
	patterns := []string{"^scratch$", "^temp", "-wip$", "node_modules"}
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			sessions := benchSessions(n)
			b.ReportAllocs()
			for b.Loop() {
				compiled := compileBlacklist(patterns)
				for _, index := range sessions.OrderedIndex {
					isBlacklisted(compiled, sessions.Directory[index].Name)
				}
			}
		})
	}
}

// benchInnerLister is the Lister a CachingLister decorates. Only
// GetAttachedTmuxSession is reached by applyFilters; the rest exist to satisfy
// the interface and fail loudly if a benchmark strays into them.
type benchInnerLister struct {
	attached model.SeshSession
	ok       bool
}

func (l *benchInnerLister) List(ListOptions) (model.SeshSessions, error) {
	panic("benchInnerLister.List: unexpected live fetch in benchmark")
}

func (l *benchInnerLister) ListTmuxPanes() (model.SeshSessions, error) {
	panic("benchInnerLister.ListTmuxPanes: unexpected call in benchmark")
}

func (l *benchInnerLister) GetAttachedTmuxSession() (model.SeshSession, bool) {
	return l.attached, l.ok
}

func (l *benchInnerLister) FindTmuxSession(string) (model.SeshSession, bool) {
	return model.SeshSession{}, false
}

func (l *benchInnerLister) FindTmuxSessionByBase(string) (model.SeshSession, bool) {
	return model.SeshSession{}, false
}

func (l *benchInnerLister) GetLastTmuxSession() (model.SeshSession, bool) {
	return model.SeshSession{}, false
}

func (l *benchInnerLister) FindConfigSession(string) (model.SeshSession, bool) {
	return model.SeshSession{}, false
}

func (l *benchInnerLister) FindConfigWildcard(string) (model.WildcardConfig, bool) {
	return model.WildcardConfig{}, false
}

func (l *benchInnerLister) FindZoxideSession(string) (model.SeshSession, bool) {
	return model.SeshSession{}, false
}

func (l *benchInnerLister) FindTmuxinatorConfig(string) (model.SeshSession, bool) {
	return model.SeshSession{}, false
}

// BenchmarkCachingListerApplyFilters is the steady-state path: it runs on every
// invocation, cache hit or miss, so it matters more than the cold List.
func BenchmarkCachingListerApplyFilters(b *testing.B) {
	cases := []struct {
		name string
		opts ListOptions
	}{
		{"passthrough", ListOptions{}},
		{"source-filter", ListOptions{Tmux: true, Zoxide: true}},
		{"hide-attached", ListOptions{HideAttached: true}},
		{"hide-duplicates", ListOptions{HideDuplicates: true}},
		{"picker-defaults", ListOptions{HideAttached: true, HideDuplicates: true}},
	}
	for _, n := range benchSizes {
		for _, tc := range cases {
			b.Run(fmt.Sprintf("n=%d/%s", n, tc.name), func(b *testing.B) {
				sessions := benchSessions(n)
				attached := sessions.Directory[sessions.OrderedIndex[0]]
				cl := NewCachingLister(&benchInnerLister{attached: attached, ok: true}, nil)
				b.ReportAllocs()
				for b.Loop() {
					cl.applyFilters(sessions, tc.opts)
				}
			})
		}
	}
}

// BenchmarkFindTmuxSessionByBase re-lists every tmux session per call, then
// scans it for a prefix; connect calls it once per reconnect.
func BenchmarkFindTmuxSessionByBase(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			l := benchLister(n, nil)
			b.ReportAllocs()
			for b.Loop() {
				l.FindTmuxSessionByBase("no-such-session")
			}
		})
	}
}
