package startup

import "github.com/joshmedeski/sesh/model"

func defaultConfigStrategy(s *RealStartup, session model.SeshSession) (string, error) {
	defaultConfig := s.config.DefaultSessionConfig
	if defaultConfig.StartupCommand != "" {
		return defaultConfig.StartupCommand, nil
	}
	return "", nil
}
