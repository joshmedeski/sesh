package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListClients(t *testing.T) {
	s := shell.NewMockShell(t)
	// ListCmd splits on newlines, so the trailing blank line comes through as
	// an empty entry and has to be dropped.
	s.EXPECT().ListCmd("tmux", "list-clients", "-F", clientFormat).
		Return([]string{
			"/dev/ttys001\t/dev/ttys001\t$1\t1788188899",
			"/dev/ttys002\t/dev/ttys002\t$2\t1788188999",
			"",
		}, nil)
	os := oswrap.NewMockOs(t)
	tm := NewTmux(os, s, "")
	clients, err := tm.ListClients()
	require.NoError(t, err)
	assert.Equal(t, []model.TmuxClient{
		{Name: "/dev/ttys001", TTY: "/dev/ttys001", SessionID: "$1", Activity: 1788188899},
		{Name: "/dev/ttys002", TTY: "/dev/ttys002", SessionID: "$2", Activity: 1788188999},
	}, clients)
}

func TestSwitchClientTarget(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().Cmd("tmux", "switch-client", "-c", "/dev/ttys001", "-t", "sesh").
		Return("", nil)
	os := oswrap.NewMockOs(t)
	tm := NewTmux(os, s, "")
	_, err := tm.SwitchClientTarget("/dev/ttys001", "sesh")
	require.NoError(t, err)
}
