package tmux

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/convert"
	"github.com/joshmedeski/sesh/v2/model"
)

func listWindowsFormat() string {
	variables := []string{
		"#{window_index}",
		"#{window_name}",
		"#{pane_current_path}",
		"#{window_active}",
	}
	return strings.Join(variables, separator)
}

func (t *RealTmux) ListWindows(targetSession string) ([]*model.TmuxWindow, error) {
	var args []string
	args = append(args, "list-windows")
	if targetSession != "" {
		args = append(args, "-t", targetSession)
	}
	args = append(args, "-F", listWindowsFormat())

	output, err := t.shell.ListCmd("tmux", args...)
	if err != nil {
		return nil, err
	}
	return parseTmuxWindowsOutput(output)
}

func parseTmuxWindowsOutput(rawList []string) ([]*model.TmuxWindow, error) {
	windows := make([]*model.TmuxWindow, 0, len(rawList))
	for _, line := range rawList {
		fields := strings.Split(line, separator)
		if len(fields) != 4 {
			continue
		}
		windows = append(windows, &model.TmuxWindow{
			Index:  convert.StringToInt(fields[0]),
			Name:   fields[1],
			Path:   fields[2],
			Active: convert.StringToBool(fields[3]),
		})
	}
	return windows, nil
}
