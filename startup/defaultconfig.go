package startup

import "github.com/joshmedeski/sesh/v2/model"

func defaultConfigStrategy(s *RealStartup, session model.SeshSession) (string, error) {
	if session.DisableStartupCommand {
		return "", nil
	}

	defaultConfig := s.config.DefaultSessionConfig
	if defaultConfig.StartupCommand != "" {
		replacements := map[string]string{
			"{}": session.Path,
		}

		return s.replacer.Replace(defaultConfig.StartupCommand, replacements), nil
	}

	return "", nil
}
