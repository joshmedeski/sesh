package seshcli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/github"
)

func NewStatusCommand(base *BaseDeps) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show contextual status for the current session (for the tmux status bar)",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			path := statusPath(deps)
			if path == "" {
				return nil
			}

			issue, found, _ := deps.Github.Issue(path)
			if !found {
				return nil
			}

			fmt.Print(formatStatus(issue))
			return nil
		},
	}
}

// statusPath resolves the directory to inspect: the attached tmux session's
// path when running inside tmux, otherwise the current working directory.
func statusPath(deps *Deps) string {
	if session, exists := deps.Lister.GetAttachedTmuxSession(); exists {
		return session.Path
	}
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return cwd
}

// formatStatus renders an issue as a tmux-styled status line.
func formatStatus(issue github.Issue) string {
	color := "green"
	if issue.State != "OPEN" {
		color = "red"
	}
	return fmt.Sprintf("#[fg=%s,bold]%s#[default] #%d %s", color, issue.State, issue.Number, issue.Title)
}
