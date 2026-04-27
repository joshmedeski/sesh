package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestSelectWindow(t *testing.T) {
	t.Run("calls tmux select-window", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().Cmd("tmux", "select-window", "-t", "editor").Return("", nil)
		result, err := tmux.SelectWindow("editor")
		assert.Nil(t, err)
		assert.Equal(t, "", result)
	})
}
