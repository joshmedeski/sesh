package zoxide

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/convert"
	"github.com/joshmedeski/sesh/v2/model"
)

func (z *RealZoxide) ListResults() ([]*model.ZoxideResult, error) {
	parts, err := z.shell.PrepareCmd(z.listCommand, map[string]string{})
	if err != nil {
		return nil, err
	}
	list, err := z.shell.ListCmd(parts[0], parts[1:]...)
	if err != nil {
		return nil, err
	}
	results := make([]*model.ZoxideResult, 0, len(list))
	for _, result := range list {
		trimmedResult := strings.TrimSpace(result)
		if trimmedResult == "" {
			continue
		}
		results = append(results, parseResult(trimmedResult))
	}
	return results, nil
}

// parseResult auto-detects an optional leading numeric score. When the first
// space-delimited field parses as a float it is used as the score and the
// remainder is the path (zoxide's `--score` output). Otherwise the whole line
// is treated as the path with a zero score (fasd, memy, plain paths).
func parseResult(line string) *model.ZoxideResult {
	if fields := strings.SplitN(line, " ", 2); len(fields) == 2 {
		if score, err := convert.StringToFloat(fields[0]); err == nil {
			return &model.ZoxideResult{Score: score, Path: fields[1]}
		}
	}
	return &model.ZoxideResult{Score: 0, Path: line}
}
