// exec.go
package core

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/joshmedeski/sesh/v2/execwrap"
)

// CommandRunner abstracts the external commands the dashboard sections run
// (tmux, git, docker, workmux, ssh, and user shell commands).
type CommandRunner interface {
	Run(name string, args ...string) (string, error)
	RunShell(cmd string) ([]byte, error)
}

// execCommandRunner is the production CommandRunner backed by execwrap.Exec.
type execCommandRunner struct {
	exec execwrap.Exec
}

// NewCommandRunner returns a CommandRunner that executes commands via exec.
func NewCommandRunner(exec execwrap.Exec) CommandRunner {
	return &execCommandRunner{exec: exec}
}

// Run executes name with args, returning stdout with one trailing newline
// trimmed. A "no server running on" error is treated as empty output.
func (r *execCommandRunner) Run(name string, args ...string) (string, error) {
	out, err := r.exec.Command(name, args...).Output()
	if err != nil {
		errString := strings.TrimSpace(err.Error())
		if strings.Contains(errString, "no server running on") {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}

// RunShell executes cmd through the platform shell (sh -c on Unix, cmd /c on
// Windows), returning stdout. stderr is only surfaced in the error when the
// command fails.
func (r *execCommandRunner) RunShell(cmd string) ([]byte, error) {
	name, flag := "sh", "-c"
	if runtime.GOOS == "windows" {
		name, flag = "cmd", "/c"
	}
	out, err := r.exec.Command(name, flag, cmd).Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
		return out, fmt.Errorf("%w: %s", err, strings.TrimSpace(string(exitErr.Stderr)))
	}
	return out, err
}
