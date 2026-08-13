package dashboard

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

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

// previewLoadedMsg carries the result of a tmux capture-pane for a hovered
// session. The name guards against stale results after the hover changes.
type previewLoadedMsg struct {
	name   string
	output string
}

// previewTickMsg is the periodic re-capture trigger.
type previewTickMsg struct {
	name string
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
	previewOutput       string
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

// capturePreview runs tmux capture-pane for a hovered session. Errors yield an
// empty preview (never a crash).
func (s *DetailsSection) capturePreview(name string) tea.Cmd {
	return func() tea.Msg {
		out, err := runCommand("tmux", "capture-pane", "-t", name, "-p", "-e")
		if err != nil {
			out = ""
		}
		return previewLoadedMsg{name: name, output: out}
	}
}

// previewTick schedules the next periodic re-capture.
func previewTick(name string) tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return previewTickMsg{name: name}
	})
}

func (s *DetailsSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case hoveredSessionMsg:
		// Hover cleared: stop refreshing and clear the preview.
		if msg.Name == "" {
			s.hoveredName = ""
			s.hoveredPath = ""
			s.hoveredWindows = 0
			s.previewOutput = ""
			return s, nil
		}

		// If the hovered session is the same as the current one, don't update
		if s.hoveredName == msg.Name && s.hoveredPath == msg.Path && s.hoveredWindows == msg.Windows {
			return s, nil
		}

		s.hoveredName = msg.Name
		s.hoveredPath = msg.Path
		s.hoveredWindows = msg.Windows
		s.previewOutput = ""

		return s, tea.Batch(s.WindowNames(msg.Name), s.VenvCheck(msg.Path), s.capturePreview(msg.Name), previewTick(msg.Name))
	case windowNamesLoadedMsg:
		s.hoveredWindowNames = msg.WindowNames
		s.hoveredActiveWindow = msg.ActiveWindow
		s.hoveredWindowIdx = msg.WindowIdx
	case venvLoadedMsg:
		s.hoveredVenvName = msg.name
		s.hoveredVenvActive = msg.active
	case previewLoadedMsg:
		// Stale capture for a previous hover: ignore.
		if msg.name != s.hoveredName {
			return s, nil
		}
		s.previewOutput = msg.output
	case previewTickMsg:
		// Hover moved/cleared: stop the ticker.
		if msg.name != s.hoveredName || s.hoveredName == "" {
			return s, nil
		}
		return s, tea.Batch(s.capturePreview(s.hoveredName), previewTick(s.hoveredName))
	}
	return s, nil
}

func (s *DetailsSection) ViewBorderless(width, height int, focused bool) (string, string) {
	s.viewHeight = height

	title := s.config.Title
	if title == "" {
		title = "Details"
	}

	// Guard: Minimum layout size checks
	const minWidth = 30
	if width < minWidth {
		msg := fmt.Sprintf("  Enlarge pane to see sessions (need ≥%d cols, have %d)", minWidth, width)
		return title, msg
	}

	if s.hoveredName == "" {
		return title, ""
	}

	var lines []string

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

	maxContentLines := max(height, 1)

	// Preview region fills the remaining pane height with the last lines of
	// the captured output, each ANSI-truncated to the pane width.
	remaining := max(maxContentLines-len(lines), 0)
	if remaining > 0 {
		lines = append(lines, previewLines(s.previewOutput, width, remaining)...)
	}

	if len(lines) > maxContentLines {
		lines = lines[:maxContentLines]
	}

	for len(lines) < maxContentLines {
		lines = append(lines, "")
	}

	return title, strings.Join(lines, "\n")
}

// previewLines returns n lines (bottom-aligned) of the captured output, each
// ANSI-truncated to width. Blank lines are emitted for empty output or when
// there are fewer captured lines than n.
func previewLines(output string, width, n int) []string {
	out := make([]string, n)
	if output == "" || n <= 0 {
		return out
	}
	raw := strings.Split(output, "\n")
	start := max(len(raw)-n, 0)
	taken := raw[start:]
	for i, line := range taken {
		out[i] = truncateRightANSI(line, width)
	}
	return out
}
