package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

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
