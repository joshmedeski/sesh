package namer

import (
	"regexp"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

// applySubstitutions runs every rule against name in order, feeding each rule
// the output of the one before it, and returns the result. A literal rule
// replaces all occurrences of Find; a regex rule replaces every match and may
// use $1-style references to capture groups. An invalid regex is skipped so a
// single bad rule can't break naming outright — the config loader validates
// patterns up front, so this only guards against unexpected input.
func applySubstitutions(name string, rules []model.NameSubstitution) string {
	for _, rule := range rules {
		if rule.Find == "" {
			continue
		}
		if rule.Regex {
			re, err := regexp.Compile(rule.Find)
			if err != nil {
				continue
			}
			name = re.ReplaceAllString(name, rule.Replace)
			continue
		}
		name = strings.ReplaceAll(name, rule.Find, rule.Replace)
	}
	return name
}

// nameSubstitution is a naming strategy that applies the user's
// [[name_substitution]] rules to the home-collapsed path. It returns a name
// only when a rule actually changes the path; otherwise it returns an empty
// string so the namer falls through to the git and directory strategies. This
// keeps naming byte-identical to prior versions when no rule matches.
func nameSubstitution(n *RealNamer, path string) (string, error) {
	if len(n.config.NameSubstitutions) == 0 {
		return "", nil
	}
	shortened, err := n.home.ShortenHome(path)
	if err != nil {
		return "", err
	}
	name := applySubstitutions(shortened, n.config.NameSubstitutions)
	if name == shortened || name == "" {
		return "", nil
	}
	return name, nil
}
