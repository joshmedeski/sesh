package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestSelectLayout(t *testing.T) {
	mockShell := shell.NewMockShell(t)
	tmux := &RealTmux{shell: mockShell, bin: "tmux"}
	mockShell.EXPECT().Cmd("tmux", "select-layout", "-t", "app:", "main-vertical").Return("", nil)

	_, err := tmux.SelectLayout("app", "main-vertical")

	assert.NoError(t, err)
}
