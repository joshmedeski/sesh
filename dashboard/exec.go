package dashboard

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func runShellCommand(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return runCommand("cmd", "/c", cmd)
	}
	return runCommand("sh", "-c", cmd)
}

// runCommand executes a command without connecting stdin (safe for BubbleTea).
// This prevents interference with BubbleTea's terminal control.
func runCommand(name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	command := exec.CommandContext(ctx, name, args...)
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
