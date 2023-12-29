package tmux

import (
	"os/exec"
	"sort"
	"strings"
)

// func isRunning() b ol {
// 	cmd := exec.Command("tmux", "ls")
// 	err := cmd.Run() // throws an exit code if tmux isn't running
// 	return err != nil
// }
// func isActive() bool {
// 	return len(os.Getenv("TMUX")) > 0
// }

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
