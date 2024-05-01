package seshcli

import (
	"fmt"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/execwrap"
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/oswrap"
	"github.com/joshmedeski/sesh/pathwrap"
	"github.com/joshmedeski/sesh/runtimewrap"
	"github.com/joshmedeski/sesh/session"
	"github.com/joshmedeski/sesh/shell"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
	cli "github.com/urfave/cli/v2"
)

func List() *cli.Command {
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
				Usage:   "show Nerd Font icons",
			},
		},
		Action: func(cCtx *cli.Context) error {
			// wrapper dependencies
			ew := execwrap.NewExec()
			os := oswrap.NewOs()
			p := pathwrap.NewPath()
			r := runtimewrap.NewRunTime()

			// base dependencies
			sh := shell.NewShell(ew)
			h := home.NewHome(os)

			// core dependencies
			tx := tmux.NewTmux(sh)
			z := zoxide.NewZoxide(sh)
			c := config.NewConfig(os, p, r)
			s := session.NewSession(c, h, tx, z)

			sessions, err := s.List(session.ListOptions{
				Config:       cCtx.Bool("config"),
				HideAttached: cCtx.Bool("hide-attached"),
				Icons:        cCtx.Bool("icons"),
				Json:         cCtx.Bool("json"),
				Tmux:         cCtx.Bool("tmux"),
				Zoxide:       cCtx.Bool("zoxide"),
			})
			if err != nil {
				return fmt.Errorf("couldn't list sessions: %q", err)
			}

			for _, session := range sessions {
				fmt.Println(session.Name)
			}

			return nil
		},
	}
}
