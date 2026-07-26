package picker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
)

// nerdGlyph is a nerd font icon (the Go language logo), written as an escape so
// the source stays readable. It occupies a single cell, like the source glyphs.
const nerdGlyph = "\ue627"

// iconTestHome expands paths the way the real home wrapper does, without the
// env-var handling the icon tests don't exercise.
func iconTestHome(t *testing.T) home.Home {
	mockOs := oswrap.NewMockOs(t)
	mockOs.On("UserHomeDir").Return("/home/user", nil).Maybe()
	mockOs.On("ExpandEnv", mock.AnythingOfType("string")).
		Return(func(s string) string { return s }).Maybe()
	return home.NewHome(mockOs)
}

// iconTestWildcards builds the real lister so wildcard icons are matched by the
// same code that matches startup_command and preview_command.
func iconTestWildcards(t *testing.T, config model.Config) WildcardFinder {
	return lister.NewLister(config, iconTestHome(t), nil, nil, nil)
}

func TestBuildIconResolver_NilWithoutIcons(t *testing.T) {
	config := model.Config{
		SessionConfigs:  []model.SessionConfig{{Name: "sesh", Path: "~/c/sesh"}},
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/*"}},
	}
	assert.Nil(t, buildIconResolver(config, iconTestHome(t), iconTestWildcards(t, config)),
		"a config with no icons should cost the picker no per-row lookup")
}

func TestBuildIconResolver_SessionName(t *testing.T) {
	config := model.Config{
		SessionConfigs: []model.SessionConfig{{Name: "notes", Path: "~/second-brain", Icon: "📓"}},
	}
	resolve := buildIconResolver(config, iconTestHome(t), nil)

	assert.Equal(t, "📓", resolve(model.SeshSession{Src: "tmux", Name: "notes"}),
		"a session listed under its configured name gets its icon")
	assert.Equal(t, "", resolve(model.SeshSession{Src: "tmux", Name: "other"}))
}

func TestBuildIconResolver_SessionPath(t *testing.T) {
	config := model.Config{
		SessionConfigs: []model.SessionConfig{{Name: "notes", Path: "~/second-brain", Icon: "📓"}},
	}
	resolve := buildIconResolver(config, iconTestHome(t), nil)

	assert.Equal(t, "📓", resolve(model.SeshSession{Src: "zoxide", Name: "~/second-brain", Path: "/home/user/second-brain"}),
		"the same directory found by another source gets the configured icon")
	assert.Equal(t, "📓", resolve(model.SeshSession{Src: "zoxide", Name: "brain", Path: "/home/user/second-brain/"}),
		"a trailing slash must not defeat the path match")
}

func TestBuildIconResolver_Wildcard(t *testing.T) {
	config := model.Config{
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/nu*", Icon: "🏠"}},
	}
	resolve := buildIconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "🏠", resolve(model.SeshSession{Src: "zoxide", Name: "nutiliti", Path: "/home/user/c/nutiliti"}))
	assert.Equal(t, "", resolve(model.SeshSession{Src: "zoxide", Name: "sesh", Path: "/home/user/c/sesh"}))
}

func TestBuildIconResolver_SessionBeatsWildcard(t *testing.T) {
	config := model.Config{
		SessionConfigs:  []model.SessionConfig{{Name: "sesh", Path: "~/c/sesh", Icon: "📔"}},
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/*", Icon: "🏠"}},
	}
	resolve := buildIconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "📔", resolve(model.SeshSession{Src: "tmux", Name: "sesh", Path: "/home/user/c/sesh"}),
		"the more specific [[session]] icon wins over the pattern")
	assert.Equal(t, "🏠", resolve(model.SeshSession{Src: "zoxide", Name: "other", Path: "/home/user/c/other"}),
		"a path the session block doesn't cover still gets the wildcard icon")
}

func TestBuildIconResolver_SessionWithoutIconFallsThroughToWildcard(t *testing.T) {
	config := model.Config{
		SessionConfigs:  []model.SessionConfig{{Name: "sesh", Path: "~/c/sesh"}},
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/*", Icon: "🏠"}},
	}
	resolve := buildIconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "🏠", resolve(model.SeshSession{Src: "tmux", Name: "sesh", Path: "/home/user/c/sesh"}),
		"an [[session]] with no icon of its own is not treated as an override")
}

func TestBuildIconResolver_EmptyIconIsUnset(t *testing.T) {
	config := model.Config{
		SessionConfigs: []model.SessionConfig{{Name: "sesh", Path: "~/c/sesh", Icon: ""}},
	}
	resolve := buildIconResolver(config, iconTestHome(t), nil)

	assert.Nil(t, resolve, `icon = "" is unset, so the source glyph is kept`)
}

func TestBuildIconResolver_FirstWildcardMatchWins(t *testing.T) {
	config := model.Config{
		WildcardConfigs: []model.WildcardConfig{
			{Pattern: "~/c/*", Icon: "🏠"},
			{Pattern: "~/c/nu*", Icon: "📓"},
		},
	}
	resolve := buildIconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "🏠", resolve(model.SeshSession{Src: "zoxide", Name: "nutiliti", Path: "/home/user/c/nutiliti"}),
		"wildcards resolve in config order, as they do for startup_command")
}

func TestBuildIconResolver_FirstWildcardMatchWithoutIconFallsBack(t *testing.T) {
	config := model.Config{
		WildcardConfigs: []model.WildcardConfig{
			{Pattern: "~/c/*", StartupCommand: "nvim"},
			{Pattern: "~/c/nu*", Icon: "📓"},
		},
	}
	resolve := buildIconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "", resolve(model.SeshSession{Src: "zoxide", Name: "nutiliti", Path: "/home/user/c/nutiliti"}),
		"the first matching pattern still wins, so the row keeps its source glyph")
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
