package shell

import (
	"strings"

	"github.com/joshmedeski/sesh/execwrap"
)

type Shell struct {
	exec execwrap.Exec
}

func NewShell(exec execwrap.Exec) *Shell {
	return &Shell{exec}
}

func (c *Shell) Cmd(cmd string, args ...string) (string, error) {
	command := c.exec.Command(cmd, args...)
	output, err := command.CombinedOutput()
	return string(output), err
}

func (c *Shell) ListCmd(cmd string, args ...string) ([]string, error) {
	command := c.exec.Command(cmd, args...)
	output, err := command.Output()
	return strings.Split(string(output), "\n"), err
}
