package lister

import (
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

func (l *RealLister) FindConfigWildcard(path string) (model.WildcardConfig, bool) {
	expandedPath, err := l.home.ExpandPath(path)
	if err != nil {
		return model.WildcardConfig{}, false
	}

	for _, wc := range l.config.WildcardConfigs {
		expandedPattern, err := l.home.ExpandPath(wc.Pattern)
		if err != nil {
			continue
		}
		if matchWildcard(expandedPattern, expandedPath) {
			return wc, true
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
