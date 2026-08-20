package picker

import (
	"fmt"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
)

// benchSizes are the session counts every benchmark in this file runs at.
// filterSessions runs on every keystroke, so the n=1000 numbers are the ones
// that decide whether the picker feels instant for someone with a long zoxide
// history.
var benchSizes = []int{10, 100, 1000}

// benchQueries are the query shapes a filter pass has to handle. The worst case
// is not the long query but the single character: it matches nearly everything,
// so ranking and the result slice stay full size.
var benchQueries = []struct {
	name  string
	query string
}{
	{"empty", ""},
	{"1-char", "a"},
	{"3-char", "app"},
	{"no-match", "qqzzxx"},
}

// benchSessions builds a list shaped like a real one: a few live tmux sessions,
// a few configured ones, and a long zoxide tail of paths. Names mix the
// separators (`-`, `_`, `/`) so the separator-aware normalization has work to do.
func benchSessions(n int) model.SeshSessions {
	orderedIndex := make([]string, 0, n)
	directory := make(model.SeshSessionMap, n)

	add := func(key string, session model.SeshSession) {
		orderedIndex = append(orderedIndex, key)
		directory[key] = session
	}

	tmuxCount := max(n/10, 1)
	configCount := max(n/10, 1)
	for i := 0; i < tmuxCount; i++ {
		add(fmt.Sprintf("tmux:%d", i), model.SeshSession{
			Src:     "tmux",
			Name:    fmt.Sprintf("app-server-%d", i),
			Path:    fmt.Sprintf("/home/bench/projects/app-server-%d", i),
			Windows: 3,
		})
	}
	for i := 0; i < configCount; i++ {
		add(fmt.Sprintf("config:%d", i), model.SeshSession{
			Src:  "config",
			Name: fmt.Sprintf("config_session_%d", i),
			Path: fmt.Sprintf("/home/bench/config/session-%d", i),
		})
	}
	for i := len(orderedIndex); i < n; i++ {
		add(fmt.Sprintf("zoxide:%d", i), model.SeshSession{
			Src:   "zoxide",
			Name:  fmt.Sprintf("~/code/repo-%d/packages/app-%d", i, i%7),
			Path:  fmt.Sprintf("/home/bench/code/repo-%d/packages/app-%d", i, i%7),
			Score: float64(n - i),
		})
	}

	return model.SeshSessions{OrderedIndex: orderedIndex, Directory: directory}
}

// benchModel loads n sessions into a model the way the async fetch does, so
// allItems is built by the same code path the picker uses.
func benchModel(n int, tweak func(*Options)) Model {
	sessions := benchSessions(n)
	opts := Options{
		Prompt:            "> ",
		Placeholder:       "Filter sessions...",
		AliasFilterPrefix: model.DefaultAliasFilterPrefix,
	}
	if tweak != nil {
		tweak(&opts)
	}
	m := New(func() (model.SeshSessions, error) { return sessions, nil }, opts)
	loaded, _ := m.Update(sessionsLoadedMsg{sessions: sessions})
	return loaded.(Model)
}

// benchKeystrokePrefix is the query already in the input when the benchmarked
// keystroke lands, so the measured pass is "typing the third character of app"
// rather than "typing into an empty filter".
const benchKeystrokePrefix = "ap"

// BenchmarkKeystroke is the round trip one keypress actually costs: the filter
// pass plus the render that follows it. filterSessions and View are benchmarked
// on their own below, but neither is what the user waits for — this is, so it is
// the number to watch when a feature lands on the picker.
//
// A CPU profile of this benchmark says the round trip is render-bound, not
// filter-bound: ~86% of it is Model.View, and ~84% is highlightMatches calling
// lipgloss.Style.Render once per matched rune. The filter pass this benchmark
// wraps is a rounding error next to that.
//
// Two consequences for reading the numbers. They do not rise monotonically with
// n, because the cost tracks the names that land in the visible window and how
// densely they match, not the length of the list. And each n is only comparable
// against itself — the same fixture and query — which is all a regression
// comparison needs.
//
// The input is restored to its starting value after every pass. Otherwise the
// query would accumulate a character per iteration and each pass would measure a
// different (and increasingly unrealistic) query. SetValue does not filter on its
// own — only Update does — so the restore costs one short string-to-rune
// conversion and no second pass.
func BenchmarkKeystroke(b *testing.B) {
	key := tea.KeyPressMsg{Code: 'p', Text: "p"}
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			m := benchModel(n, nil)
			sized, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			m = sized.(Model)
			m.filterInput.SetValue(benchKeystrokePrefix)

			// Guard against silently measuring a no-op if the key handling
			// changes shape: the keypress has to reach the filter input.
			probe, _ := m.Update(key)
			if got := probe.(Model).filterInput.Value(); got != benchKeystrokePrefix+"p" {
				b.Fatalf("keypress did not reach the filter input: query is %q", got)
			}

			b.ReportAllocs()
			for b.Loop() {
				updated, _ := m.Update(key)
				m = updated.(Model)
				m.View()
				m.filterInput.SetValue(benchKeystrokePrefix)
			}
		})
	}
}

// BenchmarkFilterSessions is the most latency-visible function in the app: it
// re-runs over every loaded session on every keystroke.
//
// The separator-aware pass is a sibling rather than a measure of the
// normalization overhead: normalizing a query this short is one Replacer call
// and invisible next to the match. What it does show is that turning
// `separator_aware` on changes the shape of the matching work, because
// replacing `-`, `_` and `/` with spaces moves the word boundaries the fuzzy
// scorer keys off.
func BenchmarkFilterSessions(b *testing.B) {
	modes := []struct {
		name           string
		separatorAware bool
	}{
		{name: "plain"},
		{name: "separator-aware", separatorAware: true},
	}
	for _, n := range benchSizes {
		for _, mode := range modes {
			m := benchModel(n, func(o *Options) { o.SeparatorAware = mode.separatorAware })
			for _, q := range benchQueries {
				b.Run(fmt.Sprintf("n=%d/%s/%s", n, mode.name, q.name), func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						m.filterSessions(q.query)
					}
				})
			}
		}
	}
}

// BenchmarkBuildItems runs once per fetch, not per keystroke, but it touches
// every session and resolves an icon for each one.
func BenchmarkBuildItems(b *testing.B) {
	cases := []struct {
		name           string
		separatorAware bool
		icons          bool
	}{
		{name: "plain"},
		{name: "separator-aware", separatorAware: true},
		{name: "icons", icons: true},
	}
	for _, n := range benchSizes {
		sessions := benchSessions(n)
		for _, tc := range cases {
			b.Run(fmt.Sprintf("n=%d/%s", n, tc.name), func(b *testing.B) {
				var resolveIcon IconFunc
				if tc.icons {
					resolveIcon = buildIconResolver(benchIconConfig(sessions), benchHome(b), nil)
					if resolveIcon == nil {
						b.Fatal("expected an icon resolver for a config that declares icons")
					}
				}
				b.ReportAllocs()
				for b.Loop() {
					buildItems(sessions, tc.separatorAware, resolveIcon)
				}
			})
		}
	}
}

// benchIconConfig declares an icon for every second session by name, so the
// resolver's map tier does real lookups and the rest fall through.
func benchIconConfig(sessions model.SeshSessions) model.Config {
	configs := make([]model.SessionConfig, 0, len(sessions.OrderedIndex)/2)
	for i, key := range sessions.OrderedIndex {
		if i%2 != 0 {
			continue
		}
		configs = append(configs, model.SessionConfig{
			Name: sessions.Directory[key].Name,
			Icon: "󰊤",
		})
	}
	return model.Config{SessionConfigs: configs}
}

func benchHome(b *testing.B) home.Home {
	b.Helper()
	return home.NewHome(oswrap.NewOs())
}

// benchWildcardCount is how many [[wildcard]] blocks the wildcard icon
// benchmark declares. A config with this many patterns is unremarkable, and
// FindConfigWildcard expands every one of them per session it is asked about.
const benchWildcardCount = 8

// BenchmarkIconResolverWildcard sizes the O(N x W) path expansion behind the
// resolver's wildcard fallback: one home.ExpandPath for the session plus one
// per wildcard pattern, with no memoization, for every session on screen.
//
// "miss" is the worst case (no pattern matches, so all W are expanded);
// "hit-last" is what a config whose catch-all sits at the bottom pays.
func BenchmarkIconResolverWildcard(b *testing.B) {
	cases := []struct {
		name string
		// lastPattern replaces the final wildcard so the session matches it.
		lastPattern string
	}{
		{name: "miss"},
		{name: "hit-last", lastPattern: "/home/bench/**"},
	}
	for _, n := range benchSizes {
		sessions := benchSessions(n)
		for _, tc := range cases {
			b.Run(fmt.Sprintf("n=%d/%s", n, tc.name), func(b *testing.B) {
				wildcards := make([]model.WildcardConfig, benchWildcardCount)
				for i := range wildcards {
					wildcards[i] = model.WildcardConfig{
						Pattern: fmt.Sprintf("~/unmatched-%d/**", i),
						Icon:    "󰊤",
					}
				}
				if tc.lastPattern != "" {
					wildcards[len(wildcards)-1].Pattern = tc.lastPattern
				}
				config := model.Config{WildcardConfigs: wildcards}
				h := benchHome(b)
				finder := lister.NewLister(config, h, nil, nil, nil)
				resolveIcon := buildIconResolver(config, h, finder)
				if resolveIcon == nil {
					b.Fatal("expected an icon resolver for a config that declares wildcard icons")
				}
				b.ReportAllocs()
				for b.Loop() {
					for _, key := range sessions.OrderedIndex {
						resolveIcon(sessions.Directory[key])
					}
				}
			})
		}
	}
}

// BenchmarkFilterAliases re-runs per keystroke in alias mode, and rebuilds the
// candidate list (a full pass over the loaded sessions plus a sort) each time.
func BenchmarkFilterAliases(b *testing.B) {
	queries := []string{"", "a", "al7"}
	for _, n := range benchSizes {
		sessions := benchSessions(n)
		aliases := make(map[string]Alias, len(sessions.OrderedIndex)/10+1)
		for i, key := range sessions.OrderedIndex {
			if i%10 != 0 {
				continue
			}
			alias := fmt.Sprintf("al%d", i)
			aliases[alias] = Alias{Alias: alias, Target: sessions.Directory[key].Name}
		}
		// An alias whose target isn't loaded: those trail the listed ones and
		// are sorted, which is the part that scales with the alias count.
		aliases["gone"] = Alias{Alias: "gone", Target: "not-in-the-list"}

		m := benchModel(n, func(o *Options) { o.Aliases = aliases })
		for _, query := range queries {
			name := query
			if name == "" {
				name = "empty"
			}
			b.Run(fmt.Sprintf("n=%d/%s", n, name), func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					m.filterAliases(query)
				}
			})
		}
	}
}

// BenchmarkView confirms rendering costs what the visible window costs rather
// than what the list costs: the n=10 and n=1000 numbers should be close.
//
// It renders with an empty query, so it measures rows without match
// highlighting — the cheap path. Highlighting dominates a real keystroke; see
// BenchmarkKeystroke.
func BenchmarkView(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			m := benchModel(n, nil)
			sized, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			m = sized.(Model)
			b.ReportAllocs()
			for b.Loop() {
				m.View()
			}
		})
	}
}
