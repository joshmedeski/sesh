package replacer

import (
	ahocorasick "github.com/petar-dambovaliev/aho-corasick"
)

type Replacer interface {
	Replace(command string, replacements map[string]string) string
}

type RealReplacer struct{}

func NewReplacer() Replacer {
	return &RealReplacer{}
}

func (r *RealReplacer) Replace(command string, replacements map[string]string) string {
	dict := make([]string, 0, len(replacements))
	replacementArray := make([]string, 0, len(replacements))
	for k, v := range replacements {
		dict = append(dict, k)
		replacementArray = append(replacementArray, v)
	}
	ac := getAhoCorasick(dict)
	replacer := ahocorasick.NewReplacer(ac)
	return replacer.ReplaceAll(command, replacementArray)
}

func getAhoCorasick(dictionary []string) ahocorasick.AhoCorasick {
	builder := ahocorasick.NewAhoCorasickBuilder(ahocorasick.Opts{
		AsciiCaseInsensitive: true,
		MatchOnlyWholeWords:  true,
		MatchKind:            ahocorasick.LeftMostFirstMatch,
		DFA:                  true,
	})

	return builder.Build(dictionary)
}
