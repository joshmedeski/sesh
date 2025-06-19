package seshcli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/previewer"
)

func NewPreviewCommand(p previewer.Previewer) *cobra.Command {
	return &cobra.Command{
		Use:     "preview",
		Aliases: []string{"p"},
		Short:   "Preview a session or directory",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("session name or directory is required")
			}

			name := args[0]

			output, err := p.Preview(name)
			if err != nil {
				return err
			}

			fmt.Print(output)

			return nil
		},
	}
}
