package shell

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/joshmedeski/sesh/execwrap"
)

type Shell interface {
	Cmd(cmd string, arg ...string) (string, error)
	ListCmd(cmd string, arg ...string) ([]string, error)
}

type RealShell struct {
	exec execwrap.Exec
}

func NewShell(exec execwrap.Exec) Shell {
	return &RealShell{exec}
}

func (c *RealShell) Cmd(cmd string, args ...string) (string, error) {
	foundCmd, err := c.exec.LookPath(cmd)
	if err != nil {
		return "", err
	}
	var stdout, stderr bytes.Buffer
	command := exec.Command(foundCmd, args...)
	command.Stdin = os.Stdin
	command.Stdout = &stdout
	command.Stderr = os.Stderr
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
	trimmedOutput := strings.TrimSuffix(string(stdout.String()), "\n")
	return trimmedOutput, nil
}

func (c *RealShell) ListCmd(cmd string, arg ...string) ([]string, error) {
	command := c.exec.Command(cmd, arg...)
	output, err := command.Output()
	return strings.Split(string(output), "\n"), err
}
