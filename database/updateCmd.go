package db

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func (c *SqliteDatabase) Update() *cli.Command {
	return &cli.Command{
		Name:                   "update",
		Aliases:                []string{"u"},
		Usage:                  "update a session entry (?)",
		Args:                   true,
		UseShortOptionHandling: true,
		Action: func(ctx *cli.Context) error {
			args := ctx.Args().Slice()
			if len(args) != 3 {
				return errors.New("Invalida args") // TODO improve error message
			}
			return c.UpdateEntry(args[0], args[1], args[2])
		},
	}
}
