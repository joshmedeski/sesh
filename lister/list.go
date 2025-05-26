package lister

import (
	"sort"

	"github.com/joshmedeski/sesh/v2/marker"
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
		Marked         bool
	}
	srcStrategy func(*RealLister) (model.SeshSessions, error)
)

var srcStrategies = map[string]srcStrategy{
	"tmux":       listTmux,
	"config":     listConfig,
	"tmuxinator": listTmuxinator,
	"zoxide":     listZoxide,
}

func (l *RealLister) List(opts ListOptions, marker marker.Marker) (model.SeshSessions, error) {
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

	markedSessions, err := marker.GetMarkedSessions()
	if err != nil {
		return model.SeshSessions{}, err
	}

	for _, index := range fullOrderedIndex {
		session := fullDirectory[index]
		for _, marked := range markedSessions {
			if marked.Session == session.Name && session.Src == "tmux" {
				session.Marked = true
				if len(session.MarkedWindows) == 0 {
					session.MarkedWindows = []string{}
				}
				session.MarkedWindows = append(session.MarkedWindows, marked.Window)
				if session.MarkTimestamp == 0 || marked.Timestamp > session.MarkTimestamp {
					session.MarkTimestamp = marked.Timestamp
				}
				fullDirectory[index] = session
			}
		}
	}

	if opts.Marked {
		var markedIndices []string
		for _, index := range fullOrderedIndex {
			if fullDirectory[index].Marked {
				markedIndices = append(markedIndices, index)
			}
		}
		
		sort.Slice(markedIndices, func(i, j int) bool {
			return fullDirectory[markedIndices[i]].MarkTimestamp > fullDirectory[markedIndices[j]].MarkTimestamp
		})
		
		fullOrderedIndex = markedIndices
	} else {
		var markedIndices []string
		var normalIndices []string
		
		for _, index := range fullOrderedIndex {
			if fullDirectory[index].Marked {
				markedIndices = append(markedIndices, index)
			} else {
				normalIndices = append(normalIndices, index)
			}
		}
		
		sort.Slice(markedIndices, func(i, j int) bool {
			return fullDirectory[markedIndices[i]].MarkTimestamp > fullDirectory[markedIndices[j]].MarkTimestamp
		})
		
		fullOrderedIndex = append(markedIndices, normalIndices...)
	}

	return model.SeshSessions{
		OrderedIndex: fullOrderedIndex,
		Directory:    fullDirectory,
	}, nil
}