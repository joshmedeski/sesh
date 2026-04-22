package startup

import "strings"

// posixSingleQuote wraps s in POSIX single quotes so every byte is literal.
// Safe for use as an argument to sh/bash/zsh/fish via -c.
func posixSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
