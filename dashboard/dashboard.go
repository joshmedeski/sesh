package dashboard

import (
	"charm.land/bubbles/v2/key"
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
	tooSmall      bool
	chosen        string
	quit          bool
	totalSessions int
	sectionWidths []int
	contentHeight int
}

type keyMap struct {
	Quit       key.Binding
	FocusLeft  key.Binding
	FocusRight key.Binding
	Select     key.Binding
}

var keys = keyMap{
	Quit:       key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "quit")),
	FocusLeft:  key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "focus left")),
	FocusRight: key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "focus right")),
	Select:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
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
		config:        config,
		sections:      sections,
		focused:       0,
		width:         80,
		height:        24,
		contentHeight: 20,
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

		// Check width and height constraints
		// with early return if too small
		if m.width < 20 || m.height < 5 {
			m.tooSmall = true
			return m, tea.Quit
		}
		m.tooSmall = false

		// Set footer height to 4 lines
		// and calculate available content height
		headerFooterHeight := 4
		contentHeight := max(m.height-headerFooterHeight, 1)

		// Allocate pixel widths to each section
		n := len(m.sections)
		sepCount := max(n-1, 0)
		availableWidth := m.width - sepCount // Subtract separators
		pw := make([]int, n)
		allocated := 0

		// Calculate available width for each section
		for i, s := range m.sections {
			if w := s.Width(); w > 0 {
				pw[i] = int(float64(availableWidth) * w)
				allocated += pw[i]
			}
		}

		// Calculate remaining available width
		remaining := availableWidth - allocated
		flexCount := 0
		for _, p := range pw {
			if p == 0 {
				flexCount++
			}
		}

		// Distribute remaining width to flex sections
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
					remaining--
				}
			}
		} else if remaining > 0 {
			pw[n-1] += remaining
		}

		// Update each section's dimensions directly
		m.sectionWidths = pw
		m.contentHeight = contentHeight

		return m, nil

	// Handle keypresses
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quit = true
			return m, tea.Quit
		case key.Matches(msg, keys.FocusLeft):
			m.focused--
			if m.focused < 0 {
				m.focused = len(m.sections) - 1
			}
		case key.Matches(msg, keys.FocusRight):
			m.focused++
			if m.focused >= len(m.sections) {
				m.focused = 0
			}
		case key.Matches(msg, keys.Select):
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

	// Non-keypress messages -> all sections
	default:
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
	name, path, windows := ss.HoveredSession()
	updated, _ := m.sections[dsIdx].Update(hoveredSessionMsg{Name: name, Path: path, Windows: windows})
	m.sections[dsIdx] = updated
	return m
}

func (m Model) View() tea.View {
	if m.quit {
		return tea.NewView("")
	}
	if m.tooSmall {
		return tea.NewView("Terminal too small for dashboard")
	}

	// Render Header & Footer strings
	title := m.config.Title
	if title == "" {
		title = "SESH COMMAND CENTER"
	}
	header := renderHeader(title, m.totalSessions, m.width)
	footer := renderFooter(m.width)

	// Render and join sub-sections horizontally
	var views []string
	for i, section := range m.sections {
		w := m.width
		h := m.height
		if m.sectionWidths != nil && i < len(m.sectionWidths) {
			w = m.sectionWidths[i]
		}
		if m.contentHeight > 0 {
			h = m.contentHeight
		}
		if v := section.View(w, h); v != "" {
			views = append(views, v)
		}
	}

	// Combine them side-by-side using Lipgloss
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, views...)

	// Stack everything vertically
	ui := lipgloss.JoinVertical(lipgloss.Top, header, mainContent, footer)

	// Force the layout to exactly fit the terminal window using Lipgloss padding
	finalString := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		MaxHeight(m.height).
		Render(ui)

	v := tea.NewView(finalString)
	v.AltScreen = true // Full-screen mode
	return v
}

func (m Model) Chosen() string {
	return m.chosen
}

func (m Model) Quit() bool {
	return m.quit
}
