package json

import (
	"encoding/json"
	"fmt"

	"github.com/joshmedeski/sesh/session"
)

func List(sessions []session.Session) string {
	jsonSessions, err := json.Marshal(sessions)
	if err != nil {
		fmt.Printf(
			"Couldn't list sessions as json: %s\n",
			err,
		)
	}
	return string(jsonSessions)
}
