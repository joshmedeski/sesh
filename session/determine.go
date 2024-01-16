package session

import (
	"log"
)

func Determine(choice string) Session {
	path, err := DeterminePath(choice)
	if err != nil {
		log.Fatal("Couldn't determine the session path", err)
	}

	name := DetermineName(path)
	if name == "" {
		log.Fatal("Couldn't determine the session name", err)
	}

	return Session{
		Name: name,
		Path: path,
	}
}
