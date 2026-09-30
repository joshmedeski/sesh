package dashboard

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/dashboard/sections"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/worktree"
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

func TestLoadedSessionsSortedAlphabeticallyByDefault(t *testing.T) {
	s := NewSessionsSection(model.DashboardSectionConfig{}, SectionDeps{}).(*SessionsSection)
	s.Update(sessionsLoadedMsg{section: s, sessions: model.SeshSessions{
		OrderedIndex: []string{"z", "a", "m"},
		Directory: model.SeshSessionMap{
			"z": {Name: "z"},
			"a": {Name: "a"},
			"m": {Name: "m"},
		},
	}})
	assert.Equal(t, []string{"a", "m", "z"}, sessionNames(s.visible()))
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

func TestSessionsSectionRenderItemLongBranchFullAndAligned(t *testing.T) {
	long := "claude/telemetry-worker/parse-bill-empty-attachments"
	s := &SessionsSection{
		sessions: []model.SeshSession{
			{Name: "nu/w/1", Path: "/home/u/nu/w/1", Branch: long, GitStatus: "~2"},
			{Name: "nu/w/2", Path: "/home/u/nu/w/2", Branch: "main", GitStatus: "!1"},
		},
		ListState: ListState{cursor: 2},
	}
	first := ansi.Strip(s.renderItem(0, 160))
	second := ansi.Strip(s.renderItem(1, 160))
	assert.Contains(t, first, long)
	assert.Equal(t, strings.Index(first, "~2"), strings.Index(second, "!1"))
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

func TestSessionsSectionRenderItemShowsResolvedIcon(t *testing.T) {
	s := &SessionsSection{
		sessions: []model.SeshSession{{Src: "tmux", Name: "nutiliti", Path: "/home/u/c/nu"}},
		deps: SectionDeps{IconWidth: 2, Icon: func(sess model.SeshSession) string {
			if sess.Name == "nutiliti" {
				return "🏠"
			}
			return ""
		}},
		ListState: ListState{cursor: 1},
	}
	assert.True(t, strings.HasPrefix(ansi.Strip(s.renderItem(0, 100)), "  🏠"))
}

func TestSessionsClickAtSkipsHeaderRow(t *testing.T) {
	s := &SessionsSection{sessions: []model.SeshSession{{Name: "a"}, {Name: "b"}, {Name: "c"}}}
	s.cursor = 1
	s.ClickAt(0)
	assert.Equal(t, 1, s.cursor)
	s.ClickAt(3)
	assert.Equal(t, 2, s.cursor)
	s.ClickAt(1)
	assert.Equal(t, 0, s.cursor)
}

func sourcesSection(sources model.SortOrder) *SessionsSection {
	built := BuildPages(onePage([]model.DashboardSectionConfig{
		{Type: "sources", Title: "Sources", Sources: sources},
	}), nil, SectionDeps{})
	return built.widgets()[0].(*SessionsSection)
}

func TestSourcesSectionListOptions(t *testing.T) {
	built := BuildPages(model.DashboardConfig{}, nil, SectionDeps{})
	assert.Equal(t, lister.ListOptions{Tmux: true}, built.widgets()[0].(*SessionsSection).listOptions())
	assert.Equal(t, "name", built.widgets()[0].(*SessionsSection).SortLabel())

	sources := model.SortOrder{"tmux", []any{"config", "zoxide"}}
	s := sourcesSection(sources)
	assert.Equal(t, "Sources", s.Name())
	assert.Equal(t, lister.ListOptions{
		Tmux: true, Config: true, Zoxide: true,
		HideDuplicates: true,
		SortOrder:      sources,
	}, s.listOptions())
	assert.Equal(t, "order", s.SortLabel())
}

func TestSourcesSectionKeepsListerOrder(t *testing.T) {
	s := sourcesSection(model.SortOrder{"tmux", []any{"config", "zoxide"}})
	s.Update(sessionsLoadedMsg{section: s, sessions: model.SeshSessions{
		OrderedIndex: []string{"t", "z1", "c", "z2"},
		Directory: model.SeshSessionMap{
			"t":  {Src: "tmux", Name: "zeta", Group: 0},
			"z1": {Src: "zoxide", Name: "~/hot", Group: 1},
			"c":  {Src: "config", Name: "alpha", Group: 1},
			"z2": {Src: "zoxide", Name: "~/cold", Group: 1},
		},
	}})
	order := []string{"zeta", "~/hot", "alpha", "~/cold"}
	assert.Equal(t, order, sessionNames(s.visible()))

	s.Update(pressKey("s"))
	assert.Equal(t, "name", s.SortLabel())
	assert.Equal(t, []string{"alpha", "zeta", "~/cold", "~/hot"}, sessionNames(s.visible()))

	for range 3 {
		s.Update(pressKey("s"))
	}
	assert.Equal(t, "order", s.SortLabel())
	assert.Equal(t, order, sessionNames(s.visible()))
}

func TestFooterSortFollowsFocusedPane(t *testing.T) {
	sources := sourcesSection(model.SortOrder{"tmux", []any{"config", "zoxide"}})
	m := testModel(sources, sections.NewSystemSection(model.DashboardSectionConfig{Type: "system"}, SectionDeps{}))
	m.width = 140
	m.pages[0].rows[0][0].(*SessionsSection).sortMode = "name"
	assert.Contains(t, ansi.Strip(m.View().Content), "sort:name")

	m.focus = 1
	assert.Contains(t, ansi.Strip(m.View().Content), "sort:order")

	m.focus = 2
	assert.NotContains(t, ansi.Strip(m.View().Content), "sort:")
}

func TestKillSkipsNonTmuxSessions(t *testing.T) {
	s := &SessionsSection{sessions: []model.SeshSession{{Src: "zoxide", Name: "~/hot"}}}
	assert.Nil(t, s.killSession())
}

func TestSourcesSectionColumns(t *testing.T) {
	s := NewSourcesSection(model.DashboardSectionConfig{Type: "sources", Columns: []string{"title", "git_status", "issue", "git_branch"}}, SectionDeps{}).(*SessionsSection)
	assert.Equal(t, []string{"title", "git_status", "git_branch"}, s.columns)

	s.Update(sessionsLoadedMsg{section: s, sessions: model.SeshSessions{
		OrderedIndex: []string{"a"},
		Directory:    model.SeshSessionMap{"a": {Src: "tmux", Name: "alpha", Path: "/p", Branch: "main"}},
	}})
	_, content := s.ViewBorderless(120, 5, true)
	assert.Equal(t, []string{"NAME", "STATUS", "BRANCH"}, strings.Fields(strings.Split(ansi.Strip(content), "\n")[0]))
}

func TestSourcesSectionMatchesWorktreeIssues(t *testing.T) {
	wt := worktree.NewMockWorktree(t)
	wt.EXPECT().List(model.WorktreeListOpts{Repo: "joshmedeski/sesh"}).Return([]model.WorktreeEntry{
		{Number: 358, Path: "/home/u/c/sesh/w/358/", Title: "tmux command updates", State: "OPEN"},
	}, nil).Once()
	built := BuildPages(onePage([]model.DashboardSectionConfig{
		{Type: "sources", Columns: []string{"title", "ghi_title"}},
		{Type: "sources", Columns: []string{"title"}},
	}), []model.WorktreeConfig{{Repo: "joshmedeski/sesh"}}, SectionDeps{Worktree: wt})
	s := built.widgets()[0].(*SessionsSection)
	plain := built.widgets()[1].(*SessionsSection)
	require.False(t, plain.showsIssues())

	msg := s.fetchIssues()()
	s.Update(msg)
	plain.Update(msg)
	assert.Nil(t, plain.issues, "issues are only stored by the section that asked for them")

	assert.Equal(t, 358, s.issueFor("/home/u/c/sesh/w/358").Number)
	assert.Equal(t, "tmux command updates", s.issueFor("/home/u/c/sesh/w/358/packages/app").Title)
	assert.Zero(t, s.issueFor("/home/u/c/sesh").Number)
	assert.Zero(t, s.issueFor("").Number)
}

func TestSourcesFilterMatchesIssueTitle(t *testing.T) {
	s := sourcesSection(model.SortOrder{"tmux"})
	s.Update(sessionsLoadedMsg{section: s, sessions: model.SeshSessions{
		OrderedIndex: []string{"a", "b"},
		Directory: model.SeshSessionMap{
			"a": {Src: "tmux", Name: "sesh/358", Path: "/w/358"},
			"b": {Src: "tmux", Name: "dotfiles", Path: "/c/dotfiles"},
		},
	}})
	s.Update(pressKey("/"))
	for _, r := range "tmux" {
		s.Update(pressKey(string(r)))
	}
	assert.Empty(t, sessionNames(s.visible()), "no issues loaded yet")

	s.Update(sessionIssuesLoadedMsg{section: s, entries: map[string]model.WorktreeEntry{
		"/w/358": {Number: 358, Path: "/w/358", Title: "Tmux command updates"},
	}})
	assert.Equal(t, []string{"sesh/358"}, sessionNames(s.visible()), "filter re-applies when issues arrive, case-insensitively")
}

func TestSourcesSectionsIgnoreEachOthersResults(t *testing.T) {
	built := BuildPages(onePage([]model.DashboardSectionConfig{
		{Type: "sources", Sources: model.SortOrder{"tmux"}},
		{Type: "sources", Sources: model.SortOrder{[]any{"config", "zoxide"}}},
	}), nil, SectionDeps{})
	m := Model{configured: built.Configured, pages: built.pages(), width: 120, height: 30}.withLayout()
	tmux := m.pages[0].rows[0][0].(*SessionsSection)
	others := m.pages[0].rows[1][0].(*SessionsSection)

	m, _ = m.broadcast(sessionsLoadedMsg{section: tmux, sessions: model.SeshSessions{
		OrderedIndex: []string{"t"},
		Directory:    model.SeshSessionMap{"t": {Src: "tmux", Name: "live"}},
	}})
	m, _ = m.broadcast(sessionsLoadedMsg{section: others, sessions: model.SeshSessions{
		OrderedIndex: []string{"z"},
		Directory:    model.SeshSessionMap{"z": {Src: "zoxide", Name: "~/c/dotfiles"}},
	}})

	assert.Equal(t, []string{"live"}, sessionNames(tmux.visible()))
	assert.Equal(t, []string{"~/c/dotfiles"}, sessionNames(others.visible()))
}
