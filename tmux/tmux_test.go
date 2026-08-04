package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListClients(t *testing.T) {
	s := shell.NewMockShell(t)
	s.EXPECT().ListCmd("tmux", "list-clients", "-F", "#{client_name}").
		Return([]string{"/dev/ttys001", ""}, nil)
	os := oswrap.NewMockOs(t)
	tm := NewTmux(os, s, "")
	clients, err := tm.ListClients()
	require.NoError(t, err)
	assert.Equal(t, []string{"/dev/ttys001", ""}, clients)
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
