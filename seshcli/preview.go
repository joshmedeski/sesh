package seshcli

import (
	"errors"
	"fmt"

	"github.com/joshmedeski/sesh/v2/previewer"
	cli "github.com/urfave/cli/v2"
)

func Preview(p previewer.Previewer) *cli.Command {
	return &cli.Command{
		Name:                   "preview",
		Aliases:                []string{"p"},
		Usage:                  "Preview a session or directory",
		UseShortOptionHandling: true,
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 1 {
				return errors.New("session name or directory is required")
			}

			name := cCtx.Args().First()

			output, err := p.Preview(name)
			if err != nil {
				return cli.Exit(err, 1)
			}

			fmt.Print(output)

			return nil
		},
	}
}
