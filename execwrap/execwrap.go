package execwrap

import (
	"os/exec"

	"github.com/stretchr/testify/mock"
)

type ExecCmd interface {
	CombinedOutput() ([]byte, error)
	Output() ([]byte, error)
}

type Exec interface {
	LookPath(executable string) (string, error)
	Command(name string, arg ...string) ExecCmd
}

type OsExec struct{}

func New() Exec {
	return &OsExec{}
}

func (e *OsExec) LookPath(executable string) (string, error) {
	return exec.LookPath(executable)
}

func (e *OsExec) Command(name string, arg ...string) ExecCmd {
	return exec.Command(name, arg...)
}

type MockExec struct {
	mock.Mock
}

func (m *MockExec) LookPath(executable string) (string, error) {
	args := m.Called(executable)
	return args.String(0), args.Error(1)
}

func (m *MockExec) Command(name string, arg ...string) ExecCmd {
	args := m.Called(name, arg)
	return args.Get(0).(ExecCmd)
}

type MockExecCmd struct {
	mock.Mock
}

// CombinedOutput mocks the os/exec.Cmd's CombinedOutput method
func (m *MockExecCmd) CombinedOutput() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockExecCmd) Output() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}
