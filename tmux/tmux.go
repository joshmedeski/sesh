package tmux

import (
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

type Tmux interface {
	ListSessions() ([]*model.TmuxSession, error)
	ListWindows(targetSession string) ([]*model.TmuxWindow, error)
	NewSession(sessionName string, startDir string) (string, error)
	NewWindow(startDir string, name string) (string, error)
	NewWindowInSession(name string, startDir string, targetSession string) (string, error)
	IsAttached() bool
	AttachSession(targetSession string) (string, error)
	SendKeys(name string, command string) (string, error)
	SwitchClient(targetSession string) (string, error)
	CapturePane(targetSession string) (string, error)
	NextWindow() (string, error)
	SelectWindow(targetWindow string) (string, error)
	SwitchOrAttach(name string, opts model.ConnectOpts) (string, error)
	ListTmuxPanes() ([]*model.TmuxPane, error)
	SelectPane(windowIndex int, paneIndex int) (string, error)
	GetCurrentSession() (string, error)
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
	return t.shell.Cmd(t.bin, "attach-session", "-t", targetSession)
}

func (t *RealTmux) SwitchClient(targetSession string) (string, error) {
	return t.shell.Cmd(t.bin, "switch-client", "-t", targetSession)
}

func (t *RealTmux) SendKeys(targetPane string, keys string) (string, error) {
	return t.shell.Cmd(t.bin, "send-keys", "-t", targetPane, keys, "Enter")
}

func (t *RealTmux) NewSession(sessionName string, startDir string) (string, error) {
	return t.shell.Cmd(t.bin, "new-session", "-d", "-s", sessionName, "-c", startDir)
}

func (t *RealTmux) NewWindow(startDir string, name string) (string, error) {
	return t.shell.Cmd(t.bin, "new-window", "-n", name, "-c", startDir)
}

func (t *RealTmux) CapturePane(targetSession string) (string, error) {
	return t.shell.Cmd(t.bin, "capture-pane", "-e", "-p", "-t", targetSession)
}

func (t *RealTmux) NextWindow() (string, error) {
	return t.shell.Cmd(t.bin, "next-window")
}

func (t *RealTmux) IsAttached() bool {
	return len(t.os.Getenv("TMUX")) > 0
}
