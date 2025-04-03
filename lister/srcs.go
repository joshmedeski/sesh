package lister

import (
	"cmp"
	"math"
	"slices"
	"strings"
)

// In-place sorting of 'sources' based on given order of elements 'desiredOrder'
//
// # Omitted elements are placed after given elements
//
// # Duplicate elements use last-most given order
//
// Example:
//
//	sources := []string{"a", "b", "c", "x"}
//	desiredOrder := []string{"b", "a", "c"}
//	sortSources(sources, desiredOrder)
//	// sources is now []string{"b", "a", "c", "x"}
func sortSources(sources, desiredOrder []string) {
	if desiredOrder == nil || len(desiredOrder) == 0 {
		return
	}
	m := make(map[string]int)
	for i, s := range desiredOrder {
		m[strings.ToLower(s)] = i
	}
	getOrder := func(s string) int {
		order, exists := m[strings.ToLower(s)]
		if !exists {
			return math.MaxInt
		}
		return order
	}
	slices.SortFunc(sources, func(a, b string) int {
		return cmp.Compare(getOrder(a), getOrder(b))
	})
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
	// TODO: if we have configured sorting, do so now
	sortSources(srcs, []string{"tmuxinator", "tmux", "config", "zoxide"})

	return srcs
}
