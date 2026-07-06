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
	Name    string
	Path    string
	Windows int
}

type windowNamesLoadedMsg struct {
	WindowIdx    []string
	WindowNames  []string
	ActiveWindow string
}

type DetailsSection struct {
	config              model.DashboardSectionConfig
	groups              []*group
	viewHeight          int
	hoveredName         string
	hoveredPath         string
	hoveredWindows      int
	hoveredWindowNames  []string
	hoveredActiveWindow string
	hoveredWindowIdx    []string
	// hoveredUptime tea.Cmd
}

func NewDetailsSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &DetailsSection{
		config: cfg,
	}
}

func (s *DetailsSection) Name() string     { return s.config.Title }
func (s *DetailsSection) TotalItems() int  { return 0 }
func (s *DetailsSection) Width() float64   { return s.config.Width }
func (s *DetailsSection) Chosen() string   { return "" }
func (s *DetailsSection) WindowCount() int { return s.hoveredWindows }

func (s *DetailsSection) WindowNames(name string) tea.Cmd {
	return func() tea.Msg {
		format := "#{window_index}|#{window_active}|#{pane_current_command}"
		out, err := exec.Command("tmux", "list-windows", "-t", name, "-F", format).Output()
		if err != nil {
			return windowNamesLoadedMsg{}
		}

		var names []string
		var active string
		var idx []string
		for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
			parts := strings.Split(line, "|")
			if len(parts) >= 3 {
				if parts[1] == "1" {
					active = parts[1]
				}
				names = append(names, parts[2])
				idx = append(idx, parts[0])
			}
		}
		return windowNamesLoadedMsg{WindowIdx: idx, WindowNames: names, ActiveWindow: active}
	}
}

func (s *DetailsSection) Init() tea.Cmd { return nil }

func (s *DetailsSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case hoveredSessionMsg:
		s.hoveredName = msg.Name
		s.hoveredPath = msg.Path
		s.hoveredWindows = msg.Windows
		return s, s.WindowNames(msg.Name)
	case windowNamesLoadedMsg:
		s.hoveredWindowNames = msg.WindowNames
		s.hoveredActiveWindow = msg.ActiveWindow
		s.hoveredWindowIdx = msg.WindowIdx
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
		return NewStyleBorder(internalWidth, internalWidth, internalHeight+2, internalHeight+2, 15, false, []int{0, 0, 0, 0}).Render(s.config.Title)
	}

	var b strings.Builder

	// Style Definitions
	sectionStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0})
	b.WriteString(sectionStyle.Render(s.config.Title))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	nameRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Name: "), valueStyle.Render(s.hoveredName))
	pathRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Path: "), valueStyle.Render(s.hoveredPath))
	windowsRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Windows: "))

	b.WriteString(nameRow)
	b.WriteString("\n")
	b.WriteString(pathRow)
	b.WriteString("\n")
	b.WriteString(windowsRow)
	// b.WriteString("\n")
	for i, w := range s.hoveredWindowNames {
		idx := ""
		row := ""
		if i < len(s.hoveredWindowIdx) {
			idx = s.hoveredWindowIdx[i]
		}
		if len(s.hoveredWindowNames) == 1 {
			row = lipgloss.JoinHorizontal(lipgloss.Left,
				labelStyle.Render(fmt.Sprintf("%s.", idx)),
				valueStyle.Render(fmt.Sprintf("%s", w)),
			)
		} else {
			row = lipgloss.JoinHorizontal(lipgloss.Left,
				labelStyle.Render(fmt.Sprintf("%s.", idx)),
				valueStyle.Render(fmt.Sprintf("%s | ", w)),
			)
		}
		b.WriteString(row)
	}

	// sessionNameStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0}, "Name:")
	// sessionStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0}, s.hoveredName)
	// pathNameStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0}, "Path:")
	// pathStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0}, s.hoveredPath)
	//

	// lines := strings.Count(b.String(), "\n")
	// for i := lines; i < internalHeight; i++ {
	// 	b.WriteString("\n")
	// }

	details := NewStyleBorder(internalWidth, internalWidth, internalHeight+2, internalHeight+2, 15, false, []int{0, 0, 0, 0}).Render(b.String())
	return details

	// return lipgloss.NewStyle().
	// 	Width(internalWidth).
	// 	Height(internalHeight). // Account for border
	// 	MaxHeight(height).
	// 	Border(lipgloss.RoundedBorder()).
	// 	Render(b.String())
}
