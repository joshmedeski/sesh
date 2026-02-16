package lister

import "regexp"

func compileBlacklist(patterns []string) []*regexp.Regexp {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		if re, err := regexp.Compile(pattern); err == nil {
			compiled = append(compiled, re)
		}
	}
	return compiled
}

func isBlacklisted(blacklist []*regexp.Regexp, name string) bool {
	for _, re := range blacklist {
		if re.MatchString(name) {
			return true
		}
	}
	return false
}
