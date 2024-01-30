package cmds

import (
	"errors"

	db "github.com/joshmedeski/sesh/database"
	"github.com/urfave/cli/v2"
)

func Delete(storage db.Storage) *cli.Command {
	return &cli.Command{
		Name:                   "delete",
		Aliases:                []string{"d"},
		Usage:                  "Delete a session entry",
		Args:                   true,
		UseShortOptionHandling: true,
		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			if name == "" {
				return errors.New("No name provided")
			}
			return storage.DeleteEntry(name)
		},
	}
}
