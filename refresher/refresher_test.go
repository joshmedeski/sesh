package refresher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefreshArgs(t *testing.T) {
	assert.Equal(t, []string{"status", "--refresh"}, refreshArgs(""))
	assert.Equal(t, []string{"status", "--refresh", "/repo"}, refreshArgs("/repo"))
}
