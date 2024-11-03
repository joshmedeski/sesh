package seshcli

import (
	"fmt"

	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/home"
	cli "github.com/urfave/cli/v2"
)

func Root(l lister.Lister, git git.Git, home home.Home) *cli.Command {
	return &cli.Command{
		Name:                   "root",
		Aliases:                []string{"r"},
		Usage:                  "Show the root from the active session",
		UseShortOptionHandling: true,
		Flags:                  []cli.Flag{},
		Action: func(cCtx *cli.Context) error {
			session, exists := l.GetAttachedTmuxSession()
			if !exists {
				return cli.Exit("Not attached to tmux session", 1)
			}
			_, path, err := git.GitMainWorktree(session.Path)
			if err != nil {
				return cli.Exit(err, 1)
			}
			root, err := home.ShortenHome(path)
			if err != nil {
				return cli.Exit(err, 1)
			}
			fmt.Print(root)
			return nil
		},
	}
}
