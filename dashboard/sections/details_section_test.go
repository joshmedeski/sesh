package sections

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshmedeski/sesh/v2/dashboard/core"
	"github.com/joshmedeski/sesh/v2/model"
)

// --- Live preview (DetailsSection) ---

func TestDetailsSectionHoverKicksCapture(t *testing.T) {
	ds := NewDetailsSection(model.DashboardSectionConfig{Title: "Details"}, core.SectionDeps{}).(*DetailsSection)
	_, cmd := ds.Update(core.HoveredSessionMsg{Name: "sesh", Path: "/x", Windows: 1})
	assert.NotNil(t, cmd)
}

func TestDetailsSectionPreviewLoadedUpdatesState(t *testing.T) {
	ds := &DetailsSection{hoveredName: "sesh"}
	updated, _ := ds.Update(previewLoadedMsg{name: "sesh", output: "hello\nworld"})
	assert.Equal(t, "hello\nworld", updated.(*DetailsSection).previewOutput)

	// Stale capture for a previous hover is ignored.
	updated2, _ := updated.(*DetailsSection).Update(previewLoadedMsg{name: "other", output: "stale"})
	assert.Equal(t, "hello\nworld", updated2.(*DetailsSection).previewOutput)
}

func TestDetailsSectionPreviewTickContinues(t *testing.T) {
	ds := &DetailsSection{hoveredName: "sesh"}
	_, cmd := ds.Update(previewTickMsg{name: "sesh"})
	assert.NotNil(t, cmd)

	// Stale tick (hover moved) stops the ticker.
	_, cmd = ds.Update(previewTickMsg{name: "other"})
	assert.Nil(t, cmd)

	// Empty hover stops the ticker.
	_, cmd = ds.Update(previewTickMsg{name: ""})
	assert.Nil(t, cmd)
}

func TestDetailsSectionHoverClearStopsPreview(t *testing.T) {
	ds := &DetailsSection{hoveredName: "sesh", previewOutput: "x"}
	updated, cmd := ds.Update(core.HoveredSessionMsg{Name: ""})
	assert.Nil(t, cmd)
	assert.Equal(t, "", updated.(*DetailsSection).hoveredName)
	assert.Equal(t, "", updated.(*DetailsSection).previewOutput)
}
