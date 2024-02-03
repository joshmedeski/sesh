package session

import (
	"fmt"
	"os"

	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

type Options struct {
	HideAttached  bool
	IncludeZoxide bool
	IncludeTmux   bool
}

func List(o Options) []string {
	var sessions []string

	tmuxSessions := make([]tmux.Session, 0)
	var sessionPaths []string
	if o.IncludeTmux {
		tmuxList, err := tmux.List(tmux.Options{
			HideAttached: o.HideAttached,
		})
		tmuxSessions = append(tmuxSessions, tmuxList...)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		tmuxSessionNames := make([]string, len(tmuxList))
		for i, session := range tmuxSessions {
			// TODO: allow support for connect as well (PrettyName?)
			// tmuxSessionNames[i] = session.Name + " (" +
			// convert.PathToPretty(session.Path) + ")"
			tmuxSessionNames[i] = session.Name()
			sessionPaths = append(sessionPaths, session.Path())
		}
		sessions = append(sessions, tmuxSessionNames...)
	}

	if o.IncludeZoxide {
		results, err := zoxide.List(sessionPaths)
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
