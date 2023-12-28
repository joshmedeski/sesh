package seshcli

import (
	"joshmedeski/sesh/cmds"

	"github.com/urfave/cli/v2"
)

func App() cli.App {
	return cli.App{
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
			cmds.ListSessions(),
		},
	}
}
