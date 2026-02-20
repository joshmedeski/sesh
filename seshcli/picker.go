package seshcli

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/picker"
)

func NewPickerCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "picker",
		Aliases: []string{"pick", "pk"},
		Short:   "Interactive session picker",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}
			if deps.CachingLister != nil {
				defer deps.CachingLister.Wait()
			}

			config, _ := cmd.Flags().GetBool("config")
			tmux, _ := cmd.Flags().GetBool("tmux")
			zoxide, _ := cmd.Flags().GetBool("zoxide")
			hideAttached, _ := cmd.Flags().GetBool("hide-attached")
			icons, _ := cmd.Flags().GetBool("icons")
			tmuxinator, _ := cmd.Flags().GetBool("tmuxinator")
			hideDuplicates, _ := cmd.Flags().GetBool("hide-duplicates")

			separatorAware := deps.Config.SeparatorAware
			if cmd.Flags().Changed("separator-aware") {
				separatorAware, _ = cmd.Flags().GetBool("separator-aware")
			}

			opts := lister.ListOptions{
				Config:         config,
				HideAttached:   hideAttached,
				Icons:          icons,
				Tmux:           tmux,
				Zoxide:         zoxide,
				Tmuxinator:     tmuxinator,
				HideDuplicates: hideDuplicates,
			}
			fetchFunc := func() (model.SeshSessions, error) {
				return deps.Lister.List(opts)
			}

			m := picker.New(fetchFunc, icons, separatorAware)
			p := tea.NewProgram(m)
			result, err := p.Run()
			if err != nil {
				return fmt.Errorf("picker error: %w", err)
			}

			pickerModel, ok := result.(picker.Model)
			if !ok {
				return fmt.Errorf("unexpected model type")
			}

			if pickerModel.LoadErr() != nil {
				return fmt.Errorf("couldn't list sessions: %w", pickerModel.LoadErr())
			}

			if pickerModel.Quit() || pickerModel.Chosen() == "" {
				return nil
			}

			if _, err := deps.Connector.Connect(pickerModel.Chosen(), model.ConnectOpts{}); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolP("config", "c", false, "show configured sessions")
	cmd.Flags().BoolP("tmux", "t", false, "show tmux sessions")
	cmd.Flags().BoolP("zoxide", "z", false, "show zoxide results")
	cmd.Flags().BoolP("hide-attached", "H", false, "don't show currently attached sessions")
	cmd.Flags().BoolP("icons", "i", false, "show icons")
	cmd.Flags().BoolP("tmuxinator", "T", false, "show tmuxinator configs")
	cmd.Flags().BoolP("hide-duplicates", "d", false, "hide duplicate entries")
	cmd.Flags().BoolP("separator-aware", "s", false, "match spaces to separators (-_/\\)")

	return cmd
}
