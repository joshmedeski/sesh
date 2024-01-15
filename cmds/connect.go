package cmds

import (
	"github.com/joshmedeski/sesh/connect"

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
				Name:    "switch",
				Aliases: []string{"s"},
				Usage:   "Always switch the session (and never attach). This is useful for third-party tools like Raycast.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			session := cCtx.Args().First()
			alwaysSwitch := cCtx.Bool("switch")
			if session == "" {
				return cli.Exit("No session provided", 0)
			}
			connect.Connect(session, alwaysSwitch)
			return nil
		},
	}
}
