package cmds

import (
	"fmt"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/joshmedeski/sesh/icons"
	"github.com/joshmedeski/sesh/json"
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
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "output as json",
			},
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
			&cli.BoolFlag{
				Name:    "icons",
				Aliases: []string{"i"},
				Usage:   "show Nerd Font icons",
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

			useIcons := cCtx.Bool("icons")
			result := make([]string, len(sessions))
			for i, session := range sessions {
				if useIcons {
					result[i] = icons.PrintWithIcon(session)
				} else {
					result[i] = session.Name
				}
			}

			useJson := cCtx.Bool("json")
			if useJson {
				fmt.Println(json.List(sessions))
			} else {
				fmt.Println(strings.Join(result, "\n"))
			}
			return nil
		},
	}
}
