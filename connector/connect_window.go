package connector

import (
	"fmt"
	"path/filepath"

	"github.com/joshmedeski/sesh/v2/model"
)

// ConnectWindow selects the named window in the target session, creating it
// when it isn't there. It returns the window's "session:index" target.
//
// An unnamed window has nothing to match on, so it is always created. So is one
// asked for with opts.New, which allows windows to share a name; a later
// connect to that name reuses the lowest-indexed one.
func (c *RealConnector) ConnectWindow(opts model.WindowConnectOpts) (string, error) {
	session, err := c.targetSession(opts.Session)
	if err != nil {
		return "", err
	}

	name, startDir := opts.Name, opts.Path
	if startDir == "" {
		// A directory argument names the window after its basename and roots
		// it there, so `sesh window ~/c/sesh` opens a window called sesh. The
		// name is derived before matching, so running it twice selects the
		// window the first run created instead of opening a duplicate.
		if dirName, dirPath, ok := c.windowDir(name); ok {
			name, startDir = dirName, dirPath
		} else {
			startDir = session.Path
		}
	}

	if name != "" && !opts.New {
		windows, err := c.tmux.ListWindows(session.Name)
		if err != nil {
			return "", fmt.Errorf("failed to list windows in '%s': %w", session.Name, err)
		}
		for _, window := range windows {
			if window.Name != name {
				continue
			}
			// Target by index: a name shared with a session resolves
			// ambiguously (see #280), an index never does.
			target := fmt.Sprintf("%s:%d", session.Name, window.Index)
			return c.reuseWindow(session.Name, target, opts)
		}
	}

	target, err := c.tmux.NewWindowInSession(model.TmuxWindowOpts{
		Name:          name,
		StartDir:      startDir,
		TargetSession: session.Name,
		Command:       opts.Command,
		Background:    opts.Background,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create window: %w", err)
	}
	if opts.Background {
		return target, nil
	}
	// new-window already made it the session's active window.
	if _, err := c.focusSession(session.Name, opts.Switch); err != nil {
		return "", err
	}
	return target, nil
}

// reuseWindow drops a command into a window that already exists. The window
// already has a running process, so the command goes in as keystrokes rather
// than becoming the pane's process the way it does on create.
func (c *RealConnector) reuseWindow(sessionName string, target string, opts model.WindowConnectOpts) (string, error) {
	if opts.Command != "" {
		if _, err := c.tmux.SendKeys(target, opts.Command); err != nil {
			return "", fmt.Errorf("failed to send command to '%s': %w", target, err)
		}
	}
	if opts.Background {
		return target, nil
	}
	if _, err := c.tmux.SelectWindow(target); err != nil {
		return "", fmt.Errorf("failed to select window '%s': %w", target, err)
	}
	if _, err := c.focusSession(sessionName, opts.Switch); err != nil {
		return "", err
	}
	return target, nil
}

// windowDir reads a window argument as a directory, reporting the window name
// and start directory it implies. Anything that isn't a directory is left to be
// used as a plain window name.
func (c *RealConnector) windowDir(name string) (string, string, bool) {
	if name == "" {
		return "", "", false
	}
	expanded, err := c.home.ExpandPath(name)
	if err != nil {
		return "", "", false
	}
	isDir, absPath := c.dir.Dir(expanded)
	if !isDir {
		return "", "", false
	}
	return filepath.Base(absPath), absPath, true
}
