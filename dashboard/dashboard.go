package dashboard

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type Model struct {
	config        model.DashboardConfig
	sections      []Section
	focused       int
	width         int
	height        int
	chosen        string
	quit          bool
	totalSessions int
}

func New(config model.DashboardConfig, tmux tmux.Tmux, lister lister.Lister, git git.Git, connector connector.Connector, homeDir string) Model {
	deps := SectionDeps{
		Tmux:      tmux,
		Lister:    lister,
		Git:       git,
		Connector: connector,
		HomeDir:   homeDir,
	}

	// build sections based on config and dependencies
	sections := BuildSections(config, deps)

	return Model{
		config:   config,
		sections: sections,
		focused:  0,
		width:    80,
		height:   24,
	}
}

func (m Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.sections))
	for i, s := range m.sections {
		cmds[i] = s.Init()
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "left", "h":
			m.focused--
			if m.focused < 0 {
				m.focused = len(m.sections) - 1
			}
		case "right", "l":
			m.focused++
			if m.focused >= len(m.sections) {
				m.focused = 0
			}
		case "enter":
			if len(m.sections) > 0 {
				m.sections[m.focused], cmd = m.sections[m.focused].Update(msg)
				if chosen := m.sections[m.focused].Chosen(); chosen != "" {
					m.chosen = chosen
					return m, tea.Quit
				}
			}
		default:
			if len(m.sections) > 0 {
				m.sections[m.focused], cmd = m.sections[m.focused].Update(msg)
			}
		}
		m = m.syncHoveredSession()
		return m, cmd

	default:
		// Non-keypress messages -> all sections
		if len(m.sections) > 0 {
			var cmds []tea.Cmd
			for i := range m.sections {
				var c tea.Cmd
				m.sections[i], c = m.sections[i].Update(msg)
				if c != nil {
					cmds = append(cmds, c)
				}
			}
			cmd = tea.Batch(cmds...)
		}
		m = m.syncHoveredSession()
		return m, cmd
	}

	return m, nil
}

func (m Model) syncHoveredSession() Model {
	var ssIdx, dsIdx int = -1, -1
	for i, s := range m.sections {
		switch s.(type) {
		case *SessionsSection:
			ssIdx = i
		case *DetailsSection:
			dsIdx = i
		}
	}
	if ssIdx < 0 || dsIdx < 0 {
		return m
	}
	ss := m.sections[ssIdx].(*SessionsSection)
	name, path := ss.HoveredSession()
	updated, _ := m.sections[dsIdx].Update(hoveredSessionMsg{Name: name, Path: path})
	m.sections[dsIdx] = updated
	return m
}

func (m Model) View() tea.View {
	if m.quit {
		return tea.NewView("")
	}

	var b strings.Builder

	m.totalSessions = 0
	for _, s := range m.sections {
		m.totalSessions += s.TotalItems()
	}

	if m.width < 20 || m.height < 5 {
		return tea.NewView("Terminal too small for dashboard")
	}

	w := m.width
	title := m.config.Title
	if title == "" {
		title = "SESH COMMAND CENTER"
	}
	header := renderHeader(title, m.totalSessions, w)
	b.WriteString(header)
	b.WriteString("\n\n")

	contentHeight := max(m.height-5, 1)

	// Allocate pixel widths to each section
	n := len(m.sections)
	sepCount := max(n-1, 0)
	availableWidth := m.width - sepCount
	pw := make([]int, n)
	allocated := 0
	for i, s := range m.sections {
		if w := s.Width(); w > 0 {
			pw[i] = int(float64(availableWidth) * w)
			allocated += pw[i]
		}
	}
	remaining := availableWidth - allocated
	flexCount := 0
	for _, p := range pw {
		if p == 0 {
			flexCount++
		}
	}
	if flexCount > 0 {
		each := remaining / flexCount
		for i := range pw {
			if pw[i] == 0 {
				pw[i] = each
				remaining -= each
			}
		}
		for i := n - 1; i >= 0 && remaining > 0; i-- {
			if pw[i] == each {
				pw[i] += remaining
			}
		}
	} else if remaining > 0 {
		pw[n-1] += remaining
	}

	// Render sections side by side
	var views []string
	for i, section := range m.sections {
		view := section.View(pw[i], contentHeight)
		if view != "" {
			views = append(views, view)
		}
	}
	if len(views) > 0 {
		b.WriteString(joinViews(views, pw))
	}

	b.WriteString(renderFooter(w))

	// Pad to full terminal height
	lines := strings.Count(b.String(), "\n")
	for i := lines; i < m.height; i++ {
		b.WriteString("\n")
	}

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}

func (m Model) Chosen() string {
	return m.chosen
}

func (m Model) Quit() bool {
	return m.quit
}

func joinViews(views []string, widths []int) string {
	if len(views) == 0 {
		return ""
	}
	if len(views) == 1 {
		return views[0]
	}
	var allLines [][]string
	maxLines := 0
	for i, v := range views {
		lines := strings.Split(v, "\n")
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		for j, line := range lines {
			if w := lipgloss.Width(line); w < widths[i] {
				lines[j] = line + strings.Repeat(" ", widths[i]-w)
			}
		}
		allLines = append(allLines, lines)
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}
	var b strings.Builder
	for line := 0; line < maxLines; line++ {
		for i, lines := range allLines {
			if line < len(lines) {
				b.WriteString(lines[line])
			} else {
				b.WriteString(strings.Repeat(" ", widths[i]))
			}
			if i < len(views)-1 {
				separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8)).Faint(true)
				b.WriteString(separatorStyle.Render("│"))
			}
		}
		if line < maxLines-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
