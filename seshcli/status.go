package seshcli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/statuscache"
)

func NewStatusCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show contextual status for the current session (for the tmux status bar)",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			path := statusPath(deps)
			if len(args) > 0 && args[0] != "" {
				path = args[0]
			}
			if path == "" {
				return nil
			}

			refresh, _ := cmd.Flags().GetBool("refresh")
			if refresh {
				return runRefresh(deps, path)
			}

			ttl := deps.Config.Github.EffectiveTTL()
			out, spawn := computeStatus(deps, ttl, path)
			if out != "" {
				fmt.Print(out)
			}
			if spawn {
				_ = deps.Refresher.Spawn(path)
			}
			return nil
		},
	}
	cmd.Flags().Bool("refresh", false, "Internal: fetch live data and update the status cache (used for background refresh)")
	_ = cmd.Flags().MarkHidden("refresh")
	return cmd
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

// computeStatus decides what to print and whether a background refresh should
// be spawned. With ttl==0 the cache is bypassed and gh is queried live.
func computeStatus(deps *Deps, ttl int, path string) (string, bool) {
	if ttl == 0 {
		issue, found, _ := deps.Github.Issue(path)
		if found {
			return formatStatus(issue), false
		}
		return "", false
	}

	ref, ok := deps.Github.Resolve(path)
	if !ok {
		return "", false
	}

	key := statuscache.Key(ref.RepoRoot, ref.Branch)
	entry, found, _ := deps.StatusCache.Read(key)

	out := ""
	if found {
		if r, ok := entry.Preferred(); ok {
			out = formatStatus(toIssue(*r))
		}
	}
	stale := !found || time.Since(entry.Timestamp) > time.Duration(ttl)*time.Second
	return out, stale
}

// runRefresh performs the live gh fetch and always writes a cache entry
// (a negative entry when there is nothing to show), so the reader respects
// the TTL instead of re-spawning every tick.
func runRefresh(deps *Deps, path string) error {
	ref, ok := deps.Github.Resolve(path)
	if !ok {
		return nil // not a repo — nothing to cache
	}

	var entry statuscache.Entry
	if ref.HasNumber {
		if issue, found, _ := deps.Github.Issue(path); found {
			entry.Issue = &statuscache.Ref{Number: issue.Number, Title: issue.Title, State: issue.State}
		}
	}
	entry.Timestamp = time.Now()
	return deps.StatusCache.Write(statuscache.Key(ref.RepoRoot, ref.Branch), entry)
}

// formatStatus renders an issue as a tmux-styled status line.
func formatStatus(issue github.Issue) string {
	color := "green"
	if issue.State != "OPEN" {
		color = "red"
	}
	return fmt.Sprintf("#[fg=%s,bold]%s#[default] #[fg=magenta]Issue #%d#[default] %s", color, issue.State, issue.Number, issue.Title)
}

// toIssue adapts a cache Ref into the github.Issue value formatStatus consumes.
func toIssue(ref statuscache.Ref) github.Issue {
	return github.Issue{Number: ref.Number, Title: ref.Title, State: ref.State}
}
