package dashboard

import (
	"fmt"
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
		out, err := runCommand("tmux", "list-windows", "-t", name, "-F", format)
		if err != nil {
			return windowNamesLoadedMsg{}
		}

		var names []string
		var active string
		var idx []string
		for line := range strings.SplitSeq(strings.TrimSpace(out), "\n") {
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 3 {
				continue
			}
			if parts[1] == "1" {
				active = parts[2]
			}
			names = append(names, parts[2])
			idx = append(idx, parts[0])
		}

		return windowNamesLoadedMsg{WindowIdx: idx, WindowNames: names, ActiveWindow: active}
	}
}

func (s *DetailsSection) VenvCheck(path string) tea.Cmd {
	return func() tea.Msg {
		if path == "" {
			return venvLoadedMsg{active: "no", name: "none"}
		}

		targetPath := path
		if strings.HasPrefix(targetPath, "~") {
			targetPath = filepath.Join(s.deps.HomeDir, targetPath[1:])
		}

		venvDirs := []string{".venv", "venv", "env"}

		for _, dir := range venvDirs {
			fullPath := filepath.Join(targetPath, dir)
			info, err := filepath.Glob(fullPath)
			if err == nil && len(info) > 0 {
				return venvLoadedMsg{
					active: "yes",
					name:   dir,
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

func (s *DetailsSection) View(width, height int, focused bool) string {
	s.viewHeight = height

	// Guard: Minimum layout size checks
	const minWidth = 30
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
		return NewStyleBorder(internalWidth, internalWidth, internalHeight+2, internalHeight+2, 15, false, []int{0, 0, 0, 1}, focused).Render(s.config.Title)
	}

	var lines []string

	lines = append(lines, s.config.Title)
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  Name: %s", s.hoveredName))
	lines = append(lines, fmt.Sprintf("  Path: %s", s.hoveredPath))

	winLabel := "  Windows:"
	if len(s.hoveredWindowNames) == 0 {
		lines = append(lines, winLabel)
	} else if len(s.hoveredWindowNames) == 1 {
		idx := ""
		if len(s.hoveredWindowIdx) > 0 {
			idx = s.hoveredWindowIdx[0]
		}
		lines = append(lines, fmt.Sprintf("%s %s.%s", winLabel, idx, s.hoveredWindowNames[0]))
	} else {
		lines = append(lines, winLabel)
		for i, w := range s.hoveredWindowNames {
			idx := ""
			if i < len(s.hoveredWindowIdx) {
				idx = s.hoveredWindowIdx[i]
			}
			lines = append(lines, fmt.Sprintf("    %s.%s", idx, w))
		}
	}

	lines = append(lines, fmt.Sprintf("  Venv: %s", s.hoveredVenvName))

	maxContentLines := internalHeight
	if len(lines) > maxContentLines {
		lines = lines[:maxContentLines]
	}

	for len(lines) < maxContentLines {
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")
	return NewStyleBorder(internalWidth, internalWidth, internalHeight+2, internalHeight+2, 15, false, []int{0, 0, 0, 1}, focused).Render(content)
}
