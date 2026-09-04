package lister

import (
	"cmp"
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

	srcGroups := groupSources(srcs(opts), l.config.SortOrder)
	srcsOrderedIndex := slices.Concat(srcGroups...)

	resultsChan := make(chan strategyResult, len(srcsOrderedIndex))
	var wg sync.WaitGroup

	// Window names are display-only and fetched once for all sessions, in
	// parallel with the source strategies, so they never cost a tmux call per
	// session.
	var windowNames map[string][]string
	if l.config.TUI.ShowWindows && slices.Contains(srcsOrderedIndex, "tmux") {
		wg.Add(1)
		go func() {
			defer wg.Done()
			windowNames = tmuxWindowNames(l)
		}()
	}

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

	// The directory is filled first so a merged group can look a session up by
	// key while it orders the group.
	for _, src := range srcsOrderedIndex {
		sessions := resultsMap[src]
		for _, i := range sessions.OrderedIndex {
			fullDirectory[i] = sessions.Directory[i]
		}
	}

	scores := zoxideScores(resultsMap["zoxide"])
	for group, sources := range srcGroups {
		for _, key := range groupOrder(sources, resultsMap, fullDirectory, scores) {
			session := fullDirectory[key]
			session.Group = group
			fullDirectory[key] = session
			fullOrderedIndex = append(fullOrderedIndex, key)
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

	sessions := model.SeshSessions{
		OrderedIndex: fullOrderedIndex,
		Directory:    fullDirectory,
	}
	attachWindowNames(sessions, windowNames)

	return sessions, nil
}

// groupOrder lays out one sort_order group's sessions. A group of one keeps the
// order its source returned. A merged group is reordered by zoxide score,
// highest first, so the block reads as "what am I working on right now" rather
// than "which source did this come from".
func groupOrder(sources []string, results map[string]model.SeshSessions, directory model.SeshSessionMap, scores map[string]float64) []string {
	index := make([]string, 0)
	for _, src := range sources {
		index = append(index, results[src].OrderedIndex...)
	}
	if len(sources) < 2 {
		return index
	}
	// Stable, so the sessions zoxide has never seen keep their source order
	// behind the scored ones instead of being shuffled among themselves.
	slices.SortStableFunc(index, func(a, b string) int {
		return cmp.Compare(sessionScore(directory[b], scores), sessionScore(directory[a], scores))
	})
	return index
}

// zoxideScores indexes the scores already fetched for the zoxide source by
// path, so scoring the other sources in a merged group costs no extra lookups.
func zoxideScores(sessions model.SeshSessions) map[string]float64 {
	scores := make(map[string]float64, len(sessions.OrderedIndex))
	for _, key := range sessions.OrderedIndex {
		session := sessions.Directory[key]
		scores[session.Path] = session.Score
	}
	return scores
}

// sessionScore is the score a session sorts by inside a merged group: its own
// when it came from zoxide, otherwise whatever zoxide has for its path. A path
// zoxide has never seen scores 0 and trails the group, which is the point —
// somewhere never actually visited doesn't belong at the top of a frecency
// ordered list.
func sessionScore(session model.SeshSession, scores map[string]float64) float64 {
	if session.Score != 0 {
		return session.Score
	}
	return scores[session.Path]
}
