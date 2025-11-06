package seshcli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/namer"
)

func NewRootSessionCommand(l lister.Lister, n namer.Namer) *cobra.Command {
	return &cobra.Command{
		Use:     "root",
		Aliases: []string{"r"},
		Short:   "Show the root from the active session",
		RunE: func(cmd *cobra.Command, args []string) error {
			session, exists := l.GetAttachedTmuxSession()
			if !exists {
				return fmt.Errorf("No root found for session")
			}
			root, err := n.RootName(session.Path)
			if err != nil {
				return err
			}
			fmt.Print(root)
			return nil
		},
	}
}
