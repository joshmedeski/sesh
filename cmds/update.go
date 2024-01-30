package cmds

import (
	"errors"

	db "github.com/joshmedeski/sesh/database"
	"github.com/urfave/cli/v2"
)

func Update(storage db.Storage) *cli.Command {
	return &cli.Command{
		Name:                   "update",
		Aliases:                []string{"u"},
		Usage:                  "Update a session entry",
		Args:                   true,
		UseShortOptionHandling: true,
		Action: func(ctx *cli.Context) error {
			args := ctx.Args().Slice()
			if len(args) != 3 {
				return errors.New("Invalida args") // TODO improve error message
			}
			return storage.UpdateEntry(args[0], args[1], args[2])
		},
	}
}
