package db

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func (c *SqliteDatabase) Add() *cli.Command {
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
			return c.CreateEntry(&Entry{
				Name: args[0],
				Path: args[1],
			})
		},
	}
}
