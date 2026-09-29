package picker

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshmedeski/sesh/v2/model"
)

// nerdGlyph is a nerd font icon (the Go language logo), written as an escape so
// the source stays readable. It occupies a single cell, like the source glyphs.
const nerdGlyph = "\ue627"

func TestIconColWidth(t *testing.T) {
	assert.Equal(t, 1, iconColWidth(model.Config{}),
		"the source glyphs need a single cell")

	assert.Equal(t, 1, iconColWidth(model.Config{
		SessionConfigs: []model.SessionConfig{{Name: "sesh", Icon: nerdGlyph}},
	}), "a single-width nerd font glyph keeps the column as it was")

	assert.Equal(t, 2, iconColWidth(model.Config{
		SessionConfigs: []model.SessionConfig{{Name: "notes", Icon: "📓"}},
	}), "a double-width emoji widens the column")

	assert.Equal(t, 2, iconColWidth(model.Config{
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/*", Icon: "🏠"}},
	}), "an emoji on a wildcard widens the column too")
}
