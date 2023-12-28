package utils

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

func ListSessions() {
	sessions, err := getTmuxSessions()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions found")
	} else {
		fmt.Println(strings.Join(sessions, "\n"))
	}

	// TODO: list zoxide sesssions
}

func getTmuxSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_last_attached} #{session_name}")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	sessionList := strings.TrimSpace(string(output))
	sessionItems := strings.Split(sessionList, "\n")

	// Custom sorting by session_last_attached in reverse order
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
