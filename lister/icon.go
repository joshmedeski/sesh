package lister

import (
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
)

// WildcardFinder resolves a path to the [[wildcard]] block that matches it. The
// icon resolver borrows the lister's matching so an icon covers exactly the
// paths that startup_command and preview_command already do, first match in
// config order included.
type WildcardFinder interface {
	FindConfigWildcard(path string) (model.WildcardConfig, bool)
}

// IconResolver indexes the icons declared in config into a lookup shared by the
// picker and `sesh list --json`. It returns nil when nothing declares one, so
// the common case costs nothing per row.
//
// The most specific match wins: an exact [[session]] name, then a [[session]]
// path, then a path under a [[worktree]] root, then a [[wildcard]] pattern. Sessions are indexed by path as well as by
// name because the same directory is often listed by another source under a
// derived name — a zoxide entry for a configured session's path still gets its
// icon.
func IconResolver(config model.Config, h home.Home, wildcards WildcardFinder) func(model.SeshSession) string {
	byName := make(map[string]string)
	byPath := make(map[string]string)
	for _, session := range config.SessionConfigs {
		if session.Icon == "" {
			continue
		}
		if session.Name != "" {
			byName[session.Name] = session.Icon
		}
		if session.Path == "" {
			continue
		}
		// An unexpandable path is skipped rather than fatal: the icon is
		// cosmetic, and the lister already reports a broken path.
		if path, err := h.ExpandPath(session.Path); err == nil {
			byPath[filepath.Clean(path)] = session.Icon
		}
	}

	var worktrees []worktreeRoot
	for _, root := range worktreeRoots(config, h) {
		if root.config.Icon != "" {
			worktrees = append(worktrees, root)
		}
	}

	hasWildcardIcon := false
	for _, wildcard := range config.WildcardConfigs {
		if wildcard.Icon != "" {
			hasWildcardIcon = true
			break
		}
	}

	if len(byName) == 0 && len(byPath) == 0 && len(worktrees) == 0 && !hasWildcardIcon {
		return nil
	}

	return func(session model.SeshSession) string {
		if icn, ok := byName[session.Name]; ok {
			return icn
		}
		if session.Path == "" {
			return ""
		}
		path := filepath.Clean(session.Path)
		if icn, ok := byPath[path]; ok {
			return icn
		}
		for _, worktree := range worktrees {
			if strings.HasPrefix(path, worktree.prefix) {
				return worktree.config.Icon
			}
		}
		if !hasWildcardIcon || wildcards == nil {
			return ""
		}
		// A first match that declares no icon still wins, matching how the rest
		// of the wildcard config resolves: the session falls back to its source
		// glyph rather than searching on for a later pattern with an icon.
		if wildcard, ok := wildcards.FindConfigWildcard(session.Path); ok {
			return wildcard.Icon
		}
		return ""
	}
}
