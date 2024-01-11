package session

import (
	"fmt"
	"joshmedeski/sesh/convert"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
	"os"
)

func List(srcs Srcs) []string {
	var sessions []string
	anySrcs := checkAnyTrue(srcs)

	tmuxSessions := make([]*tmux.TmuxSession, 0)
	if !anySrcs || srcs.Tmux {
		tmuxList, err := tmux.List()
		tmuxSessions = append(tmuxSessions, tmuxList...)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		tmuxSessionNames := make([]string, len(tmuxList))
		for i, session := range tmuxSessions {
			tmuxSessionNames[i] = session.Name + " (" + convert.PathToPretty(session.Path) + ")"
		}
		sessions = append(sessions, tmuxSessionNames...)
	}

	if !anySrcs || srcs.Zoxide {
		results, err := zoxide.List(tmuxSessions)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		zoxideResultNames := make([]string, len(results))
		for i, result := range results {
			zoxideResultNames[i] = result.Name
		}
		sessions = append(sessions, zoxideResultNames...)
	}

	return sessions
}
