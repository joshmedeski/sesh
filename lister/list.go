package lister

import (
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

const (
	srcOffset = 1000000
	srcFactor = 10000
)

func (l *RealLister) List(opts ListOptions) (model.SeshSessions, error) {
	allSessions := lo.FlatMap(srcs(opts), func(src string, i int) []model.SeshSession {
		sessions, err := srcStrategies[src](l)
		if err != nil {
			return nil
		}

		return lo.Map(lo.Values(sessions.Directory), func(session model.SeshSession, j int) model.SeshSession {
			if session.Src != "zoxide" {
				srcDownrank := float64(i) * srcFactor
				sessionDownrank := float64(j)
				session.Score = session.Score + srcOffset - srcDownrank - sessionDownrank
			}
			return session
		})
	})

	if opts.HideAttached {
		attachedSession, _ := GetAttachedTmuxSession(l)
		allSessions = lo.Filter(allSessions, func(s model.SeshSession, _ int) bool {
			return s.Name != attachedSession.Name
		})
	}

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
