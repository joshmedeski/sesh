package seshcli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/tmux"
)

func NewLastCommand(l lister.Lister, t tmux.Tmux) *cobra.Command {
	return &cobra.Command{
		Use:     "last",
		Aliases: []string{"L"},
		Short:   "Connect to the last tmux session",
		RunE: func(cmd *cobra.Command, args []string) error {
			lastSession, exists := l.GetLastTmuxSession()
			if !exists {
				// TODO: silently fail?
				return fmt.Errorf("No last session found")
			}
			t.SwitchClient(lastSession.Name)
			return nil
		},
	}
}
