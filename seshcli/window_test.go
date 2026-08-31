package seshcli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTargetFlag(t *testing.T) {
	t.Run("reads --target", func(t *testing.T) {
		cmd := NewWindowCommand(nil)
		require.NoError(t, cmd.Flags().Parse([]string{"--target", "brain"}))
		assert.Equal(t, "brain", targetFlag(cmd))
	})

	t.Run("still accepts the deprecated --session alias", func(t *testing.T) {
		cmd := NewWindowCommand(nil)
		cmd.SetErr(&bytes.Buffer{})
		require.NoError(t, cmd.Flags().Parse([]string{"--session", "brain"}))
		assert.Equal(t, "brain", targetFlag(cmd))
	})

	t.Run("prefers --target when both are given", func(t *testing.T) {
		cmd := NewWindowCommand(nil)
		cmd.SetErr(&bytes.Buffer{})
		require.NoError(t, cmd.Flags().Parse([]string{"--session", "old", "--target", "new"}))
		assert.Equal(t, "new", targetFlag(cmd))
	})

	t.Run("is empty when neither is given", func(t *testing.T) {
		cmd := NewWindowCommand(nil)
		require.NoError(t, cmd.Flags().Parse(nil))
		assert.Equal(t, "", targetFlag(cmd))
	})
}

// -s used to mean --session on `sesh window`. It now means --switch everywhere,
// so the parent command deliberately has no -s shorthand: the old form has to
// fail loudly rather than quietly connect to a window named like the session.
func TestParentWindowCommandRejectsSessionShorthand(t *testing.T) {
	cmd := NewWindowCommand(nil)
	err := cmd.Flags().Parse([]string{"-s", "brain"})
	assert.ErrorContains(t, err, "unknown shorthand flag: 's'")
}

func TestWindowConnectAcceptsSwitchShorthand(t *testing.T) {
	cmd := NewWindowCommand(nil)
	connect, _, err := cmd.Find([]string{"connect"})
	require.NoError(t, err)
	require.NoError(t, connect.Flags().Parse([]string{"-s"}))

	switched, err := connect.Flags().GetBool("switch")
	require.NoError(t, err)
	assert.True(t, switched)
}

// `sesh window <name>` predates the connect subcommand and has to keep working.
func TestWindowArgsStayWithTheParentCommand(t *testing.T) {
	cmd := NewWindowCommand(nil)

	found, args, err := cmd.Find([]string{"notes"})
	require.NoError(t, err)
	assert.Equal(t, "window", found.Name())
	assert.Equal(t, []string{"notes"}, args)

	found, _, err = cmd.Find([]string{"connect", "notes"})
	require.NoError(t, err)
	assert.Equal(t, "connect", found.Name())
}
