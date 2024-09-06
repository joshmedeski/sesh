package lister

import (
	"github.com/joshmedeski/sesh/model"
)

type (
	ListOptions struct {
		Config       bool
		HideAttached bool
		Icons        bool
		Json         bool
		Tmux         bool
		Zoxide       bool
	}
	srcStrategy func(*RealLister) (model.SeshSessions, error)
)

var srcStrategies = map[string]srcStrategy{
	"tmux":   listTmux,
	"config": listConfig,
	"zoxide": listZoxide,
}

func (l *RealLister) List(opts ListOptions) (model.SeshSessions, error) {
	fullDirectory := make(model.SeshSessionMap)
	fullOrderedIndex := make([]string, 0)

	srcsOrderedIndex := srcs(opts)

	for _, src := range srcsOrderedIndex {
		sessions, err := srcStrategies[src](l)
		if err != nil {
			return model.SeshSessions{}, err
		}
		filteredIndex := []string{}
		for _, i := range sessions.OrderedIndex {
			if opts.HideAttached && sessions.Directory[i].Attached == 1 {
				// TODO: remove the item from the fullOrderedIndex
				continue
			}
			filteredIndex = append(filteredIndex, i)
			fullDirectory[i] = sessions.Directory[i]
		}
		if opts.HideAttached {
			fullOrderedIndex = append(fullOrderedIndex, filteredIndex...)
		} else {
			fullOrderedIndex = append(fullOrderedIndex, sessions.OrderedIndex...)
		}
	}

	return model.SeshSessions{
		OrderedIndex: fullOrderedIndex,
		Directory:    fullDirectory,
	}, nil
}
