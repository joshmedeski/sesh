package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func tmuxCmd(args []string) ([]byte, error) {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(tmux, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}

func isAttached() bool {
	return len(os.Getenv("TMUX")) > 0
}

func Sessions() ([]string, error) {
	output, err := tmuxCmd([]string{"list-sessions", "-F", "#{session_last_attached} #{session_name}"})
	if err != nil {
		return nil, nil
	}

	sessionList := strings.TrimSpace(string(output))
	sessionItems := strings.Split(sessionList, "\n")
	sort.SliceStable(sessionItems, func(i, j int) bool {
		return sessionItems[i] > sessionItems[j]
	})
	sessions := make([]string, len(sessionItems))
	for i, item := range sessionItems {
		fields := strings.Fields(item)
		if len(fields) >= 2 {
			sessions[i] = fields[1]
		} else {
			sessions[i] = fields[0]
		}
	}
	return sessions, nil
}

func IsSession(session string) bool {
	sessions, err := Sessions()
	if err != nil {
		return false
	}

	for _, s := range sessions {
		if s == session {
			return true
		}
	}
	return false
}

func attachSession(session string) ([]byte, error) {
	output, err := tmuxCmd([]string{"attach", "-t", session})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func switchSession(session string) ([]byte, error) {
	output, err := tmuxCmd([]string{"switch", "-t", session})
	if err != nil {
		return nil, err
	}
	return output, nil
}

type TmuxSession struct {
	Name           string
	StartDirectory string
}

func NewSession(s TmuxSession) ([]byte, error) {
	output, err := tmuxCmd([]string{"new-session", "-d", "-s", s.Name, "-c", s.StartDirectory})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func Connect(s TmuxSession, alwaysSwitch bool) error {
	isSession := IsSession(s.Name)
	if !isSession {
		print("make session", "\n")
		output, err := NewSession(s)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(output))
	}
	isAttached := isAttached()
	if isAttached || alwaysSwitch {
		switchSession(s.Name)
	} else {
		attachSession(s.Name)
	}
	return nil
}
