package dashboard

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type Model struct {
	sections      []Section
	focused       int
	width         int
	height        int
	chosen        string
	quit          bool
	totalSessions int
}

func New(
	config model.DashboardConfig,
	tmux tmux.Tmux,
	lister lister.Lister,
	git git.Git,
	connector connector.Connector,
	homeDir string,
) Model {
	deps := SectionDeps{
		Tmux:      tmux,
		Lister:    lister,
		Git:       git,
		Connector: connector,
		HomeDir:   homeDir,
	}

	sections := BuildSections(config, deps)

	return Model{
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "enter":
			if len(m.sections) > 0 {
				chosen := m.sections[m.focused].Chosen()
				if chosen != "" {
					m.chosen = chosen
					return m, tea.Quit
				}
			}
			return m, nil
		}
	}

	// Forward all other messages to the focused section
	if len(m.sections) > 0 {
		var cmd tea.Cmd
		m.sections[m.focused], cmd = m.sections[m.focused].Update(msg)
		return m, cmd
	}

	return m, nil
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

	w := max(m.width, 20)
	sep := strings.Repeat("─", w-2)
	b.WriteString(fmt.Sprintf("┌%s┐\n", sep))
	header := renderHeader("SESH COMMAND CENTER", m.totalSessions, w)
	b.WriteString(header)
	b.WriteString("\n")

	contentHeight := max(m.height-5, 1)

	for i, section := range m.sections {
		if i == m.focused {
			view := section.View(w, contentHeight)
			if view != "" {
				b.WriteString(view)
			}
		}
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
