package connector

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
	"github.com/stretchr/testify/assert"
)

func TestDirStrategy(t *testing.T) {
	newConnector := func(l lister.Lister, n namer.Namer, h home.Home, d dir.Dir) *RealConnector {
		return &RealConnector{lister: l, namer: n, home: h, dir: d}
	}

	t.Run("reattaches to an enriched session instead of creating a new one", func(t *testing.T) {
		mockHome := new(home.MockHome)
		mockDir := new(dir.MockDir)
		mockNamer := new(namer.MockNamer)
		mockLister := new(lister.MockLister)
		mockHome.On("ExpandPath", "/p/repo").Return("/p/repo", nil)
		mockDir.On("Dir", "/p/repo").Return(true, "/p/repo")
		mockNamer.On("Name", "/p/repo").Return("repo", nil)
		mockLister.On("FindTmuxSessionByBase", "repo").
			Return(model.SeshSession{Name: "repo — fix the bug", Path: "/p/repo"}, true)

		c := newConnector(mockLister, mockNamer, mockHome, mockDir)
		conn, err := dirStrategy(c, "/p/repo")

		assert.NoError(t, err)
		assert.True(t, conn.Found)
		assert.False(t, conn.New)
		assert.Equal(t, "repo — fix the bug", conn.Session.Name)
	})

	t.Run("creates a new session when no base match exists", func(t *testing.T) {
		mockHome := new(home.MockHome)
		mockDir := new(dir.MockDir)
		mockNamer := new(namer.MockNamer)
		mockLister := new(lister.MockLister)
		mockHome.On("ExpandPath", "/p/repo").Return("/p/repo", nil)
		mockDir.On("Dir", "/p/repo").Return(true, "/p/repo")
		mockNamer.On("Name", "/p/repo").Return("repo", nil)
		mockLister.On("FindTmuxSessionByBase", "repo").Return(model.SeshSession{}, false)

		c := newConnector(mockLister, mockNamer, mockHome, mockDir)
		conn, err := dirStrategy(c, "/p/repo")

		assert.NoError(t, err)
		assert.True(t, conn.Found)
		assert.True(t, conn.New)
		assert.Equal(t, "repo", conn.Session.Name)
	})

	t.Run("not a directory returns not found", func(t *testing.T) {
		mockHome := new(home.MockHome)
		mockDir := new(dir.MockDir)
		mockHome.On("ExpandPath", "/p/x").Return("/p/x", nil)
		mockDir.On("Dir", "/p/x").Return(false, "")

		c := newConnector(nil, nil, mockHome, mockDir)
		conn, err := dirStrategy(c, "/p/x")

		assert.NoError(t, err)
		assert.False(t, conn.Found)
	})
}
