package seshcli

import (
	"github.com/joshmedeski/sesh/cmds"

	db "github.com/joshmedeski/sesh/database"
	"github.com/urfave/cli/v2"
)

func App(version string, storage db.Storage) cli.App {
	return cli.App{
		Name:    "sesh",
		Version: version,
		Usage:   "Smart session manager for the terminal",
		Commands: []*cli.Command{
			cmds.Add(storage),
			cmds.Delete(storage),
			cmds.Update(storage),
			cmds.List(),
			cmds.Choose(),
			cmds.Connect(),
			cmds.Clone(),
		},
	}
}
