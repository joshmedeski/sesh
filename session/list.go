package session

import (
	"fmt"

	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
)

type ListOptions struct {
	Config       bool
	HideAttached bool
	Icons        bool
	Json         bool
	Tmux         bool
	Zoxide       bool
}

func (s *RealSession) List(opts ListOptions) ([]model.SeshSession, error) {
	list := []model.SeshSession{}
	srcs := srcs(opts)

	if srcs["tmux"] {
		tmuxList, err := listTmuxSessions(s.tmux)
		if err != nil {
			return nil, err
		}
		list = append(list, tmuxList...)
	}

	if srcs["zoxide"] {
		zoxideList, err := listZoxideResults(s.zoxide, s.home)
		if err != nil {
			return nil, err
		}
		list = append(list, zoxideList...)
	}

	return list, nil
}

func srcs(opts ListOptions) map[string]bool {
	if !opts.Config && !opts.Tmux && !opts.Zoxide {
		// show all sources by default
		return map[string]bool{"config": true, "tmux": true, "zoxide": true}
	} else {
		return map[string]bool{"config": opts.Config, "tmux": opts.Tmux, "zoxide": opts.Zoxide}
	}
}

func listTmuxSessions(t tmux.Tmux) ([]model.SeshSession, error) {
	tmuxSessions, err := t.ListSessions()
	if err != nil {
		return nil, fmt.Errorf("couldn't list tmux sessions: %q", err)
	}
	sessions := make([]model.SeshSession, len(tmuxSessions))
	for i, session := range tmuxSessions {
		sessions[i] = model.SeshSession{
			Src: "tmux",
			// TODO: prepend icon if configured
			Name:     session.Name,
			Path:     session.Path,
			Attached: session.Attached,
			Windows:  session.Windows,
		}
	}
	return sessions, nil
}

func listZoxideResults(z zoxide.Zoxide, h home.Home) ([]model.SeshSession, error) {
	zoxideResults, err := z.ListResults()
	if err != nil {
		return nil, fmt.Errorf("couldn't list zoxide results: %q", err)
	}
	sessions := make([]model.SeshSession, len(zoxideResults))
	for i, r := range zoxideResults {
		name, err := h.ShortenHome(r.Path)
		if err != nil {
			return nil, fmt.Errorf("couldn't shorten path: %q", err)
		}
		sessions[i] = model.SeshSession{
			Src: "zoxide",
			// TODO: prepend icon if configured
			Name:  name,
			Path:  r.Path,
			Score: r.Score,
		}
	}
	return sessions, nil
}
