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

func IsSession(session string) bool {
	sessions, err := List()
	if err != nil {
		return false
	}

	for _, s := range sessions {
		if s.Name == session {
			return true
		}
	}
	return false
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
	out, err := tmuxCmd([]string{"new-session", "-d", "-s", s.Name, "-c", s.Path})
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
	_, err = tmuxCmd([]string{"send-keys", "-t", name, cmd, "Enter"})
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
	s TmuxSession,
	alwaysSwitch bool,
	command string,
	sessionPath string,
	config *config.Config,
) error {
	isSession := IsSession(s.Name)
	if !isSession {
		_, err := NewSession(s)
		if err != nil {
			log.Fatal(err)
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
