package tmux

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestDisplayMessage(t *testing.T) {
	t.Run("shows the escaped message until a key is pressed", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "display-message", "-d", "0", "sesh: pane ##1 is not defined in config").
			Return("", nil)

		_, err := tmux.DisplayMessage("sesh: pane #1 is not defined in config")

		assert.NoError(t, err)
	})

	t.Run("returns the tmux error", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().Cmd("tmux", "display-message", "-d", "0", "hi").Return("", errors.New("no client"))

		_, err := tmux.DisplayMessage("hi")

		assert.EqualError(t, err, "no client")
	})
}
