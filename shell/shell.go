package shell

import (
	"strings"

	"github.com/joshmedeski/sesh/execwrap"
	"github.com/stretchr/testify/mock"
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
	return string(output), err
}

func (c *RealShell) ListCmd(cmd string, arg ...string) ([]string, error) {
	command := c.exec.Command(cmd, arg...)
	output, err := command.Output()
	return strings.Split(string(output), "\n"), err
}

type MockShell struct {
	mock.Mock
}

func (m *MockShell) Cmd(cmd string, arg ...string) (string, error) {
	args := m.Called(cmd, arg)
	return args.String(0), args.Error(1)
}

func (m *MockShell) ListCmd(name string, arg ...string) ([]string, error) {
	args := m.Called(name, arg)
	return args.Get(0).([]string), args.Error(1)
}
