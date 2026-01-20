package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func listProjects(l *RealLister) (model.SeshSessions, error) {
	projectResults, err := l.projects.List()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list projects: %q", err)
	}

	numProjects := len(projectResults)
	orderedIndex := make([]string, numProjects)
	directory := make(model.SeshSessionMap)

	for i, p := range projectResults {
		key := fmt.Sprintf("projects:%s", p.Path)
		orderedIndex[i] = key
		directory[key] = p
	}

	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}
