package seshcli

import (
	cli "github.com/urfave/cli/v2"
)

func Clone() *cli.Command {
	return &cli.Command{
		Name:                   "clone",
		Aliases:                []string{"cl"},
		Usage:                  "Clone a git repo and connect to it as a session",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cmdDir",
				Aliases: []string{"d"},
				Usage:   "The directory to run the git command in",
			},
		},
		Action: func(cCtx *cli.Context) error {
			return nil
		},
	}
}
