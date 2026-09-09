// Package mux is the seam between sesh and the terminal multiplexer it drives.
//
// Everything that lists, creates, focuses or types into sessions goes through
// Mux rather than a concrete tmux client, so a second backend (herdr) can slot
// in behind the connector, lister, startup and previewer without those
// packages knowing which multiplexer is running. tmux is the only
// implementation today; its own interface is a superset of these.
//
// Operations that only make sense on one multiplexer live on the optional
// capability interfaces below. Callers type-assert for them and degrade when
// the backend doesn't offer one, instead of every backend growing stubs.
package mux

import "github.com/joshmedeski/sesh/v2/model"

// Mux is what every backend must provide. Names and targets are the
// multiplexer's own identifiers (a tmux session name, a herdr workspace ID);
// sesh passes back what the backend handed out.
type Mux interface {
	// Name identifies the backend ("tmux", "herdr"). Sources and strategies
	// that only make sense on one multiplexer check it.
	Name() string

	// IsAttached reports whether sesh is running inside the multiplexer, which
	// decides between switching the current client and attaching a new one.
	IsAttached() bool

	ListSessions() ([]*model.TmuxSession, error)
	NewSession(sessionName string, startDir string) (string, error)
	RenameSession(target string, newName string) (string, error)
	// SwitchOrAttach brings the user to a session, switching the current client
	// when inside the multiplexer or when opts.Switch asks for it, attaching
	// otherwise.
	SwitchOrAttach(name string, opts model.ConnectOpts) (string, error)
	SwitchClient(targetSession string) (string, error)

	// SendKeys types command into the target followed by Enter.
	SendKeys(target string, command string) (string, error)
	// CapturePane returns the visible contents of the session's active pane,
	// with ANSI escapes, for previews.
	CapturePane(targetSession string) (string, error)

	ListWindows(targetSession string) ([]*model.TmuxWindow, error)
	// ListAllWindowNames returns window names keyed by session name, in one
	// call. format is a tmux format string; backends without a format language
	// ignore it.
	ListAllWindowNames(format string) (map[string][]string, error)
	NewWindowInSession(opts model.TmuxWindowOpts) (string, error)
	NextWindowInSession(targetSession string) (string, error)
	SelectWindow(targetWindow string) (string, error)
}

// PaneNavigator addresses individual panes across the attached session. It
// backs `sesh list --panes` and connecting to a "window/pane" name.
type PaneNavigator interface {
	ListTmuxPanes() ([]*model.TmuxPane, error)
	SelectPane(windowIndex int, paneIndex int) (string, error)
}

// ClientSwitcher moves a specific attached client, which is how sesh switches
// sessions when it is launched from outside the multiplexer (a GUI launcher,
// an osascript binding) and there is no current client to act on.
type ClientSwitcher interface {
	ListClients() ([]model.TmuxClient, error)
	// ResolveClient picks the client a switch should move, or "" when there is
	// no better answer than the multiplexer's default.
	ResolveClient() string
	SwitchClientTarget(client string, targetSession string) (string, error)
}
