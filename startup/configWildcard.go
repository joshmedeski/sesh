package startup

import (
	"github.com/joshmedeski/sesh/v2/model"
)

func configWildcardStrategy(s *RealStartup, session model.SeshSession) (string, error) {
	wildcard, exists := s.lister.FindConfigWildcard(session.Path)
	if exists && wildcard.StartupCommand != "" {
		return wildcard.StartupCommand, nil
	}
	return "", nil
}
