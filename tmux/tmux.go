package tmux

import (
	"joshmedeski/sesh/session"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func IsRunning() bool {
	return len(os.Getenv("TMUX")) > 0
}

type TmuxSession struct {
	session.Session
	tmuxSessionName string
}

func Sessions() ([]string, error) {
	isRunning := IsRunning()
	if !isRunning {
		return nil, nil
	}

	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_last_attached} #{session_name}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
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
