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

			session, exists := deps.Lister.GetAttachedTmuxSession()
			if !exists {
				return fmt.Errorf("No root found for session")
			}
			root, err := deps.Namer.RootName(session.Path)
			if err != nil {
				return err
			}
			fmt.Print(root)
			return nil
		},
	}
}
