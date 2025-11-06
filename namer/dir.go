package namer

import (
	"path/filepath"
	"strings"
)

// Gets the name from a directory
func dirName(n *RealNamer, path string) (string, error) {
	dirLength := n.config.DirLength

	if dirLength <= 1 {
		return n.pathwrap.Base(path), nil
	}

	cleanPath := filepath.Clean(path)

	parts := make([]string, 0, dirLength)
	current := cleanPath

	// Collect path components backwards
	for len(parts) < dirLength && current != "/" && current != "." {
		base := filepath.Base(current)
		if base == "" || base == "." {
			break
		}

		parts = append(parts, base)
		current = filepath.Dir(current)
	}

	if len(parts) == 0 {
		return n.pathwrap.Base(path), nil
	}

	// Reverse the parts
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, "/"), nil
}
