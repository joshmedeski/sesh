package startup

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
)

func TestConfigWildcardStartupStrategy(t *testing.T) {
	mockLister := new(lister.MockLister)
	mockTmux := new(tmux.MockTmux)
	mockHome := new(home.MockHome)
	mockReplacer := new(replacer.MockReplacer)

	s := &RealStartup{
		lister:   mockLister,
		tmux:     mockTmux,
		config:   model.Config{},
		home:     mockHome,
		replacer: mockReplacer,
	}

	t.Run("should return startup command from matching wildcard", func(t *testing.T) {
		session := model.SeshSession{
			Name: "myapp",
			Path: "/Users/test/projects/myapp",
		}
		mockLister.On("FindConfigWildcard", "/Users/test/projects/myapp").Return(model.WildcardConfig{
			Pattern:        "~/projects/*",
			StartupCommand: "nvim",
		}, true)
		mockReplacer.On("Replace", "nvim", map[string]string{"{}": "/Users/test/projects/myapp"}).Return("nvim")

		cmd, err := configWildcardStartupStrategy(s, session)
		assert.Nil(t, err)
		assert.Equal(t, "nvim", cmd)
	})

	t.Run("should return empty when no wildcard matches", func(t *testing.T) {
		session := model.SeshSession{
			Name: "other",
			Path: "/Users/test/other",
		}
		mockLister.On("FindConfigWildcard", "/Users/test/other").Return(model.WildcardConfig{}, false)

		cmd, err := configWildcardStartupStrategy(s, session)
		assert.Nil(t, err)
		assert.Equal(t, "", cmd)
	})

	t.Run("should return empty when wildcard has disable_startup_command", func(t *testing.T) {
		session := model.SeshSession{
			Name: "disabled",
			Path: "/Users/test/projects/disabled",
		}
		mockLister.On("FindConfigWildcard", "/Users/test/projects/disabled").Return(model.WildcardConfig{
			Pattern:             "~/projects/*",
			StartupCommand:      "nvim",
			DisableStartCommand: true,
		}, true)

		cmd, err := configWildcardStartupStrategy(s, session)
		assert.Nil(t, err)
		assert.Equal(t, "", cmd)
	})
}
