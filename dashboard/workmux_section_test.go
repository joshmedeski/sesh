package dashboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/model"
)

func u64p(v uint64) *uint64 { return &v }
func strp(s string) *string { return &s }

func TestParseWorkmuxStatus(t *testing.T) {
	fixture := `{"agents":[
		{"worktree":"~/code/alpha","branch":"main","status":"working","elapsed_secs":120,"title":"fix bug","session":"alpha:1","agent_kind":"window","git":{"has_staged":true,"has_unstaged":false,"has_unmerged_commits":false}},
		{"worktree":"~/code/beta","branch":"dev","status":"waiting","elapsed_secs":3600,"session":null,"git":{"has_staged":false,"has_unstaged":true,"has_unmerged_commits":false}},
		{"worktree":"~/code/gamma","branch":"main","status":"done","elapsed_secs":259200,"title":null,"session":"gamma:1","git":null},
		{"worktree":"~/code/delta","branch":"main","status":"-","git":{"has_staged":false,"has_unstaged":false,"has_unmerged_commits":true}}
	]}`
	agents, err := parseWorkmuxStatus(fixture)
	require.NoError(t, err)
	require.Len(t, agents, 4)

	assert.Equal(t, "working", agents[0].Status)
	require.NotNil(t, agents[0].ElapsedSecs)
	assert.Equal(t, uint64(120), *agents[0].ElapsedSecs)
	require.NotNil(t, agents[0].Session)
	assert.Equal(t, "alpha:1", *agents[0].Session)
	require.NotNil(t, agents[0].Title)
	assert.Equal(t, "fix bug", *agents[0].Title)
	require.NotNil(t, agents[0].Git)
	assert.True(t, agents[0].Git.HasStaged)

	// Window-mode agent: nil session.
	assert.Nil(t, agents[1].Session)
	require.NotNil(t, agents[1].Git)
	assert.True(t, agents[1].Git.HasUnstaged)

	assert.Equal(t, "done", agents[2].Status)
	assert.Nil(t, agents[2].Git)
	assert.Nil(t, agents[2].Title)

	assert.Equal(t, "-", agents[3].Status)
	require.NotNil(t, agents[3].Git)
	assert.True(t, agents[3].Git.HasUnmergedCommits)
}

func TestParseWorkmuxStatusLenientUnknownFields(t *testing.T) {
	fixture := `{"agents":[{"worktree":"~/x","status":"done","extra":"ignored","nested":{"a":1}}],"extra":"top"}`
	agents, err := parseWorkmuxStatus(fixture)
	require.NoError(t, err)
	require.Len(t, agents, 1)
	assert.Equal(t, "~/x", agents[0].Worktree)
}

func TestParseWorkmuxStatusInvalidJSON(t *testing.T) {
	_, err := parseWorkmuxStatus(`not json`)
	assert.Error(t, err)
}

func TestWorkmuxStateGlyph(t *testing.T) {
	assert.Contains(t, wmStateGlyph("working"), "●")
	assert.Contains(t, wmStateGlyph("waiting"), "○")
	assert.Contains(t, wmStateGlyph("done"), "✓")
	assert.Contains(t, wmStateGlyph("-"), "-")
	assert.Contains(t, wmStateGlyph(""), "-")
}

func TestWorkmuxGitCell(t *testing.T) {
	cell := wmGitCell(&wmGit{HasStaged: true, HasUnstaged: true, HasUnmergedCommits: true})
	assert.Contains(t, cell, "+")
	assert.Contains(t, cell, "~")
	assert.Contains(t, cell, "u")

	assert.Equal(t, "", wmGitCell(nil))
	assert.Equal(t, "", wmGitCell(&wmGit{}))
}

func TestWorkmuxElapsed(t *testing.T) {
	assert.Equal(t, "", wmElapsed(nil))
	assert.Equal(t, "42s", wmElapsed(u64p(42)))
	assert.Equal(t, "2m", wmElapsed(u64p(120)))
	assert.Equal(t, "1h", wmElapsed(u64p(3600)))
	assert.Equal(t, "3d", wmElapsed(u64p(259200)))
}

func TestWorkmuxRowFullColumns(t *testing.T) {
	a := wmAgent{
		Worktree:    "~/code/alpha",
		Branch:      "main",
		Status:      "working",
		ElapsedSecs: u64p(120),
		Title:       strp("fix bug"),
		Git:         &wmGit{HasStaged: true},
	}
	row := renderWorkmuxRow(100, false, a)
	assert.Contains(t, row, "alpha")
	assert.Contains(t, row, "(main)")
	assert.Contains(t, row, "fix bug")
	assert.Contains(t, row, "2m")
	assert.Contains(t, row, "●")
	assert.Contains(t, row, "+")
}

func TestWorkmuxRowProgressiveDrops(t *testing.T) {
	a := wmAgent{
		Worktree:    "~/code/alpha",
		Branch:      "main",
		Status:      "done",
		ElapsedSecs: u64p(120),
		Title:       strp("a title"),
		Git:         &wmGit{HasStaged: true},
	}

	narrow := renderWorkmuxRow(80, false, a) // <90 drops title
	assert.Contains(t, narrow, "(main)")
	assert.Contains(t, narrow, "2m")
	assert.NotContains(t, narrow, "a title")

	veryNarrow := renderWorkmuxRow(60, false, a) // <70 drops branch+git+elapsed
	assert.NotContains(t, veryNarrow, "(main)")
	assert.NotContains(t, veryNarrow, "2m")
	assert.NotContains(t, veryNarrow, "+")
	assert.Contains(t, veryNarrow, "alpha")
}

func TestWorkmuxRowSelectionBackground(t *testing.T) {
	a := wmAgent{Worktree: "~/code/alpha", Status: "working"}
	row := renderWorkmuxRow(100, true, a)
	// Selection applies the highlight background (colour 8 = 48;5;8).
	assert.Contains(t, row, "\x1b[48;5;8m")
}

func TestWorkmuxSectionView(t *testing.T) {
	s := &WorkmuxSection{
		config: model.DashboardSectionConfig{Title: "Workmux"},
		agents: []wmAgent{
			{Worktree: "~/code/alpha", Branch: "main", Status: "working", ElapsedSecs: u64p(120), Title: strp("fix bug"), Git: &wmGit{HasStaged: true}},
			{Worktree: "~/code/beta", Status: "done", Git: &wmGit{HasUnmergedCommits: true}},
		},
	}
	title, content := s.ViewBorderless(100, 10, true)
	assert.Equal(t, "Workmux", title)
	assert.Contains(t, content, "alpha")
	assert.Contains(t, content, "fix bug")
	assert.Contains(t, content, "●") // working
	assert.Contains(t, content, "✓") // done
	assert.Contains(t, content, "2m")
	assert.Contains(t, content, "+") // staged
	assert.Contains(t, content, "u") // unmerged
}

func TestWorkmuxSectionEnterSetsChosen(t *testing.T) {
	session := "alpha:1"
	s := &WorkmuxSection{
		agents: []wmAgent{
			{Worktree: "a", Session: &session},
			{Worktree: "b", Session: nil},
		},
	}
	updated, _ := s.Update(pressKey("enter"))
	assert.Equal(t, "alpha:1", updated.(*WorkmuxSection).Chosen())

	// Move to the window-mode agent (nil session) → enter no-ops.
	updated, _ = updated.Update(pressKey("j"))
	assert.Equal(t, 1, updated.(*WorkmuxSection).cursor)
	updated, _ = updated.Update(pressKey("enter"))
	assert.Equal(t, "alpha:1", updated.(*WorkmuxSection).Chosen())
}

func TestWorkmuxSectionNav(t *testing.T) {
	s := &WorkmuxSection{agents: []wmAgent{{}, {}, {}}}
	updated, _ := s.Update(pressKey("j"))
	assert.Equal(t, 1, updated.(*WorkmuxSection).cursor)
	updated, _ = updated.Update(pressKey("k"))
	assert.Equal(t, 0, updated.(*WorkmuxSection).cursor)
}

func TestWorkmuxSectionRefresh(t *testing.T) {
	s := &WorkmuxSection{loading: false, agents: []wmAgent{{}}}
	updated, cmd := s.Update(pressKey("r"))
	assert.NotNil(t, cmd)
	assert.True(t, updated.(*WorkmuxSection).loading)
	assert.Equal(t, "", updated.(*WorkmuxSection).errorMsg)
}

func TestWorkmuxSectionStates(t *testing.T) {
	cfg := model.DashboardSectionConfig{Title: "Workmux"}

	loading := &WorkmuxSection{config: cfg, loading: true}
	_, content := loading.ViewBorderless(80, 10, true)
	assert.Contains(t, content, "Loading")

	empty := &WorkmuxSection{config: cfg, loading: false}
	_, content = empty.ViewBorderless(80, 10, true)
	assert.Contains(t, content, "No agents found")

	err := &WorkmuxSection{config: cfg, loading: false, errorMsg: "Failed to parse workmux status"}
	_, content = err.ViewBorderless(80, 10, true)
	assert.Contains(t, content, "Failed to parse workmux status")
}
