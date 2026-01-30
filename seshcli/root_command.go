package seshcli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(version string) *cobra.Command {
	base := NewBaseDeps()

	rootCmd := &cobra.Command{
		Use:              "sesh",
		Version:          version,
		Short:            "Smart session manager for the terminal",
		Long:             "Sesh is a smart terminal session manager that helps you create and manage tmux sessions quickly and easily using zoxide.",
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
	)

	return rootCmd
}
