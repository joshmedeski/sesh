package connector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshmedeski/sesh/v2/model"
)

func TestResolveAlias(t *testing.T) {
	c := &RealConnector{config: model.Config{
		SessionConfigs: []model.SessionConfig{
			{Name: "wallpaper", Alias: "wp"},
			{Name: "dotfiles", Alias: "DOT"},
			{Name: "notes"},
		},
	}}

	t.Run("resolves an alias to its session name", func(t *testing.T) {
		assert.Equal(t, "wallpaper", c.resolveAlias("wp"))
	})

	t.Run("matches aliases case-insensitively", func(t *testing.T) {
		assert.Equal(t, "dotfiles", c.resolveAlias("dot"))
		assert.Equal(t, "wallpaper", c.resolveAlias("WP"))
	})

	t.Run("passes unknown names through untouched", func(t *testing.T) {
		assert.Equal(t, "my-project", c.resolveAlias("my-project"))
		assert.Equal(t, "", c.resolveAlias(""))
	})

	t.Run("does not match a partial alias", func(t *testing.T) {
		assert.Equal(t, "w", c.resolveAlias("w"))
	})

	t.Run("leaves session names alone when they have no alias", func(t *testing.T) {
		assert.Equal(t, "notes", c.resolveAlias("notes"))
	})
}
