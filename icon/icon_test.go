package icon

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestAddIconNoColor(t *testing.T) {
	live := model.SeshSession{Src: "tmux", Name: "sesh"}

	t.Run("tmux sessions get the tmux glyph", func(t *testing.T) {
		i := NewIcon(model.Config{}, "tmux")
		assert.Equal(t, tmuxIcon+" sesh", i.AddIconNoColor(live))
	})

	t.Run("live sessions on herdr get the herdr glyph", func(t *testing.T) {
		i := NewIcon(model.Config{}, "herdr")
		assert.Equal(t, herdrIcon+" sesh", i.AddIconNoColor(live))
	})

	t.Run("other sources are unaffected by the backend", func(t *testing.T) {
		i := NewIcon(model.Config{}, "herdr")
		z := model.SeshSession{Src: "zoxide", Name: "dir"}
		assert.Equal(t, zoxideIcon+" dir", i.AddIconNoColor(z))
	})

	t.Run("unknown source has no icon", func(t *testing.T) {
		i := NewIcon(model.Config{}, "tmux")
		assert.Equal(t, "x", i.AddIconNoColor(model.SeshSession{Src: "nope", Name: "x"}))
	})
}

func TestRemoveIcon(t *testing.T) {
	i := NewIcon(model.Config{}, "herdr")
	for _, g := range Glyphs {
		assert.Equal(t, "my session", i.RemoveIcon(g.Icon+" my session"))
	}
	assert.Equal(t, "plain", i.RemoveIcon("plain"))
}
