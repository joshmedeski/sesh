package connector

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/focuser"
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

	// New session path
	mTmux.EXPECT().NewSession("nutiliti/2345", "/repo/w/2345").Return("", nil)
	// No opts.Command => startup.Exec is called; startup is a no-op mock here
	mStartup := startup.NewMockStartup(t)
	mStartup.EXPECT().Exec(mock.Anything).Return("", nil).Maybe()

	// Detached: not attached, --switch set
	mTmux.EXPECT().IsAttached().Return(false)
	mTmux.EXPECT().ListClients().Return([]model.TmuxClient{
		{Name: "/dev/ttys002", TTY: "/dev/ttys002", SessionID: "$2"},
	}, nil)
	mTmux.EXPECT().ResolveClient().Return("/dev/ttys002")
	mTmux.EXPECT().SwitchClientTarget("/dev/ttys002", "nutiliti/2345").Return("", nil)
	mFocuser.EXPECT().Activate("wezterm").Return(true, nil)

	c := NewConnector(
		model.Config{Terminal: "wezterm"},
		nil, nil, nil, nil, mStartup, mTmux, nil, nil, mFocuser,
	).(*RealConnector)

	conn := model.Connection{
		Found: true, New: true,
		Session: model.SeshSession{Src: "dir", Name: "nutiliti/2345", Path: "/repo/w/2345"},
	}
	_, err := connectToTmux(c, conn, model.ConnectOpts{Switch: true})
	require.NoError(t, err)
}

func TestConnectToTmuxDetachedWithNoClientStillFocuses(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	mFocuser := focuser.NewMockFocuser(t)

	// Nothing is attached to the server, so there is no client to switch.
	mTmux.EXPECT().IsAttached().Return(false)
	mTmux.EXPECT().ListClients().Return(nil, nil)
	mFocuser.EXPECT().Activate("wezterm").Return(true, nil)

	c := NewConnector(
		model.Config{Terminal: "wezterm"},
		nil, nil, nil, nil, nil, mTmux, nil, nil, mFocuser,
	).(*RealConnector)

	conn := model.Connection{
		Found: true, New: false,
		Session: model.SeshSession{Src: "tmux", Name: "nutiliti/2345", Path: "/repo/w/2345"},
	}
	msg, err := connectToTmux(c, conn, model.ConnectOpts{Switch: true})
	require.NoError(t, err)
	assert.Equal(t, "connected to tmux session: nutiliti/2345", msg)
}

// A GUI launcher runs sesh with a bare PATH, so tmux isn't reachable while
// Activate's osascript still is. That used to focus the terminal, change
// nothing, and report success.
func TestConnectToTmuxDetachedReportsUnreachableTmux(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	mFocuser := focuser.NewMockFocuser(t)

	mTmux.EXPECT().IsAttached().Return(false)
	mTmux.EXPECT().ListClients().
		Return(nil, errors.New(`exec: "tmux": executable file not found in $PATH`))

	c := NewConnector(
		model.Config{Terminal: "wezterm"},
		nil, nil, nil, nil, nil, mTmux, nil, nil, mFocuser,
	).(*RealConnector)

	conn := model.Connection{
		Found: true, New: false,
		Session: model.SeshSession{Src: "tmux", Name: "nutiliti/2345", Path: "/repo/w/2345"},
	}
	_, err := connectToTmux(c, conn, model.ConnectOpts{Switch: true})
	require.ErrorContains(t, err, "couldn't reach tmux to switch to 'nutiliti/2345'")
}
