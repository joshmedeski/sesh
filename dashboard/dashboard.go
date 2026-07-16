package dashboard

import (
	"slices"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type Model struct {
	config             model.DashboardConfig
	sections           []Section
	focused            int
	width              int
	height             int
	tooSmall           bool
	chosen             string
	quit               bool
	totalSessions      int
	sectionWidths      []int
	contentHeight      int
	rows               [][]int // rows[rowIdx] = []section indices
	lastHoveredSession string
}

type keyMap struct {
	Quit   key.Binding
	Next   key.Binding
	Prev   key.Binding
	Select key.Binding
}

var keys = keyMap{
	Quit:   key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "quit")),
	Next:   key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next section")),
	Prev:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev section")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
}

func New(config model.DashboardConfig, tmux tmux.Tmux, lister lister.Lister, git git.Git, connector connector.Connector, sh shell.Shell, homeDir string) Model {
	deps := SectionDeps{
		Tmux:      tmux,
		Lister:    lister,
		Git:       git,
		Connector: connector,
		Shell:     sh,
		HomeDir:   homeDir,
	}

	sections := BuildSections(config, deps)

	// Build row assignments from config
	sectionRows := make([]int, len(sections))
	for i, sc := range config.Sections {
		if i >= len(sections) {
			break
		}
		sectionRows[i] = sc.Row
	}
	// For default sections (no config), all go to row 0
	for i := len(config.Sections); i < len(sections); i++ {
		sectionRows[i] = 0
	}

	// Group sections by row
	rowMap := make(map[int][]int)
	for i, row := range sectionRows {
		rowMap[row] = append(rowMap[row], i)
	}

	// Sort rows by key for stable ordering
	rowKeys := make([]int, 0, len(rowMap))
	for r := range rowMap {
		rowKeys = append(rowKeys, r)
	}
	slices.Sort(rowKeys)

	var rows [][]int
	for _, r := range rowKeys {
		rows = append(rows, rowMap[r])
	}

	return Model{
		config:        config,
		sections:      sections,
		focused:       0,
		width:         80,
		height:        24,
		contentHeight: 20,
		rows:          rows,
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

		if m.width < 20 || m.height < 5 {
			m.tooSmall = true
			return m, tea.Quit
		}
		m.tooSmall = false

		headerFooterHeight := 4
		contentHeight := max(m.height-headerFooterHeight, 1)
		m.contentHeight = contentHeight

		// Calculate heights per row (equal distribution)
		numRows := len(m.rows)
		if numRows == 0 {
			numRows = 1
		}

		// Calculate widths per-row
		m.sectionWidths = make([]int, len(m.sections))
		for _, rowSections := range m.rows {
			n := len(rowSections)
			sepCount := max(n-1, 0)
			availableWidth := m.width - sepCount
			pw := make([]int, n)

			flexCount := 0
			totalFraction := 0.0
			for _, si := range rowSections {
				if w := m.sections[si].Width(); w > 0 {
					totalFraction += w
				} else {
					flexCount++
				}
			}

			scale := 1.0
			if totalFraction > 1.0 {
				scale = 1.0 / totalFraction
			}

			allocated := 0
			for j, si := range rowSections {
				if w := m.sections[si].Width(); w > 0 {
					pw[j] = int(float64(availableWidth) * w * scale)
					allocated += pw[j]
				}
			}

			remaining := availableWidth - allocated
			if flexCount > 0 {
				each := remaining / flexCount
				for j := range pw {
					if pw[j] == 0 {
						pw[j] = each
						remaining -= each
					}
				}
				for j := n - 1; j >= 0 && remaining > 0; j-- {
					if pw[j] == each {
						pw[j]++
						remaining--
					}
				}
			} else if remaining > 0 {
				pw[n-1] += remaining
			}

			for j, si := range rowSections {
				m.sectionWidths[si] = pw[j]
			}
		}
		var cmds []tea.Cmd
		for i := range m.sections {
			var c tea.Cmd
			m.sections[i], c = m.sections[i].Update(msg)
			if c != nil {
				cmds = append(cmds, c)
			}
		}
		return m, tea.Batch(cmds...)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quit = true
			return m, tea.Quit

		case key.Matches(msg, keys.Next):
			if len(m.sections) > 0 {
				m.focused = (m.focused + 1) % len(m.sections)
			}

		case key.Matches(msg, keys.Prev):
			if len(m.sections) > 0 {
				m.focused--
				if m.focused < 0 {
					m.focused = len(m.sections) - 1
				}
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
			for i := range m.sections {
				if ss, ok := m.sections[i].(*SessionsSection); ok {
					m.totalSessions = ss.totalSessions
				}
			}
		}
		m, syncCmd := m.syncHoveredSession()
		return m, tea.Batch(cmd, syncCmd)

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
		m, syncCmd := m.syncHoveredSession()
		return m, tea.Batch(cmd, syncCmd)
	}
}

func (m Model) syncHoveredSession() (Model, tea.Cmd) {
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
		return m, nil
	}
	ss := m.sections[ssIdx].(*SessionsSection)
	name, path, windows := ss.HoveredSession()

	// ONLY update details and execute background commands if the selection changed!
	if name == m.lastHoveredSession {
		return m, nil
	}
	m.lastHoveredSession = name

	updated, cmd := m.sections[dsIdx].Update(hoveredSessionMsg{Name: name, Path: path, Windows: windows})
	m.sections[dsIdx] = updated
	return m, cmd
}

func (m Model) View() tea.View {
	if m.quit {
		return tea.NewView("")
	}
	if m.tooSmall {
		return tea.NewView("Terminal too small for dashboard")
	}

	title := m.config.Title
	if title == "" {
		title = "SESH COMMAND CENTER"
	}
	header := renderHeader(title, m.totalSessions, m.width)
	footer := renderFooter(m.width)

	numRows := len(m.rows)
	if numRows == 0 {
		numRows = 1
	}
	sepLines := max(numRows-1, 0)
	availableForRows := max(m.contentHeight-sepLines, numRows)
	baseHeight := availableForRows / numRows
	extra := availableForRows - baseHeight*numRows

	// Render each row horizontally, then stack rows vertically
	var rowViews []string
	for ri, rowSections := range m.rows {
		h := baseHeight
		if ri < extra {
			h++
		}
		var views []string
		for _, si := range rowSections {
			w := m.width
			if m.sectionWidths != nil && si < len(m.sectionWidths) && m.sectionWidths[si] > 0 {
				w = m.sectionWidths[si]
			}
			if v := m.sections[si].View(w, h, si == m.focused); v != "" {
				views = append(views, v)
			}
		}
		if len(views) > 0 {
			rowViews = append(rowViews, lipgloss.JoinHorizontal(lipgloss.Top, views...))
		}
	}

	mainContent := lipgloss.JoinVertical(lipgloss.Top, rowViews...)

	ui := lipgloss.JoinVertical(lipgloss.Top, header, mainContent, footer)

	finalString := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(ui)

	v := tea.NewView(finalString)
	v.AltScreen = true
	return v
}

func (m Model) Chosen() string {
	return m.chosen
}

func (m Model) Quit() bool {
	return m.quit
}
