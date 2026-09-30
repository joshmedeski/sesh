package dashboard

import (
	"fmt"
	"log/slog"
	"slices"
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

// BuiltPage is one dashboard page: its tab title and rows of panes.
type BuiltPage struct {
	Title string
	Rows  [][]Section
}

// Built is the result of BuildPages: the Configured tab's list plus the
// dashboard pages.
type Built struct {
	Configured *ConfiguredSection
	Pages      []BuiltPage
}

func (b Built) pages() []dashPage {
	pages := make([]dashPage, len(b.Pages))
	for i, p := range b.Pages {
		pages[i] = dashPage{title: p.Title, rows: p.Rows}
	}
	return pages
}

// BuildPages builds the Configured list and one dashboard page per
// `[[dashboard.page]]`, each from its rows of section tables. Unknown section
// types and unmatched worktree repos are logged and skipped, and so are rows
// and pages left empty. With no usable page there is a single "Dashboard"
// page holding a tmux sessions list.
func BuildPages(cfg model.DashboardConfig, worktrees []model.WorktreeConfig, deps SectionDeps) Built {
	var pages []BuiltPage
	for _, pc := range cfg.Pages {
		var rows [][]Section
		for _, rowCfg := range pc.Sections {
			var row []Section
			for _, sc := range rowCfg {
				if sec, ok := buildSection(sc, worktrees, deps); ok {
					row = append(row, sec)
				}
			}
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}
		if len(rows) == 0 {
			slog.Warn("dashboard page has no usable sections", "title", pc.Title)
			continue
		}
		pages = append(pages, BuiltPage{Title: pc.Title, Rows: rows})
	}
	if len(pages) == 0 {
		sessions, _ := buildSection(model.DashboardSectionConfig{Type: "sources", Title: "Sessions"}, worktrees, deps)
		pages = []BuiltPage{{Rows: [][]Section{{sessions}}}}
	}
	for i := range pages {
		if pages[i].Title != "" {
			continue
		}
		pages[i].Title = "Dashboard"
		if i > 0 {
			pages[i].Title = fmt.Sprintf("Dashboard %d", i+1)
		}
	}

	configured := NewConfiguredSection(
		model.DashboardSectionConfig{Type: "configured", Title: "Configured"},
		deps,
	).(*ConfiguredSection)

	return Built{Configured: configured, Pages: pages}
}

func buildSection(sc model.DashboardSectionConfig, worktrees []model.WorktreeConfig, deps SectionDeps) (Section, bool) {
	switch sc.Type {
	case "":
		slog.Warn("unknown dashboard section type")
		return nil, false
	case "worktree":
		wc, ok := findWorktreeConfig(worktrees, sc.Repo)
		if !ok {
			slog.Warn("dashboard worktree section has no matching [[worktree]] block", "repo", sc.Repo)
			return nil, false
		}
		ws := NewWorktreeSection(wc, deps)
		ws.title = sc.Title
		if len(sc.Columns) > 0 {
			ws.columns = resolveColumns(sc.Columns, render.WorktreeColumns, "worktree")
		}
		return ws, true
	}
	factory, ok := registry[sc.Type]
	if !ok {
		if sc.Type == "aiagent" {
			slog.Warn("unknown dashboard section type", "type", sc.Type, "hint", "aiagent is deprecated; use workmux")
		} else {
			slog.Warn("unknown dashboard section type", "type", sc.Type)
		}
		return nil, false
	}
	if sc.Groups != nil {
		slog.Warn("dashboard section groups are no longer applied", "type", sc.Type)
	}
	sec := factory(sc, deps)
	if s, ok := sec.(*SessionsSection); ok {
		s.worktrees = worktrees
	}
	return sec, true
}

func resolveColumns(requested, supported []string, sectionType string) []string {
	if len(requested) == 0 {
		return nil
	}
	columns := make([]string, 0, len(requested))
	for _, id := range requested {
		if !slices.Contains(supported, id) {
			slog.Warn("unknown dashboard column", "section", sectionType, "column", id, "supported", supported)
			continue
		}
		columns = append(columns, id)
	}
	return columns
}

func findWorktreeConfig(worktrees []model.WorktreeConfig, repo string) (model.WorktreeConfig, bool) {
	for _, wc := range worktrees {
		if repo != "" && strings.EqualFold(wc.Repo, repo) {
			return wc, true
		}
	}
	return model.WorktreeConfig{}, false
}
