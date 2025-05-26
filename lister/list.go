package lister

import (
	"fmt"
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

func (l *RealLister) getWindowName(session, window string) string {
	name, err := l.tmux.GetWindowName(session, window)
	if err != nil {
		return window
	}
	return name
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

	// Add individual marked windows as separate entries
	for _, marked := range markedSessions {
		if !marked.Marked {
			continue
		}
		
		windowName := l.getWindowName(marked.Session, marked.Window)
		markedKey := fmt.Sprintf("%s:%s(%s)", marked.Session, windowName, marked.Window)
		
		// Find the original session to get its path
		originalSession, exists := fullDirectory[marked.Session]
		if !exists {
			// If original session not found, create a basic one
			originalSession = model.SeshSession{
				Src:  "tmux",
				Name: marked.Session,
				Path: "",
			}
		}
		
		// Create a new session entry for this marked window
		markedSession := model.SeshSession{
			Src:           "tmux",
			Name:          markedKey,
			Path:          originalSession.Path,
			Marked:        true,
			MarkedWindows: []string{marked.Window},
			MarkTimestamp: marked.Timestamp,
			AlertLevel:    marker.GetAlertLevel(marked.Session, marked.Window),
		}
		
		fullDirectory[markedKey] = markedSession
		fullOrderedIndex = append(fullOrderedIndex, markedKey)
	}

	if opts.Marked {
		var markedIndices []string
		for _, index := range fullOrderedIndex {
			if fullDirectory[index].Marked {
				markedIndices = append(markedIndices, index)
			}
		}
		
		sort.Slice(markedIndices, func(i, j int) bool {
			sessionI := fullDirectory[markedIndices[i]]
			sessionJ := fullDirectory[markedIndices[j]]
			
			if sessionI.AlertLevel != sessionJ.AlertLevel {
				return sessionI.AlertLevel > sessionJ.AlertLevel
			}
			return sessionI.MarkTimestamp > sessionJ.MarkTimestamp
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
			sessionI := fullDirectory[markedIndices[i]]
			sessionJ := fullDirectory[markedIndices[j]]
			
			if sessionI.AlertLevel != sessionJ.AlertLevel {
				return sessionI.AlertLevel > sessionJ.AlertLevel
			}
			return sessionI.MarkTimestamp > sessionJ.MarkTimestamp
		})
		
		fullOrderedIndex = append(markedIndices, normalIndices...)
	}

	return model.SeshSessions{
		OrderedIndex: fullOrderedIndex,
		Directory:    fullDirectory,
	}, nil
}