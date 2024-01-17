package session

import (
	"log"

	"github.com/joshmedeski/sesh/config"
)

func Determine(choice string, config *config.Config) Session {
	path, err := DeterminePath(choice)
	if err != nil {
		log.Fatal("Couldn't determine the session path", err)
	}

	name := DetermineName(path, config)
	if name == "" {
		log.Fatal("Couldn't determine the session name", err)
	}

	return Session{
		Name: name,
		Path: path,
	}
}
