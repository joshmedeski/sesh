package dashboard

import (
	"log/slog"

	"github.com/joshmedeski/sesh/v2/dashboard/core"
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

type SectionFactory func(cfg model.DashboardSectionConfig, deps SectionDeps) Section

type Registry map[string]SectionFactory

// registry maps configurable widget types to their factories. The "sessions"
// type is now implicit (always built) and is therefore not part of the
// addable widget registry.
var registry = Registry{
	// "details": sections.NewDetailsSection,
	"system":  sections.NewSystemSection,
	"ssh":     sections.NewSSHSection,
	"git":     sections.NewGitSection,
	"custom":  sections.NewCustomSection,
	"docker":  sections.NewDockerSection,
	"workmux": sections.NewWorkmuxSection,
}

// BuiltSections is the result of BuildSections: the two permanent lists plus
// the optional user-configured widgets.
type BuiltSections struct {
	Sessions   *SessionsSection
	Configured *ConfiguredSection
	Widgets    []Section
}

// BuildSections always builds the Open (sessions) and Configured lists. Config
// `[dashboard.sections]` entries are treated as optional widgets only; the
// "sessions" type and unknown types are logged and skipped. If a legacy
// "sessions" entry exists, its Title is carried over to the implicit sessions
// list; Groups are parsed but no longer applied (grouping was removed).
func BuildSections(cfg model.DashboardConfig, deps SectionDeps) BuiltSections {
	sessionsCfg := model.DashboardSectionConfig{Type: "sessions", Title: "Sessions"}

	var widgets []Section
	for _, sc := range cfg.Sections {
		switch sc.Type {
		case "":
			slog.Warn("unknown dashboard section type")
			continue
		case "sessions":
			slog.Warn("dashboard section type \"sessions\" is now implicit; ignoring entry")
			if sc.Groups != nil {
				slog.Warn("dashboard \"sessions\" groups are no longer applied; sessions render as a flat list")
			}
			if sc.Title != "" {
				sessionsCfg.Title = sc.Title
			}
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

	sessions := NewSessionsSection(sessionsCfg, deps).(*SessionsSection)
	configured := NewConfiguredSection(
		model.DashboardSectionConfig{Type: "configured", Title: "Configured"},
		deps,
	).(*ConfiguredSection)

	return BuiltSections{
		Sessions:   sessions,
		Configured: configured,
		Widgets:    widgets,
	}
}
