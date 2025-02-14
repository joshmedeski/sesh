package seshcli

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/namer"
	cli "github.com/urfave/cli/v2"
)

func Root(l lister.Lister, n namer.Namer) *cli.Command {
	return &cli.Command{
		Name:                   "root",
		Aliases:                []string{"r"},
		Usage:                  "Show the root from the active session",
		UseShortOptionHandling: true,
		Flags:                  []cli.Flag{},
		Action: func(cCtx *cli.Context) error {
			session, exists := l.GetAttachedTmuxSession()
			if !exists {
				return cli.Exit("No root found for session", 1)
			}
			root, err := n.RootName(session.Path)
			if err != nil {
				return cli.Exit(err, 1)
			}
			fmt.Print(root)
			return nil
		},
	}
}
