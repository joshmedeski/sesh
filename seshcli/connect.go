package seshcli

import (
	"errors"
	"strings"

	"github.com/joshmedeski/sesh/connector"
	"github.com/joshmedeski/sesh/icon"
	"github.com/joshmedeski/sesh/model"
	cli "github.com/urfave/cli/v2"
)

func Connect(c connector.Connector, i icon.Icon) *cli.Command {
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
			&cli.BoolFlag{
				Name:    "tmuxinator",
				Aliases: []string{"T"},
				Usage:   "Use tmuxinator to start session if it doesnt exist",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() == 0 {
				return errors.New("please provide a session name")
			}
			name := strings.Join(cCtx.Args().Slice(), " ")
			if name == "" {
				return nil
			}
      opts := model.ConnectOpts{Switch: cCtx.Bool("switch"), Command: cCtx.String("command"), Tmuxinator: cCtx.Bool("tmuxinator")}
			trimmedName := i.RemoveIcon(name)
			if _, err := c.Connect(trimmedName, opts); err != nil {
				// TODO: add to logging
				return err
			} else {
				// TODO: add to logging
				return nil
			}
		},
	}
}
