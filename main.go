package main

import (
	"joshmedeski/sesh/utils"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "lang",
				Value: "english",
				Usage: "language for the greeting",
			},
		},
		Name:  "sesh",
		Usage: "Smart session manager for the terminal",
		Commands: []*cli.Command{
			{
				Name:    "isTmuxRunning",
				Aliases: []string{"itr"},
				Usage:   "Determines if tmux is running",
				Action: func(*cli.Context) error {
					isTmuxRunning, err := utils.IsTmuxRunning()
					if err != nil {
						return err
					}
					if isTmuxRunning {
						log.Println("Tmux is running")
					} else {
						log.Println("Tmux is not running")
					}
					return nil
				},
			},
			{
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
				Action: func(ctx *cli.Context) error {
					utils.ListSessions(ctx)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
