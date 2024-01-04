package cmds

import (
	"bytes"
	"joshmedeski/sesh/connect"
	"joshmedeski/sesh/session"
	"log"
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

			var cmdOutput bytes.Buffer
			cmd.Stdout = &cmdOutput

			stdin, err := cmd.StdinPipe()
			if err != nil {
				return err
			}
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

			err = cmd.Wait()
			if err != nil {
				log.Fatal(err)
			}
			choice := cmdOutput.String()
			connect.Connect(choice)
			return nil
		},
	}
}
