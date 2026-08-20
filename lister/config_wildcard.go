package lister

import (
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

// expandedWildcard pairs a [[wildcard]] block with its home-expanded pattern so
// the expansion happens once per lister instead of once per lookup.
type expandedWildcard struct {
	config  model.WildcardConfig
	pattern string
}

// expandedWildcards expands every configured wildcard pattern the first time it
// is asked for and caches the result. The patterns come from config, so they
// cannot change during a run, but the picker asks about every session on screen
// — expanding per call made that O(sessions x patterns).
//
// Expansion is lazy rather than done in NewLister so a lister built for a
// command that never resolves a wildcard pays nothing.
//
// A pattern that fails to expand is dropped, which keeps the previous
// per-lookup `continue` behavior, and the surviving patterns stay in config
// order so first-match-wins is unchanged.
func (l *RealLister) expandedWildcards() []expandedWildcard {
	l.wildcardsOnce.Do(func() {
		expanded := make([]expandedWildcard, 0, len(l.config.WildcardConfigs))
		for _, wc := range l.config.WildcardConfigs {
			pattern, err := l.home.ExpandPath(wc.Pattern)
			if err != nil {
				continue
			}
			expanded = append(expanded, expandedWildcard{config: wc, pattern: pattern})
		}
		l.wildcards = expanded
	})
	return l.wildcards
}

func (l *RealLister) FindConfigWildcard(path string) (model.WildcardConfig, bool) {
	wildcards := l.expandedWildcards()
	if len(wildcards) == 0 {
		return model.WildcardConfig{}, false
	}

	expandedPath, err := l.home.ExpandPath(path)
	if err != nil {
		return model.WildcardConfig{}, false
	}

	for _, wc := range wildcards {
		if matchWildcard(wc.pattern, expandedPath) {
			return wc.config, true
		}
	}
	return model.WildcardConfig{}, false
}

func matchWildcard(pattern, path string) bool {
	cleanPath := filepath.Clean(path)

	if strings.HasSuffix(pattern, "/**") {
		prefix := strings.TrimSuffix(pattern, "/**")
		if !strings.HasPrefix(cleanPath, prefix+"/") {
			return false
		}
		return len(cleanPath) > len(prefix)+1
	}

	matched, err := filepath.Match(pattern, cleanPath)
	return err == nil && matched
}
