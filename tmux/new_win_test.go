package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestNewWindowInSession(t *testing.T) {
	const format = "#{session_name}:#{window_index}"

	t.Run("targets the session unambiguously with a trailing colon", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "new-window", "-P", "-F", format, "-n", "agent", "-c", "/home/dev/project", "-t", "project:").
			Return("project:2\n", nil)

		result, err := tmux.NewWindowInSession(model.TmuxWindowOpts{
			Name: "agent", StartDir: "/home/dev/project", TargetSession: "project",
		})

		assert.Nil(t, err)
		assert.Equal(t, "project:2", result)
	})

	t.Run("omits the target when no session is provided", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "new-window", "-P", "-F", format, "-n", "agent", "-c", "/home/dev/project").
			Return("", nil)

		result, err := tmux.NewWindowInSession(model.TmuxWindowOpts{
			Name: "agent", StartDir: "/home/dev/project",
		})

		assert.Nil(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("passes the command after -- so it becomes the pane process", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "new-window", "-P", "-F", format, "-c", "/home/dev/project", "-d", "-t", "project:", "--", `claude "do the thing"`).
			Return("project:3", nil)

		result, err := tmux.NewWindowInSession(model.TmuxWindowOpts{
			StartDir: "/home/dev/project", TargetSession: "project",
			Command: `claude "do the thing"`, Background: true,
		})

		assert.Nil(t, err)
		assert.Equal(t, "project:3", result)
	})

	t.Run("uses the configured tmux binary", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "/opt/bin/tmux"}
		mockShell.EXPECT().
			Cmd("/opt/bin/tmux", "new-window", "-P", "-F", format).
			Return("", nil)

		_, err := tmux.NewWindowInSession(model.TmuxWindowOpts{})

		assert.Nil(t, err)
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
