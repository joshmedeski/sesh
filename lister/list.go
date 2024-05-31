package lister

import (
	"github.com/joshmedeski/sesh/model"
)

type ListOptions struct {
	Config       bool
	HideAttached bool
	Icons        bool
	Json         bool
	Tmux         bool
	Zoxide       bool
}

func (s *RealLister) List(opts ListOptions) (model.SeshSessions, error) {
	fullDirectory := make(model.SeshSessionMap)
	fullOrderedIndex := make([]string, 0)
	srcs := srcs(opts)

	if srcs["tmux"] {
		tmuxSessions, err := listTmuxSessions(s.tmux)
		if err != nil {
			return model.SeshSessions{}, err
		}
		fullOrderedIndex = append(fullOrderedIndex, tmuxSessions.OrderedIndex...)
		for _, i := range tmuxSessions.OrderedIndex {
			fullDirectory[i] = tmuxSessions.Directory[i]
		}
	}

	if srcs["config"] {
		configSessions := listConfigSessions(s.config)
		fullOrderedIndex = append(fullOrderedIndex, configSessions.OrderedIndex...)
		for _, i := range configSessions.OrderedIndex {
			fullDirectory[i] = configSessions.Directory[i]
		}
	}

	if srcs["zoxide"] {
		sessions, err := listZoxideSessions(s.zoxide, s.home)
		if err != nil {
			return model.SeshSessions{}, err
		}
		for k, s := range sessions {
			fullDirectory[k] = s
		}
	}

	return model.SeshSessions{
		OrderedIndex: fullOrderedIndex,
		Directory:    fullDirectory,
	}, nil
}

func srcs(opts ListOptions) map[string]bool {
	if !opts.Config && !opts.Tmux && !opts.Zoxide {
		// show all sources by default
		return map[string]bool{"config": true, "tmux": true, "zoxide": true}
	} else {
		return map[string]bool{"config": opts.Config, "tmux": opts.Tmux, "zoxide": opts.Zoxide}
	}
}
