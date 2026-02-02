package seshcli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewLastCommand(base *BaseDeps) *cobra.Command {
	return &cobra.Command{
		Use:     "last",
		Aliases: []string{"L"},
		Short:   "Connect to the last tmux session",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			lastSession, exists := deps.Lister.GetLastTmuxSession()
			if !exists {
				// TODO: silently fail?
				return fmt.Errorf("No last session found")
			}
			base.Tmux.SwitchClient(lastSession.Name)
			return nil
		},
	}
}
