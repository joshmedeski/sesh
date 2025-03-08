package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func listZoxide(l *RealLister) (model.SeshSessions, error) {
	zoxideResults, err := l.zoxide.ListResults()
	numZoxideResults := len(zoxideResults)
	orderedIndex := make([]string, numZoxideResults)
	directory := make(model.SeshSessionMap)
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list zoxide sessions: %q", err)
	}
	for i, r := range zoxideResults {
		name, err := l.home.ShortenHome(r.Path)
		if err != nil {
			return model.SeshSessions{}, fmt.Errorf("couldn't shorten path: %q", err)
		}
		key := fmt.Sprintf("zoxide:%s", name)
		orderedIndex[i] = key
		directory[key] = model.SeshSession{
			Src:   "zoxide",
			Name:  name,
			Path:  r.Path,
			Score: r.Score,
		}
	}
	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) FindZoxideSession(path string) (model.SeshSession, bool) {
	result, err := l.zoxide.Query(path)
	if err != nil {
		return model.SeshSession{}, false
	}
	return model.SeshSession{
		Src:   "zoxide",
		Name:  result.Path,
		Path:  result.Path,
		Score: result.Score,
	}, true
}
