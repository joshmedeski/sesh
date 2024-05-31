package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/model"
	"github.com/stretchr/testify/assert"
)

func TestListConfigSessions(t *testing.T) {
	t.Run("should list config sessions", func(t *testing.T) {
		config := model.Config{
			SessionConfigs: []model.SessionConfig{
				{
					Name: "sesh config",
					Path: "/Users/joshmedeski/.config/sesh",
				},
			},
		}
		sessions := listConfigSessions(config)
		assert.Equal(t, "config:sesh config", sessions.OrderedIndex[0])
		assert.Equal(t, "/Users/joshmedeski/.config/sesh", sessions.Directory["config:sesh config"].Path)
	})
}
