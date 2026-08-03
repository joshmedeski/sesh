package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestNewWindowInSession(t *testing.T) {
	t.Run("targets the session unambiguously with a trailing colon", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "new-window", "-n", "agent", "-c", "/home/dev/project", "-t", "project:").
			Return("", nil)

		result, err := tmux.NewWindowInSession("agent", "/home/dev/project", "project")

		assert.Nil(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("omits the target when no session is provided", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "new-window", "-n", "agent", "-c", "/home/dev/project").
			Return("", nil)

		result, err := tmux.NewWindowInSession("agent", "/home/dev/project", "")

		assert.Nil(t, err)
		assert.Equal(t, "", result)
	})
}

func TestNextWindowInSession(t *testing.T) {
	t.Run("targets the session unambiguously with a trailing colon", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "next-window", "-t", "project:").
			Return("", nil)

		result, err := tmux.NextWindowInSession("project")

		assert.Nil(t, err)
		assert.Equal(t, "", result)
	})
}
