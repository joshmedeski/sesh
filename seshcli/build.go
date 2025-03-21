package seshcli

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/builder"
	"github.com/joshmedeski/sesh/v2/lister"
	cli "github.com/urfave/cli/v2"
)

func Build(l lister.Lister, b builder.Builder) *cli.Command {
	return &cli.Command{
		Name:                   "build",
		Aliases:                []string{"b"},
		Usage:                  "Builds the current session (Experimental)",
		UseShortOptionHandling: true,
		Flags:                  []cli.Flag{},
		Action: func(cCtx *cli.Context) error {
			session, exists := l.GetAttachedTmuxSession()
			if !exists {
				return cli.Exit("No tmux session currently attached", 1)
			}
			out, err := b.Build(session)
			if err != nil {
				return cli.Exit(err, 1)
			}
			fmt.Print(out)
			return nil
		},
	}
}
