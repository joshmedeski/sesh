package icon

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/ansi"
	"github.com/joshmedeski/sesh/v2/model"
)

func TestSource(t *testing.T) {
	glyphs, err := SourceGlyphs(nil)
	require.NoError(t, err)

	for _, src := range []string{"tmux", "config", "zoxide", "tmuxinator", "tmux-pane"} {
		icn, clr := glyphs.Source(src)
		assert.NotEmpty(t, strings.TrimSpace(icn), "icon for %s should not be empty", src)
		assert.NotEqual(t, "? ", icn, "icon for %s should not be fallback", src)
		assert.NotNil(t, clr, "color for %s should not be nil", src)
	}

	icn, clr := glyphs.Source("other")
	assert.Equal(t, "? ", icn)
	assert.NotNil(t, clr)
}

func TestSourceGlyphs_Overrides(t *testing.T) {
	glyphs, err := SourceGlyphs(map[string]model.SourceIcon{
		"tmux":   {Icon: "T", Color: "magenta"},
		"zoxide": {Icon: "📁"},
		"config": {Color: "Bright-Yellow"},
	})
	require.NoError(t, err)

	assert.Equal(t, Glyph{Icon: "T", ColorCode: ansi.Magenta}, glyphs["tmux"])
	assert.Equal(t, Glyph{Icon: "📁", ColorCode: ansi.Cyan}, glyphs["zoxide"])
	assert.Equal(t, Glyph{Icon: defaultGlyphs["config"].Icon, ColorCode: ansi.BrightYellow}, glyphs["config"])
	assert.Equal(t, defaultGlyphs["tmuxinator"], glyphs["tmuxinator"])

	icn, clr := glyphs.Source("tmux")
	assert.Equal(t, "T ", icn)
	assert.Equal(t, lipgloss.ANSIColor(5), clr)
}

func TestSourceGlyphs_DoesNotMutateDefaults(t *testing.T) {
	before := defaultGlyphs["tmux"]
	_, err := SourceGlyphs(map[string]model.SourceIcon{"tmux": {Icon: "T", Color: "red"}})
	require.NoError(t, err)
	assert.Equal(t, before, defaultGlyphs["tmux"])
}

func TestSourceGlyphs_Errors(t *testing.T) {
	_, err := SourceGlyphs(map[string]model.SourceIcon{"zoxdie": {Icon: "Z"}})
	assert.ErrorContains(t, err, `icons: unknown source "zoxdie"`)

	glyphs, err := SourceGlyphs(map[string]model.SourceIcon{"tmux": {Icon: "T", Color: "purple"}})
	assert.ErrorContains(t, err, `icons.tmux.color: unknown color "purple"`)
	assert.Equal(t, defaultGlyphs["tmux"], glyphs["tmux"])
}

func TestAddIcon_ConfiguredColor(t *testing.T) {
	ic := NewIcon(model.Config{Icons: map[string]model.SourceIcon{"tmux": {Icon: "T", Color: "red"}}})
	session := model.SeshSession{Src: "tmux", Name: "dotfiles"}

	assert.Equal(t, "\033[31mT\033[39m dotfiles", ic.AddIcon(session))
	assert.Equal(t, "T dotfiles", ic.AddIconNoColor(session))
}

func TestRemoveIcon(t *testing.T) {
	ic := NewIcon(model.Config{Icons: map[string]model.SourceIcon{"zoxide": {Icon: "📁"}}})

	assert.Equal(t, "dotfiles", ic.RemoveIcon("📁 dotfiles"))
	assert.Equal(t, "dotfiles", ic.RemoveIcon("\uebc8 dotfiles"))
	assert.Equal(t, "dotfiles", ic.RemoveIcon("dotfiles"))
}

func TestColumnWidth_SourceIcons(t *testing.T) {
	assert.Equal(t, 1, ColumnWidth(model.Config{}))
	assert.Equal(t, 2, ColumnWidth(model.Config{Icons: map[string]model.SourceIcon{"zoxide": {Icon: "📁"}}}))
}
