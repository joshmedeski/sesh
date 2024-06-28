package shell

import (
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

func (c *RealShell) Cmd(cmd string, arg ...string) (string, error) {
	command := c.exec.Command(cmd, arg...)
	output, err := command.CombinedOutput()
	trimmedOutput := strings.TrimSuffix(string(output), "\n")
	return trimmedOutput, err
}

func (c *RealShell) ListCmd(cmd string, arg ...string) ([]string, error) {
	command := c.exec.Command(cmd, arg...)
	output, err := command.Output()
	return strings.Split(string(output), "\n"), err
}
