package tmux

import (
	"fmt"
	"os"
	"os/exec"
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

// func SessionPath(session string) *string {
// 	sessions, err := Sessions()
// 	if err != nil {
// 		return nil
// 	}
//
// 	for _, s := range sessions {
// 		if s == session {
// 			if len(fields) >= 3 {
// 				sessions[i] = fields[1]
// 			} else {
// 				sessions[i] = fields[0]
// 			}
// 		}
// 	}
// 	return false
//
// 	output, err := tmuxCmd([]string{"display-message", "-p", "-F", "#{session_path}"})
// }

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

func NewSession(s TmuxSession) ([]byte, error) {
	output, err := tmuxCmd([]string{"new-session", "-d", "-s", s.Name, "-c", s.Path})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func Connect(s TmuxSession, alwaysSwitch bool) error {
	isSession := IsSession(s.Name)
	if !isSession {
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
