package session

import (
	"fmt"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
	"os"
)

func List(srcs Srcs) []string {
	var sessions []string
	anySrcs := checkAnyTrue(srcs)

	if !anySrcs || srcs.Tmux {
		tmuxSessions, err := tmux.List()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		tmuxSessionNames := make([]string, len(tmuxSessions))
		for i, session := range tmuxSessions {
			tmuxSessionNames[i] = session.Name
		}

		sessions = append(sessions, tmuxSessionNames...)
	}

	if !anySrcs || srcs.Zoxide {
		dirs, err := zoxide.Dirs()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		sessions = append(sessions, dirs...)
	}
	return sessions
}
