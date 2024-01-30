package cmds

import (
	"errors"

	db "github.com/joshmedeski/sesh/database"
	"github.com/urfave/cli/v2"
)

func Add(storage db.Storage) *cli.Command {
	return &cli.Command{
		Name:                   "add",
		Aliases:                []string{"a"},
		Usage:                  "Add a new session",
		Args:                   true,
		UseShortOptionHandling: true,
		Action: func(ctx *cli.Context) error {
			args := ctx.Args().Slice()
			if len(args) != 2 {
				return errors.New("Name and Path needed")
			}
			return storage.CreateEntry(&db.Entry{
				Name: args[0],
				Path: args[1],
			})
		},
	}
}
