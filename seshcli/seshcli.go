package seshcli

import (
	"github.com/urfave/cli/v2"
)

func App(version string) cli.App {
	return cli.App{
		Name:    "sesh",
		Version: version,
		Usage:   "Smart session manager for the terminal",
		Commands: []*cli.Command{
			List(),
			Connect(),
			Clone(),
		},
	}
}
