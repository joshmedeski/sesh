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
		Tmuxinator   bool
	}
	srcStrategy func(*RealLister) (model.SeshSessions, error)
)

var srcStrategies = map[string]srcStrategy{
	"tmux":       listTmux,
	"config":     listConfig,
	"zoxide":     listZoxide,
	"tmuxinator": listTmuxinator,
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
		fullOrderedIndex = append(fullOrderedIndex, sessions.OrderedIndex...)
		filteredIndex := fullOrderedIndex[:0] // Create a slice with the same underlying array but length 0
		for _, i := range sessions.OrderedIndex {
			if opts.HideAttached && sessions.Directory[i].Attached == 1 {
				// TODO: remove the item from the fullOrderedIndex
				continue
			}
			filteredIndex = append(filteredIndex, i)
			fullDirectory[i] = sessions.Directory[i]
		}
		fullOrderedIndex = filteredIndex
	}

	return model.SeshSessions{
		OrderedIndex: fullOrderedIndex,
		Directory:    fullDirectory,
	}, nil
}
