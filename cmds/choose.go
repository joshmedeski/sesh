package cmds

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/joshmedeski/sesh/connect"
	"github.com/joshmedeski/sesh/session"
)

func Choose() *cli.Command {
	return &cli.Command{
		Name:                   "choose",
		Aliases:                []string{"ch"},
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
			&cli.BoolFlag{
				Name:    "hide-attached",
				Aliases: []string{"H"},
				Usage:   "don't show currently attached sessions",
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

			o := session.Options{
				HideAttached: cCtx.Bool("hide-attached"),
			}
			sessions := session.List(o, session.Srcs{
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
			choice := strings.TrimSpace(cmdOutput.String())
			// TODO: get choice from Session structs array
			return connect.Connect(choice, false, "")
		},
	}
}
