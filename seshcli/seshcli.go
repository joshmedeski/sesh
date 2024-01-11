package seshcli

import (
	"joshmedeski/sesh/cmds"

	"github.com/urfave/cli/v2"
)

func App() cli.App {
	return cli.App{
		Name:  "sesh",
		Usage: "Smart session manager for the terminal",
		Commands: []*cli.Command{
			cmds.List(),
			cmds.Choose(),
			cmds.Connect(),
		},
	}
}
