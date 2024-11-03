package seshcli

import (
	"fmt"

	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/namer"
	"github.com/joshmedeski/sesh/git"
	cli "github.com/urfave/cli/v2"
)

func Root(l lister.Lister, n namer.Namer, git git.Git) *cli.Command {
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
			_, path, err := git.GitMainWorktree(session.Path)
			root, err := n.RootName(path)
			if err != nil {
				return cli.Exit(err, 1)
			}
			fmt.Print(root)
			return nil
		},
	}
}
