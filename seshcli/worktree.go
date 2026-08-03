package seshcli

import (
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
	return cmd
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
