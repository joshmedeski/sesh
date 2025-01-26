package lister

import (
	"path/filepath"
	"slices"
	"sort"

	"github.com/joshmedeski/sesh/model"
	"github.com/samber/lo"
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
	"tmuxinator": listTmuxinator,
	"zoxide":     listZoxide,
}

const zoxideFactor = 100

func (l *RealLister) List(opts ListOptions) (model.SeshSessions, error) {
	allSessions := sessionFromSources(opts, l)

	allSessions = removeConfigDuplicates(allSessions)
	allSessions = removeZoxideDuplicates(allSessions)

	if opts.HideAttached {
		allSessions = removeAttachedSession(l, allSessions)
	}

	return createSeshSessions(allSessions)
}

func sessionFromSources(opts ListOptions, l *RealLister) []model.SeshSession {
	allSessions := lo.FlatMap(srcs(opts), func(src string, i int) []model.SeshSession {
		sessions, err := srcStrategies[src](l)
		if err != nil {
			return nil
		}

		return lo.Map(lo.Values(sessions.Directory), func(session model.SeshSession, j int) model.SeshSession {
			if session.Src != "zoxide" {
				session.Score = session.Score * zoxideFactor
			}
			return session
		})
	})

	return allSessions
}

func removeConfigDuplicates(allSessions []model.SeshSession) []model.SeshSession {
	configDuplicates := lo.GroupBy(allSessions, func(s model.SeshSession) string {
		return s.Name
	})

	return lo.MapToSlice(configDuplicates, func(_ string, sessions []model.SeshSession) model.SeshSession {
		return lo.MaxBy(sessions, func(a, b model.SeshSession) bool {
			return a.Score > b.Score
		})
	})
}

func removeZoxideDuplicates(allSessions []model.SeshSession) []model.SeshSession {
	runningSessionNames := lo.FilterMap(allSessions, func(s model.SeshSession, _ int) (string, bool) {
		return s.Name, s.Src == "tmux"
	})

	return lo.Filter(allSessions, func(s model.SeshSession, _ int) bool {
		return s.Src != "zoxide" || !slices.Contains(runningSessionNames, filepath.Base(s.Path))
	})
}

func removeAttachedSession(l *RealLister, allSessions []model.SeshSession) []model.SeshSession {
	attachedSession, _ := GetAttachedTmuxSession(l)
	allSessions = lo.Filter(allSessions, func(s model.SeshSession, _ int) bool {
		return s.Name != attachedSession.Name
	})
	return allSessions
}

func createSeshSessions(allSessions []model.SeshSession) (model.SeshSessions, error) {
	sort.Slice(allSessions, func(i, j int) bool {
		return allSessions[i].Score > allSessions[j].Score
	})

	orderedIndex := lo.Map(allSessions, func(s model.SeshSession, _ int) string {
		return s.Src + s.Name
	})
	directory := lo.KeyBy(allSessions, func(s model.SeshSession) string {
		return s.Src + s.Name
	})

	return model.SeshSessions{
		OrderedIndex: orderedIndex,
		Directory:    directory,
	}, nil
}
