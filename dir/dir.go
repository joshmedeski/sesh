package dir

import (
	"os"
	"strings"
)

func PrettyPath(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1), nil
	}

	return path, nil
}
