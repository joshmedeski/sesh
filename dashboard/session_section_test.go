package dashboard

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/model"
)

// sessionNames returns the names of ss in order.
func sessionNames(ss []model.SeshSession) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = s.Name
	}
	return out
}

// apSessionsSection returns a section whose "ap" query filters 5 sessions down
// to [api, app, ape] (cursor on the first match).
func apSessionsSection() *SessionsSection {
	s := &SessionsSection{
		sessions: []model.SeshSession{
			{Name: "api"}, {Name: "app"}, {Name: "ape"}, {Name: "zoo"}, {Name: "yak"},
		},
		ListState: ListState{
			filtering:   true,
			filterQuery: "ap",
		},
	}
	s.applyFilter()
	return s
}

// --- Sessions section (flat list) ---

func TestFlattenSessionsSortedAlphabetically(t *testing.T) {
	sessions := model.SeshSessions{
		OrderedIndex: []string{"z", "a", "m"},
		Directory: model.SeshSessionMap{
			"z": {Name: "z"},
			"a": {Name: "a"},
			"m": {Name: "m"},
		},
	}
	flat := flattenSessions(sessions)
	require.Len(t, flat, 3)
	assert.Equal(t, []string{"a", "m", "z"}, []string{flat[0].Name, flat[1].Name, flat[2].Name})
}

func TestSessionsSectionTKeyIsNoop(t *testing.T) {
	s := &SessionsSection{sessions: []model.SeshSession{{Name: "a"}, {Name: "b"}}}
	updated, cmd := s.handleKey(pressKey("t"))
	assert.Nil(t, cmd)
	// The list is flat; `t` no longer collapses/expands anything, so the
	// cursor and list are untouched.
	assert.Equal(t, 0, updated.cursor)
	assert.Len(t, updated.sessions, 2)
}

// --- Sort modes ---

func TestSessionsSectionSortModeCycle(t *testing.T) {
	now := time.Now()
	t1 := now.Add(-1 * time.Hour)
	t2 := now.Add(-2 * time.Hour)
	t3 := now.Add(-3 * time.Hour)
	c1 := now.Add(-10 * time.Hour)
	c2 := now.Add(-20 * time.Hour)
	c3 := now.Add(-30 * time.Hour)

	s := &SessionsSection{
		sessions: []model.SeshSession{
			{Name: "b", LastAttached: &t1, Created: &c2},
			{Name: "a", LastAttached: &t2, Created: &c1},
			{Name: "c", LastAttached: &t3, Created: &c3},
		},
		sortMode: "name",
	}
	s.applySort()
	assert.Equal(t, []string{"a", "b", "c"}, sessionNames(s.sessions))

	s.cycleSortMode() // name → recent
	assert.Equal(t, "recent", s.SortLabel())
	assert.Equal(t, []string{"b", "a", "c"}, sessionNames(s.sessions))

	s.cycleSortMode() // recent → created
	assert.Equal(t, "created", s.SortLabel())
	assert.Equal(t, []string{"a", "b", "c"}, sessionNames(s.sessions))

	s.cycleSortMode() // created → name
	assert.Equal(t, "name", s.SortLabel())
	assert.Equal(t, []string{"a", "b", "c"}, sessionNames(s.sessions))
}

func TestSessionsSectionSKeyCyclesSort(t *testing.T) {
	s := &SessionsSection{
		sessions: []model.SeshSession{{Name: "b"}, {Name: "a"}},
		sortMode: "name",
	}
	updated, _ := s.handleKey(pressKey("s"))
	assert.Equal(t, "recent", updated.sortMode)
	updated, _ = updated.handleKey(pressKey("s"))
	assert.Equal(t, "created", updated.sortMode)
	updated, _ = updated.handleKey(pressKey("s"))
	assert.Equal(t, "name", updated.sortMode)
}

// --- Type-to-filter ---

func TestFilterMatchesCaseInsensitive(t *testing.T) {
	s := &SessionsSection{
		sessions: []model.SeshSession{{Name: "Alpha"}, {Name: "BETA"}, {Name: "alpine"}},
		ListState: ListState{
			filtering: true,
		},
	}
	s.filterQuery = "ALP"
	s.applyFilter()
	require.Len(t, s.filtered, 2)
	assert.Equal(t, "Alpha", s.filtered[0].Name)
	assert.Equal(t, "alpine", s.filtered[1].Name)
}

// --- Filter navigation / enter / esc (regression) ---

func TestSessionsFilterNavigationMovesThroughResults(t *testing.T) {
	s := apSessionsSection()
	require.Len(t, s.filtered, 3)

	// j/k and the arrow keys move the cursor through the filtered results.
	step := func(key string) *SessionsSection {
		updated, _ := s.handleKey(pressKey(key))
		return updated
	}
	s = step("j")
	assert.Equal(t, 1, s.cursor)
	s = step("down")
	assert.Equal(t, 2, s.cursor)
	s = step("j") // clamped at the last filtered item
	assert.Equal(t, 2, s.cursor)
	s = step("k")
	assert.Equal(t, 1, s.cursor)
	s = step("up")
	assert.Equal(t, 0, s.cursor)
	s = step("k") // clamped at the first filtered item
	assert.Equal(t, 0, s.cursor)

	// Navigation never mutates the query or exits filtering.
	assert.True(t, s.filtering)
	assert.Equal(t, "ap", s.filterQuery)
	assert.Equal(t, []string{"api", "app", "ape"}, sessionNames(s.filtered))
}

func TestSessionsFilterEnterSelectsHighlightedAndExits(t *testing.T) {
	s := apSessionsSection()
	s.cursor = 2 // highlight "ape"

	updated, _ := s.handleKey(pressKey("enter"))
	assert.Equal(t, "ape", updated.chosen)
	assert.False(t, updated.filtering)
	assert.Equal(t, "", updated.filterQuery)
}

func TestSessionsFilterEscCancelsWithoutSelecting(t *testing.T) {
	s := apSessionsSection()
	s.cursor = 2 // a filtered item is highlighted, but esc must not select it

	updated, _ := s.handleKey(pressKey("esc"))
	assert.False(t, updated.filtering)
	assert.Equal(t, "", updated.filterQuery)
	assert.Equal(t, "", updated.chosen)
}

// --- Current-session highlight ---

func TestSessionsSectionRenderItemCurrentHighlight(t *testing.T) {
	s := &SessionsSection{
		sessions:    []model.SeshSession{{Name: "active", Path: "/home/u/active"}},
		currentName: "active",
		ListState: ListState{
			cursor: 1, // row 0 is not selected, so only the name accent shows
		},
	}
	row := s.renderItem(0, 100)
	assert.Contains(t, row, "\x1b[1;38;5;14m") // bold cyan accent on the name
}

func TestSessionsSectionRenderItemNonCurrentNotAccent(t *testing.T) {
	s := &SessionsSection{
		sessions:    []model.SeshSession{{Name: "active", Path: "/home/u/active"}},
		currentName: "other",
		ListState: ListState{
			cursor: 1, // row 0 not selected → no marker accent either
		},
	}
	row := s.renderItem(0, 100)
	assert.NotContains(t, row, "\x1b[38;5;14m")
}

// --- Alias support (Open sessions) ---

func TestFilterMatchesAlias(t *testing.T) {
	s := &SessionsSection{
		sessions: []model.SeshSession{
			{Name: "wallpaper", Alias: "wp"},
			{Name: "dotfiles", Alias: "dot"},
			{Name: "notes"},
		},
		ListState: ListState{
			filtering:   true,
			filterQuery: "wp",
		},
	}
	s.applyFilter()
	require.Len(t, s.filtered, 1)
	assert.Equal(t, "wallpaper", s.filtered[0].Name)

	// Case-insensitive alias matching.
	s.filterQuery = "DOT"
	s.applyFilter()
	require.Len(t, s.filtered, 1)
	assert.Equal(t, "dotfiles", s.filtered[0].Name)

	// A query that matches neither name nor alias returns nothing.
	s.filterQuery = "xyz"
	s.applyFilter()
	assert.Empty(t, s.filtered)
}

func TestSelectByAliasReturnsSessionName(t *testing.T) {
	s := &SessionsSection{
		sessions: []model.SeshSession{
			{Name: "wallpaper", Alias: "wp"},
			{Name: "dotfiles"},
		},
		ListState: ListState{
			filtering:   true,
			filterQuery: "wp",
		},
	}
	s.applyFilter()
	require.Len(t, s.filtered, 1)
	s.selectItem()
	assert.Equal(t, "wallpaper", s.chosen)
}
