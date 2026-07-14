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
	cmd.AddCommand(newWorktreeCreateCommand(base))
	return cmd
}

func newWorktreeCreateCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <number>",
		Aliases: []string{"c"},
		Short:   "Create (or reconnect to) a worktree for an issue/PR and connect to it",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("worktree number must be an integer: %q", args[0])
			}

			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			repo, _ := cmd.Flags().GetString("repo")
			pr, _ := cmd.Flags().GetBool("pr")
			switchFlag, _ := cmd.Flags().GetBool("switch")

			_, err = deps.Worktree.Create(model.WorktreeCreateOpts{
				Number: number,
				Repo:   repo,
				Pr:     pr,
				Switch: switchFlag,
			})
			return err
		},
	}

	cmd.Flags().StringP("repo", "r", "", "Target repo as org/repo (overrides cwd detection)")
	cmd.Flags().BoolP("pr", "p", false, "Treat the number as a pull request")
	cmd.Flags().BoolP("switch", "s", false, "Switch the session rather than attach (for use outside tmux)")

	return cmd
}
