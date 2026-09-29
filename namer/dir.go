package namer

import (
	// Since sesh uses forward slashes for all internal path handling,
	// using slash normalization and the path package for path splitting
	// provides better Windows compatibility than filepath here.
	pathpkg "path"
	"path/filepath"
	"strings"
)

// lastNComponents returns the last n path components joined by "/".
// For n <= 1 it returns the basename. Empty string in, empty string out.
func lastNComponents(path string, n int) string {
	if path == "" {
		return ""
	}
	cleanPath := filepath.ToSlash(filepath.Clean(path))
	if n <= 1 {
		return pathpkg.Base(cleanPath)
	}

	parts := make([]string, 0, n)
	current := cleanPath

	for len(parts) < n && current != "/" && current != "." {
		base := pathpkg.Base(current)
		if base == "" || base == "." {
			break
		}
		parts = append(parts, base)
		parent := pathpkg.Dir(current)

		// Extra check for Windows compatibility (breaks when driver letter reached).
		if parent == current {
			break
		}
		current = parent
	}

	if len(parts) == 0 {
		return pathpkg.Base(cleanPath)
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
