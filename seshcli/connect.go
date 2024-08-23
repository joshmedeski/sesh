package seshcli

import (
	"errors"
	"fmt"

	"github.com/joshmedeski/sesh/connector"
	"github.com/joshmedeski/sesh/icon"
	"github.com/joshmedeski/sesh/model"
	cli "github.com/urfave/cli/v2"
)

//	func tmuxCmd(args []string) (string, error) {
//		tmux, err := exec.LookPath("tmux")
//		if err != nil {
//			return "", err
//		}
//		var stdout, stderr bytes.Buffer
//		cmd := exec.Command(tmux, args...)
//		cmd.Stdin = os.Stdin
//		cmd.Stdout = &stdout
//		cmd.Stderr = os.Stderr
//		cmd.Stderr = &stderr
//		if err := cmd.Start(); err != nil {
//			return "", err
//		}
//		if err := cmd.Wait(); err != nil {
//			errString := strings.TrimSpace(stderr.String())
//			if strings.HasPrefix(errString, "no server running on") {
//				return "", nil
//			}
//			return "", err
//		}
//		return stdout.String(), nil
//	}
//
//	func attachSession(session string) error {
//		if _, err := tmuxCmd([]string{"attach", "-t", session}); err != nil {
//			return err
//		}
//		return nil
//	}
func Connect(c connector.Connector, i icon.Icon) *cli.Command {
	return &cli.Command{
		Name:                   "connect",
		Aliases:                []string{"cn"},
		Usage:                  "Connect to the given session",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "switch",
				Aliases: []string{"s"},
				Usage:   "Switch the session (rather than attach). This is useful for actions triggered outside the terminal.",
			},
			&cli.StringFlag{
				Name:    "command",
				Aliases: []string{"c"},
				Usage:   "Execute a command when connecting to a new session. Will be ignored if the session exists.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() == 0 {
				return errors.New("please provide a session name")
			}
			name := cCtx.Args().First()
			if name == "" {
				return nil
			}
			opts := model.ConnectOpts{Switch: cCtx.Bool("switch"), Command: cCtx.String("command")}
			trimmedName := i.RemoveIcon(name)
			fmt.Println(trimmedName)
			if connection, err := c.Connect(trimmedName, opts); err != nil {
				// TODO: print to logs?
				return err
			} else {
				// TODO: create a message that is helpful to the end user
				fmt.Println(connection)
				return nil
			}

			// attachSession("sesh/v2")
			// return nil
		},
	}
}
