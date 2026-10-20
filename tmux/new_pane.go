package tmux

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

func (t *RealTmux) NewPaneInWindow(opts model.TmuxPaneOpts) (string, error) {
	args := []string{"split-window", "-d", "-P", "-F", "#{pane_id}", "-t", sessionTarget(opts.TargetWindow)}
	if opts.StartDir != "" {
		args = append(args, "-c", escapeFormat(opts.StartDir))
	}

	paneID, err := t.shell.Cmd(t.bin, args...)
	if err != nil {
		return "", err
	}
	paneID = strings.TrimSpace(paneID)

	if opts.Command != "" {
		if _, err := t.shell.Cmd(t.bin, "send-keys", "-t", paneID, opts.Command, "Enter"); err != nil {
			return paneID, err
		}
	}
	if opts.Layout != "" {
		if _, err := t.SelectLayout(opts.TargetWindow, opts.Layout); err != nil {
			return paneID, err
		}
	}
	return paneID, nil
}
