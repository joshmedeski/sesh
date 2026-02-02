package seshcli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

func NewListCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			config, _ := cmd.Flags().GetBool("config")
			jsonOutput, _ := cmd.Flags().GetBool("json")
			tmux, _ := cmd.Flags().GetBool("tmux")
			zoxide, _ := cmd.Flags().GetBool("zoxide")
			hideAttached, _ := cmd.Flags().GetBool("hide-attached")
			icons, _ := cmd.Flags().GetBool("icons")
			tmuxinator, _ := cmd.Flags().GetBool("tmuxinator")
			hideDuplicates, _ := cmd.Flags().GetBool("hide-duplicates")

			sessions, err := deps.Lister.List(lister.ListOptions{
				Config:         config,
				HideAttached:   hideAttached,
				Icons:          icons,
				Json:           jsonOutput,
				Tmux:           tmux,
				Zoxide:         zoxide,
				Tmuxinator:     tmuxinator,
				HideDuplicates: hideDuplicates,
			})
			if err != nil {
				return fmt.Errorf("couldn't list sessions: %q", err)
			}

			if jsonOutput {
				var sessionsArray []model.SeshSession
				for _, i := range sessions.OrderedIndex {
					sessionsArray = append(sessionsArray, sessions.Directory[i])
				}
				fmt.Println(base.Json.EncodeSessions(sessionsArray))
				return nil
			}

			for _, i := range sessions.OrderedIndex {
				name := sessions.Directory[i].Name
				if icons {
					name = deps.Icon.AddIcon(sessions.Directory[i])
				}
				fmt.Println(name)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("config", "c", false, "show configured sessions")
	cmd.Flags().BoolP("json", "j", false, "output as json")
	cmd.Flags().BoolP("tmux", "t", false, "show tmux sessions")
	cmd.Flags().BoolP("zoxide", "z", false, "show zoxide results")
	cmd.Flags().BoolP("hide-attached", "H", false, "don't show currently attached sessions")
	cmd.Flags().BoolP("icons", "i", false, "show icons")
	cmd.Flags().BoolP("tmuxinator", "T", false, "show tmuxinator configs")
	cmd.Flags().BoolP("hide-duplicates", "d", false, "hide duplicate entries")

	return cmd
}
