package seshcli

import (
	"errors"

	"github.com/joshmedeski/sesh/v2/cloner"
	"github.com/joshmedeski/sesh/v2/model"
	cli "github.com/urfave/cli/v2"
)

func Clone(c cloner.Cloner) *cli.Command {
	return &cli.Command{
		Name:                   "clone",
		Aliases:                []string{"cl"},
		Usage:                  "Clone a git repo and connect to it as a session",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cmdDir",
				Aliases: []string{"c"},
				Usage:   "The directory to run the git command in",
			},
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   "The name of the directory that git is creating",
			},
		},
		Action: func(cCtx *cli.Context) error {

			if cCtx.NArg() != 1 {
				return errors.New("please provide url to clone")
			}
			repo := cCtx.Args().First()

			opts := model.GitCloneOptions{CmdDir: cCtx.String("cmdDir"), Repo: repo, Dir: cCtx.String("dir")}
			if _, err := c.Clone(opts); err != nil {
				return err
			} else {
				return nil
			}
		},
	}
}
