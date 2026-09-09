package browser

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActiveTabURLOnMacOSDefaultCommand(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("darwin")
	s := shell.NewMockShell(t)
	s.EXPECT().
		Cmd("osascript", "-e", `tell application "Helium" to return URL of active tab of front window`).
		Return("https://github.com/joshmedeski/sesh/issues/409", nil)

	b := NewBrowser(rt, s, model.BrowserConfig{Application: "Helium"})
	url, ok, err := b.ActiveTabURL()
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "https://github.com/joshmedeski/sesh/issues/409", url)
}

func TestActiveTabURLWithURLCommandOverride(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("darwin")
	s := shell.NewMockShell(t)
	s.EXPECT().
		Cmd("osascript", "-e", `tell application "Safari" to return URL of current tab of front window`).
		Return("https://github.com/joshmedeski/sesh/pull/410", nil)

	b := NewBrowser(rt, s, model.BrowserConfig{
		Application: "Safari",
		URLCommand:  "URL of current tab of front window",
	})
	url, ok, err := b.ActiveTabURL()
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "https://github.com/joshmedeski/sesh/pull/410", url)
}

func TestActiveTabURLSkippedOnLinux(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("linux")
	s := shell.NewMockShell(t) // no shell calls expected
	b := NewBrowser(rt, s, model.BrowserConfig{Application: "Helium"})
	url, ok, err := b.ActiveTabURL()
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, "", url)
}

func TestActiveTabURLSkippedWhenNoApplication(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	s := shell.NewMockShell(t) // no runtime/shell calls expected
	b := NewBrowser(rt, s, model.BrowserConfig{})
	url, ok, err := b.ActiveTabURL()
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, "", url)
}

func TestActiveTabURLPropagatesShellError(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("darwin")
	s := shell.NewMockShell(t)
	s.EXPECT().
		Cmd("osascript", "-e", `tell application "Helium" to return URL of active tab of front window`).
		Return("", errors.New("no window"))

	b := NewBrowser(rt, s, model.BrowserConfig{Application: "Helium"})
	_, ok, err := b.ActiveTabURL()
	require.Error(t, err)
	assert.False(t, ok)
}
