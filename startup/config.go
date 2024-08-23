package startup

import "github.com/joshmedeski/sesh/model"

func configStrategy(s *RealStartup, session model.SeshSession) (string, error) {
	config, exists := s.lister.FindConfigSession(session.Name)
	if exists && config.StartupCommand != "" {
		return config.StartupCommand, nil
	}
	return "", nil
}
