package zoxide

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	t.Run("substitutes the path into the default zoxide command", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{})
		mockShell.EXPECT().
			PrepareCmd("zoxide remove {}", map[string]string{"{}": "/Users/joshmedeski/c/sesh/v2"}).
			Return([]string{"zoxide", "remove", "/Users/joshmedeski/c/sesh/v2"}, nil)
		mockShell.EXPECT().Cmd("zoxide", "remove", "/Users/joshmedeski/c/sesh/v2").Return("", nil)
		err := zoxide.Remove("/Users/joshmedeski/c/sesh/v2")
		assert.Nil(t, err)
	})

	t.Run("substitutes the path into a custom command", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{
			RemoveCommand: "fasd -D {}",
		})
		mockShell.EXPECT().
			PrepareCmd("fasd -D {}", map[string]string{"{}": "/Users/joshmedeski/c/sesh/v2"}).
			Return([]string{"fasd", "-D", "/Users/joshmedeski/c/sesh/v2"}, nil)
		mockShell.EXPECT().Cmd("fasd", "-D", "/Users/joshmedeski/c/sesh/v2").Return("", nil)
		err := zoxide.Remove("/Users/joshmedeski/c/sesh/v2")
		assert.Nil(t, err)
	})

	t.Run("returns the error when the command fails", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{})
		mockShell.EXPECT().
			PrepareCmd("zoxide remove {}", map[string]string{"{}": "/tmp/gone"}).
			Return([]string{"zoxide", "remove", "/tmp/gone"}, nil)
		mockShell.EXPECT().
			Cmd("zoxide", "remove", "/tmp/gone").
			Return("", errors.New("exit status 1"))
		err := zoxide.Remove("/tmp/gone")
		assert.Error(t, err)
	})
}
