package lister

import (
	"cmp"
	"math"
	"slices"
	"strings"
)

// returns a sorted list of sources based on the provided sort order.
func sortSources(sources, sortOrder []string) []string {
	if sortOrder == nil || len(sortOrder) == 0 {
		return sources
	}
	m := make(map[string]int)
	for i, s := range sortOrder {
		m[strings.ToLower(s)] = i
	}
	getOrder := func(s string) int {
		if order, exists := m[strings.ToLower(s)]; exists {
			return order
		} else {
			return math.MaxInt
		}
	}
	result := slices.Clone(sources)
	slices.SortStableFunc(result, func(a, b string) int {
		return cmp.Compare(getOrder(a), getOrder(b))
	})
	return result
}

func srcs(opts ListOptions) []string {
	var srcs []string
	count := 0
	if opts.Tmux {
		count++
	}
	if opts.Config {
		count++
	}
	if opts.Tmuxinator {
		count++
	}
	if opts.Zoxide {
		count++
	}
	if count == 0 {
		return []string{"tmux", "config", "tmuxinator", "zoxide"}
	}
	srcs = make([]string, count)
	i := 0
	if opts.Tmux {
		srcs[i] = "tmux"
		i++
	}
	if opts.Config {
		srcs[i] = "config"
		i++
	}
	if opts.Tmuxinator {
		srcs[i] = "tmuxinator"
		i++
	}
	if opts.Zoxide {
		srcs[i] = "zoxide"
		i++
	}
	return srcs
}
