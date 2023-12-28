package zoxide

import (
	"joshmedeski/sesh/dir"
	"os/exec"
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
