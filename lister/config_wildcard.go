package lister

import (
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

func (l *RealLister) FindConfigWildcard(path string) (model.WildcardConfig, bool) {
	expandedPath, err := l.home.ExpandHome(path)
	if err != nil {
		return model.WildcardConfig{}, false
	}

	for _, wc := range l.config.WildcardConfigs {
		expandedPattern, err := l.home.ExpandHome(wc.Pattern)
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
	if strings.HasSuffix(pattern, "/**") {
		prefix := strings.TrimSuffix(pattern, "/**")
		return strings.HasPrefix(path, prefix+"/")
	}
	matched, err := filepath.Match(pattern, path)
	return err == nil && matched
}
