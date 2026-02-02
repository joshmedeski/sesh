package seshcli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func NewPreviewCommand(base *BaseDeps) *cobra.Command {
	return &cobra.Command{
		Use:     "preview",
		Aliases: []string{"p"},
		Short:   "Preview a session or directory",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("session name or directory is required")
			}

			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			name := args[0]

			output, err := deps.Previewer.Preview(name)
			if err != nil {
				return err
			}

			fmt.Print(output)

			return nil
		},
	}
}
