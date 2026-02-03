package lister

import (
	"path/filepath"

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
		matched, err := filepath.Match(expandedPattern, expandedPath)
		if err != nil {
			continue
		}
		if matched {
			return wc, true
		}
	}
	return model.WildcardConfig{}, false
}
