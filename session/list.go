package session

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

type Options struct {
	HideAttached bool
	Json         bool
}

func checkAnyTrue(s interface{}) bool {
	val := reflect.ValueOf(s)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Bool && field.Bool() {
			return true
		}
	}
	return false
}

func List(o Options, srcs Srcs) []Session {
	var sessions []Session
	anySrcs := checkAnyTrue(srcs)

	tmuxSessions := make([]*tmux.TmuxSession, 0)
	if !anySrcs || srcs.Tmux {
		tmuxList, err := tmux.List(tmux.Options{
			HideAttached: o.HideAttached,
		})
		tmuxSessions = append(tmuxSessions, tmuxList...)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		tmuxSessionNames := make([]Session, len(tmuxList))
		for i, session := range tmuxSessions {
			// TODO: allow support for connect as well (PrettyName?)
			// tmuxSessionNames[i] = session.Name + " (" + convert.PathToPretty(session.Path) + ")"
			tmuxSessionNames[i] = Session{
				Src:      "tmux",
				Name:     session.Name,
				Path:     session.Path,
				Attached: session.Attached,
				Windows:  session.Windows,
			}
		}
		sessions = append(sessions, tmuxSessionNames...)
	}

	if !anySrcs || srcs.Zoxide {
		results, err := zoxide.List(tmuxSessions)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		zoxideResultNames := make([]Session, len(results))
		for i, result := range results {
			zoxideResultNames[i] = Session{
				Src:   "zoxide",
				Name:  result.Name,
				Path:  result.Path,
				Score: result.Score,
			}
		}
		sessions = append(sessions, zoxideResultNames...)
	}

	return sessions
}
