package seshcli

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/json"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	cli "github.com/urfave/cli/v2"
)

func List(icon icon.Icon, json json.Json, list lister.Lister) *cli.Command {
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
			&cli.BoolFlag{
				Name:    "tmuxinator",
				Aliases: []string{"T"},
				Usage:   "show tmuxinator configs",
			},
			&cli.BoolFlag{
				Name:    "hide-duplicates",
				Aliases: []string{"d"},
				Usage:   "hide duplicate entries",
			},
		},
		Action: func(cCtx *cli.Context) error {
			sessions, err := list.List(lister.ListOptions{
				Config:         cCtx.Bool("config"),
				HideAttached:   cCtx.Bool("hide-attached"),
				Icons:          cCtx.Bool("icons"),
				Json:           cCtx.Bool("json"),
				Tmux:           cCtx.Bool("tmux"),
				Zoxide:         cCtx.Bool("zoxide"),
				Tmuxinator:     cCtx.Bool("tmuxinator"),
				HideDuplicates: cCtx.Bool("hide-duplicates"),
			})
			if err != nil {
				return fmt.Errorf("couldn't list sessions: %q", err)
			}

			if cCtx.Bool("json") {
				var sessionsArray []model.SeshSession
				for _, i := range sessions.OrderedIndex {
					sessionsArray = append(sessionsArray, sessions.Directory[i])
				}
				fmt.Println(json.EncodeSessions(sessionsArray))
				return nil
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
