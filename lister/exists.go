package lister

import "github.com/joshmedeski/sesh/v2/model"

func exists(key string, sessions map[string]model.SeshSession) (model.SeshSession, bool) {
	if session, exists := sessions[key]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
