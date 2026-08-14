package dashboard

import (
	"log/slog"

	tea "charm.land/bubbletea/v2"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type Section interface {
	Name() string
	Init() tea.Cmd
	Update(msg tea.Msg) (Section, tea.Cmd)
	// ViewBorderless returns the pane title (drawn on the shared frame border)
	// and the borderless content (already clipped/padded to the given height).
	ViewBorderless(width, height int, focused bool) (title string, content string)
	Chosen() string
	TotalItems() int
	Width() float64
}

// Sorter is implemented by sections whose list can be re-sorted by the `s`
// key. SortLabel returns the current sort mode's name.
type Sorter interface {
	SortLabel() string
}

// Filterer is implemented by sections that support type-to-filter. While
// Filtering() is true the Model routes every key (except ctrl+c) to the pane
// so quit/switch/nav keys don't fire mid-typing.
type Filterer interface {
	Filtering() bool
	FilterQuery() string
}

// Clicker is implemented by list sections that move their selection to a
// specific row in response to a mouse click. row is 0-based and
// view-relative (the line within the pane's content area).
type Clicker interface {
	ClickAt(row int)
}

type SectionDeps struct {
	Tmux      tmux.Tmux
	Lister    lister.Lister
	Git       git.Git
	Connector connector.Connector
	Shell     shell.Shell
	HomeDir   string
}

type SectionFactory func(cfg model.DashboardSectionConfig, deps SectionDeps) Section

type Registry map[string]SectionFactory

// registry maps configurable widget types to their factories. The "sessions"
// type is now implicit (always built) and is therefore not part of the
// addable widget registry.
var registry = Registry{
	// "details": NewDetailsSection,
	"system":  NewSystemSection,
	"ssh":     NewSSHSection,
	"git":     NewGitSection,
	"custom":  NewCustomSection,
	"docker":  NewDockerSection,
	"workmux": NewWorkmuxSection,
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
