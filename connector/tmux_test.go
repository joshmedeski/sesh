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
	mock "github.com/stretchr/testify/mock"
)

func TestEstablishTmuxConnection(t *testing.T) {
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
	mockTmux.On("AttachSession", mock.Anything).Return("attaching", nil)
	mockZoxide.On("Add", mock.Anything).Return(nil)

	t.Run("should attach to tmux session", func(t *testing.T) {
		mockTmux.On("IsAttached").Return(false)
		mockLister.On("FindTmuxSession", "dotfiles").Return(model.SeshSession{
			Name: "dotfiles",
			Path: "/Users/joshmedeski/c/dotfiles",
		}, true)
		connection, err := tmuxStrategy(c, "dotfiles")
		assert.Nil(t, err)
		assert.Equal(t, "dotfiles", connection.Session.Name)
	})

	t.Run("should switch to tmux session", func(t *testing.T) {
		mockTmux.On("IsAttached").Return(true)
		mockLister.On("FindTmuxSession", "dotfiles").Return(model.SeshSession{
			Name: "dotfiles",
			Path: "/Users/joshmedeski/c/dotfiles",
		}, true)
		connection, err := tmuxStrategy(c, "dotfiles")
		assert.Nil(t, err)
		assert.Equal(t, "dotfiles", connection.Session.Name)
	})
}
