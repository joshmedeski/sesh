package seshcli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFormattingFlags(t *testing.T) {
	t.Run("format defaults to empty", func(t *testing.T) {
		cmd := NewListCommand(nil)

		format, err := cmd.Flags().GetString("format")

		require.NoError(t, err)
		assert.Empty(t, format)
	})

	t.Run("icons exclude defaults to empty", func(t *testing.T) {
		cmd := NewListCommand(nil)

		excluded, err := cmd.Flags().GetStringSlice("icons-exclude")

		require.NoError(t, err)
		assert.Empty(t, excluded)
	})

	t.Run("accepts format and repeated icon exclusions", func(t *testing.T) {
		cmd := NewListCommand(nil)

		require.NoError(t, cmd.Flags().Parse([]string{
			"--format",
			"{session}\x1f{active_window_name_prefix}{name}",
			"--icons-exclude",
			"tmux",
			"--icons-exclude",
			"config",
		}))

		format, err := cmd.Flags().GetString("format")
		require.NoError(t, err)
		assert.Equal(
			t,
			"{session}\x1f{active_window_name_prefix}{name}",
			format,
		)

		excluded, err := cmd.Flags().GetStringSlice("icons-exclude")
		require.NoError(t, err)
		assert.Equal(t, []string{"tmux", "config"}, excluded)
	})
}
