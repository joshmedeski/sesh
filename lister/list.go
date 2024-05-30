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

func (s *RealLister) List(opts ListOptions) (model.SeshSessionMap, error) {
	allSessions := make(model.SeshSessionMap)
	srcs := srcs(opts)

	if srcs["tmux"] {
		sessions, err := listTmuxSessions(s.tmux)
		if err != nil {
			return nil, err
		}
		for k, s := range sessions {
			allSessions[k] = s
		}
	}

	if srcs["config"] {
		sessions := listConfigSessions(s.config)
		for k, s := range sessions {
			if s.Name != "" {
				allSessions[k] = s
			}
		}
	}

	if srcs["zoxide"] {
		sessions, err := listZoxideSessions(s.zoxide, s.home)
		if err != nil {
			return nil, err
		}
		for k, s := range sessions {
			allSessions[k] = s
		}
	}

	return allSessions, nil
}

func srcs(opts ListOptions) map[string]bool {
	if !opts.Config && !opts.Tmux && !opts.Zoxide {
		// show all sources by default
		return map[string]bool{"config": true, "tmux": true, "zoxide": true}
	} else {
		return map[string]bool{"config": opts.Config, "tmux": opts.Tmux, "zoxide": opts.Zoxide}
	}
}
