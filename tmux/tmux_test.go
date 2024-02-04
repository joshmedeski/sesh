package tmux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommand_Run(t *testing.T) {
	c := Command{
		execFunc: func(string, []string) (string, error) {
			return "stub", nil
		},
	}

	res, err := c.Run([]string{"arg1", "arg2"})
	require.NoError(t, err)
	require.Equal(t, "stub", res)
}
