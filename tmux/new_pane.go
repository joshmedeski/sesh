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

// SelectPaneByID makes the pane active in its window without switching the client's window.
func (t *RealTmux) SelectPaneByID(paneID string) (string, error) {
	return t.shell.Cmd(t.bin, "select-pane", "-t", paneID)
}

// ZoomPane zooms (fullscreen) the target pane. A window target zooms its active pane.
func (t *RealTmux) ZoomPane(target string) (string, error) {
	return t.shell.Cmd(t.bin, "resize-pane", "-Z", "-t", target)
}
