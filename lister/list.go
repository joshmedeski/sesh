package lister

import (
	"github.com/joshmedeski/sesh/v2/model"
	"slices"
)

type (
	ListOptions struct {
		Config         bool
		HideAttached   bool
		Icons          bool
		Json           bool
		Tmux           bool
		Zoxide         bool
		Tmuxinator     bool
		HideDuplicates bool
	}
	srcStrategy func(*RealLister) (model.SeshSessions, error)
)

var srcStrategies = map[string]srcStrategy{
	"tmux":       listTmux,
	"config":     listConfig,
	"tmuxinator": listTmuxinator,
	"zoxide":     listZoxide,
}

func (l *RealLister) List(opts ListOptions) (model.SeshSessions, error) {
	fullDirectory := make(model.SeshSessionMap)
	fullOrderedIndex := make([]string, 0)

	srcsOrderedIndex := srcs(opts)
	srcsOrderedIndex = sortSources(srcsOrderedIndex, l.config.SortOrder)

	for _, src := range srcsOrderedIndex {
		sessions, err := srcStrategies[src](l)
		if err != nil {
			return model.SeshSessions{}, err
		}
		fullOrderedIndex = append(fullOrderedIndex, sessions.OrderedIndex...)
		for _, i := range sessions.OrderedIndex {
			fullDirectory[i] = sessions.Directory[i]
		}
	}

	if opts.HideDuplicates {
		directoryHash := make(map[string]int)
		nameHash := make(map[string]int)
		destIndex := 0
		for _, index := range fullOrderedIndex {
			session := fullDirectory[index]
			nameIsDuplicate := nameHash[session.Name] != 0
			pathIsDuplicate := session.Path != "" && directoryHash[session.Path] != 0
			if !nameIsDuplicate && !pathIsDuplicate {
				fullOrderedIndex[destIndex] = index
				directoryHash[session.Path] = 1
				nameHash[session.Name] = 1
				destIndex = destIndex + 1
			}
		}
		fullOrderedIndex = fullOrderedIndex[:destIndex]
	}

	if opts.HideAttached {
		attachedSession, _ := GetAttachedTmuxSession(l)
		for i, index := range fullOrderedIndex {
			if fullDirectory[index].Name == attachedSession.Name {
				fullOrderedIndex = slices.Delete(fullOrderedIndex, i, i+1)
				break
			}
		}
	}

	return model.SeshSessions{
		OrderedIndex: fullOrderedIndex,
		Directory:    fullDirectory,
	}, nil
}
