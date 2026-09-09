package lister

import (
	"log"
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestListConfigSessions(t *testing.T) {
	mockHome := new(home.MockHome)
	mockHome.On("ExpandPath", "/Users/joshmedeski/.config/sesh").Return("/Users/joshmedeski/.config/sesh", nil)
	mockZoxide := new(zoxide.MockZoxide)
	mockTmux := new(tmux.MockTmux)
	mockTmuxinator := new(tmuxinator.MockTmuxinator)
	config := model.Config{
		SessionConfigs: []model.SessionConfig{
			{
				Name: "sesh config",
				Path: "/Users/joshmedeski/.config/sesh",
			},
		},
	}
	lister := NewLister(config, mockHome, mockTmux, mockZoxide, mockTmuxinator)

	realLister, ok := lister.(*RealLister)
	if !ok {
		log.Fatal("Cannot convert lister to *RealLister")
	}

	// TODO: make sure Path has home expanded
	t.Run("should list config sessions", func(t *testing.T) {
		sessions, err := listConfig(realLister)
		assert.Nil(t, err)
		assert.Equal(t, "config:sesh config", sessions.OrderedIndex[0])
		assert.Equal(t, "/Users/joshmedeski/.config/sesh", sessions.Directory["config:sesh config"].Path)
		assert.Equal(t, "sesh config", sessions.Directory["config:sesh config"].Name)
	})

	t.Run("should find config session", func(t *testing.T) {
		sessions, exists := lister.FindConfigSession("sesh config")
		assert.Equal(t, true, exists)
		assert.Equal(t, "sesh config", sessions.Name)
	})
}

func TestListConfigIncludesAlias(t *testing.T) {
	mockHome := new(home.MockHome)
	mockHome.On("ExpandPath", "/home/user/wallpaper").Return("/home/user/wallpaper", nil)
	mockHome.On("ExpandPath", "/home/user/notes").Return("/home/user/notes", nil)
	mockTmux := new(tmux.MockTmux)
	mockZoxide := new(zoxide.MockZoxide)
	mockTmuxinator := new(tmuxinator.MockTmuxinator)
	config := model.Config{
		SessionConfigs: []model.SessionConfig{
			{Name: "wallpaper", Path: "/home/user/wallpaper", Alias: "wp"},
			{Name: "notes", Path: "/home/user/notes"},
		},
	}
	lister := NewLister(config, mockHome, mockTmux, mockZoxide, mockTmuxinator)

	realLister, ok := lister.(*RealLister)
	if !ok {
		log.Fatal("Cannot convert lister to *RealLister")
	}

	sessions, err := listConfig(realLister)
	assert.Nil(t, err)
	assert.Equal(t, "wp", sessions.Directory["config:wallpaper"].Alias)
	assert.Equal(t, "", sessions.Directory["config:notes"].Alias)
}
