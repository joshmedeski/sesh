package db

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func (c *SqliteDatabase) Delete() *cli.Command {
	return &cli.Command{
		Name:                   "delete",
		Aliases:                []string{"d"},
		Usage:                  "delete a session entry (?)",
		Args:                   true,
		UseShortOptionHandling: true,
		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			if name == "" {
				return errors.New("No name provided")
			}
			return c.DeleteEntry(name)
		},
	}
}
