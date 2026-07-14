package namer

import (
	"path/filepath"
	"strings"
)

// lastNComponents returns the last n path components joined by "/".
// For n <= 1 it returns the basename. Empty string in, empty string out.
func lastNComponents(path string, n int) string {
	if path == "" {
		return ""
	}
	if n <= 1 {
		return filepath.Base(path)
	}

	cleanPath := filepath.Clean(path)
	parts := make([]string, 0, n)
	current := cleanPath

	for len(parts) < n && current != "/" && current != "." {
		base := filepath.Base(current)
		if base == "" || base == "." {
			break
		}
		parts = append(parts, base)
		current = filepath.Dir(current)
	}

	if len(parts) == 0 {
		return filepath.Base(path)
	}

	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, "/")
}

// Gets the name from a directory
func dirName(n *RealNamer, path string) (string, error) {
	return lastNComponents(path, n.config.DirLength), nil
}
