package picker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
)

// nerdGlyph is a nerd font icon (the Go language logo), written as an escape so
// the source stays readable. It occupies a single cell, like the source glyphs.
const nerdGlyph = "\ue627"

func forgeResolver(f forge.Forge) forge.Resolver {
	return forge.NewResolver(
		func(cfg model.WorktreeConfig) (forge.Forge, error) {
			return f, nil
		},
		func(path string) (forge.Forge, error) {
			return f, nil
		},
	)
}

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

func TestWorktreeResolver_OffByDefault(t *testing.T) {
	mockOs := oswrap.NewMockOs(t)
	mockOs.On("UserHomeDir").Return("/home/user", nil).Maybe()
	mockOs.On("ExpandEnv", mock.AnythingOfType("string")).Return(func(s string) string { return s }).Maybe()
	h := home.NewHome(mockOs)
	config := model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{Repo: "org/repo", Path: "~/c/repo", WorktreeDir: "w"}},
	}

	assert.Nil(t, worktreeResolver(config, h, nil, nil))

	mockF := forge.NewMockForge(t)
	mockF.EXPECT().WorktreeNumber("7").Return(7, false, nil)

	config.TUI.ShowWorktreeTitles = true
	resolve := worktreeResolver(config, h, forgeResolver(mockF), nil)
	require.NotNil(t, resolve)
	assert.Equal(t, 7, resolve("/home/user/c/repo/w/7").Number)
}
