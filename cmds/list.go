package cmds

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/session"

	"github.com/urfave/cli/v2"
)

func List() *cli.Command {
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
			sessions := session.List(session.Srcs{
				Tmux:   cCtx.Bool("tmux"),
				Zoxide: cCtx.Bool("zoxide"),
			})
			fmt.Println(strings.Join(sessions, "\n"))
			return nil
		},
	}
}
