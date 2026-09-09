package seshcli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func windowSubcommand(t *testing.T, name string) *cobra.Command {
	t.Helper()
	cmd, _, err := NewWindowCommand(nil).Find([]string{name})
	require.NoError(t, err)
	require.Equal(t, name, cmd.Name())
	return cmd
}

func TestTargetFlag(t *testing.T) {
	t.Run("reads --target", func(t *testing.T) {
		cmd := windowSubcommand(t, "connect")
		require.NoError(t, cmd.Flags().Parse([]string{"--target", "brain"}))
		assert.Equal(t, "brain", targetFlag(cmd))
	})

	t.Run("still accepts the deprecated --session alias", func(t *testing.T) {
		cmd := windowSubcommand(t, "list")
		cmd.SetErr(&bytes.Buffer{})
		require.NoError(t, cmd.Flags().Parse([]string{"--session", "brain"}))
		assert.Equal(t, "brain", targetFlag(cmd))
	})

	t.Run("prefers --target when both are given", func(t *testing.T) {
		cmd := windowSubcommand(t, "connect")
		cmd.SetErr(&bytes.Buffer{})
		require.NoError(t, cmd.Flags().Parse([]string{"--session", "old", "--target", "new"}))
		assert.Equal(t, "new", targetFlag(cmd))
	})

	t.Run("is empty when neither is given", func(t *testing.T) {
		cmd := windowSubcommand(t, "connect")
		require.NoError(t, cmd.Flags().Parse(nil))
		assert.Equal(t, "", targetFlag(cmd))
	})
}

// `sesh window <name>` used to select-or-create that window. It is a
// subcommand entrypoint now, so the old form has to fail with a message that
// names its replacement rather than a bare "unknown command".
func TestWindowRejectsTheLegacyPositionalForm(t *testing.T) {
	cmd := NewWindowCommand(nil)

	err := cmd.RunE(cmd, []string{"claude"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), `unknown command "claude"`)
	assert.Contains(t, err.Error(), "sesh window connect claude")
}

func TestWindowWithNoArgsPrintsHelp(t *testing.T) {
	cmd := NewWindowCommand(nil)
	out := &bytes.Buffer{}
	cmd.SetOut(out)

	require.NoError(t, cmd.RunE(cmd, nil))
	assert.Contains(t, out.String(), "connect")
	assert.Contains(t, out.String(), "list")
}

// -s used to mean --session on `sesh window`. It means --switch everywhere now,
// and the entrypoint has no flags of its own for the old form to land on.
func TestParentWindowCommandRejectsSessionShorthand(t *testing.T) {
	cmd := NewWindowCommand(nil)
	err := cmd.Flags().Parse([]string{"-s", "brain"})
	assert.ErrorContains(t, err, "unknown shorthand flag: 's'")
}

func TestWindowConnectAcceptsSwitchShorthand(t *testing.T) {
	cmd := windowSubcommand(t, "connect")
	require.NoError(t, cmd.Flags().Parse([]string{"-s"}))

	switched, err := cmd.Flags().GetBool("switch")
	require.NoError(t, err)
	assert.True(t, switched)
}

func TestWindowSubcommandRouting(t *testing.T) {
	cmd := NewWindowCommand(nil)

	for _, tc := range []struct{ arg, want string }{
		{"connect", "connect"},
		{"cn", "connect"},
		{"list", "list"},
		{"ls", "list"},
	} {
		found, _, err := cmd.Find([]string{tc.arg})
		require.NoError(t, err)
		assert.Equal(t, tc.want, found.Name(), "routing %q", tc.arg)
	}
}
