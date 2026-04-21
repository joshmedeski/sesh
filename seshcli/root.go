package seshcli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootSessionCommand(base *BaseDeps) *cobra.Command {
	return &cobra.Command{
		Use:     "root",
		Aliases: []string{"r"},
		Short:   "Show the root from the active session",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			var sessionPath string
			var exists bool

			if deps.Config.Backend == "wezterm" {
				session, ok := deps.Lister.GetActiveWeztermWorkspace()
				exists = ok
				sessionPath = session.Path
			} else {
				session, ok := deps.Lister.GetAttachedTmuxSession()
				exists = ok
				sessionPath = session.Path
			}

			if !exists {
				return fmt.Errorf("No root found for session")
			}
			root, err := deps.Namer.RootName(sessionPath)
			if err != nil {
				return err
			}
			fmt.Print(root)
			return nil
		},
	}
}
