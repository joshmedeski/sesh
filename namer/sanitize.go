package namer

import "strings"

// SanitizeTitle prepares a GitHub issue title for use inside a tmux session
// name: it keeps spaces and original casing but replaces the two characters
// tmux forbids in session names ('.' and ':') with a space, then collapses
// the resulting runs of whitespace and trims.
func SanitizeTitle(title string) string {
	replaced := strings.NewReplacer(".", " ", ":", " ").Replace(title)
	return strings.Join(strings.Fields(replaced), " ")
}
