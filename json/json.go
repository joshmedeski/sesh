package json

import (
	"encoding/json"
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

type Json interface {
	EncodeSessions(sessions []model.SeshSession) string
}

type RealJson struct{}

func NewJson() Json {
	return &RealJson{}
}

func (j *RealJson) EncodeSessions(sessions []model.SeshSession) string {
	jsonSessions, err := json.Marshal(sessions)
	if err != nil {
		fmt.Printf(
			"Couldn't list sessions as json: %s\n",
			err,
		)
	}
	return string(jsonSessions)
}
