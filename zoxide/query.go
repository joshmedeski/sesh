package zoxide

import "github.com/joshmedeski/sesh/v2/model"

func (z *RealZoxide) Query(query string) (*model.ZoxideResult, error) {
	parts, err := z.shell.PrepareCmd(z.queryCommand, map[string]string{"{}": query})
	if err != nil {
		return nil, err
	}
	// TODO: handle no result found
	result, err := z.shell.Cmd(parts[0], parts[1:]...)
	if err != nil {
		return nil, err
	}
	return &model.ZoxideResult{
		Score: 0,
		Path:  result,
	}, nil
}
