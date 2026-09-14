package seshcli

import (
	"errors"
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
			if deps.CachingLister != nil {
				defer deps.CachingLister.Wait()
			}

			config, _ := cmd.Flags().GetBool("config")
			jsonOutput, _ := cmd.Flags().GetBool("json")
			tmux, _ := cmd.Flags().GetBool("tmux")
			zoxide, _ := cmd.Flags().GetBool("zoxide")
			hideAttached, _ := cmd.Flags().GetBool("hide-attached")
			icons, _ := cmd.Flags().GetBool("icons")
			noColor, _ := cmd.Flags().GetBool("no-color")
			tmuxinator, _ := cmd.Flags().GetBool("tmuxinator")
			hideDuplicates, _ := cmd.Flags().GetBool("hide-duplicates")
			panes, _ := cmd.Flags().GetBool("panes")
			blacklisted, _ := cmd.Flags().GetBool("blacklisted")
			listFormat, _ := cmd.Flags().GetString("format")
			formatChanged := cmd.Flags().Changed("format")
			iconExcludes, _ := cmd.Flags().GetStringSlice("icons-exclude")

			if jsonOutput && formatChanged {
				return errors.New("--format cannot be used with --json")
			}

			if panes && !deps.Tmux.IsAttached() {
				return errors.New("--panes requires being inside a tmux session")
			}

			sessions, err := deps.Lister.List(lister.ListOptions{
				Config:         config,
				HideAttached:   hideAttached,
				Icons:          icons,
				NoColor:        noColor,
				Json:           jsonOutput,
				Tmux:           tmux,
				Zoxide:         zoxide,
				Tmuxinator:     tmuxinator,
				HideDuplicates: hideDuplicates,
				Panes:          panes,
				Blacklisted:    blacklisted,
				Format:         listFormat,
				FormatSet:      formatChanged,
				IconExcludes:   iconExcludes,
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
				fmt.Println(sessions.Directory[i].Name)
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
	cmd.Flags().BoolP("no-color", "n", false, "show icons without color (requires --icons)")
	cmd.Flags().BoolP("tmuxinator", "T", false, "show tmuxinator configs")
	cmd.Flags().BoolP("hide-duplicates", "d", false, "hide duplicate entries")
	cmd.Flags().BoolP("panes", "p", false, "show panes in current session")
	cmd.Flags().BoolP("blacklisted", "b", false, "show blacklisted sessions")
	cmd.Flags().String("format", "", "format each session row ({name}, {session}, {source}, {path}, {active_window_name}, {active_window_name_prefix}; colors: {fg:yellow}...{/fg})")
	cmd.Flags().StringSlice("icons-exclude", nil, "don't show source icons for these session sources (for example: tmux)")

	return cmd
}
