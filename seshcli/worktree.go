package seshcli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/model"
)

func NewWorktreeCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "worktree",
		Aliases: []string{"wt"},
		Short:   "Manage git worktrees as sesh sessions",
	}
	cmd.AddCommand(newWorktreeConnectCommand(base))
	cmd.AddCommand(newWorktreeListCommand(base))
	return cmd
}

func newWorktreeListCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List worktrees with their issue titles",
		Long: `List every worktree for a repo, showing the issue title beside its number.

Titles are cached under $XDG_CACHE_HOME/sesh, so only numbers that are new or
past their TTL cost a request — and those are fetched in one batch rather than
one call per worktree.

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

			entries, err := deps.Worktree.List(model.WorktreeListOpts{Path: path, Repo: repo})
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
