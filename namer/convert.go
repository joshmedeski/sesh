package namer

import "strings"

func convertToValidName(name string, backend string) string {
	if backend == "wezterm" {
		return name
	}
	// tmux does not allow . or : in session names.
	validName := strings.ReplaceAll(name, ".", "_")
	validName = strings.ReplaceAll(validName, ":", "_")
	return validName
}
