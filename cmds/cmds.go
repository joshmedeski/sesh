package cmds

import (
	"fmt"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func ListSessions() *cli.Command {
	return &cli.Command{
		Name:                   "list",
		Aliases:                []string{"l"},
		Usage:                  "List sessions",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "tmux",
				Aliases: []string{"t"},
				Usage:   "show tmux sessions",
			},
			&cli.BoolFlag{
				Name:    "zoxide",
				Aliases: []string{"z"},
				Usage:   "show zoxide results",
			},
		},
		Action: func(cCtx *cli.Context) error {
			var sessions []string
			hasFlags := cCtx.Bool("tmux") || cCtx.Bool("zoxide")

			if !hasFlags || cCtx.Bool("tmux") {
				tmuxSessions, err := tmux.Sessions()
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				sessions = append(sessions, tmuxSessions...)
			}

			if !hasFlags || cCtx.Bool("zoxide") {
				dirs, err := zoxide.Dirs()
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				sessions = append(sessions, dirs...)
			}

			fmt.Println(strings.Join(sessions, "\n"))
			return nil
		},
	}
}
