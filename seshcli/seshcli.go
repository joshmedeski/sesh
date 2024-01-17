package seshcli

import (
	"github.com/joshmedeski/sesh/cmds"

	"github.com/urfave/cli/v2"
)

var version = "dev"

func App() cli.App {
	return cli.App{
		Name:    "sesh",
		Version: version,
		Usage:   "Smart session manager for the terminal",
		Commands: []*cli.Command{
			cmds.List(),
			cmds.Choose(),
			cmds.Connect(),
			cmds.Clone(),
		},
	}
}
