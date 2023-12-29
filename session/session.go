package session

import (
	"fmt"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
	"os"
	"reflect"
)

type session struct {
	Type        string
	Value       string
	DisplayName string
}

type Srcs struct {
	Tmux   bool
	Zoxide bool
}

func checkAnyTrue(s interface{}) bool {
	val := reflect.ValueOf(s)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Ptr && !field.IsNil() && field.Elem().Bool() {
			return true
		}
	}
	return false
}

func Sessions(srcs Srcs) []string {
	var sessions []string
	anySrcs := checkAnyTrue(srcs)

	if !anySrcs || srcs.Tmux {
		tmuxSessions, err := tmux.Sessions()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		sessions = append(sessions, tmuxSessions...)
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
