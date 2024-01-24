package cmds

import (
	cli "github.com/urfave/cli/v2"

	"github.com/joshmedeski/sesh/connect"
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
			&cli.StringFlag{
				Name:    "command",
				Aliases: []string{"c"},
				Usage:   "Execute a command when connecting to a new session. Will be ignored if the session exists.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			session := cCtx.Args().First()
			alwaysSwitch := cCtx.Bool("switch")
			command := cCtx.String("command")
			if session == "" {
				return cli.Exit("No session provided", 0)
			}
			return connect.Connect(session, alwaysSwitch, command)
		},
	}
}
