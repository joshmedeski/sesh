package dashboard

import (
	"bytes"
	"os/exec"
	"strings"
)

// runCommand executes a command without connecting stdin (safe for BubbleTea).
// This prevents interference with BubbleTea's terminal control.
func runCommand(name string, args ...string) (string, error) {
	command := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		return "", err
	}
	if err := command.Wait(); err != nil {
		errString := strings.TrimSpace(stderr.String())
		if strings.HasPrefix(errString, "no server running on") {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSuffix(stdout.String(), "\n"), nil
}
