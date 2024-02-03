package tmux

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/dir"
)

var (
	command *Command
	once    sync.Once
)

func init() {
	once.Do(func() {
		var err error
		command, err = NewCommand()
		if err != nil {
			log.Fatal(err)
		}
	})
}

type Error struct{ msg string }

func (e Error) Error() string { return e.msg }

var ErrNotRunning = Error{"no server running"}

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

func GetSession(s string) (Session, error) {
	sessionList, err := List(Options{})
	if err != nil {
		return Session{}, fmt.Errorf("unable to get tmux sessions: %w", err)
	}

	altPath := dir.AlternatePath(s)

	for _, session := range sessionList {
		if session.Name == s {
			return *session, nil
		}

		if session.Path == s {
			return *session, nil
		}

		if altPath != "" && session.Path == altPath {
			return *session, nil
		}
	}

	return Session{}, fmt.Errorf(
		"no tmux session found with name or path matching %q",
		s,
	)
}

func tmuxCmd(args []string) (string, error) {
	return command.Run(args)
}

func isAttached() bool {
	return len(os.Getenv("TMUX")) > 0
}

func IsSession(session string) (bool, string) {
	sessions, err := List(Options{})
	if err != nil {
		return false, ""
	}

	for _, s := range sessions {
		if s.Name == session {
			return true, s.Path
		}
	}
	return false, ""
}

func attachSession(session string) error {
	if _, err := tmuxCmd([]string{"attach", "-t", session}); err != nil {
		return err
	}
	return nil
}

func switchSession(session string) error {
	if _, err := tmuxCmd([]string{"switch-client", "-t", session}); err != nil {
		return err
	}
	return nil
}

func runPersistentCommand(session string, command string) error {
	finalCmd := []string{"send-keys", "-t", session, command, "Enter"}
	if _, err := tmuxCmd(finalCmd); err != nil {
		return err
	}
	return nil
}

func NewSession(s Session) (string, error) {
	out, err := tmuxCmd(
		[]string{"new-session", "-d", "-s", s.Name, "-c", s.Path},
	)
	if err != nil {
		return "", err
	}
	return out, nil
}

func execStartupScript(name string, scriptPath string) error {
	bash, err := exec.LookPath("bash")
	if err != nil {
		return err
	}
	cmd := strings.Join(
		[]string{bash, "-c", fmt.Sprintf("\"source %s\"", scriptPath)},
		" ",
	)
	err = runPersistentCommand(name, cmd)
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

func Connect(
	s Session,
	alwaysSwitch bool,
	command string,
	sessionPath string,
	config *config.Config,
) error {
	isSession, _ := IsSession(s.Name)
	if !isSession {
		_, err := NewSession(s)
		if err != nil {
			return fmt.Errorf(
				"unable to connect to tmux session %q: %w",
				s.Name,
				err,
			)
		}
		if command != "" {
			runPersistentCommand(s.Name, command)
		} else if scriptPath := getStartupScript(sessionPath, config); scriptPath != "" {
			err := execStartupScript(s.Name, scriptPath)
			if err != nil {
				log.Fatal(err)
			}
		} else if config.DefaultStartupScript != "" {
			err := execStartupScript(s.Name, config.DefaultStartupScript)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	isAttached := isAttached()
	if isAttached || alwaysSwitch {
		switchSession(s.Name)
	} else {
		attachSession(s.Name)
	}
	return nil
}
