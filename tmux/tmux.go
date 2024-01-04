package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"
)

// func isRunning() b ol {
// 	cmd := exec.Command("tmux", "ls")
// 	err := cmd.Run() // throws an exit code if tmux isn't running
// 	return err != nil
// }

func isAttached() bool {
	return len(os.Getenv("TMUX")) > 0
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
		sessions[i] = fields[1]
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
	cmd := exec.Command("tmux", "switch-client", "-t", session)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch to session '%s', error: '%v', output: '%s'", session, err, output)
	}
	return nil
}

func Connect(session string) error {
	if isAttached() {
		switchSession(session)
	} else {
		print("Attaching")
		attachSession(session)
	}
	return nil
}
