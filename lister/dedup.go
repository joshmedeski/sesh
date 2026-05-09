package lister

import (
	"github.com/joshmedeski/sesh/v2/model"
)

type dedupKey string

const (
	keyName dedupKey = "name"
	keyPath dedupKey = "path"
)

type dedupRule struct {
	against []string
	key     dedupKey
}

// Per-source rules: each source declares which other sources cause it to be
// dropped on a given key, and which key is used for the comparison.
// Sources without a rule (tmux, tmux-pane) are always kept.
var dedupRules = map[string]dedupRule{
	"config":     {against: []string{"tmux"}, key: keyName},
	"tmuxinator": {against: []string{"tmux"}, key: keyName},
	// zoxide is not deduped against tmuxinator: tmuxinator entries do not
	// carry a Path (only Name), so a path-keyed rule would never match.
	// Revisit if tmuxinator's `root:` becomes available on SeshSession.
	"zoxide": {against: []string{"tmux", "config"}, key: keyPath},
}

// sourceTier orders sources by priority: tier 0 always wins, tier 1 may lose
// only to tier 0 sources, tier 2 may lose to tier 0 or 1. applyDedup walks
// tiers in order and only checks against values from KEPT earlier-tier
// sessions, so a dropped session never causes another session to be dropped.
// Unknown sources default to the highest tier (lowest priority); see
// tierOf.
var sourceTier = map[string]int{
	"tmux":       0,
	"tmux-pane":  0,
	"config":     1,
	"tmuxinator": 1,
	"zoxide":     2,
}

func tierOf(src string) int {
	if t, ok := sourceTier[src]; ok {
		return t
	}
	return maxDedupTier
}

var maxDedupTier = func() int {
	max := 0
	for _, t := range sourceTier {
		if t > max {
			max = t
		}
	}
	return max
}()

type srcKey struct {
	src string
	key dedupKey
}

func valueOf(k dedupKey, s model.SeshSession) string {
	switch k {
	case keyName:
		return s.Name
	case keyPath:
		return s.Path
	}
	return ""
}

func applyDedup(in model.SeshSessions) []string {
	keptSets := make(map[srcKey]map[string]struct{})
	keep := make(map[string]bool, len(in.OrderedIndex))

	for tier := 0; tier <= maxDedupTier; tier++ {
		for _, idx := range in.OrderedIndex {
			s := in.Directory[idx]
			if tierOf(s.Src) != tier {
				continue
			}
			if !shouldKeep(s, keptSets) {
				continue
			}
			keep[idx] = true
			recordValues(s, keptSets)
		}
	}

	out := make([]string, 0, len(keep))
	for _, idx := range in.OrderedIndex {
		if keep[idx] {
			out = append(out, idx)
		}
	}
	return out
}

func shouldKeep(s model.SeshSession, keptSets map[srcKey]map[string]struct{}) bool {
	rule, hasRule := dedupRules[s.Src]
	if !hasRule {
		return true
	}
	v := valueOf(rule.key, s)
	if v == "" {
		return true
	}
	for _, against := range rule.against {
		if _, found := keptSets[srcKey{against, rule.key}][v]; found {
			return false
		}
	}
	return true
}

func recordValues(s model.SeshSession, keptSets map[srcKey]map[string]struct{}) {
	for _, k := range []dedupKey{keyName, keyPath} {
		v := valueOf(k, s)
		if v == "" {
			continue
		}
		sk := srcKey{s.Src, k}
		if keptSets[sk] == nil {
			keptSets[sk] = make(map[string]struct{})
		}
		keptSets[sk][v] = struct{}{}
	}
}
