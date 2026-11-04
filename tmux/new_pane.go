package tmux

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

func (t *RealTmux) NewPaneInWindow(opts model.TmuxPaneOpts) (string, error) {
	target := cmp.Or(opts.TargetPane, sessionTarget(opts.TargetWindow))
	args := []string{"split-window", "-d", "-P", "-F", "#{pane_id}", "-t", target}
	if opts.StartDir != "" {
		args = append(args, "-c", escapeFormat(opts.StartDir))
	}

	paneID, err := t.shell.Cmd(t.bin, args...)
	if err != nil {
		return "", fmt.Errorf("split window: %w", err)
	}
	return strings.TrimSpace(paneID), nil
}
