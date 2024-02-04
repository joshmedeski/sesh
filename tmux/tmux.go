package tmux

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/dir"
)

type Error struct{ msg string }

func (e Error) Error() string { return e.msg }

var (
	ErrNotRunning = Error{"no server running"}
	ErrNotFound   = Error{"no tmux session found"}
)

func executeCommand(command string, args []string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		if strings.Contains(stderr.String(), "no server running on") {
			return "", ErrNotRunning
		}

		return "", err
	}

	out := strings.TrimSpace(stdout.String())
	if strings.Contains(out, "no server running on") {
		return "", ErrNotRunning
	}

	return out, nil
}

type Command struct {
	cliPath  string
	execFunc func(string, []string) (string, error)
}

func NewCommand() (c *Command, err error) {
	c = new(Command)

	c.cliPath, err = exec.LookPath("tmux")
	if err != nil {
		return nil, err
	}

	c.execFunc = executeCommand

	return c, nil
}

func (c *Command) Run(args []string) (string, error) {
	return c.execFunc(c.cliPath, args)
}

func (c *Command) GetSession(s string) (Session, error) {
	sessionList, err := c.List(Options{})
	if err != nil {
		return Session{}, fmt.Errorf("unable to get tmux sessions: %w", err)
	}

	altPath := dir.AlternatePath(s)

	for _, session := range sessionList {
		if session.Name() == s {
			return session, nil
		}

		if session.Path() == s {
			return session, nil
		}

		if altPath != "" && session.Path() == altPath {
			return session, nil
		}
	}

	return Session{}, fmt.Errorf(
		"%w with name or path matching %q",
		ErrNotFound,
		s,
	)
}

func isAttached() bool {
	return len(os.Getenv("TMUX")) > 0
}

func (c *Command) IsSession(session string) (bool, string) {
	sessions, err := c.List(Options{})
	if err != nil {
		return false, ""
	}

	for _, s := range sessions {
		if s.name == session {
			return true, s.path
		}
	}
	return false, ""
}

func (c *Command) attachSession(session string) error {
	if _, err := c.Run([]string{"attach", "-t", session}); err != nil {
		return err
	}
	return nil
}

func (c *Command) switchSession(session string) error {
	if _, err := c.Run([]string{"switch-client", "-t", session}); err != nil {
		return err
	}
	return nil
}

func (c *Command) runPersistentCommand(session string, command string) error {
	finalCmd := []string{"send-keys", "-t", session, command, "Enter"}
	if _, err := c.Run(finalCmd); err != nil {
		return err
	}
	return nil
}

func (c *Command) NewSession(sessionName, sessionPath string) (string, error) {
	out, err := c.Run(
		[]string{"new-session", "-d", "-s", sessionName, "-c", sessionPath},
	)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (c *Command) execStartupScript(name string, scriptPath string) error {
	bash, err := exec.LookPath("bash")
	if err != nil {
		return err
	}
	cmd := strings.Join(
		[]string{bash, "-c", fmt.Sprintf("\"source %s\"", scriptPath)},
		" ",
	)
	err = c.runPersistentCommand(name, cmd)
	if err != nil {
		return err
	}
	return nil
}

func getStartupScript(sessionPath string, config *config.Config) string {
	for _, script := range config.StartupScripts {
		if dir.FullPath(script.SessionPath) == sessionPath {
			return dir.FullPath(script.ScriptPath)
		}
	}
	return ""
}

func (c *Command) Connect(
	sessionName string,
	alwaysSwitch bool,
	cmd string,
	sessionPath string,
	config *config.Config,
) error {
	isSession, _ := c.IsSession(sessionName)
	if !isSession {
		_, err := c.NewSession(sessionName, sessionPath)
		if err != nil {
			return fmt.Errorf(
				"unable to connect to tmux session %q: %w",
				sessionName,
				err,
			)
		}
		if cmd != "" {
			c.runPersistentCommand(sessionName, cmd)
		} else if scriptPath := getStartupScript(sessionPath, config); scriptPath != "" {
			err := c.execStartupScript(sessionName, scriptPath)
			if err != nil {
				log.Fatal(err)
			}
		} else if config.DefaultStartupScript != "" {
			err := c.execStartupScript(sessionName, config.DefaultStartupScript)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	isAttached := isAttached()
	if isAttached || alwaysSwitch {
		c.switchSession(sessionName)
	} else {
		c.attachSession(sessionName)
	}
	return nil
}
