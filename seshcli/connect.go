package seshcli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

func NewConnectCommand(base *BaseDeps) *cobra.Command {
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

			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			switchFlag, _ := cmd.Flags().GetBool("switch")
			command, _ := cmd.Flags().GetString("command")
			tmuxinator, _ := cmd.Flags().GetBool("tmuxinator")
			root, _ := cmd.Flags().GetBool("root")
			background, _ := cmd.Flags().GetBool("background")

			if root {
				hasRootDir, rootDir := base.Dir.RootDir(name)
				if hasRootDir {
					name = rootDir
				}
			}

			opts := model.ConnectOpts{Switch: switchFlag, Command: command, Tmuxinator: tmuxinator, Background: background}
			trimmedName := deps.Icon.RemoveIcon(name)
			session, err := deps.Connector.Connect(trimmedName, opts)
			if err != nil {
				// TODO: add to logging
				return err
			}
			if background {
				fmt.Println(session)
			}
			// Refresh cache in background so next sesh list has fresh data
			if deps.CachingLister != nil {
				deps.CachingLister.RefreshCache(lister.ListOptions{})
				deps.CachingLister.Wait()
			}
			return nil
		},
	}

	cmd.Flags().BoolP("switch", "s", false, "Switch the session (rather than attach). This is useful for actions triggered outside the terminal.")
	cmd.Flags().StringP("command", "c", "", "Execute a command when connecting to a new session. Will be ignored if the session exists.")
	cmd.Flags().BoolP("tmuxinator", "T", false, "Use tmuxinator to start session if it doesnt exist")
	cmd.Flags().BoolP("root", "r", false, "Switches to the root of the current session")
	cmd.Flags().BoolP("background", "b", false, "create without selecting it or moving the client")

	return cmd
}
