package dashboard

import (
	"log/slog"
	"strings"

	"github.com/joshmedeski/sesh/v2/dashboard/core"
	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/dashboard/sections"
	"github.com/joshmedeski/sesh/v2/model"
)

// Section, Sorter, Filterer, Clicker, and SectionDeps are defined in
// dashboard/core and re-exported here so existing callers keep working.
// CommandRunner is re-exported from exec.go.
type Section = core.Section
type Sorter = core.Sorter
type Filterer = core.Filterer
type Clicker = core.Clicker
type SectionDeps = core.SectionDeps

func iconCol(deps SectionDeps, sess model.SeshSession, fallbackSrc string, selected bool) render.Col {
	custom := ""
	if deps.Icon != nil {
		custom = deps.Icon(sess)
	}
	src := sess.Src
	if fallbackSrc != "" {
		src = fallbackSrc
	}
	return render.IconCol(custom, src, max(deps.IconWidth, 1), selected)
}

type SectionFactory func(cfg model.DashboardSectionConfig, deps SectionDeps) Section

type Registry map[string]SectionFactory

// registry maps configurable section types to their factories. "sessions" is
// kept as an alias of an unconfigured "sources" section (tmux sessions only).
var registry = Registry{
	// "details": sections.NewDetailsSection,
	"system":   sections.NewSystemSection,
	"ssh":      sections.NewSSHSection,
	"git":      sections.NewGitSection,
	"custom":   sections.NewCustomSection,
	"docker":   sections.NewDockerSection,
	"workmux":  sections.NewWorkmuxSection,
	"sources":  NewSourcesSection,
	"sessions": NewSourcesSection,
}

// BuiltSections is the result of BuildSections: the Configured tab's list plus
// the first page's sections.
type BuiltSections struct {
	Configured *ConfiguredSection
	Widgets    []Section
}

// BuildSections builds the Configured list and the first page from the
// `[[dashboard.section]]` entries, in config order. A "worktree" entry lists
// the [[worktree]] block whose repo matches its own. Unknown types and
// unmatched repos are logged and skipped. With no usable entries the first
// page is a single tmux sessions list.
func BuildSections(cfg model.DashboardConfig, worktrees []model.WorktreeConfig, deps SectionDeps) BuiltSections {
	var widgets []Section
	for _, sc := range cfg.Sections {
		if sc.Type == "" {
			slog.Warn("unknown dashboard section type")
			continue
		}
		if sc.Type == "worktree" {
			wc, ok := findWorktreeConfig(worktrees, sc.Repo)
			if !ok {
				slog.Warn("dashboard worktree section has no matching [[worktree]] block", "repo", sc.Repo)
				continue
			}
			ws := NewWorktreeSection(wc, deps)
			ws.title = sc.Title
			widgets = append(widgets, ws)
			continue
		}
		factory, ok := registry[sc.Type]
		if !ok {
			if sc.Type == "aiagent" {
				slog.Warn("unknown dashboard section type", "type", sc.Type, "hint", "aiagent is deprecated; use workmux")
			} else {
				slog.Warn("unknown dashboard section type", "type", sc.Type)
			}
			continue
		}
		if sc.Groups != nil {
			slog.Warn("dashboard section groups are no longer applied", "type", sc.Type)
		}
		widgets = append(widgets, factory(sc, deps))
	}
	if len(widgets) == 0 {
		widgets = []Section{NewSourcesSection(model.DashboardSectionConfig{Type: "sources", Title: "Sessions"}, deps)}
	}

	configured := NewConfiguredSection(
		model.DashboardSectionConfig{Type: "configured", Title: "Configured"},
		deps,
	).(*ConfiguredSection)

	return BuiltSections{
		Configured: configured,
		Widgets:    widgets,
	}
}

func findWorktreeConfig(worktrees []model.WorktreeConfig, repo string) (model.WorktreeConfig, bool) {
	for _, wc := range worktrees {
		if repo != "" && strings.EqualFold(wc.Repo, repo) {
			return wc, true
		}
	}
	return model.WorktreeConfig{}, false
}
