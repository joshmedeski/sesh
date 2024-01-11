package zoxide

import (
	"fmt"
	"joshmedeski/sesh/convert"
	"joshmedeski/sesh/tmux"
	"os"
	"strings"
)

type ZoxideResult struct {
	Name  string
	Path  string
	Score float64
}

func List(tmuxSessions []*tmux.TmuxSession) ([]*ZoxideResult, error) {
	output, err := zoxideCmd([]string{"query", "-ls"})
	if err != nil {
		return nil, err
	}
	cleanOutput := strings.TrimSpace(string(output))
	list := strings.Split(cleanOutput, "\n")

	results := make([]*ZoxideResult, 0, len(list))
	tmuxSessionPaths := make(map[string]struct{})
	for _, session := range tmuxSessions {
		tmuxSessionPaths[session.Path] = struct{}{}
	}

	for _, line := range list {
		trimmed := strings.Trim(line, "[]")
		trimmed = strings.Trim(trimmed, " ")
		fields := strings.SplitN(trimmed, " ", 2)
		if len(fields) != 2 {
			fmt.Println("Zoxide entry has invalid number of fields (expected 2)")
			os.Exit(1)
		}
		path := fields[1]
		if _, exists := tmuxSessionPaths[path]; !exists {
			results = append(results, &ZoxideResult{
				Score: convert.StringToFloat(fields[0]),
				Name:  convert.PathToPretty(path),
				Path:  path,
			})
		}
	}

	return results, nil
}
