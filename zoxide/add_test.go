package zoxide

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Run("substitutes the path into the default zoxide command", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{})
		mockShell.EXPECT().
			PrepareCmd("zoxide add {}", map[string]string{"{}": "/Users/joshmedeski/c/sesh/v2"}).
			Return([]string{"zoxide", "add", "/Users/joshmedeski/c/sesh/v2"}, nil)
		mockShell.EXPECT().Cmd("zoxide", "add", "/Users/joshmedeski/c/sesh/v2").Return("", nil)
		err := zoxide.Add("/Users/joshmedeski/c/sesh/v2")
		assert.Nil(t, err)
	})

	t.Run("substitutes the path into a custom command", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{
			AddCommand: "fasd -A {}",
		})
		mockShell.EXPECT().
			PrepareCmd("fasd -A {}", map[string]string{"{}": "/Users/joshmedeski/c/sesh/v2"}).
			Return([]string{"fasd", "-A", "/Users/joshmedeski/c/sesh/v2"}, nil)
		mockShell.EXPECT().Cmd("fasd", "-A", "/Users/joshmedeski/c/sesh/v2").Return("", nil)
		err := zoxide.Add("/Users/joshmedeski/c/sesh/v2")
		assert.Nil(t, err)
	})
}
