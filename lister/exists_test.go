package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	sessions := map[string]model.SeshSession{
		"session1": {},
	}
	_, session1Exists := exists("session1", sessions)
	assert.Equal(t, true, session1Exists)
	_, session3Exists := exists("session3", sessions)
	assert.Equal(t, false, session3Exists)
}
