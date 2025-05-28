package lister

import "regexp"

func isBlacklisted(blacklist []string, name string) bool {
	for _, blacklistedName := range blacklist {
		if regexp.MustCompile(blacklistedName).MatchString(name) {
			return true
		}
	}
	return false
}
