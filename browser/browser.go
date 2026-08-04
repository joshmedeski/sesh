package browser

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

// Browser reads the active tab URL from a browser's front window. Currently
// macOS-only via osascript; a no-op elsewhere.
type Browser interface {
	// ActiveTabURL returns the front window's active-tab URL.
	// Returns ("", false, nil) when skipped: non-macOS or no application
	// configured.
	ActiveTabURL() (url string, ok bool, err error)
}

// defaultURLCommand is the Chrome-family AppleScript fragment (Helium, Chrome,
// Arc, Brave, Edge). Safari uses "URL of current tab of front window".
const defaultURLCommand = "URL of active tab of front window"

type RealBrowser struct {
	runtime runtimewrap.Runtime
	shell   shell.Shell
	config  model.BrowserConfig
}

func NewBrowser(runtime runtimewrap.Runtime, shell shell.Shell, config model.BrowserConfig) Browser {
	return &RealBrowser{runtime, shell, config}
}

func (b *RealBrowser) ActiveTabURL() (string, bool, error) {
	if b.config.Application == "" {
		return "", false, nil
	}
	if b.runtime.GOOS() != "darwin" {
		return "", false, nil
	}
	urlCommand := b.config.URLCommand
	if urlCommand == "" {
		urlCommand = defaultURLCommand
	}
	script := fmt.Sprintf("tell application %q to return %s", b.config.Application, urlCommand)
	url, err := b.shell.Cmd("osascript", "-e", script)
	if err != nil {
		return "", false, err
	}
	return url, true, nil
}
