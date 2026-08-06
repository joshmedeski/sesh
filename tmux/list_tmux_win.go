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

func listAllWindowNamesFormat() string {
	variables := []string{
		"#{session_name}",
		"#{pane_title}",
	}
	return strings.Join(variables, separator)
}

// ListAllWindowNames returns the window names of every session, keyed by
// session name, in a single tmux invocation. The picker uses this so showing
// window names doesn't cost one shell call per session.
func (t *RealTmux) ListAllWindowNames() (map[string][]string, error) {
	output, err := t.shell.ListCmd(t.bin, "list-windows", "-a", "-F", listAllWindowNamesFormat())
	if err != nil {
		return nil, err
	}
	return parseAllWindowNamesOutput(output), nil
}

func parseAllWindowNamesOutput(rawList []string) map[string][]string {
	windowNames := make(map[string][]string)
	for _, line := range rawList {
		// Split on the first separator only: window names are more likely to
		// contain it than session names (which come from directory names).
		session, name, found := strings.Cut(line, separator)
		if !found {
			continue
		}
		if session == "" || name == "" {
			continue
		}
		windowNames[session] = append(windowNames[session], name)
	}
	return windowNames
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
