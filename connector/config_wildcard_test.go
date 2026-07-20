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
		mockHome.On("ExpandPath", "~/projects/myapp").Return("/Users/test/projects/myapp", nil)
		mockDir.On("Dir", "/Users/test/projects/myapp").Return(true, "/Users/test/projects/myapp")
		mockLister.On("FindConfigWildcard", "/Users/test/projects/myapp").Return(model.WildcardConfig{
			Pattern:        "~/projects/*",
			StartupCommand: "nvim",
			Windows:        []string{"code", "server"},
		}, true)
		mockNamer.On("Name", "/Users/test/projects/myapp").Return("myapp", nil)

		connection, err := configWildcardStrategy(c, "~/projects/myapp")
		assert.Nil(t, err)
		assert.True(t, connection.Found)
		assert.True(t, connection.New)
		assert.True(t, connection.AddToZoxide)
		assert.Equal(t, "myapp", connection.Session.Name)
		assert.Equal(t, "/Users/test/projects/myapp", connection.Session.Path)
		assert.Equal(t, "config_wildcard", connection.Session.Src)
		assert.Equal(t, []string{"code", "server"}, connection.Session.WindowNames)
	})

	t.Run("should match wildcard using absolute path for relative connect arg", func(t *testing.T) {
		// user runs `sesh connect dev/my-app-dir` from ~ (relative arg)
		mockHome.On("ExpandPath", "dev/my-app-dir").Return("dev/my-app-dir", nil)
		mockDir.On("Dir", "dev/my-app-dir").Return(true, "/Users/test/dev/my-app-dir")
		mockLister.On("FindConfigWildcard", "/Users/test/dev/my-app-dir").Return(model.WildcardConfig{
			Pattern: "~/dev/*",
			Windows: []string{"agent"},
		}, true)
		mockNamer.On("Name", "/Users/test/dev/my-app-dir").Return("my-app-dir", nil)

		connection, err := configWildcardStrategy(c, "dev/my-app-dir")
		assert.Nil(t, err)
		assert.True(t, connection.Found)
		assert.Equal(t, "/Users/test/dev/my-app-dir", connection.Session.Path)
		assert.Equal(t, []string{"agent"}, connection.Session.WindowNames)
	})

	t.Run("should propagate DisableStartupCommand from wildcard config", func(t *testing.T) {
		mockHome.On("ExpandPath", "~/projects/quiet").Return("/Users/test/projects/quiet", nil)
		mockDir.On("Dir", "/Users/test/projects/quiet").Return(true, "/Users/test/projects/quiet")
		mockLister.On("FindConfigWildcard", "/Users/test/projects/quiet").Return(model.WildcardConfig{
			Pattern:             "~/projects/*",
			DisableStartCommand: true,
		}, true)
		mockNamer.On("Name", "/Users/test/projects/quiet").Return("quiet", nil)

		connection, err := configWildcardStrategy(c, "~/projects/quiet")
		assert.Nil(t, err)
		assert.True(t, connection.Found)
		assert.True(t, connection.Session.DisableStartupCommand)
	})

	t.Run("should return not found when no wildcard matches", func(t *testing.T) {
		mockHome.On("ExpandPath", "/other/path").Return("/other/path", nil)
		mockDir.On("Dir", "/other/path").Return(true, "/other/path")
		mockLister.On("FindConfigWildcard", "/other/path").Return(model.WildcardConfig{}, false)

		connection, err := configWildcardStrategy(c, "/other/path")
		assert.Nil(t, err)
		assert.False(t, connection.Found)
	})

	t.Run("should return not found when path is not a directory", func(t *testing.T) {
		mockLister.On("FindConfigWildcard", "~/projects/notadir").Return(model.WildcardConfig{
			Pattern: "~/projects/*",
		}, true)
		mockHome.On("ExpandPath", "~/projects/notadir").Return("/Users/test/projects/notadir", nil)
		mockDir.On("Dir", "/Users/test/projects/notadir").Return(false, "")

		connection, err := configWildcardStrategy(c, "~/projects/notadir")
		assert.Nil(t, err)
		assert.False(t, connection.Found)
	})
}
