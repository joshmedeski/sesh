package startup

import "github.com/joshmedeski/sesh/v2/model"

func configWildcardStartupStrategy(s *RealStartup, session model.SeshSession) (string, error) {
	wc, found := s.lister.FindConfigWildcard(session.Path)
	if !found {
		return "", nil
	}

	if wc.DisableStartCommand {
		return "", nil
	}

	if wc.StartupCommand != "" {
		replacements := map[string]string{
			"{}": session.Path,
		}
		return s.replacer.Replace(wc.StartupCommand, replacements), nil
	}

	return "", nil
}
