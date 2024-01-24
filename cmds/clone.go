package cmds

import (
	cli "github.com/urfave/cli/v2"

	"github.com/joshmedeski/sesh/connect"
	"github.com/joshmedeski/sesh/git"
)

func Clone() *cli.Command {
	return &cli.Command{
		Name:                   "clone",
		Aliases:                []string{"cl"},
		Usage:                  "Clone a git repo and connect to it as a session",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cmdDir",
				Aliases: []string{"d"},
				Usage:   "The directory to run the git command in",
			},
		},
		Action: func(cCtx *cli.Context) error {
			repo := cCtx.Args().First()
			dir := cCtx.Args().Get(1)
			cmdDir := cCtx.String("cmdDir")
			c, err := git.Clone(git.CloneOptions{
				Dir:    &dir,
				CmdDir: &cmdDir,
				Repo:   repo,
			})
			if err != nil {
				return cli.Exit(err, 1)
			}

			return connect.Connect(c.Path, false, "")
		},
	}
}
