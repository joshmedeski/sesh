package seshcli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(version string) *cobra.Command {
	base := NewBaseDeps()

	rootCmd := &cobra.Command{
		Use:              "cc-sesh",
		Version:          version,
		Short:            "Smart session manager for the terminal, with Claude Code awareness",
		Long:             "cc-sesh is a fork of sesh that adds Claude Code state badges (busy / idle / needs-input / sub-agent) to the session picker, on top of all original sesh functionality.",
		TraverseChildren: true,
	}

	rootCmd.PersistentFlags().StringP("config", "C", "", "path to config file")

	rootCmd.AddCommand(
		NewListCommand(base),
		NewLastCommand(base),
		NewConnectCommand(base),
		NewCloneCommand(base),
		NewRootSessionCommand(base),
		NewPreviewCommand(base),
		NewPickerCommand(base),
		NewWindowCommand(base),
	)

	return rootCmd
}
