package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joshmedeski/sesh/v2/execwrap"
	"github.com/joshmedeski/sesh/v2/home"
)

type Shell interface {
	Cmd(cmd string, arg ...string) (string, error)
	CmdInDir(dir string, cmd string, arg ...string) (string, error)
	CmdWithOutput(cmd string, arg ...string) (string, error)
	// CmdCapture runs cmd and returns whatever it wrote to stdout even when
	// it exits non-zero. Some tools report partial success that way: `gh api
	// graphql` prints a complete response body containing both resolved data
	// and an "errors" array, then exits 1. Cmd discards that body; callers
	// that need to salvage it use this instead and decide for themselves
	// whether the output is usable.
	CmdCapture(cmd string, arg ...string) (string, error)
	ListCmd(cmd string, arg ...string) ([]string, error)
	PrepareCmd(cmd string, replacements map[string]string) ([]string, error)
	// ShellCmd runs cmd through the user's shell (honoring $SHELL, falling
	// back to /bin/sh) instead of splitting it into argv and exec'ing it
	// directly. This lets config commands such as preview_command use
	// shell operators (&&, ||, |, ;) and other shell syntax that a naive
	// space-split + exec.Command call cannot support. Each key in
	// replacements is substituted in cmd with its value, shell-quoted so
	// that values containing spaces or shell metacharacters are treated
	// as a single literal argument.
	ShellCmd(cmd string, replacements map[string]string) (string, error)
}

type RealShell struct {
	exec execwrap.Exec
	home home.Home
}

func NewShell(exec execwrap.Exec, home home.Home) Shell {
	return &RealShell{exec, home}
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

func (c *RealShell) CmdInDir(dir string, cmd string, args ...string) (string, error) {
	foundCmd, err := c.exec.LookPath(cmd)
	if err != nil {
		return "", err
	}
	var stdout, stderr bytes.Buffer
	command := exec.Command(foundCmd, args...)
	command.Dir = dir
	command.Stdin = os.Stdin
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return "", err
	}
	return strings.TrimSuffix(stdout.String(), "\n"), nil
}

func (c *RealShell) CmdCapture(cmd string, args ...string) (string, error) {
	foundCmd, err := c.exec.LookPath(cmd)
	if err != nil {
		return "", err
	}
	var stdout, stderr bytes.Buffer
	command := exec.Command(foundCmd, args...)
	command.Stdin = os.Stdin
	command.Stdout = &stdout
	command.Stderr = &stderr

	runErr := command.Run()
	out := strings.TrimSuffix(stdout.String(), "\n")
	if runErr != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return out, fmt.Errorf("%s: %w: %s", cmd, runErr, msg)
		}
		return out, fmt.Errorf("%s: %w", cmd, runErr)
	}
	return out, nil
}

func (c *RealShell) CmdWithOutput(cmd string, args ...string) (string, error) {
	foundCmd, err := c.exec.LookPath(cmd)
	if err != nil {
		return "", err
	}
	command := exec.Command(foundCmd, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return "", err
	}
	return "", nil
}

func (c *RealShell) ListCmd(cmd string, arg ...string) ([]string, error) {
	command := c.exec.Command(cmd, arg...)
	output, err := command.Output()
	return strings.Split(string(output), "\n"), err
}

func (c *RealShell) PrepareCmd(cmd string, replacements map[string]string) ([]string, error) {
	cmdParts := strings.Split(cmd, " ")
	result := make([]string, len(cmdParts))

	for i, arg := range cmdParts {
		expanded, err := c.home.ExpandPath(arg)
		if err != nil {
			return nil, err
		}
		if replacement, ok := replacements[expanded]; ok {
			result[i] = replacement
		} else {
			result[i] = expanded
		}
	}

	return result, nil
}

// userShell returns the shell binary to run script commands with. It
// prefers $SHELL (the user's configured login shell) so aliases, functions
// and syntax the user expects are honored, and falls back to /bin/sh when
// $SHELL is unset (e.g. minimal/CI environments).
func userShell() string {
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh
	}
	return "/bin/sh"
}

// shellQuote wraps s in single quotes so it is treated as one literal shell
// argument, escaping any single quotes it contains. This is the standard
// POSIX-shell-safe quoting technique: close the quote, emit an escaped
// quote, reopen the quote.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func (c *RealShell) ShellCmd(cmd string, replacements map[string]string) (string, error) {
	for placeholder, value := range replacements {
		cmd = strings.ReplaceAll(cmd, placeholder, shellQuote(value))
	}

	sh := userShell()
	foundCmd, err := c.exec.LookPath(sh)
	if err != nil {
		return "", err
	}

	var stdout, stderr bytes.Buffer
	command := exec.Command(foundCmd, "-c", cmd)
	command.Stdin = os.Stdin
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
	trimmedOutput := strings.TrimSuffix(stdout.String(), "\n")
	return trimmedOutput, nil
}
