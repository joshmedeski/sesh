package lister

import (
	"fmt"
	"path/filepath"

	"github.com/joshmedeski/sesh/v2/model"
)

func configWildcardKey(name string) string {
	return fmt.Sprintf("config_wildcard:%s", name)
}

func listConfigWildcards(l *RealLister) (model.SeshWildcards, error) {
	orderedIndex := make([]string, 0)
	directory := make(model.SeshWildcardMap)
	for _, session := range l.config.WildcardConfigs {
		if session.Pattern != "" {
			key := configWildcardKey(session.Pattern)
			orderedIndex = append(orderedIndex, key)
			directory[key] = model.SeshWildcard{
				Src:                   "config",
				Pattern:               session.Pattern,
				StartupCommand:        session.StartupCommand,
				PreviewCommand:        session.PreviewCommand,
				DisableStartupCommand: session.DisableStartCommand,
			}
		}
	}
	return model.SeshWildcards{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) FindConfigWildcard(pattern string) (model.SeshWildcard, bool) {
	path, err := l.home.ExpandHome(pattern)
	if err != nil {
		return model.SeshWildcard{}, false
	}

	wildcards, _ := listConfigWildcards(l)
	for _, wildcard := range wildcards.Directory {
		expandedWildcard, err := l.home.ExpandHome(wildcard.Pattern)
		if err != nil {
			return model.SeshWildcard{}, false
		}
		if matched, _ := filepath.Match(expandedWildcard, path); matched {
			return wildcard, true
		}
	}
	return model.SeshWildcard{}, false
}
