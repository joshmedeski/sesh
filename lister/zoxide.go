package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/zoxide"
)

func zoxideKey(name string) string {
	return fmt.Sprintf("zoxide:%s", name)
}

func listZoxideSessions(z zoxide.Zoxide, h home.Home) (model.SeshSessionMap, error) {
	zoxideResults, err := z.ListResults()
	if err != nil {
		return nil, fmt.Errorf("couldn't list zoxide sessions: %q", err)
	}
	sessions := make(model.SeshSessionMap)
	for _, r := range zoxideResults {
		name, err := h.ShortenHome(r.Path)
		if err != nil {
			return nil, fmt.Errorf("couldn't shorten path: %q", err)
		}
		key := zoxideKey(name)
		sessions[key] = model.SeshSession{
			Src:   "zoxide",
			Name:  name,
			Path:  r.Path,
			Score: r.Score,
		}
	}
	return sessions, nil
}

func (l *RealLister) FindZoxideSession(name string) (model.SeshSession, bool) {
	sessions, err := listZoxideSessions(l.zoxide, l.home)
	if err != nil {
		return model.SeshSession{}, false
	}
	if session, exists := sessions[name]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
