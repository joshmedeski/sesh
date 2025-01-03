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


func (l *RealLister) List(opts ListOptions) (model.SeshSessions, error) {
	allSessions := lo.FlatMap(srcs(opts), func(src string, i int) []model.SeshSession {
		sessions, err := srcStrategies[src](l)
		if err != nil {
			return nil
		}

		return lo.Map(lo.Values(sessions.Directory), func(session model.SeshSession, j int) model.SeshSession {
			return session
		})
	})

	if opts.HideAttached {
		attachedSession, _ := GetAttachedTmuxSession(l)
		allSessions = lo.Filter(allSessions, func(s model.SeshSession, _ int) bool {
			return s.Name != attachedSession.Name
		})
	}

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
