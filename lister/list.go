package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/model"
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

func (s *RealLister) List(opts ListOptions) (model.SeshSessionMap, error) {
	allSessions := make(model.SeshSessionMap)
	srcs := srcs(opts)

	if srcs["tmux"] {
		tmuxSessions, err := listTmuxSessions(s.tmux)
		if err != nil {
			return nil, err
		}
		for _, s := range tmuxSessions {
			key := fmt.Sprintf("tmux:%s", s.Name)
			allSessions[key] = s
		}
	}

	if srcs["config"] {
		configList, err := listConfigSessions(s.config)
		if err != nil {
			return nil, err
		}
		for _, s := range configList {
			if s.Name != "" {
				key := fmt.Sprintf("config:%s", s.Name)
				allSessions[key] = s
			}
		}
	}

	if srcs["zoxide"] {
		zoxideList, err := listZoxideResults(s.zoxide, s.home)
		if err != nil {
			return nil, err
		}
		for _, s := range zoxideList {
			key := fmt.Sprintf("zoxide:%s", s.Name)
			allSessions[key] = s
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

func listConfigSessions(c model.Config) ([]model.SeshSession, error) {
	var configSessions []model.SeshSession
	for _, session := range c.SessionConfigs {
		if session.Name != "" {
			configSessions = append(configSessions, model.SeshSession{
				Src:  "config",
				Name: session.Name,
				Path: session.Path,
			})
		}
	}
	return configSessions, nil
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
