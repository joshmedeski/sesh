package session

import (
	"fmt"
	"log"

	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/name"
)

func Determine(choice string, config *config.Config) (s Session, err error) {
	path, err := DeterminePath(choice)
	if err != nil {
		return s, fmt.Errorf(
			"couldn't determine the path for %q: %w",
			choice,
			err,
		)
	}
	s.Path = path

	sessionName := name.DetermineName(path)
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
