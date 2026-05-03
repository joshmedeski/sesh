package startup

import (
	"testing"

	"github.com/Wingsdh/cc-sesh/v2/home"
	"github.com/Wingsdh/cc-sesh/v2/lister"
	"github.com/Wingsdh/cc-sesh/v2/model"
	"github.com/Wingsdh/cc-sesh/v2/oswrap"
	"github.com/Wingsdh/cc-sesh/v2/replacer"
	"github.com/Wingsdh/cc-sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
)

func TestResolveCommand(t *testing.T) {
	newStartup := func(config model.Config) (*RealStartup, *lister.MockLister, *replacer.MockReplacer) {
		mockLister := new(lister.MockLister)
		mockReplacer := new(replacer.MockReplacer)
		s := &RealStartup{
			os:       new(oswrap.MockOs),
			lister:   mockLister,
			tmux:     new(tmux.MockTmux),
			config:   config,
			home:     new(home.MockHome),
			replacer: mockReplacer,
		}
		return s, mockLister, mockReplacer
	}

	session := model.SeshSession{Name: "proj", Path: "/p"}

	t.Run("returns session-level config command first", func(t *testing.T) {
		s, mockLister, mockReplacer := newStartup(model.Config{})
		mockLister.On("FindConfigSession", "proj").Return(model.SeshSession{
			Name:           "proj",
			Path:           "/p",
			StartupCommand: "nvim",
		}, true)
		mockReplacer.On("Replace", "nvim", map[string]string{"{}": "/p"}).Return("nvim")

		cmd, err := s.ResolveCommand(session)
		assert.Nil(t, err)
		assert.Equal(t, "nvim", cmd)
	})

	t.Run("falls through to wildcard when no session config", func(t *testing.T) {
		s, mockLister, mockReplacer := newStartup(model.Config{})
		mockLister.On("FindConfigSession", "proj").Return(model.SeshSession{}, false)
		mockLister.On("FindConfigWildcard", "/p").Return(model.WildcardConfig{
			Pattern:        "/*",
			StartupCommand: "htop",
		}, true)
		mockReplacer.On("Replace", "htop", map[string]string{"{}": "/p"}).Return("htop")

		cmd, err := s.ResolveCommand(session)
		assert.Nil(t, err)
		assert.Equal(t, "htop", cmd)
	})

	t.Run("falls through to default when no session/wildcard match", func(t *testing.T) {
		s, mockLister, mockReplacer := newStartup(model.Config{
			DefaultSessionConfig: model.DefaultSessionConfig{StartupCommand: "lazygit"},
		})
		mockLister.On("FindConfigSession", "proj").Return(model.SeshSession{}, false)
		mockLister.On("FindConfigWildcard", "/p").Return(model.WildcardConfig{}, false)
		mockReplacer.On("Replace", "lazygit", map[string]string{"{}": "/p"}).Return("lazygit")

		cmd, err := s.ResolveCommand(session)
		assert.Nil(t, err)
		assert.Equal(t, "lazygit", cmd)
	})

	t.Run("returns empty when no strategy yields a command", func(t *testing.T) {
		s, mockLister, _ := newStartup(model.Config{})
		mockLister.On("FindConfigSession", "proj").Return(model.SeshSession{}, false)
		mockLister.On("FindConfigWildcard", "/p").Return(model.WildcardConfig{}, false)

		cmd, err := s.ResolveCommand(session)
		assert.Nil(t, err)
		assert.Equal(t, "", cmd)
	})
}
