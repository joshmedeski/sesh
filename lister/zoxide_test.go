package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestListZoxideSessions(t *testing.T) {
	t.Run("should list zoxide sessions", func(t *testing.T) {
		mockZoxide := new(zoxide.MockZoxide)
		mockZoxide.On("ListResults").Return([]model.ZoxideResult{
			{
				Score: 0.3,
				Path:  "/Users/joshmedeski/.config/fish",
			},
			{
				Score: 0.5,
				Path:  "/Users/joshmedeski/.config/sesh",
			},
		}, nil)
		mockHome := new(home.MockHome)
		sessions, err := listZoxideSessions(mockZoxide, mockHome)
		assert.Equal(t, "zoxide:sesh/main", sessions.OrderedIndex[0])
		assert.Equal(t, "Score", sessions.Directory["tmux:sesh/main"].Name)
		assert.Equal(t, "zoxide:sesh/v2", sessions.OrderedIndex[1])
		assert.Equal(t, nil, err)
	})
}
