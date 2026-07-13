package connector

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
	"github.com/stretchr/testify/assert"
)

func TestZoxideStrategy(t *testing.T) {
	t.Run("reattaches to an enriched session", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockNamer := new(namer.MockNamer)
		mockLister.On("FindZoxideSession", "repo").
			Return(model.SeshSession{Name: "repo", Path: "/p/repo"}, true)
		mockNamer.On("Name", "/p/repo").Return("repo", nil)
		mockLister.On("FindTmuxSessionByBase", "repo").
			Return(model.SeshSession{Name: "repo — fix the bug", Path: "/p/repo"}, true)

		c := &RealConnector{lister: mockLister, namer: mockNamer}
		conn, err := zoxideStrategy(c, "repo")

		assert.NoError(t, err)
		assert.True(t, conn.Found)
		assert.False(t, conn.New)
		assert.Equal(t, "repo — fix the bug", conn.Session.Name)
	})

	t.Run("creates a new session when no base match", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockNamer := new(namer.MockNamer)
		mockLister.On("FindZoxideSession", "repo").
			Return(model.SeshSession{Name: "repo", Path: "/p/repo"}, true)
		mockNamer.On("Name", "/p/repo").Return("repo", nil)
		mockLister.On("FindTmuxSessionByBase", "repo").Return(model.SeshSession{}, false)

		c := &RealConnector{lister: mockLister, namer: mockNamer}
		conn, err := zoxideStrategy(c, "repo")

		assert.NoError(t, err)
		assert.True(t, conn.Found)
		assert.True(t, conn.New)
		assert.Equal(t, "repo", conn.Session.Name)
	})

	t.Run("zoxide miss returns not found", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockLister.On("FindZoxideSession", "nope").Return(model.SeshSession{}, false)

		c := &RealConnector{lister: mockLister}
		conn, err := zoxideStrategy(c, "nope")

		assert.NoError(t, err)
		assert.False(t, conn.Found)
	})
}
