package dashboard

import (
	"fmt"
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
	config      model.DashboardSectionConfig
	viewHeight  int
	hoveredName string
	hoveredPath string
	// hoveredUptime tea.Cmd
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
		// s.hoveredUptime = tea.Cmd(func() tea.Msg {
		// 	pipeline := fmt.Sprintf("tmux ls | awk -F: '$1 == \"%s\" {print $2}'", s.hoveredName)
		// 	c := exec.Command("sh", "-c", pipeline)
		// 	out, err := c.Output()
		// 	if err != nil {
		// 		return hoveredSessionMsg{Name: "", Path: ""}
		// 	}
		// 	return hoveredSessionMsg{Name: s.hoveredName, Path: strings.TrimSpace(string(out))}
		// })
	}
	return s, nil
}

func (s *DetailsSection) View(width, height int) string {
	s.viewHeight = height

	// Guard: Minimum layout size checks
	const minWidth = 64
	if width < minWidth {
		msg := fmt.Sprintf("  Enlarge pane to see sessions (need ≥%d cols, have %d)", minWidth, width)
		return lipgloss.NewStyle().Faint(true).Width(width).Height(height).Render(msg)
	}

	// Calculate internal view height
	internalWidth := max(width-2, 1)
	internalHeight := max(height-2, 1)

	// Calculate active available viewing rows
	chrome := 2 // Accounts for title header line space
	available := height - chrome
	if available < 1 {
		available = 5
	}

	if s.hoveredName == "" {
		return NewStyleBorder(internalWidth, internalWidth, internalHeight, internalHeight, 15, false, []int{0, 0, 0, 0}).Render("")
	}

	var b strings.Builder

	// Style Definitions
	sectionStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0})
	b.WriteString(sectionStyle.Render(s.config.Title))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	nameRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Name:"), valueStyle.Render(s.hoveredName))
	pathRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Path:"), valueStyle.Render(s.hoveredPath))

	b.WriteString(nameRow)
	b.WriteString("\n\n")
	b.WriteString(pathRow)
	b.WriteString("\n")

	// sessionNameStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0}, "Name:")
	// sessionStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0}, s.hoveredName)
	// pathNameStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0}, "Path:")
	// pathStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0}, s.hoveredPath)
	//

	// lines := strings.Count(b.String(), "\n")
	// for i := lines; i < internalHeight; i++ {
	// 	b.WriteString("\n")
	// }

	return lipgloss.NewStyle().
		Width(internalWidth).
		Height(internalHeight). // Account for border
		MaxHeight(height).
		Border(lipgloss.RoundedBorder()).
		Render(b.String())
}
