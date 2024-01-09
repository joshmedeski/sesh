package cmds

import (
	"joshmedeski/sesh/connect"

	"github.com/urfave/cli/v2"
)

func Connect() *cli.Command {
	return &cli.Command{
		Name:                   "connect",
		Aliases:                []string{"cn"},
		Usage:                  "Connect to the given session",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "tmux",
				Aliases: []string{"t"},
				Usage:   "connect to tmux session",
			},
		},
		Action: func(cCtx *cli.Context) error {
			session := cCtx.Args().First()
			if session == "" {
				return cli.Exit("No session provided", 0)
			}
			connect.Connect(session)
			return nil
		},
	}
}
