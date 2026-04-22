package tmuxinator

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	t.Run("runs tmuxinator with --no-attach and --name matching the project", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		tmuxinator := &RealTmuxinator{shell: mockShell}
		mockShell.EXPECT().
			Cmd("tmuxinator", "start", "--no-attach", "--name", "sys", "sys").
			Return("", nil)

		_, err := tmuxinator.Start("sys")

		assert.Nil(t, err)
		mockShell.AssertExpectations(t)
	})
}
