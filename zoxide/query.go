package zoxide

import "github.com/joshmedeski/sesh/v2/model"

func (z *RealZoxide) Query(query string) (*model.ZoxideResult, error) {
	result, err := z.shell.Cmd("zoxide", "query", query)
	// TODO: handle no result found
	if err != nil {
		return nil, err
	}
	return &model.ZoxideResult{
		Score: 0,
		Path:  result,
	}, nil
}
