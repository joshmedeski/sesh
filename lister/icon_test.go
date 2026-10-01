package lister

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
)

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
	return NewLister(config, iconTestHome(t), nil, nil, nil, nil)
}

func TestIconResolver_NilWithoutIcons(t *testing.T) {
	config := model.Config{
		SessionConfigs:  []model.SessionConfig{{Name: "sesh", Path: "~/c/sesh"}},
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/*"}},
	}
	assert.Nil(t, IconResolver(config, iconTestHome(t), iconTestWildcards(t, config)),
		"a config with no icons should cost no per-row lookup")
}

func TestIconResolver_SessionName(t *testing.T) {
	config := model.Config{
		SessionConfigs: []model.SessionConfig{{Name: "notes", Path: "~/second-brain", Icon: "📓"}},
	}
	resolve := IconResolver(config, iconTestHome(t), nil)

	assert.Equal(t, "📓", resolve(model.SeshSession{Src: "tmux", Name: "notes"}),
		"a session listed under its configured name gets its icon")
	assert.Equal(t, "", resolve(model.SeshSession{Src: "tmux", Name: "other"}))
}

func TestIconResolver_SessionPath(t *testing.T) {
	config := model.Config{
		SessionConfigs: []model.SessionConfig{{Name: "notes", Path: "~/second-brain", Icon: "📓"}},
	}
	resolve := IconResolver(config, iconTestHome(t), nil)

	assert.Equal(t, "📓", resolve(model.SeshSession{Src: "zoxide", Name: "~/second-brain", Path: "/home/user/second-brain"}),
		"the same directory found by another source gets the configured icon")
	assert.Equal(t, "📓", resolve(model.SeshSession{Src: "zoxide", Name: "brain", Path: "/home/user/second-brain/"}),
		"a trailing slash must not defeat the path match")
}

func TestIconResolver_Wildcard(t *testing.T) {
	config := model.Config{
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/nu*", Icon: "🏠"}},
	}
	resolve := IconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "🏠", resolve(model.SeshSession{Src: "zoxide", Name: "nutiliti", Path: "/home/user/c/nutiliti"}))
	assert.Equal(t, "", resolve(model.SeshSession{Src: "zoxide", Name: "sesh", Path: "/home/user/c/sesh"}))
}

func TestIconResolver_SessionBeatsWildcard(t *testing.T) {
	config := model.Config{
		SessionConfigs:  []model.SessionConfig{{Name: "sesh", Path: "~/c/sesh", Icon: "📔"}},
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/*", Icon: "🏠"}},
	}
	resolve := IconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "📔", resolve(model.SeshSession{Src: "tmux", Name: "sesh", Path: "/home/user/c/sesh"}),
		"the more specific [[session]] icon wins over the pattern")
	assert.Equal(t, "🏠", resolve(model.SeshSession{Src: "zoxide", Name: "other", Path: "/home/user/c/other"}),
		"a path the session block doesn't cover still gets the wildcard icon")
}

func TestIconResolver_SessionWithoutIconFallsThroughToWildcard(t *testing.T) {
	config := model.Config{
		SessionConfigs:  []model.SessionConfig{{Name: "sesh", Path: "~/c/sesh"}},
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/*", Icon: "🏠"}},
	}
	resolve := IconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "🏠", resolve(model.SeshSession{Src: "tmux", Name: "sesh", Path: "/home/user/c/sesh"}),
		"an [[session]] with no icon of its own is not treated as an override")
}

func TestIconResolver_EmptyIconIsUnset(t *testing.T) {
	config := model.Config{
		SessionConfigs: []model.SessionConfig{{Name: "sesh", Path: "~/c/sesh", Icon: ""}},
	}
	resolve := IconResolver(config, iconTestHome(t), nil)

	assert.Nil(t, resolve, `icon = "" is unset, so the source glyph is kept`)
}

func TestIconResolver_FirstWildcardMatchWins(t *testing.T) {
	config := model.Config{
		WildcardConfigs: []model.WildcardConfig{
			{Pattern: "~/c/*", Icon: "🏠"},
			{Pattern: "~/c/nu*", Icon: "📓"},
		},
	}
	resolve := IconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "🏠", resolve(model.SeshSession{Src: "zoxide", Name: "nutiliti", Path: "/home/user/c/nutiliti"}),
		"wildcards resolve in config order, as they do for startup_command")
}

func TestIconResolver_FirstWildcardMatchWithoutIconFallsBack(t *testing.T) {
	config := model.Config{
		WildcardConfigs: []model.WildcardConfig{
			{Pattern: "~/c/*", StartupCommand: "nvim"},
			{Pattern: "~/c/nu*", Icon: "📓"},
		},
	}
	resolve := IconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "", resolve(model.SeshSession{Src: "zoxide", Name: "nutiliti", Path: "/home/user/c/nutiliti"}),
		"the first matching pattern still wins, so the row keeps its source glyph")
}

func TestIconResolver_Worktree(t *testing.T) {
	config := model.Config{
		SessionConfigs:  []model.SessionConfig{{Name: "nutiliti", Path: "~/c/nu", Icon: "🏠"}},
		WorktreeConfigs: []model.WorktreeConfig{{Repo: "Nutiliti/nutiliti", Path: "~/c/nu", WorktreeDir: "w", Icon: "🌳"}},
		WildcardConfigs: []model.WildcardConfig{{Pattern: "~/c/**", Icon: "🚀"}},
	}
	resolve := IconResolver(config, iconTestHome(t), iconTestWildcards(t, config))

	assert.Equal(t, "🌳", resolve(model.SeshSession{Src: "tmux", Name: "nu/w/123", Path: "/home/user/c/nu/w/123"}),
		"a worktree session gets its [[worktree]] icon over a matching wildcard")
	assert.Equal(t, "🌳", resolve(model.SeshSession{Src: "zoxide", Name: "web", Path: "/home/user/c/nu/w/123/apps/web"}))
	assert.Equal(t, "🏠", resolve(model.SeshSession{Src: "tmux", Name: "nu", Path: "/home/user/c/nu"}),
		"the repo root itself is not a worktree")
	assert.Equal(t, "🚀", resolve(model.SeshSession{Src: "zoxide", Name: "nu/wiki", Path: "/home/user/c/nu/wiki"}),
		"a sibling sharing the root's prefix is not under it")
}

func TestIconResolver_WorktreeDefaultDir(t *testing.T) {
	config := model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{Repo: "joshmedeski/sesh", Path: "~/c/sesh", Icon: "⚡"}},
	}
	resolve := IconResolver(config, iconTestHome(t), nil)

	assert.Equal(t, "⚡", resolve(model.SeshSession{Src: "tmux", Name: "sesh/.wk/7", Path: "/home/user/c/sesh/.wk/7"}))
}
