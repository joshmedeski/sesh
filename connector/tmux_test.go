package connector

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/focuser"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/startup"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
		nil,
		nil,
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

func TestConnectToTmuxDetachedSwitchesClientAndFocuses(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	mFocuser := focuser.NewMockFocuser(t)
	mRuntime := runtimewrap.NewMockRuntime(t)

	// New session path
	mTmux.EXPECT().NewSession("nutiliti/2345", "/repo/w/2345").Return("", nil)
	// No opts.Command => startup.Exec is called; startup is a no-op mock here
	mStartup := startup.NewMockStartup(t)
	mStartup.EXPECT().Exec(mock.Anything).Return("", nil).Maybe()

	// Detached: not attached, --switch set
	mTmux.EXPECT().IsAttached().Return(false)
	mTmux.EXPECT().ListClients().Return([]string{"", "/dev/ttys002"}, nil)
	mTmux.EXPECT().SwitchClientTarget("/dev/ttys002", "nutiliti/2345").Return("", nil)
	mFocuser.EXPECT().Activate("wezterm").Return(true, nil)

	c := NewConnector(
		model.Config{Terminal: "wezterm"},
		nil, nil, nil, nil, mStartup, mTmux, nil, nil, mRuntime, mFocuser,
	).(*RealConnector)

	conn := model.Connection{
		Found: true, New: true,
		Session: model.SeshSession{Src: "dir", Name: "nutiliti/2345", Path: "/repo/w/2345"},
	}
	_, err := connectToTmux(c, conn, model.ConnectOpts{Switch: true})
	require.NoError(t, err)
}
