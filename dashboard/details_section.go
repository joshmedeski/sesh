package dashboard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

type venvLoadedMsg struct {
	active string
	name   string
}

type DetailsSection struct {
	config              model.DashboardSectionConfig
	deps                SectionDeps
	viewHeight          int
	hoveredName         string
	hoveredPath         string
	hoveredWindows      int
	hoveredWindowNames  []string
	hoveredActiveWindow string
	hoveredWindowIdx    []string
	hoveredVenvName     string
	hoveredVenvActive   string
	// hoveredUptime tea.Cmd
}

func NewDetailsSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &DetailsSection{
		config: cfg,
		deps:   deps,
	}
}

func (s *DetailsSection) Name() string    { return s.config.Title }
func (s *DetailsSection) TotalItems() int { return 0 }
func (s *DetailsSection) Width() float64  { return s.config.Width }
func (s *DetailsSection) Chosen() string  { return "" }

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

func (s *DetailsSection) VenvCheck(path string) tea.Cmd {
	return func() tea.Msg {
		if path == "" {
			return venvLoadedMsg{active: "no", name: "none"}
		}

		// Expand paths like "~" if necessary, otherwise use as-is
		targetPath := path
		if strings.HasPrefix(targetPath, "~") {
			home, err := os.UserHomeDir()
			if err == nil {
				targetPath = filepath.Join(home, targetPath[1:])
			}
		}

		// Common virtual environment folder names
		venvDirs := []string{".venv", "venv", "env"}

		for _, dir := range venvDirs {
			fullPath := filepath.Join(targetPath, dir)
			info, err := os.Stat(fullPath)

			// If the directory exists, we found an inactive/available venv for this path
			if err == nil && info.IsDir() {
				return venvLoadedMsg{
					active: "yes",
					name:   dir, // returns ".venv", "venv", etc.
				}
			}
		}

		return venvLoadedMsg{active: "no", name: "none"}
	}
}

func (s *DetailsSection) Init() tea.Cmd { return nil }

func (s *DetailsSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case hoveredSessionMsg:
		// If the hovered session is the same as the current one, don't update
		if s.hoveredName == msg.Name && s.hoveredPath == msg.Path && s.hoveredWindows == msg.Windows {
			return s, nil
		}

		s.hoveredName = msg.Name
		s.hoveredPath = msg.Path
		s.hoveredWindows = msg.Windows

		return s, tea.Batch(s.WindowNames(msg.Name), s.VenvCheck(msg.Path))
	case windowNamesLoadedMsg:
		s.hoveredWindowNames = msg.WindowNames
		s.hoveredActiveWindow = msg.ActiveWindow
		s.hoveredWindowIdx = msg.WindowIdx
	case venvLoadedMsg:
		s.hoveredVenvName = msg.name
		s.hoveredVenvActive = msg.active
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
		return NewStyleBorder(internalWidth, internalWidth, internalHeight+2, internalHeight+2, 15, false, []int{0, 0, 0, 1}).Render(s.config.Title)
	}

	var b strings.Builder

	// Style Definitions
	sectionStyle := NewStyle(internalWidth, internalWidth, 1, 1, 15, false, []int{0, 0, 0, 0})
	b.WriteString(sectionStyle.Render(s.config.Title))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true).Padding(0, 0, 0, 1)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	nameRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Name: "), valueStyle.Render(s.hoveredName))
	pathRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Path: "), valueStyle.Render(s.hoveredPath))
	windowsRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Windows: "))
	venvCheck := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Venv: "), valueStyle.Render(s.hoveredVenvName))

	b.WriteString(nameRow)
	b.WriteString("\n")
	b.WriteString(pathRow)
	b.WriteString("\n")
	b.WriteString(windowsRow)
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
	b.WriteString("\n")
	b.WriteString(venvCheck)

	lines := strings.Count(b.String(), "\n")
	for i := lines + 1; i < internalHeight; i++ {
		b.WriteString("\n")
	}

	details := NewStyleBorder(internalWidth, internalWidth, internalHeight+2, internalHeight+2, 15, false, []int{0, 0, 0, 1}).Render(b.String())
	return details
}
