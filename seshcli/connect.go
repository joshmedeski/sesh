package seshcli

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/model"
)

func NewConnectCommand(c connector.Connector, i icon.Icon, d dir.Dir) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connect",
		Aliases: []string{"cn"},
		Short:   "Connect to the given session",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("please provide a session name")
			}
			name := strings.Join(args, " ")
			if name == "" {
				return nil
			}

			switchFlag, _ := cmd.Flags().GetBool("switch")
			command, _ := cmd.Flags().GetString("command")
			tmuxinator, _ := cmd.Flags().GetBool("tmuxinator")
			root, _ := cmd.Flags().GetBool("root")

			if root {
				hasRootDir, rootDir := d.RootDir(name)
				if hasRootDir {
					name = rootDir
				}
			}

			opts := model.ConnectOpts{Switch: switchFlag, Command: command, Tmuxinator: tmuxinator}
			trimmedName := i.RemoveIcon(name)
			if _, err := c.Connect(trimmedName, opts); err != nil {
				// TODO: add to logging
				return err
			} else {
				// TODO: add to logging
				return nil
			}
		},
	}

	cmd.Flags().BoolP("switch", "s", false, "Switch the session (rather than attach). This is useful for actions triggered outside the terminal.")
	cmd.Flags().StringP("command", "c", "", "Execute a command when connecting to a new session. Will be ignored if the session exists.")
	cmd.Flags().BoolP("tmuxinator", "T", false, "Use tmuxinator to start session if it doesnt exist")
	cmd.Flags().BoolP("root", "r", false, "Switches to the root of the current session")

	return cmd
}
