package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Header / footer rendering ---

func TestRenderHeaderContainsTabsAndCount(t *testing.T) {
	h := RenderHeader(0, 3, 80)
	assert.Contains(t, h, "Open")
	assert.Contains(t, h, "Configured")
	assert.Contains(t, h, "3 active")
}

func TestRenderHeaderSmallDropsCount(t *testing.T) {
	h := RenderHeader(0, 3, 40)
	assert.NotContains(t, h, "active")
}

func TestRenderFooterPage0(t *testing.T) {
	f := RenderFooter(0, 120, "name", false, "")
	assert.Contains(t, f, "ctrl+d")
	assert.Contains(t, f, "1-9")
	assert.Contains(t, f, "panes")
	assert.Contains(t, f, "sort:name")
	assert.NotContains(t, f, "widgets")
}

func TestRenderFooterPage0DropsLabelsWhenOverflow(t *testing.T) {
	// The labeled tab-1 footer (after dropping the `t group` bind) is 94 cols,
	// so at 90 cols labels are dropped (keys only) rather than wrapping.
	f := RenderFooter(0, 90, "name", false, "")
	assert.Contains(t, f, "1-9")
	assert.NotContains(t, f, "panes")
}

func TestRenderFooterPage1(t *testing.T) {
	f := RenderFooter(1, 100, "name", false, "")
	assert.NotContains(t, f, "ctrl+d")
	assert.Contains(t, f, "filter")
	assert.Contains(t, f, "refresh")
}

func TestRenderFooterNarrowKeysOnly(t *testing.T) {
	f := RenderFooter(0, 50, "name", false, "")
	assert.Contains(t, f, "ctrl+d")
	assert.NotContains(t, f, "kill")
}

func TestRenderFooterTiny(t *testing.T) {
	f := RenderFooter(0, 20, "name", false, "")
	assert.Contains(t, f, "tab")
	assert.Contains(t, f, "j/k")
	assert.Contains(t, f, "enter")
	assert.NotContains(t, f, "ctrl+d")
}

func TestRenderFooterFiltering(t *testing.T) {
	f := RenderFooter(0, 120, "name", true, "foo")
	assert.Contains(t, f, "filter:")
	assert.Contains(t, f, "foo")
	assert.Contains(t, f, "esc")
	assert.Contains(t, f, "enter")
	// The filter line replaces all binds, so no quit/help/panes binds remain.
	assert.NotContains(t, f, "quit")
	assert.NotContains(t, f, "help")
	assert.NotContains(t, f, "ctrl+d")
}
