package dashboard

import (
	"fmt"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

type hoveredSessionMsg struct {
	Name string
	Path string
}

type DetailsSection struct {
	config        model.DashboardSectionConfig
	hoveredName   string
	hoveredPath   string
	hoveredUptime tea.Cmd
}

func NewDetailsSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &DetailsSection{
		config: cfg,
	}
}

func (s *DetailsSection) Name() string    { return s.config.Title }
func (s *DetailsSection) TotalItems() int { return 0 }
func (s *DetailsSection) Width() float64  { return s.config.Width }
func (s *DetailsSection) Chosen() string  { return "" }

func (s *DetailsSection) Init() tea.Cmd { return nil }

func (s *DetailsSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case hoveredSessionMsg:
		s.hoveredName = msg.Name
		s.hoveredPath = msg.Path
		// function to get uptime of hovered tmux session
		s.hoveredUptime = tea.Cmd(func() tea.Msg {
			pipeline := fmt.Sprintf("tmux ls | awk -F: '$1 == \"%s\" {print $2}'", s.hoveredName)
			c := exec.Command("sh", "-c", pipeline)
			out, err := c.Output()
			if err != nil {
				return hoveredSessionMsg{Name: "", Path: ""}
			}
			return hoveredSessionMsg{Name: s.hoveredName, Path: strings.TrimSpace(string(out))}
		})
	}
	return s, nil
}

func (s *DetailsSection) View(width, height int) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().Bold(true).Padding(0, 2)
	b.WriteString(sectionStyle.Render("" + s.config.Title))
	b.WriteString("\n\n")

	if s.hoveredName != "" {
		sessionNameStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Padding(0, 1)
		sessionStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15))
		pathNameStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Padding(0, 1)
		pathStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15))
		fmt.Fprintf(&b, "%s%s\n", sessionNameStyle.Render("Name:"), sessionStyle.Render(s.hoveredName))
		if s.hoveredPath != "" {
			fmt.Fprintf(&b, "%s%s\n", pathNameStyle.Render("Path:"), pathStyle.Render(s.hoveredPath))
		}
		// uptimeNameStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Padding(0, 1)
		// uptimeStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Width(55)
		// fmt.Fprintf(&b, "%s%s\n", uptimeNameStyle.Render("Uptime:"), uptimeStyle.Render(s.hoveredUptime))
	} else {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render(""))
	}

	lines := strings.Count(b.String(), "\n")
	for i := lines; i < height; i++ {
		b.WriteString("\n")
	}

	return b.String()
}
