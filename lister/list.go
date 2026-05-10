package lister

import (
	"slices"
	"sync"

	"github.com/joshmedeski/sesh/v2/model"
)

type (
	ListOptions struct {
		Config         bool
		HideAttached   bool
		Icons          bool
		NoColor        bool
		Json           bool
		Tmux           bool
		Zoxide         bool
		Tmuxinator     bool
		HideDuplicates bool
		Panes          bool
		Blacklisted    bool
	}
	srcStrategy func(*RealLister) (model.SeshSessions, error)
)

type strategyResult struct {
	source   string
	sessions model.SeshSessions
	err      error
}

var srcStrategies = map[string]srcStrategy{
	"tmux":       listTmux,
	"config":     listConfig,
	"tmuxinator": listTmuxinator,
	"zoxide":     listZoxide,
	"tmux-pane":  listTmuxPanes,
}

func (l *RealLister) List(opts ListOptions) (model.SeshSessions, error) {
	fullDirectory := make(model.SeshSessionMap)
	fullOrderedIndex := make([]string, 0)

	srcsOrderedIndex := srcs(opts)
	srcsOrderedIndex = sortSources(srcsOrderedIndex, l.config.SortOrder)

	resultsChan := make(chan strategyResult, len(srcsOrderedIndex))
	var wg sync.WaitGroup

	for _, src := range srcsOrderedIndex {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			sessions, err := srcStrategies[s](l)
			resultsChan <- strategyResult{source: s, sessions: sessions, err: err}
		}(src)
	}

	wg.Wait()
	close(resultsChan)

	// Collect results into a map for easy lookup
	resultsMap := make(map[string]model.SeshSessions)
	for res := range resultsChan {
		if res.err != nil {
			return model.SeshSessions{}, res.err
		}
		resultsMap[res.source] = res.sessions
	}

	for _, src := range srcsOrderedIndex {
		sessions := resultsMap[src]
		fullOrderedIndex = append(fullOrderedIndex, sessions.OrderedIndex...)
		for _, i := range sessions.OrderedIndex {
			fullDirectory[i] = sessions.Directory[i]
		}
	}

	if len(l.config.Blacklist) > 0 || opts.Blacklisted {
		compiled := compileBlacklist(l.config.Blacklist)
		filteredIndex := make([]string, 0, len(fullOrderedIndex))
		filteredDirectory := make(model.SeshSessionMap)
		for _, index := range fullOrderedIndex {
			session := fullDirectory[index]
			if isBlacklisted(compiled, session.Name) == opts.Blacklisted {
				filteredIndex = append(filteredIndex, index)
				filteredDirectory[index] = session
			}
		}
		fullOrderedIndex = filteredIndex
		fullDirectory = filteredDirectory
	}

	// HideAttached runs before HideDuplicates so the attached tmux session
	// is removed from the dedup input. Otherwise tmux would win the dedup
	// against a same-named config/tmuxinator entry and then be hidden,
	// leaving no entry for the user to pick.
	if opts.HideAttached {
		if attachedSession, ok := GetAttachedTmuxSession(l); ok {
			for i, index := range fullOrderedIndex {
				s := fullDirectory[index]
				if s.Src == "tmux" && s.Name == attachedSession.Name {
					fullOrderedIndex = slices.Delete(fullOrderedIndex, i, i+1)
					break
				}
			}
		}
	}

	if opts.HideDuplicates {
		fullOrderedIndex = applyDedup(model.SeshSessions{
			OrderedIndex: fullOrderedIndex,
			Directory:    fullDirectory,
		})
	}

	return model.SeshSessions{
		OrderedIndex: fullOrderedIndex,
		Directory:    fullDirectory,
	}, nil
}
