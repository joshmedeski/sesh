package seshcli

import (
	"fmt"

	"github.com/joshmedeski/sesh/icon"
	"github.com/joshmedeski/sesh/lister"
	cli "github.com/urfave/cli/v2"
)

func List(icon icon.Icon, list lister.Lister) *cli.Command {
	return &cli.Command{
		Name:                   "list",
		Aliases:                []string{"l"},
		Usage:                  "List sessions",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "show configured sessions",
			},
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
				Name:    "tmuxinator",
				Aliases: []string{"T"},
				Usage:   "show tmuxinator configs",
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
				Usage:   "show icons",
			},
		},
		Action: func(cCtx *cli.Context) error {
			sessions, err := list.List(lister.ListOptions{
				Config:       cCtx.Bool("config"),
				HideAttached: cCtx.Bool("hide-attached"),
				Icons:        cCtx.Bool("icons"),
				Json:         cCtx.Bool("json"),
				Tmux:         cCtx.Bool("tmux"),
				Tmuxinator:   cCtx.Bool("tmuxinator"),
				Zoxide:       cCtx.Bool("zoxide"),
			})
			if err != nil {
				return fmt.Errorf("couldn't list sessions: %q", err)
			}

			for _, i := range sessions.OrderedIndex {
				name := sessions.Directory[i].Name
				if cCtx.Bool("icons") {
					name = icon.AddIcon(sessions.Directory[i])
				}
				fmt.Println(name)
			}

			return nil
		},
	}
}
