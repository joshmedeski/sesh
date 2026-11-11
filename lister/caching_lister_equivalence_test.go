package lister_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

// A list served from a warm cache must equal what the uncached lister returns
// for the same machine, whatever the options and the config.
func TestCachingLister_WarmHitEqualsTheUncachedList(t *testing.T) {
	blacklist := func(c *model.Config) { c.Blacklist = []string{"^work$", "^~/code$"} }
	showWindows := func(c *model.Config) { c.TUI.ShowWindows = true }

	cases := []struct {
		name  string
		opts  lister.ListOptions
		tweak func(*model.Config)
	}{
		{name: "default"},
		{name: "tmux only", opts: lister.ListOptions{Tmux: true}},
		{name: "config only", opts: lister.ListOptions{Config: true}},
		{name: "zoxide only", opts: lister.ListOptions{Zoxide: true}},
		{name: "tmuxinator only", opts: lister.ListOptions{Tmuxinator: true}},
		{name: "tmux and zoxide", opts: lister.ListOptions{Tmux: true, Zoxide: true}},
		{name: "panes", opts: lister.ListOptions{Panes: true}},
		{name: "panes and zoxide", opts: lister.ListOptions{Panes: true, Zoxide: true}},
		{name: "hide attached", opts: lister.ListOptions{HideAttached: true}},
		{name: "hide duplicates", opts: lister.ListOptions{HideDuplicates: true}},
		{name: "hide attached and duplicates", opts: lister.ListOptions{HideAttached: true, HideDuplicates: true}},
		{name: "config blacklist", tweak: blacklist},
		{name: "blacklisted only", opts: lister.ListOptions{Blacklisted: true}, tweak: blacklist},
		{
			name: "reversed flat sort order",
			tweak: func(c *model.Config) {
				c.SortOrder = model.SortOrder{"zoxide", "tmuxinator", "config", "tmux"}
			},
		},
		{
			name: "merged sort order group",
			tweak: func(c *model.Config) {
				c.SortOrder = model.SortOrder{"tmux", []string{"config", "zoxide"}}
			},
		},
		{
			name: "sort order from the options",
			opts: lister.ListOptions{SortOrder: model.SortOrder{"config", "tmux"}},
		},
		{name: "show windows", tweak: showWindows},
		{name: "icons", opts: lister.ListOptions{Icons: true}},
		{name: "icons without color", opts: lister.ListOptions{Icons: true, NoColor: true}},
		{name: "json", opts: lister.ListOptions{Json: true}},
		{name: "format", opts: lister.ListOptions{Format: "{source}:{name}", FormatSet: true}},
		{name: "format with the active window", opts: lister.ListOptions{Format: "{active_window_name_prefix}{name}", FormatSet: true}},
	}

	// Each case runs against every way the cache can serve a list: nothing
	// cached yet, a fresh cache, a stale one, and a read that is skipped.
	modes := []struct {
		name string
		list func(t *testing.T, f *fixture, opts lister.ListOptions) model.SeshSessions
	}{
		{"cold", func(t *testing.T, f *fixture, opts lister.ListOptions) model.SeshSessions {
			cl, _ := f.cached()
			got, err := cl.List(opts)
			require.NoError(t, err)
			cl.Wait()
			return got
		}},
		{"fresh", func(t *testing.T, f *fixture, opts lister.ListOptions) model.SeshSessions {
			cl, _ := f.cached()
			_, err := cl.List(lister.ListOptions{})
			require.NoError(t, err, "fills the cache")
			got, err := cl.List(opts)
			require.NoError(t, err)
			cl.Wait()
			return got
		}},
		{"stale", func(t *testing.T, f *fixture, opts lister.ListOptions) model.SeshSessions {
			cl, fc := f.cached()
			_, err := cl.List(lister.ListOptions{})
			require.NoError(t, err, "fills the cache")
			filled, err := fc.Read()
			require.NoError(t, err)
			writeCacheFile(t, f.cachePath, filled.Sessions, longAgo)
			got, err := listWithin(t, cl, opts)
			require.NoError(t, err)
			cl.Wait()
			return got
		}},
		{"skip read", func(t *testing.T, f *fixture, opts lister.ListOptions) model.SeshSessions {
			cl, _ := f.cached()
			cl.SkipRead()
			got, err := cl.List(opts)
			require.NoError(t, err)
			cl.Wait()
			return got
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, mode := range modes {
				t.Run(mode.name, func(t *testing.T) {
					f := newFixture(t, tc.tweak)

					want, err := f.real().List(tc.opts)
					require.NoError(t, err)
					got := mode.list(t, f, tc.opts)

					// What a caller sees is the ordered rows. Directory may hold rows
					// the order leaves out, so it is not compared.
					assert.Equal(t, rows(want), rows(got))
				})
			}
		})
	}
}

func rows(sessions model.SeshSessions) []model.SeshSession {
	out := make([]model.SeshSession, 0, len(sessions.OrderedIndex))
	for _, key := range sessions.OrderedIndex {
		out = append(out, sessions.Directory[key])
	}
	return out
}
