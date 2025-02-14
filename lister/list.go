package lister

import (
	"github.com/joshmedeski/sesh/v2/model"
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

	for _, src := range srcsOrderedIndex {
		sessions, err := srcStrategies[src](l)
		if err != nil {
			return model.SeshSessions{}, err
		}
		if opts.HideAttached {
			attachedSession, _ := GetAttachedTmuxSession(l)
			sessionsCopy := sessions.OrderedIndex
			for i, ses := range sessionsCopy {
				if attachedSession.Name == sessions.Directory[ses].Name {
					sessions.OrderedIndex = append(sessions.OrderedIndex[:i],
						sessions.OrderedIndex[i+1:]...)
				}
			}
		}
		fullOrderedIndex = append(fullOrderedIndex, sessions.OrderedIndex...)
		for _, i := range sessions.OrderedIndex {
			fullDirectory[i] = sessions.Directory[i]
		}
	}

	if opts.HideDuplicates {
		directoryHash := make(map[string]int)
		destIndex := 0
		for _, index := range fullOrderedIndex {
			directoryPath := fullDirectory[index].Path
			if _, exists := directoryHash[directoryPath]; !exists {
				fullOrderedIndex[destIndex] = index
				directoryHash[directoryPath] = 1
				destIndex = destIndex + 1
			}
		}
		fullOrderedIndex = fullOrderedIndex[:destIndex]
	}

	return model.SeshSessions{
		OrderedIndex: fullOrderedIndex,
		Directory:    fullDirectory,
	}, nil
}
