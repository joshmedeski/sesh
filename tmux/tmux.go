package tmux

import (
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"

	"github.com/mattn/go-isatty"
)

// func isRunning() b ol {
// 	cmd := exec.Command("tmux", "ls")
// 	err := cmd.Run() // throws an exit code if tmux isn't running
// 	return err != nil
// }

func isAttached() bool {
	return len(os.Getenv("TMUX")) > 0
}

func isTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

func Sessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_last_attached} #{session_name}")
	output, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	sessionList := strings.TrimSpace(string(output))
	sessionItems := strings.Split(sessionList, "\n")
	sort.SliceStable(sessionItems, func(i, j int) bool {
		return sessionItems[i] > sessionItems[j]
	})
	sessions := make([]string, len(sessionItems))
	for i, item := range sessionItems {
		fields := strings.Fields(item)
		if len(fields) >= 2 {
			sessions[i] = fields[1]
		} else {
			// handle case when there is no whitespace - replace with your behavior in this case
			sessions[i] = ""
		}
	}
	return sessions, nil
}

func IsSession(session string) bool {
	sessions, err := Sessions()
	if err != nil {
		return false
	}

	for _, s := range sessions {
		if s == session {
			return true
		}
	}
	return false
}

func attachSession(session string) error {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}

	args := append([]string{tmux}, "attach", "-t", session)
	if err := syscall.Exec(tmux, args, os.Environ()); err != nil {
		return err
	}

	return nil
}

func switchSession(session string) error {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}

	args := append([]string{tmux}, "switch", "-t", session)
	if err := syscall.Exec(tmux, args, os.Environ()); err != nil {
		return err
	}

	return nil
}

type TmuxSession struct {
	Name           string
	StartDirectory string
}

func NewSession(s TmuxSession) error {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}

	args := append([]string{tmux}, "new-session", "-d", "-s", s.Name, "-c", s.StartDirectory)
	if err := syscall.Exec(tmux, args, os.Environ()); err != nil {
		return err
	}

	return nil
}

func Connect(s TmuxSession) error {
	isSession := IsSession(s.Name)
	if !isSession {
		NewSession(s)
	}
	isAttached := isAttached()
	if isAttached || !isTerminal() {
		switchSession(s.Name)
	} else {
		attachSession(s.Name)
	}
	return nil
}
