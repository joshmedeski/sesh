package seshcli

import (
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/tmux"
	cli "github.com/urfave/cli/v2"
)

func Last(l lister.Lister, t tmux.Tmux) *cli.Command {
	return &cli.Command{
		Name:                   "last",
		Aliases:                []string{"L"},
		Usage:                  "Connect to the last tmux session",
		UseShortOptionHandling: true,
		Flags:                  []cli.Flag{},
		Action: func(cCtx *cli.Context) error {
			lastSession, exists := l.GetLastTmuxSession()
			if !exists {
				// TODO: silently fail?
				return cli.Exit("No last session found", 1)
			}
			t.SwitchClient(lastSession.Name)
			return nil
		},
	}
}
