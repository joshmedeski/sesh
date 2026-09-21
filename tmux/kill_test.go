package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestKillSession(t *testing.T) {
	t.Run("calls tmux kill-session with a session-scoped target", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "kill-session", "-t", "my-project:").
			Return("", nil)

		result, err := tmux.KillSession("my-project")

		assert.Nil(t, err)
		assert.Equal(t, "", result)
	})
}
