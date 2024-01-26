package session

import (
	"fmt"
	"log"

	"github.com/joshmedeski/sesh/config"
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

	name := DetermineName(path, config)
	if name == "" {
		log.Fatal("Couldn't determine the session name", err)
		return s, fmt.Errorf(
			"couldn't determine the session name for %q",
			choice,
		)
	}
	s.Name = name

	return s, nil
}
