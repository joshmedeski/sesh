package seshcli

import (
	"os"

	"github.com/joshmedeski/sesh/cmds"
	db "github.com/joshmedeski/sesh/database"

	"github.com/urfave/cli/v2"
)

func App(version string) cli.App {
	sqlPath := os.ExpandEnv("$HOME/.local/share/sesh/sesh.db")
	storage := db.NewSqliteDatabase(sqlPath)

	return cli.App{
		Name:    "sesh",
		Version: version,
		Usage:   "Smart session manager for the terminal",
		Commands: []*cli.Command{
			storage.Add(),
			storage.Delete(),
			cmds.List(),
			cmds.Choose(),
			cmds.Connect(),
			cmds.Clone(),
		},
	}
}
