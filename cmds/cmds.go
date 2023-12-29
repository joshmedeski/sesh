package cmds

import (
	"fmt"
	"joshmedeski/sesh/tmux"
	"joshmedeski/sesh/zoxide"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func ListSessions() *cli.Command {
	return &cli.Command{
		Name:                   "list",
		Aliases:                []string{"l"},
		Usage:                  "List sessions",
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
			fmt.Println(
				strings.Join(requestedSessions(cCtx), "\n"),
			)
			return nil
		},
	}
}

func ChooseSession() *cli.Command {
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

			stdin.Write([]byte(
				strings.Join(requestedSessions(cCtx), "\n"),
			))
			err = stdin.Close()
			if err != nil {
				return err
			}

			return cmd.Wait()
		},
	}
}

func requestedSessions(cCtx *cli.Context) []string {
	var sessions []string
	hasFlags := cCtx.Bool("tmux") || cCtx.Bool("zoxide")

	if !hasFlags || cCtx.Bool("tmux") {
		tmuxSessions, err := tmux.Sessions()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		sessions = append(sessions, tmuxSessions...)
	}

	if !hasFlags || cCtx.Bool("zoxide") {
		dirs, err := zoxide.Dirs()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		sessions = append(sessions, dirs...)
	}
	return sessions
}
