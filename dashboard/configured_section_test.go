package dashboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/model"
)

// apConfiguredSection returns a section whose "ap" query filters 5 sessions
// down to [api, app, ape] (cursor on the first match).
func apConfiguredSection() *ConfiguredSection {
	s := &ConfiguredSection{
		sessions: []model.SeshSession{
			{Name: "api"}, {Name: "app"}, {Name: "ape"}, {Name: "zoo"}, {Name: "yak"},
		},
		running: map[string]bool{},
		ListState: ListState{
			filtering:   true,
			filterQuery: "ap",
		},
	}
	s.applyFilter()
	return s
}

// --- Configured section ---

func TestConfiguredSectionBasic(t *testing.T) {
	cs := NewConfiguredSection(model.DashboardSectionConfig{Type: "configured", Title: "Configured"}, SectionDeps{})
	assert.Equal(t, "Configured", cs.Name())
	assert.Equal(t, "", cs.Chosen())
	assert.Equal(t, 0, cs.TotalItems())
}

func TestConfiguredSectionViewLoading(t *testing.T) {
	cs := NewConfiguredSection(model.DashboardSectionConfig{Type: "configured", Title: "Configured"}, SectionDeps{})
	title, content := cs.ViewBorderless(80, 10, true)
	assert.Equal(t, "Configured", title)
	assert.Contains(t, content, "Loading")
}

func TestConfiguredSectionSelectItem(t *testing.T) {
	cs := &ConfiguredSection{sessions: []model.SeshSession{{Name: "a"}, {Name: "b"}}}
	cs.cursor = 1
	cs.selectItem()
	assert.Equal(t, "b", cs.chosen)
}

// --- Type-to-filter ---

func TestConfiguredFilterMatchesCaseInsensitive(t *testing.T) {
	s := &ConfiguredSection{
		sessions: []model.SeshSession{{Name: "Alpha"}, {Name: "beta"}},
		ListState: ListState{
			filtering: true,
		},
	}
	s.filterQuery = "ALP"
	s.applyFilter()
	require.Len(t, s.filtered, 1)
	assert.Equal(t, "Alpha", s.filtered[0].Name)
}

// --- Filter navigation / enter / esc (regression) ---

func TestConfiguredFilterNavigationMovesThroughResults(t *testing.T) {
	s := apConfiguredSection()
	require.Len(t, s.filtered, 3)

	step := func(key string) *ConfiguredSection {
		updated, _ := s.handleKey(pressKey(key))
		return updated.(*ConfiguredSection)
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

	assert.True(t, s.filtering)
	assert.Equal(t, "ap", s.filterQuery)
	assert.Equal(t, []string{"api", "app", "ape"}, sessionNames(s.filtered))
}

func TestConfiguredFilterEnterSelectsHighlightedAndExits(t *testing.T) {
	s := apConfiguredSection()
	s.cursor = 1 // highlight "app"

	updated, _ := s.handleKey(pressKey("enter"))
	cfg := updated.(*ConfiguredSection)
	assert.Equal(t, "app", cfg.chosen)
	assert.False(t, cfg.filtering)
	assert.Equal(t, "", cfg.filterQuery)
}

func TestConfiguredFilterEscCancelsWithoutSelecting(t *testing.T) {
	s := apConfiguredSection()
	s.cursor = 1 // a filtered item is highlighted, but esc must not select it

	updated, _ := s.handleKey(pressKey("esc"))
	cfg := updated.(*ConfiguredSection)
	assert.False(t, cfg.filtering)
	assert.Equal(t, "", cfg.filterQuery)
	assert.Equal(t, "", cfg.chosen)
}
