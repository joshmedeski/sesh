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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
