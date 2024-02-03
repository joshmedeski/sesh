package zoxide

import (
	"fmt"
	"os"
	"strings"

	"github.com/joshmedeski/sesh/convert"
)

func Query(dir string) (ZoxideResult, error) {
	output, err := zoxideCmd([]string{"query", "-s", dir})
	if err != nil {
		return ZoxideResult{}, err
	}
	cleanOutput := strings.TrimSpace(string(output))
	list := strings.Split(cleanOutput, "\n")
	listLen := len(list)
	if listLen == 1 && list[0] == "" {
		return ZoxideResult{}, fmt.Errorf("no results found")
	}
	results := make([]ZoxideResult, 0, listLen)
	for _, line := range list {
		trimmed := strings.Trim(line, "[]")
		trimmed = strings.Trim(trimmed, " ")
		fields := strings.SplitN(trimmed, " ", 2)
		if len(fields) != 2 {
			fmt.Println(
				"Zoxide entry has invalid number of fields (expected 2)",
			)
			os.Exit(1)
		}
		path := fields[1]
		results = append(results, ZoxideResult{
			Score: convert.StringToFloat(fields[0]),
			Name:  convert.PathToPretty(path),
			Path:  path,
		})
	}
	if len(results) == 0 {
		return ZoxideResult{}, err
	}
	return results[0], nil
}
