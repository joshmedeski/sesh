package startup

import "github.com/joshmedeski/sesh/v2/model"

func defaultConfigStrategy(s *RealStartup, session model.SeshSession) (string, error) {
	if session.DisableStartupCommand {
		return "", nil
	}

	defaultConfig := s.config.DefaultSessionConfig
	if defaultConfig.StartupCommand != "" {
		return defaultConfig.StartupCommand, nil
	}

	return "", nil
}
