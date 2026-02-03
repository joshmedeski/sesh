package connector

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
	"github.com/joshmedeski/sesh/v2/startup"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestConfigWildcardStrategy(t *testing.T) {
	mockDir := new(dir.MockDir)
	mockHome := new(home.MockHome)
	mockLister := new(lister.MockLister)
	mockNamer := new(namer.MockNamer)
	mockStartup := new(startup.MockStartup)
	mockTmux := new(tmux.MockTmux)
	mockZoxide := new(zoxide.MockZoxide)
	mockTmuxinator := new(tmuxinator.MockTmuxinator)

	c := &RealConnector{
		model.Config{},
		mockDir,
		mockHome,
		mockLister,
		mockNamer,
		mockStartup,
		mockTmux,
		mockZoxide,
		mockTmuxinator,
	}

	t.Run("should connect via wildcard when pattern matches a directory", func(t *testing.T) {
		mockLister.On("FindConfigWildcard", "~/projects/myapp").Return(model.WildcardConfig{
			Pattern:        "~/projects/*",
			StartupCommand: "nvim",
		}, true)
		mockHome.On("ExpandHome", "~/projects/myapp").Return("/Users/test/projects/myapp", nil)
		mockDir.On("Dir", "/Users/test/projects/myapp").Return(true, "/Users/test/projects/myapp")
		mockNamer.On("Name", "/Users/test/projects/myapp").Return("myapp", nil)

		connection, err := configWildcardStrategy(c, "~/projects/myapp")
		assert.Nil(t, err)
		assert.True(t, connection.Found)
		assert.True(t, connection.New)
		assert.True(t, connection.AddToZoxide)
		assert.Equal(t, "myapp", connection.Session.Name)
		assert.Equal(t, "/Users/test/projects/myapp", connection.Session.Path)
		assert.Equal(t, "config_wildcard", connection.Session.Src)
	})

	t.Run("should return not found when no wildcard matches", func(t *testing.T) {
		mockLister.On("FindConfigWildcard", "/other/path").Return(model.WildcardConfig{}, false)

		connection, err := configWildcardStrategy(c, "/other/path")
		assert.Nil(t, err)
		assert.False(t, connection.Found)
	})

	t.Run("should return not found when path is not a directory", func(t *testing.T) {
		mockLister.On("FindConfigWildcard", "~/projects/notadir").Return(model.WildcardConfig{
			Pattern: "~/projects/*",
		}, true)
		mockHome.On("ExpandHome", "~/projects/notadir").Return("/Users/test/projects/notadir", nil)
		mockDir.On("Dir", "/Users/test/projects/notadir").Return(false, "")

		connection, err := configWildcardStrategy(c, "~/projects/notadir")
		assert.Nil(t, err)
		assert.False(t, connection.Found)
	})
}
