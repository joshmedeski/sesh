package connector

import (
	"testing"

	"github.com/joshmedeski/sesh/dir"
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestEstablishTmuxConnection(t *testing.T) {
	mockDir := new(dir.MockDir)
	mockHome := new(home.MockHome)
	mockLister := new(lister.MockLister)
	mockTmux := new(tmux.MockTmux)
	mockZoxide := new(zoxide.MockZoxide)

	c := &RealConnector{
		mockDir,
		mockHome,
		mockLister,
		mockTmux,
		mockZoxide,
		model.Config{},
	}
	mockTmux.On("AttachSession", mock.Anything).Return("attaching", nil)
	mockZoxide.On("Add", mock.Anything).Return(nil)

	t.Run("should attach to tmux session", func(t *testing.T) {
		mockTmux.On("IsAttached").Return(false)
		mockLister.On("FindTmuxSession", "dotfiles").Return(model.SeshSession{
			Name: "dotfiles",
			Path: "/Users/joshmedeski/c/dotfiles",
		}, true)
		connection, err := establishTmuxConnection(c, "dotfiles", model.ConnectOpts{})
		assert.Equal(t, nil, err)
		assert.Equal(t, "attaching to existing tmux session: dotfiles", connection)
	})

	t.Run("should switch to tmux session", func(t *testing.T) {
		mockTmux.On("IsAttached").Return(true)
		mockLister.On("FindTmuxSession", "dotfiles").Return(model.SeshSession{
			Name: "dotfiles",
			Path: "/Users/joshmedeski/c/dotfiles",
		}, true)
		connection, err := establishTmuxConnection(c, "dotfiles", model.ConnectOpts{})
		assert.Equal(t, nil, err)
		assert.Equal(t, "switching to existing tmux session: dotfiles", connection)
	})
}
