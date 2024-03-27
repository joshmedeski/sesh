package tmux

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/dir"
)

func GetSession(s string) (TmuxSession, error) {
	sessionList, err := List(Options{})
	if err != nil {
		return TmuxSession{}, fmt.Errorf("unable to get tmux sessions: %w", err)
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

	return TmuxSession{}, fmt.Errorf(
		"no tmux session found with name or path matching %q",
		s,
	)
}

func tmuxCmd(args []string) (string, error) {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return "", err
	}
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(tmux, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		errString := strings.TrimSpace(stderr.String())
		if strings.HasPrefix(errString, "no server running on") {
			return "", nil
		}
		return "", err
	}
	return stdout.String(), nil
}

func isAttached() bool {
	return len(os.Getenv("TMUX")) > 0
}

func FindSession(session string) (*TmuxSession, error) {
	sessions, err := List(Options{})
	if err != nil {
		return nil, err
	}

	for _, s := range sessions {
		if s.Name == session {
			return s, nil
		}
	}
	return nil, nil
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

func NewSession(s TmuxSession) (string, error) {
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

func execStartupCommand(name string, command string) error {
	err := runPersistentCommand(name, command)
	if err != nil {
		return err
	}
	return nil
}

func getStartupScript(sessionPath string, config *config.Config) string {
	for _, sessionConfig := range config.SessionConfigs {
		scriptFullPath := dir.FullPath(sessionConfig.Path)
		match, _ := filepath.Match(scriptFullPath, sessionPath)
		if match {
			return sessionConfig.StartupScript
		}
	}
	return ""
}

func getStartupCommand(sessionPath string, config *config.Config) string {
	for _, sessionConfig := range config.SessionConfigs {
		scriptFullPath := dir.FullPath(sessionConfig.Path)
		match, _ := filepath.Match(scriptFullPath, sessionPath)
		if match {
			return sessionConfig.StartupCommand
		}
	}
	return ""
}

func Connect(
	s TmuxSession,
	alwaysSwitch bool,
	command string,
	sessionPath string,
	config *config.Config,
) error {
	session, _ := FindSession(s.Name)
	if session == nil {
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
		} else if startupScript := getStartupScript(sessionPath, config); startupScript != "" {
			err := execStartupScript(s.Name, startupScript)
			if err != nil {
				log.Fatal(err)
			}
		} else if startupCommand := getStartupCommand(sessionPath, config); startupCommand != "" {
			err := execStartupCommand(s.Name, startupCommand)
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
