package session

import (
	"fmt"
	"log"
)

func Determine(choice string) (s Session, err error) {
	path, err := DeterminePath(choice)
	if err != nil {
		return s, fmt.Errorf(
			"couldn't determine the path for %q: %w",
			choice,
			err,
		)
	}
	s.Path = path

	name := DetermineName(path)
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
