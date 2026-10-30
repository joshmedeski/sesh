package tmux

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

// sessionTarget adds a trailing colon so tmux reads a dotted name as a session, not a pane
// an existing colon already marks a window or pane target, so that target stays unchanged
func sessionTarget(name string) string {
	if strings.Contains(name, ":") || isPaneID(name) {
		return name
	}
	return name + ":"
}

func isPaneID(name string) bool {
	id, ok := strings.CutPrefix(name, "%")
	if !ok || id == "" {
		return false
	}
	return strings.Trim(id, "0123456789") == ""
}

func escapeFormat(s string) string {
	return strings.ReplaceAll(s, "#", "##")
}

type Tmux interface {
	ListSessions() ([]*model.TmuxSession, error)
	ListWindows(targetSession string) ([]*model.TmuxWindow, error)
	ListAllWindowNames(format string) (map[string][]string, error)
	ListAllWindows() (map[string][]model.TmuxWindow, error)
	NewSession(sessionName string, startDir string) (string, error)
	NewWindowInSession(opts model.TmuxWindowOpts) (string, error)
	NewPaneInWindow(opts model.TmuxPaneOpts) (string, error)
	IsAttached() bool
	AttachSession(targetSession string) (string, error)
	SendKeys(name string, command string) (string, error)
	SwitchClient(targetSession string) (string, error)
	CapturePane(targetSession string) (string, error)
	NextWindowInSession(targetSession string) (string, error)
	SelectWindow(targetWindow string) (string, error)
	SelectLayout(targetWindow string, layout string) (string, error)
	SwitchOrAttach(name string, opts model.ConnectOpts) (string, error)
	ListTmuxPanes() ([]*model.TmuxPane, error)
	SelectPane(windowIndex int, paneIndex int) (string, error)
	SelectPaneByID(paneID string) (string, error)
	ZoomPane(target string) (string, error)
	GetCurrentSession() (string, error)
	ListClients() ([]model.TmuxClient, error)
	ResolveClient() string
	FocusedSessionPath() (string, error)
	SwitchClientTarget(client string, targetSession string) (string, error)
	RenameSession(target string, newName string) (string, error)
	KillSession(target string) (string, error)
	DisplayMessage(message string) (string, error)
}

type RealTmux struct {
	os    oswrap.Os
	shell shell.Shell
	bin   string
}

func NewTmux(os oswrap.Os, shell shell.Shell, bin string) Tmux {
	if bin == "" {
		bin = "tmux"
	}
	return &RealTmux{os, shell, bin}
}

func (t *RealTmux) AttachSession(targetSession string) (string, error) {
	// Using t.shell.Cmd() here (capturing stdout/stderr in buffers) would block psmux on Windows.
	return t.shell.CmdWithOutput(t.bin, "attach-session", "-t", sessionTarget(targetSession))
}

func (t *RealTmux) SwitchClient(targetSession string) (string, error) {
	return t.shell.Cmd(t.bin, "switch-client", "-t", sessionTarget(targetSession))
}

func (t *RealTmux) SendKeys(targetPane string, keys string) (string, error) {
	return t.shell.Cmd(t.bin, "send-keys", "-t", sessionTarget(targetPane), keys, "Enter")
}

func (t *RealTmux) NewSession(sessionName string, startDir string) (string, error) {
	return t.shell.Cmd(t.bin, "new-session", "-d", "-s", escapeFormat(sessionName), "-c", escapeFormat(startDir))
}

func (t *RealTmux) CapturePane(targetSession string) (string, error) {
	return t.shell.Cmd(t.bin, "capture-pane", "-e", "-p", "-t", sessionTarget(targetSession))
}

func (t *RealTmux) NextWindowInSession(targetSession string) (string, error) {
	return t.shell.Cmd(t.bin, "next-window", "-t", sessionTarget(targetSession))
}

func (t *RealTmux) IsAttached() bool {
	return len(t.os.Getenv("TMUX")) > 0
}

func (t *RealTmux) SwitchClientTarget(client string, targetSession string) (string, error) {
	return t.shell.Cmd(t.bin, "switch-client", "-c", client, "-t", sessionTarget(targetSession))
}

// DisplayMessage shows the message in the status line until a key is pressed.
func (t *RealTmux) DisplayMessage(message string) (string, error) {
	return t.shell.Cmd(t.bin, "display-message", "-d", "0", escapeFormat(message))
}

func (t *RealTmux) SelectLayout(targetWindow string, layout string) (string, error) {
	return t.shell.Cmd(t.bin, "select-layout", "-t", sessionTarget(targetWindow), layout)
}
