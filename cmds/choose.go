package cmds

import (
	"joshmedeski/sesh/session"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func Choose() *cli.Command {
	return &cli.Command{
		Name:                   "choose",
		Aliases:                []string{"c"},
		Usage:                  "Select session",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "tmux",
				Aliases: []string{"t"},
				Usage:   "show tmux sessions",
			},
			&cli.BoolFlag{
				Name:    "zoxide",
				Aliases: []string{"z"},
				Usage:   "show zoxide results",
			},
		},
		Action: func(cCtx *cli.Context) error {
			cmd := exec.Command("fzf")
			stdin, err := cmd.StdinPipe()
			if err != nil {
				return err
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Start()
			if err != nil {
				return err
			}

			sessions := session.Sessions(session.Srcs{
				Tmux:   cCtx.Bool("tmux"),
				Zoxide: cCtx.Bool("zoxide"),
			})
			stdin.Write([]byte(
				strings.Join(sessions, "\n"),
			))
			err = stdin.Close()
			if err != nil {
				return err
			}

			return cmd.Wait()
		},
	}
}
