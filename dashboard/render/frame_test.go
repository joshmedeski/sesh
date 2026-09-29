package render

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Frame rendering ---

func TestRenderFrame_ContainsTitlesAndJunctions(t *testing.T) {
	panes := []FramePane{
		{Title: "Sessions", Content: "row1\nrow2", Width: 16, Focused: true},
		{Title: "Details", Content: "x", Width: 16, Focused: false},
	}
	out := RenderFrame(panes, 6)
	assert.Contains(t, out, "Sessions")
	assert.Contains(t, out, "Details")
	for _, j := range []string{"┬", "┴", "┌", "┐", "└", "┘"} {
		assert.Contains(t, out, j)
	}
}

func TestRenderFrame_TopBorderTitlesInOrder(t *testing.T) {
	panes := []FramePane{
		{Title: "Sessions", Content: "", Width: 16, Focused: false},
		{Title: "Details", Content: "", Width: 16, Focused: false},
		{Title: "Git", Content: "", Width: 16, Focused: false},
	}
	out := RenderFrame(panes, 4)
	top := strings.Split(out, "\n")[0]
	iS := strings.Index(top, "Sessions")
	iD := strings.Index(top, "Details")
	iG := strings.Index(top, "Git")
	assert.True(t, iS >= 0 && iD > iS && iG > iD)
}

func TestRenderFrame_FocusedTitleAccent(t *testing.T) {
	panes := []FramePane{
		{Title: "Sessions", Content: "", Width: 16, Focused: true},
		{Title: "Details", Content: "", Width: 16, Focused: false},
	}
	out := RenderFrame(panes, 4)
	assert.Contains(t, out, accentStyle().Render(" Sessions "))
}

func TestRenderFrame_ContentHeight(t *testing.T) {
	panes := []FramePane{{Title: "A", Content: "line", Width: 6, Focused: false}}
	out := RenderFrame(panes, 5)
	lines := strings.Split(out, "\n")
	// 1 top + 3 content + 1 bottom = 5 lines
	assert.Len(t, lines, 5)
}
