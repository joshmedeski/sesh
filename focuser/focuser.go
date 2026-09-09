package focuser

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

// Focuser brings a terminal emulator to the foreground. Currently macOS-only
// via osascript; a no-op elsewhere.
type Focuser interface {
	// Activate focuses the given app. Returns (true, nil) if it ran,
	// (false, nil) if skipped (empty app or non-macOS).
	Activate(app string) (bool, error)
}

type RealFocuser struct {
	runtime runtimewrap.Runtime
	shell   shell.Shell
}

func NewFocuser(runtime runtimewrap.Runtime, shell shell.Shell) Focuser {
	return &RealFocuser{runtime, shell}
}

func (f *RealFocuser) Activate(app string) (bool, error) {
	if app == "" || f.runtime.GOOS() != "darwin" {
		return false, nil
	}
	script := fmt.Sprintf("tell application %q to activate", app)
	if _, err := f.shell.Cmd("osascript", "-e", script); err != nil {
		return false, err
	}
	return true, nil
}
