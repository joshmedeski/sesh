// gitutil.go
package render

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/git"
)

// FormatGitStatus renders a git.StatusSummary as the "+n ~n -n !n" compact
// string used in the status column (shared with the sessions section). Each
// part is styled with a distinct ANSI colour so the counts are distinguishable
// at a glance: staged green, unstaged yellow, deleted red, untracked magenta.
//
// The parts are emitted as one continuous ANSI run with only a final reset, so
// an outer cursor/hover background applied by RenderRow is not cancelled by
// per-part reset sequences. Spaces between parts are left unstyled so they
// inherit the row's highlight background cleanly.
func FormatGitStatus(status git.StatusSummary) string {
	parts := make([]string, 0, 4)
	if status.Staged > 0 {
		parts = append(parts, ansiFg(colorStaged)+fmt.Sprintf("+%d", status.Staged))
	}
	if status.Unstaged > 0 {
		parts = append(parts, ansiFg(colorUnstaged)+fmt.Sprintf("~%d", status.Unstaged))
	}
	if status.Deleted > 0 {
		parts = append(parts, ansiFg(ColorDeleted)+fmt.Sprintf("-%d", status.Deleted))
	}
	if status.Untracked > 0 {
		parts = append(parts, ansiFg(colorUntracked)+fmt.Sprintf("!%d", status.Untracked))
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ") + "\x1b[0m"
}

// ansiFg returns an ANSI 256-colour foreground sequence for c.
func ansiFg(c lipgloss.ANSIColor) string {
	return fmt.Sprintf("\x1b[38;5;%dm", int(c))
}
