package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestRenameSession(t *testing.T) {
	t.Run("calls tmux rename-session with target and new name", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "rename-session", "-t", "400-status", "400-status — warm the cache").
			Return("", nil)

		result, err := tmux.RenameSession("400-status", "400-status — warm the cache")

		assert.Nil(t, err)
		assert.Equal(t, "", result)
	})
}
