package seshcli

import (
	"errors"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/dashboard"
	"github.com/joshmedeski/sesh/v2/model"
)

func NewDashboardCommand(base *BaseDeps) *cobra.Command {
	return &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"dash", "d"},
		Short:   "Full-screen session dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			// check if we're inside a tmux session
			if !deps.Tmux.IsAttached() {
				return errors.New("dashboard requires being inside a tmux session")
			}

			// check user home dir
			homeDir, err := deps.Os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("couldn't get home directory: %w", err)
			}

			m := dashboard.New(deps.Config.Dashboard, deps.Tmux, deps.Lister, deps.Git, deps.Connector, deps.Shell, homeDir)
			prog := tea.NewProgram(m)
			result, err := prog.Run()
			if err != nil {
				return fmt.Errorf("dashboard error: %w", err)
			}

			dashModel, ok := result.(dashboard.Model)
			if !ok {
				return errors.New("unexpected model type")
			}

			if dashModel.Quit() {
				return nil
			}

			if chosen := dashModel.Chosen(); chosen != "" {
				if _, err := deps.Connector.Connect(chosen, model.ConnectOpts{}); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
