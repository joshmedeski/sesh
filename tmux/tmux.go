package tmux

import (
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/shell"
)

type Tmux interface {
	ListSessions() ([]*model.TmuxSession, error)
	AttachSession(targetSession string) (string, error)
	SwitchClient(targetSession string) (string, error)
}

type RealTmux struct {
	shell shell.Shell
}

func NewTmux(shell shell.Shell) Tmux {
	return &RealTmux{shell}
}

func (t *RealTmux) AttachSession(targetSession string) (string, error) {
	return t.shell.Cmd("tmux", "attach-session", "-t", targetSession)
}

func (t *RealTmux) SwitchClient(targetSession string) (string, error) {
	return t.shell.Cmd("tmux", "switch-client", "-t", targetSession)
}

func (t *RealTmux) SendKeys(targetPane string, keys string) (string, error) {
	return t.shell.Cmd("tmux", "send-keys", "-t", targetPane, keys)
}

func (t *RealTmux) NewSession(sessionName string, startDir string) (string, error) {
	return t.shell.Cmd("tmux", "new-session", "-s", sessionName, "-d", startDir)
}
