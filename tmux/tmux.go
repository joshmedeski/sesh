package tmux

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

type Tmux interface {
	ListSessions() ([]*model.TmuxSession, error)
	NewSession(sessionName string, startDir string) (string, error)
	IsAttached() bool
	AttachSession(targetSession string) (string, error)
	SendKeys(name string, command string) (string, error)
	SwitchClient(targetSession string) (string, error)
	CapturePane(targetSession string) (string, error)
	SwitchOrAttach(name string, opts model.ConnectOpts) (string, error)
	GetWindowName(session, window string) (string, error)
}

type RealTmux struct {
	os    oswrap.Os
	shell shell.Shell
}

func NewTmux(os oswrap.Os, shell shell.Shell) Tmux {
	return &RealTmux{os, shell}
}

func (t *RealTmux) AttachSession(targetSession string) (string, error) {
	return t.shell.Cmd("tmux", "attach-session", "-t", targetSession)
}

func (t *RealTmux) SwitchClient(targetSession string) (string, error) {
	return t.shell.Cmd("tmux", "switch-client", "-t", targetSession)
}

func (t *RealTmux) SendKeys(targetPane string, keys string) (string, error) {
	return t.shell.Cmd("tmux", "send-keys", "-t", targetPane, keys, "Enter")
}

func (t *RealTmux) NewSession(sessionName string, startDir string) (string, error) {
	return t.shell.Cmd("tmux", "new-session", "-d", "-s", sessionName, "-c", startDir)
}

func (t *RealTmux) CapturePane(targetSession string) (string, error) {
	return t.shell.Cmd("tmux", "capture-pane", "-e", "-p", "-t", targetSession)
}

func (t *RealTmux) IsAttached() bool {
	return len(t.os.Getenv("TMUX")) > 0
}

func (t *RealTmux) GetWindowName(session, window string) (string, error) {
	result, err := t.shell.Cmd("tmux", "list-windows", "-t", session, "-F", "#{window_index}:#{window_name}")
	if err != nil {
		return window, err
	}
	
	lines := strings.Split(strings.TrimSpace(result), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[0] == window {
			return parts[1], nil
		}
	}
	return window, nil
}
