package zoxide

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	t.Run("substitutes the query into the default zoxide command", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{})
		mockShell.EXPECT().
			PrepareCmd("zoxide query {}", map[string]string{"{}": "sesh"}).
			Return([]string{"zoxide", "query", "sesh"}, nil)
		mockShell.EXPECT().Cmd("zoxide", "query", "sesh").
			Return("/Users/joshmedeski/c/sesh/v2", nil)
		result, err := zoxide.Query("sesh")
		assert.Nil(t, err)
		assert.Equal(t, &model.ZoxideResult{Path: "/Users/joshmedeski/c/sesh/v2", Score: 0}, result)
	})

	t.Run("substitutes the query into a custom command", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{
			QueryCommand: "fasd -d {}",
		})
		mockShell.EXPECT().
			PrepareCmd("fasd -d {}", map[string]string{"{}": "sesh"}).
			Return([]string{"fasd", "-d", "sesh"}, nil)
		mockShell.EXPECT().Cmd("fasd", "-d", "sesh").
			Return("/Users/joshmedeski/c/sesh/v2", nil)
		result, err := zoxide.Query("sesh")
		assert.Nil(t, err)
		assert.Equal(t, &model.ZoxideResult{Path: "/Users/joshmedeski/c/sesh/v2", Score: 0}, result)
	})
}
