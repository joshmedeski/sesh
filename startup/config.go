package startup

import "github.com/joshmedeski/sesh/v2/model"

func configStrategy(s *RealStartup, session model.SeshSession) (string, error) {
	config, exists := s.lister.FindConfigSession(session.Name)

	if exists && config.Tmuxinator != "" {
		return config.Tmuxinator, nil
	}

	if exists && config.StartupCommand != "" {
		return config.StartupCommand, nil
	}
	return "", nil
}
