package cmds

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/connect"
	"github.com/joshmedeski/sesh/session"

	"github.com/urfave/cli/v2"
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

			sessions := session.List(session.Srcs{
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
			config := config.ParseConfigFile()
			connect.Connect(choice, false, "", &config)
			return nil
		},
	}
}
