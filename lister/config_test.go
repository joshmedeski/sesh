package lister

import (
	"log"
	"testing"

	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestListConfigSessions(t *testing.T) {
	t.Run("should list config sessions", func(t *testing.T) {
		mockHome := new(home.MockHome)
		mockZoxide := new(zoxide.MockZoxide)
		mockTmux := new(tmux.MockTmux)
		config := model.Config{
			SessionConfigs: []model.SessionConfig{
				{
					Name: "sesh config",
					Path: "/Users/joshmedeski/.config/sesh",
				},
			},
		}
		lister := NewLister(config, mockHome, mockTmux, mockZoxide)

		realLister, ok := lister.(*RealLister)
		if !ok {
			log.Fatal("Cannot convert lister to *RealLister")
		}
		sessions, err := listConfig(realLister)
		assert.Nil(t, err)
		assert.Equal(t, "config:sesh config", sessions.OrderedIndex[0])
		assert.Equal(t, "/Users/joshmedeski/.config/sesh", sessions.Directory["config:sesh config"].Path)
	})
}
