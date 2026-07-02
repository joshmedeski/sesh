package dashboard

import (
	tea "charm.land/bubbletea/v2"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type Section interface {
	Name() string
	Init() tea.Cmd
	Update(msg tea.Msg) (Section, tea.Cmd)
	View(width, height int) string
	Chosen() string
	TotalItems() int
	Width() float64
}

type SectionDeps struct {
	Tmux      tmux.Tmux
	Lister    lister.Lister
	Git       git.Git
	Connector connector.Connector
	HomeDir   string
}

type SectionFactory func(cfg model.DashboardSectionConfig, deps SectionDeps) Section

type Registry map[string]SectionFactory

// TODO: add more sections eg. system, ssh, etc.
var registry = Registry{
	"sessions": NewSessionsSection,
	"details":  NewDetailsSection,
	// "system":   NewSystemSection,
	// "ssh":      NewSshSection,
	// "git":      NewGitSection,
	// "aiagent":  NewAiAgentSection,
	// "custom":   NewCustomSection,
	// "docker":   NewDockerSection,
}

func BuildSections(cfg model.DashboardConfig, deps SectionDeps) []Section {
	var sections []Section
	for _, sc := range cfg.Sections {
		if factory, ok := registry[sc.Type]; ok {
			sections = append(sections, factory(sc, deps))
		}
	}
	if len(sections) == 0 {
		sections = append(sections, NewSessionsSection(
			model.DashboardSectionConfig{Type: "sessions", Title: "Sessions"},
			deps,
		))
	}
	return sections
}
