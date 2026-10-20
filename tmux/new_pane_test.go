package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestNewPaneInWindow(t *testing.T) {
	t.Run("splits detached, runs the command in the new pane and applies the layout", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "split-window", "-d", "-P", "-F", "#{pane_id}", "-t", "app:dev", "-c", "/repo/app").
			Return("%5\n", nil)
		mockShell.EXPECT().Cmd("tmux", "send-keys", "-t", "%5", "lazydocker", "Enter").Return("", nil)
		mockShell.EXPECT().Cmd("tmux", "select-layout", "-t", "app:dev", "main-vertical").Return("", nil)

		paneID, err := tmux.NewPaneInWindow(model.TmuxPaneOpts{
			TargetWindow: "app:dev", StartDir: "/repo/app", Command: "lazydocker", Layout: "main-vertical",
		})

		assert.NoError(t, err)
		assert.Equal(t, "%5", paneID)
	})

	t.Run("only splits when there is no command or layout", func(t *testing.T) {
		mockShell := shell.NewMockShell(t)
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "split-window", "-d", "-P", "-F", "#{pane_id}", "-t", "app:").
			Return("%6", nil)

		paneID, err := tmux.NewPaneInWindow(model.TmuxPaneOpts{TargetWindow: "app"})

		assert.NoError(t, err)
		assert.Equal(t, "%6", paneID)
	})
}

func TestSelectAndZoomPane(t *testing.T) {
	mockShell := shell.NewMockShell(t)
	tmux := &RealTmux{shell: mockShell, bin: "tmux"}
	mockShell.EXPECT().Cmd("tmux", "select-pane", "-t", "%5").Return("", nil)
	mockShell.EXPECT().Cmd("tmux", "resize-pane", "-Z", "-t", "%5").Return("", nil)

	_, err := tmux.SelectPaneByID("%5")
	assert.NoError(t, err)
	_, err = tmux.ZoomPane("%5")
	assert.NoError(t, err)
}
