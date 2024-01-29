package dir

import (
	"os"
	"path/filepath"
	"strings"
)

func AlternatePath(s string) (altPath string) {
	if filepath.IsAbs(s) {
		return ""
	}

	if strings.HasPrefix(s, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			altPath = filepath.Join(homeDir, strings.TrimPrefix(s, "~/"))
		}
	}

	if strings.HasPrefix(s, ".") {
		if a, err := filepath.Abs(s); err == nil {
			altPath = a
		}
	}

	return altPath
}
