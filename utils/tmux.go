package utils

import (
	"os/exec"
	"strings"
)

// Determines if a tmux server is running
// Return true if tmux is running, false otherwise
func IsTmuxRunning() (bool, error) {
	output, err := exec.Command("tmux", "list-sessions").Output()
	if err != nil {
		return false, err
	}
	outputString := strings.TrimSpace(string(output))
	return len(outputString) > 0, err
}
