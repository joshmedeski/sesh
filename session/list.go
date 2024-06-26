package session

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joshmedeski/sesh/config"
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

// Takes as input a list of sessions and an injected mapper function
// Returns a map using the mapper function
// to show existing sessions
func makeSessionsMap(sessions []Session, mapper func(s Session) string) map[string]bool {
	sessionMap := make(map[string]bool, len(sessions))
	for _, session := range sessions {
		sessionMap[mapper(session)] = true
	}
	return sessionMap
}

func isInSessionMap(sessionMap map[string]bool, key string) bool {
	_, exists := sessionMap[key]
	return exists
}

func listTmuxSessions(o Options) (sessions []Session, err error) {
	tmuxList, err := tmux.List(tmux.Options{
		HideAttached: o.HideAttached,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't list tmux sessions: %q", err)
	}
	tmuxSessions := make([]Session, len(tmuxList))
	for i, session := range tmuxList {
		tmuxSessions[i] = Session{
			Src:      "tmux",
			Name:     session.Name,
			Path:     session.Path,
			Attached: session.Attached,
			Windows:  session.Windows,
		}
	}
	return tmuxSessions, nil
}

func listConfigSessions(c *config.Config, existingSessions []Session) (sessions []Session, err error) {
	var configSessions []Session
	fString := "%s-%s"

	// filter sessions by name+path combination
	sessionMap := makeSessionsMap(existingSessions, func(s Session) string { return fmt.Sprintf(fString, s.Name, s.Path) })

	for _, sessionConfig := range c.SessionConfigs {
		if !isInSessionMap(sessionMap, fmt.Sprintf(fString, sessionConfig.Name, sessionConfig.Path)) && sessionConfig.Name != "" {
			configSessions = append(configSessions, Session{
				Src:  "config",
				Name: sessionConfig.Name,
				Path: sessionConfig.Path,
			})
		}
	}
	return configSessions, nil
}

func listZoxideSessions(existingSessions []Session) (sessions []Session, err error) {
	results, err := zoxide.List()
	if err != nil {
		return nil, fmt.Errorf("couldn't list zoxide results: %q", err)
	}
	var zoxideSessions []Session
	sessionMap := makeSessionsMap(existingSessions, func(s Session) string { return s.Path })
	for _, result := range results {
		if !isInSessionMap(sessionMap, result.Path) {
			zoxideSessions = append(zoxideSessions, Session{
				Src:   "zoxide",
				Name:  result.Name,
				Path:  result.Path,
				Score: result.Score,
			})
		}
	}
	return zoxideSessions, nil
}

func List(options Options, srcs Srcs, config *config.Config) []Session {
	var sessions []Session
	anySrcs := checkAnyTrue(srcs)

	if !anySrcs || srcs.Tmux {
		tmuxSessions, err := listTmuxSessions(options)
		if err != nil {
			fmt.Println("list failed:", err)
			os.Exit(1)
		}
		sessions = append(sessions, tmuxSessions...)
	}

	if !anySrcs || srcs.Config {
		configSessions, err := listConfigSessions(config, sessions)
		if err != nil {
			fmt.Println("list failed:", err)
			os.Exit(1)
		}
		sessions = append(sessions, configSessions...)
	}

	if !anySrcs || srcs.Zoxide {
		zoxideSessions, err := listZoxideSessions(sessions)
		if err != nil {
			fmt.Println("list failed:", err)
			os.Exit(1)
		}
		sessions = append(sessions, zoxideSessions...)
	}

	return sessions
}
