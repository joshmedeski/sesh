package tmux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	c, err := NewCommand(Options{})
	require.NoError(t, err)
	require.NotNil(t, c)
}
