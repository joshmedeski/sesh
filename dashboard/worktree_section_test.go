package dashboard

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/worktree"
)

func worktreeTestModel(repos ...string) Model {
	m := testModel()
	for _, repo := range repos {
		m.worktrees = append(m.worktrees, NewWorktreeSection(model.WorktreeConfig{Repo: repo}, SectionDeps{}))
	}
	return m.withLayout()
}

func loadedWorktrees(repo string) worktreesLoadedMsg {
	return worktreesLoadedMsg{repo: repo, entries: []model.WorktreeEntry{
		{Number: 358, Path: "/r/w/358", Title: "tmux command updates", State: "OPEN"},
		{Number: 411, Path: "/r/w/411", Title: "configurable dashboard", State: "CLOSED"},
	}}
}

func TestTabCyclesThroughWorktreeTabsInConfigOrder(t *testing.T) {
	m := worktreeTestModel("joshmedeski/sesh", "Nutiliti/nutiliti")
	var pages []string
	for range m.pageCount() {
		m = updateModel(m, pressKey("tab"))
		if m.page > pageConfigured {
			pages = append(pages, m.pageSection().Name())
		}
	}
	assert.Equal(t, []string{"joshmedeski/sesh", "Nutiliti/nutiliti"}, pages)
	assert.Equal(t, pageOpen, m.page)

	m = updateModel(m, pressKey("shift+tab"))
	assert.Equal(t, "Nutiliti/nutiliti", m.pageSection().Name())
}

func TestHeaderListsWorktreeTabsByRepoName(t *testing.T) {
	m := worktreeTestModel("joshmedeski/sesh", "Nutiliti/nutiliti")
	m.width = 120
	header := ansi.Strip(m.View().Content)
	assert.Contains(t, header, "Open │ Configured │ sesh │ nutiliti")
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

func TestEnterOnWorktreeTabChoosesWorktree(t *testing.T) {
	m := worktreeTestModel("joshmedeski/sesh")
	m.page = pageConfigured + 1
	m.worktrees[0].Update(loadedWorktrees("joshmedeski/sesh"))
	m = updateModel(m, pressKey("j"))
	m = updateModel(m, pressKey("enter"))
	require.NotNil(t, m.ChosenWorktree())
	assert.Equal(t, model.WorktreeConnectOpts{Number: 411, Repo: "joshmedeski/sesh"}, *m.ChosenWorktree())
	assert.Equal(t, "", m.Chosen())
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
	assert.Contains(t, ansi.Strip(content), "#358    tmux command updates")
	assert.Contains(t, ansi.Strip(content), "jam/358-tmux-command-updates")
}

func TestWorktreeIconFallsBackToRepoRootIcon(t *testing.T) {
	deps := SectionDeps{HomeDir: "/home/u", IconWidth: 2, Icon: func(sess model.SeshSession) string {
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
