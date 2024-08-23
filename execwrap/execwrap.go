package execwrap

import (
	"os/exec"
)

type ExecCmd interface {
	CombinedOutput() ([]byte, error)
	Output() ([]byte, error)
}

type Exec interface {
	LookPath(executable string) (string, error)
	Command(name string, args ...string) ExecCmd
}

type OsExec struct{}

func NewExec() Exec {
	return &OsExec{}
}

func (e *OsExec) LookPath(executable string) (string, error) {
	return exec.LookPath(executable)
}

func (e *OsExec) Command(name string, args ...string) ExecCmd {
	return exec.Command(name, args...)
}
