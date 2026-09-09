package connector

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/focuser"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A bare mux.Mux offers neither pane navigation nor per-client switching. These
// tests pin down how the connector degrades on such a backend, which is the
// contract a non-tmux multiplexer is written against.

func TestConnectDetachedWithoutClientSwitcherFallsBackToSwitchClient(t *testing.T) {
	m := mux.NewMockMux(t)
	mFocuser := focuser.NewMockFocuser(t)
	m.EXPECT().SwitchClient("dotfiles").Return("", nil)
	mFocuser.EXPECT().Activate("wezterm").Return(true, nil)

	c := NewConnector(
		model.Config{Terminal: "wezterm"}, nil, nil, nil, nil, nil, m, nil, nil, mFocuser,
	).(*RealConnector)

	out, err := c.connectDetached("dotfiles")

	require.NoError(t, err)
	assert.Equal(t, "switching to session: dotfiles", out)
}

func TestConnectDetachedWithoutClientSwitcherStillFocusesOnFailure(t *testing.T) {
	m := mux.NewMockMux(t)
	mFocuser := focuser.NewMockFocuser(t)
	m.EXPECT().SwitchClient("dotfiles").Return("", assert.AnError)
	mFocuser.EXPECT().Activate("wezterm").Return(true, nil)

	c := NewConnector(
		model.Config{Terminal: "wezterm"}, nil, nil, nil, nil, nil, m, nil, nil, mFocuser,
	).(*RealConnector)

	_, err := c.connectDetached("dotfiles")

	require.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestConnectToPaneWithoutPaneNavigatorErrors(t *testing.T) {
	m := mux.NewMockMux(t)
	c := NewConnector(
		model.Config{}, nil, nil, lister.NewMockLister(t), nil, nil, m, nil, nil, nil,
	).(*RealConnector)

	_, err := connectToTmuxPane(c, model.Connection{
		Session: model.SeshSession{Src: "tmux-pane", Name: "editor/nvim"},
	}, model.ConnectOpts{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "editor/nvim")
}
