package focuser

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivateOnMacOS(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("darwin")
	s := shell.NewMockShell(t)
	s.EXPECT().Cmd("osascript", "-e", `tell application "wezterm" to activate`).Return("", nil)

	f := NewFocuser(rt, s)
	ran, err := f.Activate("wezterm")
	require.NoError(t, err)
	assert.True(t, ran)
}

func TestActivateSkippedOnLinux(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("linux")
	s := shell.NewMockShell(t) // no shell calls expected
	f := NewFocuser(rt, s)
	ran, err := f.Activate("wezterm")
	require.NoError(t, err)
	assert.False(t, ran)
}

func TestActivateSkippedWhenEmpty(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	s := shell.NewMockShell(t)
	f := NewFocuser(rt, s)
	ran, err := f.Activate("")
	require.NoError(t, err)
	assert.False(t, ran)
}
