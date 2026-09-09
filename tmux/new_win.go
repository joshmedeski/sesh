package tmux

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

// windowTargetFormat makes tmux print the window it just created, so callers
// can hand a concrete target to a script or a later send-keys instead of
// guessing an index.
const windowTargetFormat = "#{session_name}:#{window_index}"

// NewWindowInSession creates a window and returns its "session:index" target.
func (t *RealTmux) NewWindowInSession(opts model.TmuxWindowOpts) (string, error) {
	args := []string{"new-window", "-P", "-F", windowTargetFormat}
	if opts.Name != "" {
		args = append(args, "-n", opts.Name)
	}
	if opts.StartDir != "" {
		args = append(args, "-c", opts.StartDir)
	}
	if opts.Background {
		args = append(args, "-d")
	}
	if opts.TargetSession != "" {
		// trailing colon forces session (not window) target resolution
		args = append(args, "-t", opts.TargetSession+":")
	}
	if opts.Command != "" {
		args = append(args, "--", opts.Command)
	}
	target, err := t.shell.Cmd(t.bin, args...)
	return strings.TrimSpace(target), err
}
