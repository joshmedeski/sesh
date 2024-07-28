package tmux

import (
	"bytes"
	"fmt"
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

	if s.PathList != nil {
		var combinedErr error
		for _, path := range s.PathList {
			if err := createNewWindow(s.Name, path); err != nil {
				combinedErr = fmt.Errorf("%w; %v", combinedErr, err)
			}
		}
		if combinedErr != nil {
			return out, combinedErr
		}
	}

	return out, nil
}

func createNewWindow(sessionName, path string) error {
	fullPath := dir.FullPath(path)
	if fullPath == "" {
		return nil
	}
	info, err := os.Stat(fullPath)
	if err != nil || !info.IsDir() {
		return nil
	}
	_, err = tmuxCmd([]string{"new-window", "-t", sessionName, "-c", fullPath, "-d"})
	return err
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

func execTmuxp(name string, command string) error {
	err := runPersistentCommand(name, command)
	if err != nil {
		return err
	}
	return nil
}

func removeTrailingSlash(path string) string {
	return strings.TrimRight(path, "/")
}

func getStartupScript(sessionPath string, config *config.Config) string {
	for _, sessionConfig := range config.SessionConfigs {
		// TODO: get working with /* again
		scriptFullPath := removeTrailingSlash(dir.FullPath(sessionConfig.Path))
		match, _ := filepath.Match(scriptFullPath, sessionPath)
		if match {
			return sessionConfig.StartupScript
		}
	}
	return ""
}

func getStartupCommand(sessionPath string, config *config.Config) string {
	for _, sessionConfig := range config.SessionConfigs {
		scriptFullPath := removeTrailingSlash(dir.FullPath(sessionConfig.Path))
		match, _ := filepath.Match(scriptFullPath, sessionPath)
		if match {
			return sessionConfig.StartupCommand
		}
	}
	return ""
}

func getTmuxp(sessionPath string, config *config.Config) string {
	for _, sessionConfig := range config.SessionConfigs {
		scriptFullPath := dir.FullPath(sessionConfig.Path)
		match, _ := filepath.Match(scriptFullPath, sessionPath)
		if match {
			return sessionConfig.Tmuxp
		}
	}
	return ""
}
