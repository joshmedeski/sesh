package session

import (
	"fmt"
	"os"
)

func Determine(choice string) Session {
	path, err := DeterminPath(choice)
	if err != nil {
		fmt.Println("Couldn't determine the session path", err)
		os.Exit(1)
	}

	name := DetermineName(path)
	if name == "" {
		fmt.Println("Couldn't determine the session name", err)
		os.Exit(1)
	}

	return Session{
		Name: name,
		Path: path,
	}
}
