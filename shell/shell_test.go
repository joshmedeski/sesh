package shell

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/execwrap"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShellCmd(t *testing.T) {
	t.Run("run should succeed", func(t *testing.T) {
		mockExec := new(execwrap.MockExec)
		mockExec.On("LookPath", "echo", mock.Anything).Return("echo", nil)
		mockCmd := new(execwrap.MockExecCmd)
		shell := &RealShell{exec: mockExec}
		mockCmd.On("CombinedOutput").Return([]byte("hello"), nil)
		mockExec.On("Command", "echo", mock.Anything).Return(mockCmd)
		out, err := shell.Cmd("echo", "hello")
		assert.Nil(t, err)
		assert.Equal(t, "hello", out)
	})
}

func TestShellListCmd(t *testing.T) {
	t.Run("run should succeed", func(t *testing.T) {
		mockExec := new(execwrap.MockExec)
		mockCmd := new(execwrap.MockExecCmd)
		shell := &RealShell{exec: mockExec}
		dirListingActual := []byte(`total 9720
drwxr-xr-x  17 joshmedeski  staff      544 Apr 11 21:40 ./
drwxr-xr-x   8 joshmedeski  staff      256 Apr 11 19:05 ../
-rw-r--r--   1 joshmedeski  staff       53 Apr 11 09:00 .git`)
		mockCmd.On("Output").Return(dirListingActual, nil)
		mockExec.On("Command", "ls", mock.Anything).Return(mockCmd)
		dirListingExpected := []string{
			"total 9720",
			"drwxr-xr-x  17 joshmedeski  staff      544 Apr 11 21:40 ./",
			"drwxr-xr-x   8 joshmedeski  staff      256 Apr 11 19:05 ../",
			"-rw-r--r--   1 joshmedeski  staff       53 Apr 11 09:00 .git",
		}
		list, err := shell.ListCmd("ls", "-la")
		assert.Nil(t, err)
		assert.Equal(t, dirListingExpected, list)
	})
}

func TestShellPrepareCmd(t *testing.T) {
	t.Run("should succeed with correct replacements and expansions", func(t *testing.T) {
		mockHome := new(home.MockHome)
		shell := &RealShell{home: mockHome}
		mockHome.On("ExpandHome", "~/.local/bin/rat").Return("/home/test/.local/bin/rat", nil)
		cmdParts, err := shell.PrepareCmd("~/.local/bin/rat {}", map[string]string{"{}": "hello"})
		assert.Nil(t, err)
		assert.Equal(t, []string{"/home/test/.local/bin/rat", "hello"}, cmdParts)
	})
}
