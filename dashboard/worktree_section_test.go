package dashboard

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/worktree"
)

func loadedWorktrees(repo string) worktreesLoadedMsg {
	return worktreesLoadedMsg{repo: repo, entries: []model.WorktreeEntry{
		{Number: 358, Path: "/r/w/358", Title: "tmux command updates", State: "OPEN"},
		{Number: 411, Path: "/r/w/411", Title: "configurable dashboard", State: "CLOSED"},
	}}
}

func TestWorktreeSectionIgnoresOtherReposEntries(t *testing.T) {
	s := NewWorktreeSection(model.WorktreeConfig{Repo: "joshmedeski/sesh"}, SectionDeps{})
	s.Update(loadedWorktrees("Nutiliti/nutiliti"))
	assert.True(t, s.loading)
	s.Update(loadedWorktrees("joshmedeski/sesh"))
	assert.False(t, s.loading)
	assert.Equal(t, 2, s.TotalItems())
}

func TestWorktreeSectionFetchListsItsRepo(t *testing.T) {
	wt := worktree.NewMockWorktree(t)
	wt.EXPECT().List(model.WorktreeListOpts{Repo: "joshmedeski/sesh"}).Return(loadedWorktrees("joshmedeski/sesh").entries, nil)
	s := NewWorktreeSection(model.WorktreeConfig{Repo: "joshmedeski/sesh"}, SectionDeps{Worktree: wt})
	msg := s.Init()()
	assert.Equal(t, loadedWorktrees("joshmedeski/sesh"), msg)
}

func TestWorktreeFilterMatchesNumberAndTitle(t *testing.T) {
	s := NewWorktreeSection(model.WorktreeConfig{Repo: "joshmedeski/sesh"}, SectionDeps{})
	s.Update(loadedWorktrees("joshmedeski/sesh"))
	s.filtering = true
	for query, want := range map[string]string{"411": "configurable dashboard", "tmux": "tmux command updates"} {
		s.filterQuery = query
		s.applyFilter()
		require.Len(t, s.visible(), 1, query)
		assert.Equal(t, want, s.visible()[0].Name)
	}
}

func TestWorktreeRowShowsNumberTitleAndBranch(t *testing.T) {
	s := NewWorktreeSection(model.WorktreeConfig{Repo: "joshmedeski/sesh"}, SectionDeps{})
	s.Update(loadedWorktrees("joshmedeski/sesh"))
	s.Update(branchLoadedMsg{path: "/r/w/358", branch: "jam/358-tmux-command-updates"})
	_, content := s.ViewBorderless(120, 10, true)
	assert.Contains(t, ansi.Strip(content), "#358    open   tmux command updates")
	assert.Contains(t, ansi.Strip(content), "jam/358-tmux-command-updates")
}

func TestWorktreeIconFallsBackToRepoRootIcon(t *testing.T) {
	mockHome := new(home.MockHome)
	mockHome.On("ExpandPath", "~/c/nu").Return("/home/u/c/nu", nil)
	deps := SectionDeps{Home: mockHome, IconWidth: 2, Icon: func(sess model.SeshSession) string {
		if sess.Path == "/home/u/c/nu" {
			return "🏠"
		}
		return ""
	}}
	s := NewWorktreeSection(model.WorktreeConfig{Repo: "Nutiliti/nutiliti", Path: "~/c/nu"}, deps)
	s.Update(loadedWorktrees("Nutiliti/nutiliti"))
	_, content := s.ViewBorderless(120, 10, true)
	for _, line := range strings.Split(strings.TrimSuffix(ansi.Strip(content), "\n"), "\n")[1:] {
		assert.True(t, strings.HasPrefix(strings.TrimLeft(line, "▌ "), "🏠"), line)
	}
}

func TestWorktreeIconWithoutMatchUsesConfigGlyph(t *testing.T) {
	s := NewWorktreeSection(model.WorktreeConfig{Repo: "joshmedeski/sesh"}, SectionDeps{})
	assert.Equal(t, render.IconCol("", "config", 1, false), s.iconCol(model.SeshSession{Src: "worktree", Path: "/r/w/358"}, false))
}

func worktreeNumbers(s *WorktreeSection) []int {
	var out []int
	for _, sess := range s.visible() {
		out = append(out, s.entries[sess.Path].Number)
	}
	return out
}

func TestWorktreeSortModes(t *testing.T) {
	s := NewWorktreeSection(model.WorktreeConfig{Repo: "o/r"}, SectionDeps{})
	now := time.Now()
	s.Update(worktreesLoadedMsg{repo: "o/r", entries: []model.WorktreeEntry{
		{Number: 3, Path: "/w/3", State: "CLOSED", Created: now.Add(-1 * time.Hour)},
		{Number: 1, Path: "/w/1", State: "", Created: now.Add(-3 * time.Hour)},
		{Number: 2, Path: "/w/2", State: "MERGED"},
		{Number: 4, Path: "/w/4", State: "OPEN", Created: now.Add(-2 * time.Hour)},
	}})
	s.Update(statusLoadedMsg{path: "/w/2", changes: 5})
	s.Update(statusLoadedMsg{path: "/w/4", changes: 1})

	assert.Equal(t, "issue", s.SortLabel())
	assert.Equal(t, []int{1, 2, 3, 4}, worktreeNumbers(s))

	s.Update(pressKey("s"))
	assert.Equal(t, "age", s.SortLabel())
	assert.Equal(t, []int{3, 4, 1, 2}, worktreeNumbers(s))

	s.Update(pressKey("s"))
	assert.Equal(t, "state", s.SortLabel())
	assert.Equal(t, []int{4, 2, 3, 1}, worktreeNumbers(s))

	s.Update(pressKey("s"))
	assert.Equal(t, "status", s.SortLabel())
	assert.Equal(t, []int{2, 4, 1, 3}, worktreeNumbers(s))

	s.Update(pressKey("s"))
	assert.Equal(t, "issue", s.SortLabel())
}

func TestFooterShowsWorktreeSortLabel(t *testing.T) {
	m := New(model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{Repo: "o/r"}},
		Dashboard: onePage([]model.DashboardSectionConfig{
			{Type: "custom"},
			{Type: "worktree", Repo: "o/r"},
		}),
	}, SectionDeps{})
	m.width = 120
	assert.NotContains(t, ansi.Strip(m.View().Content), "sort:")
	m = updateModel(m, pressKey("ctrl+l"))
	m = updateModel(m, pressKey("s"))
	assert.Contains(t, ansi.Strip(m.View().Content), "sort:age")
}

func TestWorktreeConfigsAddNoTabs(t *testing.T) {
	m := New(model.Config{WorktreeConfigs: []model.WorktreeConfig{{Repo: "joshmedeski/sesh"}}}, SectionDeps{})
	assert.Equal(t, []string{"Dashboard"}, m.tabTitles())
}

func TestBuildSectionsWorktreeMatchesRepo(t *testing.T) {
	worktrees := []model.WorktreeConfig{
		{Repo: "Nutiliti/nutiliti", Path: "~/c/nu"},
		{Repo: "joshmedeski/sesh", Path: "~/c/sesh"},
	}
	built := BuildPages(onePage([]model.DashboardSectionConfig{
		{Type: "worktree", Repo: "JoshMedeski/Sesh"},
		{Type: "worktree", Repo: "nobody/missing"},
		{Type: "worktree", Repo: "nutiliti/nutiliti", Title: "Nutiliti"},
	}), worktrees, SectionDeps{})

	require.Len(t, built.widgets(), 2)
	sesh := built.widgets()[0].(*WorktreeSection)
	assert.Equal(t, worktrees[1], sesh.config)
	assert.Equal(t, "joshmedeski/sesh", sesh.Name())
	assert.Equal(t, "Nutiliti", built.widgets()[1].Name())
}

func TestDashboardWorktreePaneConnects(t *testing.T) {
	cfg := model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{Repo: "joshmedeski/sesh", Path: "~/c/sesh"}},
		Dashboard: onePage([]model.DashboardSectionConfig{
			{Type: "worktree", Repo: "joshmedeski/sesh"},
		}),
	}
	m := New(cfg, SectionDeps{})
	m = updateModel(m, loadedWorktrees("joshmedeski/sesh"))
	assert.Contains(t, ansi.Strip(m.View().Content), "tmux command updates")

	m = updateModel(m, pressKey("enter"))
	require.NotNil(t, m.ChosenWorktree())
	assert.Equal(t, model.WorktreeConnectOpts{Number: 358, Repo: "joshmedeski/sesh"}, *m.ChosenWorktree())
}

func TestWorktreeColumnsFromConfig(t *testing.T) {
	worktrees := []model.WorktreeConfig{{Repo: "joshmedeski/sesh"}}
	built := BuildPages(onePage([]model.DashboardSectionConfig{
		{Type: "worktree", Repo: "joshmedeski/sesh"},
		{Type: "worktree", Repo: "joshmedeski/sesh", Columns: []string{"ghi_title", "title", "age"}},
	}), worktrees, SectionDeps{})

	require.Len(t, built.widgets(), 2)
	assert.Nil(t, built.widgets()[0].(*WorktreeSection).columns)
	assert.Equal(t, []string{"ghi_title", "age"}, built.widgets()[1].(*WorktreeSection).columns)
}
