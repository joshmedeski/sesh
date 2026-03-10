package tmux

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/v2/convert"
	"github.com/joshmedeski/sesh/v2/model"
)

func listpanesformat() string {
	variables := []string{
		"#{window_index}",
		"#{window_name}",
		"#{pane_index}",
		"#{pane_title}",
		"#{pane_current_command}",
		"#{pane_current_path}",
		"#{pane_id}",
	}
	return strings.Join(variables, separator)
}

func (t *RealTmux) ListTmuxPanes() ([]*model.TmuxPane, error) {
	output, err := t.shell.ListCmd("tmux", "list-panes", "-s", "-F", listpanesformat())
	if err != nil {
		return []*model.TmuxPane{}, nil
	}
	return parseTmuxPanesOutput(output)
}

func parseTmuxPanesOutput(rawList []string) ([]*model.TmuxPane, error) {
	panes := make([]*model.TmuxPane, 0, len(rawList))
	for _, line := range rawList {
		fields := strings.Split(line, separator)
		if len(fields) != 7 {
			continue
		}
		panes = append(panes, &model.TmuxPane{
			WindowIndex: convert.StringToInt(fields[0]),
			WindowName:  fields[1],
			PaneIndex:   convert.StringToInt(fields[2]),
			PaneTitle:   fields[3],
			PaneCommand: fields[4],
			PanePath:    fields[5],
			PaneID:      fields[6],
		})
	}
	return panes, nil
}

func (t *RealTmux) SelectPane(windowIndex int, paneIndex int) (string, error) {
	if _, err := t.shell.Cmd("tmux", "select-window", "-t", fmt.Sprintf(":%d", windowIndex)); err != nil {
		return "", fmt.Errorf("failed to select window %d: %w", windowIndex, err)
	}
	if _, err := t.shell.Cmd("tmux", "select-pane", "-t", fmt.Sprintf(".%d", paneIndex)); err != nil {
		return "", fmt.Errorf("failed to select pane %d: %w", paneIndex, err)
	}
	return fmt.Sprintf("selected pane %d in window %d", paneIndex, windowIndex), nil
}

func (t *RealTmux) GetCurrentSession() (string, error) {
	return t.shell.Cmd("tmux", "display-message", "-p", "#{session_name}")
}
