package tmuxinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisabled(t *testing.T) {
	d := NewDisabled()

	configs, err := d.List()
	require.NoError(t, err)
	assert.Empty(t, configs)

	_, err = d.Start("anything")
	assert.Error(t, err)
}
