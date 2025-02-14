package tmuxinator

import (
	"slices"

	"github.com/joshmedeski/sesh/v2/model"
)

func (t *RealTmuxinator) List() ([]*model.TmuxinatorConfig, error) {
	res, err := t.shell.ListCmd("tmuxinator", "list", "-n")
	if err != nil {
		// NOTE: return empty list if error
		return []*model.TmuxinatorConfig{}, nil
	}
	return parseTmuxinatorConfigsOutput(res)
}

func parseTmuxinatorConfigsOutput(rawList []string) ([]*model.TmuxinatorConfig, error) {
	cleanedList := slices.Delete(rawList, 0, 1)
	cleanedList = cleanedList[:len(cleanedList)-1]
	sessions := make([]*model.TmuxinatorConfig, 0, len(cleanedList))
	for _, line := range cleanedList {
		if len(line) > 0 {
			session := &model.TmuxinatorConfig{
				Name: line,
			}
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}
