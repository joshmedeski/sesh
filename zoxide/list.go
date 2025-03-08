package zoxide

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/convert"
	"github.com/joshmedeski/sesh/v2/model"
)

func (z *RealZoxide) ListResults() ([]*model.ZoxideResult, error) {
	list, err := z.shell.ListCmd("zoxide", "query", "--list", "--score")
	if err != nil {
		return nil, err
	}
	results := make([]*model.ZoxideResult, 0, len(list))
	for _, result := range list {
		if result == "" {
			break
		}
		trimmedResult := strings.TrimSpace(result)
		fields := strings.SplitN(trimmedResult, " ", 2)
		score, err := convert.StringToFloat(fields[0])
		if err != nil {
			return nil, err
		}
		results = append(results, &model.ZoxideResult{
			Score: score,
			Path:  fields[1],
		})
	}
	return results, nil
}
