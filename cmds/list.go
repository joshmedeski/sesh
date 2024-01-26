package cmds

import (
	"fmt"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/joshmedeski/sesh/session"
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
			&cli.BoolFlag{
				Name:    "hide-attached",
				Aliases: []string{"H"},
				Usage:   "don't show currently attached sessions",
			},
		},
		Action: func(cCtx *cli.Context) error {
			o := session.Options{
				HideAttached: cCtx.Bool("hide-attached"),
			}
			sessions := session.List(o, session.Srcs{
				Tmux:   cCtx.Bool("tmux"),
				Zoxide: cCtx.Bool("zoxide"),
			})
			fmt.Println(strings.Join(sessions, "\n"))
			return nil
		},
	}
}
