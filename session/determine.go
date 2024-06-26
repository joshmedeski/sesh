package session

import (
	"fmt"
	"log"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/name"
)

func isConfigSession(choice string) *Session {
	config := config.ParseConfigFile(&config.DefaultConfigDirectoryFetcher{})
	for _, sessionConfig := range config.SessionConfigs {
		if sessionConfig.Name == choice {
			return &Session{
				Src:  "config",
				Name: sessionConfig.Name,
				Path: sessionConfig.Path,
			}
		}
	}
	return nil
}

func Determine(choice string, config *config.Config) (s Session, err error) {
	configSession := isConfigSession(choice)
	if configSession != nil {
		return *configSession, nil
	}

	path, err := DeterminePath(choice)
	if err != nil {
		return s, fmt.Errorf(
			"couldn't determine the path for %q: %w",
			choice,
			err,
		)
	}
	s.Path = path

	sessionName := name.DetermineName(choice, path)
	if sessionName == "" {
		log.Fatal("Couldn't determine the session name", err)
		return s, fmt.Errorf(
			"couldn't determine the session name for %q",
			choice,
		)
	}
	s.Name = sessionName

	return s, nil
}
