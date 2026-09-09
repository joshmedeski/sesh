package seshcli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/picker"
)

func NewWorktreeCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "worktree",
		Aliases: []string{"wt"},
		Short:   "Manage git worktrees as sesh sessions",
	}
	cmd.AddCommand(newWorktreeConnectCommand(base))
	cmd.AddCommand(newWorktreeListCommand(base))
	cmd.AddCommand(newWorktreePickerCommand(base))
	return cmd
}

func newWorktreePickerCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "picker",
		Aliases: []string{"pick", "pk", "p"},
		Short:   "Interactively pick a worktree and connect to it",
		Long: `Pick a worktree from an interactive list and connect to its session.

Each row shows the issue number, a color-coded badge for its state, and its
title. Typing filters on the number and the title together, so both "409" and
"worktree support" find the same row.

Titles and states come from the same cache ` + "`sesh worktree list`" + ` fills, so the
list is instant once it is warm, and it then refetches them behind the rows on
screen — a title edited since the last listing corrects itself a moment after the
picker opens. ctrl+r asks for that again at any time; --refresh does it on the way
in, before the first row is drawn.

--path and --repo are two ways to select the same [[worktree]] block; with
neither, the repo is detected from the current directory.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			path, _ := cmd.Flags().GetString("path")
			repo, _ := cmd.Flags().GetString("repo")
			switchFlag, _ := cmd.Flags().GetBool("switch")
			refresh, _ := cmd.Flags().GetBool("refresh")

			var pickerOpts picker.WorktreePickerOptions
			// --refresh already ignored the cache on the initial load, so the
			// automatic refresh behind it would fetch the same titles twice.
			if refresh {
				autoRefresh := false
				pickerOpts.AutoRefresh = &autoRefresh
			}
			if cmd.Flags().Changed("icons") {
				showIcons, _ := cmd.Flags().GetBool("icons")
				pickerOpts.ShowIcons = &showIcons
			}
			if cmd.Flags().Changed("prompt") {
				prompt, _ := cmd.Flags().GetString("prompt")
				pickerOpts.Prompt = &prompt
			}
			if cmd.Flags().Changed("placeholder") {
				placeholder, _ := cmd.Flags().GetString("placeholder")
				pickerOpts.Placeholder = &placeholder
			}
			if cmd.Flags().Changed("query") {
				query, _ := cmd.Flags().GetString("query")
				pickerOpts.Query = &query
			}

			listOpts := model.WorktreeListOpts{Path: path, Repo: repo, Refresh: refresh}
			return pickWorktree(deps, listOpts, switchFlag, pickerOpts)
		},
	}

	cmd.Flags().String("path", "", "Worktree root to pick from, e.g. ~/c/nu/w")
	cmd.Flags().StringP("repo", "r", "", "Target repo as org/repo (overrides cwd detection)")
	cmd.Flags().Bool("refresh", false, "Refetch titles and states, ignoring the cache (ctrl+r does the same from the picker)")
	cmd.Flags().BoolP("switch", "s", false, "Switch the session rather than attach (for use outside tmux)")
	cmd.Flags().BoolP("icons", "i", false, "round the state badges off with nerd font half circles")
	cmd.Flags().StringP("prompt", "p", "", "prompt shown in the picker TUI")
	cmd.Flags().String("placeholder", "", "placeholder text in the picker TUI")
	cmd.Flags().StringP("query", "q", "", "prefill the picker's filter with this query")

	return cmd
}

// pickWorktree runs the picker over the worktrees selected by listOpts and
// connects to whatever was chosen. Quitting the picker is not an error: it is
// how the user says they changed their mind, so nothing is connected to.
func pickWorktree(
	deps *Deps,
	listOpts model.WorktreeListOpts,
	switchSession bool,
	pickerOpts picker.WorktreePickerOptions,
) error {
	// refresh is threaded through rather than baked into listOpts so ctrl+r can
	// ask for a fetch that ignores the cache on a picker that opened from it.
	fetchFunc := func(refresh bool) ([]model.WorktreeEntry, error) {
		opts := listOpts
		opts.Refresh = opts.Refresh || refresh
		return deps.Worktree.List(opts)
	}

	entry, picked, err := deps.Picker.PickWorktree(fetchFunc, pickerOpts)
	if err != nil {
		return err
	}
	if !picked {
		return nil
	}

	// Connecting goes back through Worktree.Connect rather than straight to the
	// picked path, so it lands in exactly the session `sesh worktree connect
	// <number>` would have — startup command and PR refresh included.
	_, err = deps.Worktree.Connect(model.WorktreeConnectOpts{
		Number: entry.Number,
		Repo:   listOpts.Repo,
		Switch: switchSession,
	})
	return err
}

func newWorktreeListCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List worktrees with their issue titles",
		Long: `List every worktree for a repo, showing the issue title beside its number.

Titles are cached under $XDG_CACHE_HOME/sesh, so only numbers that are new or
past their TTL cost a request — and those are fetched in one batch rather than
one call per worktree. --refresh refetches every one of them, for a title edited
or an issue closed since the last listing.

--path and --repo are two ways to select the same [[worktree]] block; with
neither, the repo is detected from the current directory.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			path, _ := cmd.Flags().GetString("path")
			repo, _ := cmd.Flags().GetString("repo")
			jsonOutput, _ := cmd.Flags().GetBool("json")
			refresh, _ := cmd.Flags().GetBool("refresh")

			entries, err := deps.Worktree.List(model.WorktreeListOpts{Path: path, Repo: repo, Refresh: refresh})
			if err != nil {
				return err
			}

			if jsonOutput {
				encoded, err := json.MarshalIndent(entries, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(encoded))
				return nil
			}

			for _, entry := range entries {
				fmt.Println(worktreeLabel(entry))
			}
			return nil
		},
	}

	cmd.Flags().String("path", "", "Worktree root to list, e.g. ~/c/nu/w")
	cmd.Flags().StringP("repo", "r", "", "Target repo as org/repo (overrides cwd detection)")
	cmd.Flags().BoolP("json", "j", false, "output as json")
	cmd.Flags().Bool("refresh", false, "Refetch every title and state, ignoring the cache")

	return cmd
}

// worktreeLabel renders one entry for the picker and for plain output. The
// number leads so the list sorts and filters by it, with the title appended
// when it is known.
func worktreeLabel(entry model.WorktreeEntry) string {
	if entry.Title == "" {
		return strconv.Itoa(entry.Number)
	}
	return fmt.Sprintf("%d  %s", entry.Number, entry.Title)
}

func newWorktreeConnectCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connect [number]",
		Aliases: []string{"c"},
		Short:   "Connect to a worktree for an issue/PR, creating it if it doesn't exist",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromBrowser, _ := cmd.Flags().GetBool("browser")

			if fromBrowser && len(args) > 0 {
				return fmt.Errorf("provide a number or --browser, not both")
			}
			if !fromBrowser && len(args) == 0 {
				return fmt.Errorf("provide an issue/PR number or use --browser")
			}

			number := 0
			if !fromBrowser {
				n, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("worktree number must be an integer: %q", args[0])
				}
				number = n
			}

			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			repo, _ := cmd.Flags().GetString("repo")
			pr, _ := cmd.Flags().GetBool("pr")
			switchFlag, _ := cmd.Flags().GetBool("switch")

			_, err = deps.Worktree.Connect(model.WorktreeConnectOpts{
				Number:      number,
				Repo:        repo,
				Pr:          pr,
				Switch:      switchFlag,
				FromBrowser: fromBrowser,
			})
			return err
		},
	}

	cmd.Flags().StringP("repo", "r", "", "Target repo as org/repo (overrides cwd detection)")
	cmd.Flags().BoolP("pr", "p", false, "Treat the number as a pull request")
	cmd.Flags().BoolP("switch", "s", false, "Switch the session rather than attach (for use outside tmux)")
	cmd.Flags().BoolP("browser", "b", false, "Read the issue/PR from the active browser tab (macOS)")

	return cmd
}
