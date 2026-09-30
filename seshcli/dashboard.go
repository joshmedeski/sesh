package seshcli

import (
	"errors"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/dashboard"
	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/lister"
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

			m := dashboard.New(deps.Config, dashboard.SectionDeps{
				Tmux:      deps.Tmux,
				Lister:    deps.Lister,
				Git:       deps.Git,
				Connector: deps.Connector,
				Shell:     deps.Shell,
				Runner:    dashboard.NewCommandRunner(deps.Exec),
				Worktree:  deps.Worktree,
				Home:      deps.Home,
				Icon:      lister.IconResolver(deps.Config, deps.Home, deps.Lister),
				IconWidth: icon.ColumnWidth(deps.Config),
			})
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

			if opts := dashModel.ChosenWorktree(); opts != nil {
				if _, err := deps.Worktree.Connect(*opts); err != nil {
					return err
				}
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
