package main

import (
	"strconv"
	"strings"
	"unicode"
)

// A benchmark name is written for `go test -bench=`, not for someone skimming
// a pull request: "picker FilterSessions/n=1000/separator-aware/no-match" says
// nothing about what slowed down until you have read the benchmark. The maps
// below turn one into a sentence.
//
// Anything missing from them still renders — the base name is split on its
// capitals and the parameters are passed through — so a new benchmark reads
// awkwardly rather than not at all. TestEveryBenchmarkIsNamed keeps that from
// being the normal case.

// what a benchmark exercises. "{n}" takes the session count.
var what = map[string]string{
	"ApplyDedup":                "Deduplicating {n}",
	"BlacklistFilter":           "Blacklisting {n}",
	"BuildItems":                "Building picker rows for {n}",
	"CachingListerApplyFilters": "Re-filtering {n} on every run",
	"FilterAliases":             "Filtering aliases, {n}",
	"FilterSessions":            "Filtering {n} on a keystroke",
	"FindTmuxSessionByBase":     "Finding a tmux session among {n}",
	"IconResolverWildcard":      "Resolving a wildcard icon, {n}",
	"Keystroke":                 "A keystroke in the picker, {n}",
	"List":                      "Listing {n}",
	"ListBlacklist":             "Listing {n} through the blacklist",
	"ListHideDuplicates":        "Listing {n} without duplicates",
	"ListShowWindows":           "Listing {n} with windows",
	"RankMatches":               "Ranking matches from {n}",
	"View":                      "Redrawing the picker, {n}",
}

// conditions name the sub-benchmark parameters. A key may be qualified by its
// benchmark, because the same token means different things in different ones:
// "plain" is a matching mode under FilterSessions and a row style under
// BuildItems. An empty value drops the condition, for the ones that only say
// "the default".
var conditions = map[string]string{
	"1-char":                     "1 char typed",
	"3-char":                     "3 chars typed",
	"a":                          "1 char typed",
	"al7":                        "whole alias typed",
	"app":                        "3 chars typed",
	"empty":                      "nothing typed",
	"hide-attached":              "hiding attached",
	"hide-duplicates":            "hiding duplicates",
	"hit-last":                   "last pattern hit",
	"icons":                      "with icons",
	"miss":                       "no pattern hit",
	"no-match":                   "no matches",
	"passthrough":                "no filters",
	"picker-defaults":            "picker defaults",
	"plain":                      "",
	"separator-aware":            "separator-aware",
	"source-filter":              "filtered by source",
	"BuildItems/plain":           "plain rows",
	"BuildItems/separator-aware": "separator-aware rows",
}

// phrase turns a benchmark's name into something readable: what it does, with
// the session count folded in, and whatever conditions distinguish it from its
// siblings in parentheses.
func phrase(name string) string {
	base, parts := splitName(name)

	sessions, rest := countOf(parts)
	head, ok := what[base]
	if !ok {
		head = deCamel(base) + ", {n}"
	}
	head = strings.ReplaceAll(head, "{n}", sessions)

	var conds []string
	for _, p := range rest {
		c, ok := conditions[base+"/"+p]
		if !ok {
			if c, ok = conditions[p]; !ok {
				c = p
			}
		}
		if c != "" {
			conds = append(conds, c)
		}
	}
	if len(conds) == 0 {
		return head
	}
	return head + " (" + strings.Join(conds, ", ") + ")"
}

// splitName separates "FilterSessions/n=1000/plain" into its base name and its
// sub-benchmark parameters.
func splitName(name string) (string, []string) {
	parts := strings.Split(name, "/")
	return parts[0], parts[1:]
}

// countOf pulls the "n=1000" parameter out and renders it as a phrase. A
// benchmark without one still gets a subject to attach the sentence to.
func countOf(parts []string) (string, []string) {
	rest := make([]string, 0, len(parts))
	sessions := "the session list"
	for _, p := range parts {
		if n, ok := strings.CutPrefix(p, "n="); ok {
			sessions = commas(n) + " sessions"
			continue
		}
		rest = append(rest, p)
	}
	return sessions, rest
}

// commas groups a number for reading: 1000 is a size, 1,000 is a quantity, and
// the whole point here is to read it as a quantity.
func commas(digits string) string {
	if _, err := strconv.Atoi(digits); err != nil {
		return digits
	}
	var out []byte
	for i, c := range []byte(digits) {
		if i > 0 && (len(digits)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, c)
	}
	return string(out)
}

// deCamel is the fallback for a benchmark nobody has named yet:
// "FindTmuxSessionByBase" becomes "Find tmux session by base", which at least
// reads as words.
func deCamel(s string) string {
	var words []string
	start := 0
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			words = append(words, s[start:i])
			start = i
		}
	}
	words = append(words, s[start:])
	for i, w := range words {
		if i == 0 {
			continue
		}
		// An acronym-ish run stays as it is; anything else is mid-sentence.
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, " ")
}
