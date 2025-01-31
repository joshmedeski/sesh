package lister

import "strings"

func isBlacklisted(blacklist []string, name string) bool {
	for _, blacklistedName := range blacklist {
		if strings.EqualFold(blacklistedName, name) {
			return true
		}
	}
	return false
}

func createBlacklistSet(blacklist []string) map[string]struct{} {
	blacklistSet := make(map[string]struct{}, len(blacklist))
	for _, name := range blacklist {
		blacklistSet[name] = struct{}{}
	}
	return blacklistSet
}
