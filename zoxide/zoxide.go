package zoxide

import (
	"fmt"
	"joshmedeski/sesh/dir"
	"os"
	"os/exec"
	"path"
	"strings"
)

func Dirs() ([]string, error) {
	cmd := exec.Command("zoxide", "query", "-l")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	resultList := strings.TrimSpace(string(output))
	results := strings.Split(resultList, "\n")

	for i, path := range results {
		prettyPath, err := dir.PrettyPath(path)
		if err != nil {
			return nil, err
		}
		results[i] = prettyPath
	}
	return results, nil
}

func Add(result string) {
	if !path.IsAbs(result) {
		return
	}
	cmd := exec.Command("zoxide", "add", result)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
