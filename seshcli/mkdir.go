package seshcli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

func NewMkdirCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mkdir <path>",
		Aliases: []string{"md"},
		Short:   "Create a directory and connect to it as a session",
		Long: "Create a directory and connect to it as a session, combining `mkdir` and `sesh connect` into a single step.\n\n" +
			"The <path> can be relative (resolved against the current working directory) or absolute, and may use `~` for the home directory.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			if path == "" {
				return errors.New("please provide a directory path")
			}

			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			switchFlag, _ := cmd.Flags().GetBool("switch")
			command, _ := cmd.Flags().GetString("command")

			opts := model.ConnectOpts{Switch: switchFlag, Command: command}
			if _, err := deps.Mkdirer.Mkdir(path, opts); err != nil {
				// TODO: add to logging
				return err
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
	cmd.Flags().StringP("command", "c", "", "Execute a command when connecting to the new session.")

	return cmd
}
