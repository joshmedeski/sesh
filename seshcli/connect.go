package seshcli

import (
	"errors"
	"fmt"

	"github.com/joshmedeski/sesh/connector"
	"github.com/joshmedeski/sesh/model"
	cli "github.com/urfave/cli/v2"
)

func Connect(c connector.Connector) *cli.Command {
	return &cli.Command{
		Name:                   "connect",
		Aliases:                []string{"cn"},
		Usage:                  "Connect to the given session",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "switch",
				Aliases: []string{"s"},
				Usage:   "Switch the session (rather than attach). This is useful for actions triggered outside the terminal.",
			},
			&cli.StringFlag{
				Name:    "command",
				Aliases: []string{"c"},
				Usage:   "Execute a command when connecting to a new session. Will be ignored if the session exists.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() == 0 {
				return errors.New("please provide a session name")
			}
			name := cCtx.Args().First()
			if name == "" {
				return nil
			}
			opts := model.ConnectOpts{Switch: cCtx.Bool("switch"), Command: cCtx.String("command")}
			if connection, err := c.Connect(name, opts); err != nil {
				// TODO: print to logs?
				return err
			} else {
				// TODO: create a message that is helpful to the end user
				fmt.Println(connection)
				return nil
			}
		},
	}
}
