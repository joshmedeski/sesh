package dir

import (
	"fmt"
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

func FullPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if strings.HasPrefix(path, "~") {
		return strings.Replace(path, "~", home, 1)
	}
	return path
}
