package seshcli

import (
	"github.com/spf13/cobra"

	"github.com/Wingsdh/cc-sesh/v2/lister"
	"github.com/Wingsdh/cc-sesh/v2/model"
	"github.com/Wingsdh/cc-sesh/v2/picker"
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
			tmuxinator, _ := cmd.Flags().GetBool("tmuxinator")
			hideDuplicates, _ := cmd.Flags().GetBool("hide-duplicates")

			listerOpts := lister.ListOptions{
				Config:         config,
				HideAttached:   hideAttached,
				Tmux:           tmux,
				Zoxide:         zoxide,
				Tmuxinator:     tmuxinator,
				HideDuplicates: hideDuplicates,
			}
			fetchFunc := func() (model.SeshSessions, error) {
				return deps.Lister.List(listerOpts)
			}

			var pickerOpts picker.PickerOptions
			if cmd.Flags().Changed("icons") {
				showIcons := true
				pickerOpts.ShowIcons = &showIcons
			} else {
				showIcons := deps.Config.TUI.ShowIcons
				pickerOpts.ShowIcons = &showIcons
			}
			if cmd.Flags().Changed("separator-aware") {
				separatorAware := true
				pickerOpts.SeparatorAware = &separatorAware
			}
			if cmd.Flags().Changed("prompt") {
				prompt, _ := cmd.Flags().GetString("prompt")
				pickerOpts.Prompt = &prompt
			} else if deps.Config.TUI.Prompt != "" {
				pickerOpts.Prompt = &deps.Config.TUI.Prompt
			}
			if cmd.Flags().Changed("placeholder") {
				placeholder, _ := cmd.Flags().GetString("placeholder")
				pickerOpts.Placeholder = &placeholder
			} else if deps.Config.TUI.Placeholder != "" {
				pickerOpts.Placeholder = &deps.Config.TUI.Placeholder
			}

			chosen, err := deps.Picker.Pick(fetchFunc, pickerOpts)
			if err != nil {
				return err
			}

			if chosen == "" {
				return nil
			}

			if _, err := deps.Connector.Connect(chosen, model.ConnectOpts{}); err != nil {
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
	cmd.Flags().StringP("prompt", "p", "", "prompt shown in the picker TUI")
	cmd.Flags().String("placeholder", "", "placeholder text in the picker TUI")

	return cmd
}
