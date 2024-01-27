package tmux

import (
	"os"
	"path/filepath"
	"strings"
)

// alternatePath returns an altnerate string that should be check when doing
// path based comparisons.
func alternatePath(s string) (altPath string) {
	// If the path is absolute, there is no alternate path.
	if filepath.IsAbs(s) {
		return ""
	}

	// If the path starts with a ~ it's likely relative to the home directory.
	if strings.HasPrefix(s, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			altPath = filepath.Join(homeDir, strings.TrimPrefix(s, "~/"))
		}
	}

	// If the path starts with a . it's likely relative to the current
	// directory.
	if strings.HasPrefix(s, ".") {
		if a, err := filepath.Abs(s); err == nil {
			altPath = a
		}
	}

	// If we get to this point the path is likely a subdirectory.
	return altPath
}
