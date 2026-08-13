package dashboard

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

type customOutputMsg struct {
	output string
	err    error
}

type CustomSection struct {
	config  model.DashboardSectionConfig
	deps    SectionDeps
	output  string
	loading bool
}

func NewCustomSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &CustomSection{
		config:  cfg,
		deps:    deps,
		loading: true,
	}
}

func (s *CustomSection) Name() string    { return s.config.Title }
func (s *CustomSection) TotalItems() int { return 0 }
func (s *CustomSection) Width() float64  { return s.config.Width }
func (s *CustomSection) Chosen() string  { return "" }

func (s *CustomSection) Init() tea.Cmd {
	return s.fetchOutput
}

func (s *CustomSection) fetchOutput() tea.Msg {
	cmd := s.config.Custom.Command
	if cmd == "" {
		return customOutputMsg{output: "No command configured"}
	}
	out, err := runShellCommand(cmd)
	if err != nil {
		return customOutputMsg{err: err}
	}
	return customOutputMsg{output: string(out)}
}

func (s *CustomSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case customOutputMsg:
		s.loading = false
		if msg.err != nil {
			s.output = "Error: " + msg.err.Error()
		} else {
			s.output = msg.output
		}
	case tea.KeyPressMsg:
		if msg.String() == "r" {
			s.loading = true
			return s, s.fetchOutput
		}
	}
	return s, nil
}

func (s *CustomSection) ViewBorderless(width, height int, focused bool) (string, string) {
	title := s.config.Title
	if title == "" {
		title = "Custom"
	}

	const minWidth = 16
	if width < minWidth {
		return title, "  Custom"
	}

	var b strings.Builder

	if s.loading {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("  Loading..."))
	} else if s.output == "" {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("  No output"))
	} else {
		outputStyle := lipgloss.NewStyle().Width(width - 4)
		b.WriteString(outputStyle.Render(s.output))
	}

	return title, b.String()
}
