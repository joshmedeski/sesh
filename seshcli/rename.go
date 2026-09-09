package seshcli

import (
	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
)

func NewRenameCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename [name]",
		Short: "Rename a tmux session, optionally enriching it with its GitHub issue title",
		Long: "Rename a tmux session. With --enrich, the session named by [name] " +
			"(or the attached session) is renamed to '<namer name> — <issue title>' " +
			"when its branch resolves to a GitHub issue, and back to the bare namer " +
			"name when it does not. Intended to be run from a tmux session-created hook.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}
			enrich, _ := cmd.Flags().GetBool("enrich")
			if !enrich {
				return nil // only the enrich mode is implemented today
			}
			return runEnrich(deps, args)
		},
	}
	cmd.Flags().Bool("enrich", false, "Rename the session to include its GitHub issue title")
	return cmd
}

// runEnrich resolves the target session, computes its enriched name, and
// renames it when that differs from the current name. Every "nothing to do"
// path returns nil.
func runEnrich(deps *Deps, args []string) error {
	target, ok := renameTarget(deps, args)
	if !ok || target.Name == "" {
		return nil
	}
	newName := enrichedName(deps, target.Path)
	if newName == "" || newName == target.Name {
		return nil
	}
	_, err := deps.Tmux.RenameSession(target.Name, newName)
	return err
}

// renameTarget resolves which session to rename: the one named by the first
// arg, or the attached session when no arg is given.
func renameTarget(deps *Deps, args []string) (model.SeshSession, bool) {
	if len(args) > 0 && args[0] != "" {
		return deps.Lister.FindTmuxSession(args[0])
	}
	return deps.Lister.GetAttachedTmuxSession()
}

// enrichedName recomputes the base name from path (the deterministic source of
// truth, which makes re-runs idempotent) and appends the sanitized issue title
// when the branch resolves to an issue. With no issue it returns the bare base,
// which self-heals a stale suffix. Returns "" when the base cannot be computed.
func enrichedName(deps *Deps, path string) string {
	baseName, err := deps.Namer.Name(path)
	if err != nil || baseName == "" {
		return ""
	}
	issue, found, _ := deps.Github.Issue(path)
	if !found {
		return baseName
	}
	title := namer.SanitizeTitle(issue.Title)
	if title == "" {
		return baseName
	}
	return baseName + model.SessionNameSeparator + title
}
