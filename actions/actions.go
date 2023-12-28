package actions

import (
	"fmt"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
	"strings"

	"github.com/urfave/cli"
)

func Sessions(cCtx *cli.Context) {
	var sessions []string
	hasFlags := cCtx.Bool("tmux") || cCtx.Bool("zoxide")

	if !hasFlags || cCtx.Bool("tmux") {
		tmuxSessions, err := tmux.Sessions()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		sessions = append(sessions, tmuxSessions...)
	}

	if !hasFlags || cCtx.Bool("zoxide") {
		dirs, err := zoxide.Dirs()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		sessions = append(sessions, dirs...)
	}

	fmt.Println(strings.Join(sessions, "\n"))
}
