package picker

import (
	"github.com/joshmedeski/sesh/v2/model"
)

type sessionItem struct {
	session    model.SeshSession
	name       string // raw session name (no icons/ANSI)
	searchName string // normalized name used for fuzzy matching
	src        string // source type (tmux, config, zoxide, tmuxinator)
}

// sessionItems implements fuzzy.Source for fuzzy matching.
type sessionItems []sessionItem

func (s sessionItems) String(i int) string { return s[i].searchName }
func (s sessionItems) Len() int            { return len(s) }

type filteredItem struct {
	item           sessionItem
	matchedIndexes []int
}

// FetchFunc loads sessions asynchronously. It is called in a goroutine by Init().
type FetchFunc func() (model.SeshSessions, error)

// sessionsLoadedMsg carries the result of the async fetch back to Update().
type sessionsLoadedMsg struct {
	sessions model.SeshSessions
	err      error
}

type Model struct {
	allItems       sessionItems
	filtered       []filteredItem
	filterInput    textinput.Model
	cursor         int
	offset         int
	width          int
	height         int
	chosen         string
	quit           bool
	showIcons      bool
	separatorAware bool
	focusCmd       tea.Cmd
	loading        bool
	fetchFunc      FetchFunc
	loadErr        error
}

// srcIcon returns the nerd font icon and color for a session source.
func srcIcon(src string) (string, color.Color) {
	if g, ok := icon.Glyphs[src]; ok {
		var ansi int
		switch {
		case g.ColorCode >= 90 && g.ColorCode <= 97:
			ansi = g.ColorCode - 82
		case g.ColorCode >= 30 && g.ColorCode <= 37:
			ansi = g.ColorCode - 30
		default:
			ansi = g.ColorCode
		}
		return g.Icon + " ", lipgloss.ANSIColor(ansi)
	}
	return "? ", lipgloss.ANSIColor(8)
}

var separatorReplacer = strings.NewReplacer("-", " ", "_", " ", "/", " ", "\\", " ")

func normalizeSeparators(s string) string {
	return separatorReplacer.Replace(s)
}

func buildItems(sessions model.SeshSessions, separatorAware bool) sessionItems {
	items := make(sessionItems, 0, len(sessions.OrderedIndex))
	for _, key := range sessions.OrderedIndex {
		s := sessions.Directory[key]
		searchName := s.Name
		if separatorAware {
			searchName = normalizeSeparators(s.Name)
		}
		items = append(items, sessionItem{
			session:    s,
			name:       s.Name,
			searchName: searchName,
			src:        s.Src,
		})
	}
	return items
}

func New(fetchFunc FetchFunc, showIcons bool, separatorAware bool) Model {
	ti := textinput.New()
	ti.Placeholder = "Filter sessions..."
	ti.Prompt = "> "

	m := Model{
		filterInput:    ti,
		showIcons:      showIcons,
		separatorAware: separatorAware,
		loading:        true,
		fetchFunc:      fetchFunc,
	}
	m.focusCmd = m.filterInput.Focus()
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.focusCmd, m.fetchSessions())
}

func (m Model) fetchSessions() tea.Cmd {
	return func() tea.Msg {
		sessions, err := m.fetchFunc()
		return sessionsLoadedMsg{sessions: sessions, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sessionsLoadedMsg:
		if msg.err != nil {
			m.loadErr = msg.err
			return m, tea.Quit
		}
		m.loading = false
		m.allItems = buildItems(msg.sessions, m.separatorAware)
		m.applyFilter()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.filterInput.SetWidth(m.contentWidth() - 4)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if m.loading {
				return m, nil
			}
			if len(m.filtered) > 0 {
				selected := m.filtered[m.cursor]
				m.chosen = selected.item.name
			}
			return m, tea.Quit

		case "esc", "ctrl+c":
			m.quit = true
			return m, tea.Quit

		case "up", "ctrl+k":
			m.cursorUp(1)
			return m, nil

		case "down", "ctrl+j":
			m.cursorDown(1)
			return m, nil

		case "ctrl+u":
			m.cursorUp(m.visibleCount() / 2)
			return m, nil

		case "ctrl+d":
			m.cursorDown(m.visibleCount() / 2)
			return m, nil
		}
	}

	// Forward to text input
	prevValue := m.filterInput.Value()
	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)

	if m.filterInput.Value() != prevValue {
		if !m.loading {
			m.applyFilter()
		}
		m.cursor = 0
		m.offset = 0
	}

	return m, cmd
}

func (m *Model) applyFilter() {
	pattern := m.filterInput.Value()
	if pattern == "" {
		m.filtered = make([]filteredItem, len(m.allItems))
		for i, item := range m.allItems {
			m.filtered[i] = filteredItem{item: item}
		}
		return
	}

	if m.separatorAware {
		pattern = normalizeSeparators(pattern)
	}

	matches := fuzzy.FindFrom(pattern, m.allItems)
	m.filtered = make([]filteredItem, len(matches))
	for i, match := range matches {
		m.filtered[i] = filteredItem{
			item:           m.allItems[match.Index],
			matchedIndexes: match.MatchedIndexes,
		}
	}
}

func (m *Model) cursorUp(n int) {
	m.cursor -= n
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
}

func (m *Model) cursorDown(n int) {
	m.cursor += n
	max := len(m.filtered) - 1
	if max < 0 {
		max = 0
	}
	if m.cursor > max {
		m.cursor = max
	}
	visible := m.visibleCount()
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
}

func (m Model) visibleCount() int {
	// border(2) + title(1) + blank(1) + filter(1) + blank(1) + counter(1) + help(1) + blank before counter(1)
	chrome := 9
	available := m.height - chrome
	if available < 1 {
		available = 5
	}
	if available > 15 {
		available = 15
	}
	return available
}

func (m Model) contentWidth() int {
	w := m.width
	if w < 30 {
		w = 40
	}
	if w > 60 {
		w = 60
	}
	return w
}

func (m Model) View() tea.View {
	var b strings.Builder

	// Filter input
	b.WriteString("  " + m.filterInput.View())
	b.WriteString("\n\n")

	visible := m.visibleCount()

	if m.loading {
		loadingStyle := lipgloss.NewStyle().Faint(true)
		b.WriteString(loadingStyle.Render("  Loading sessions..."))
		b.WriteString("\n")
		// Pad remaining visible lines to prevent layout jump
		for i := 1; i < visible; i++ {
			b.WriteString("\n")
		}
	} else {
		// Session list
		end := m.offset + visible
		if end > len(m.filtered) {
			end = len(m.filtered)
		}

		cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)
		matchStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1)).Bold(true)
		normalStyle := lipgloss.NewStyle()

		for i := m.offset; i < end; i++ {
			item := m.filtered[i]
			prefix := "  "
			if i == m.cursor {
				prefix = cursorStyle.Render("> ")
			}

			var tag string
			if m.showIcons {
				icn, clr := srcIcon(item.item.src)
				iconStyle := lipgloss.NewStyle().Foreground(clr)
				tag = iconStyle.Render(icn)
			}
			name := highlightMatches(item.item.name, item.matchedIndexes, matchStyle, normalStyle)

			b.WriteString(fmt.Sprintf("%s%s%s\n", prefix, tag, name))
		}

		// Pad remaining visible lines
		for i := end - m.offset; i < visible; i++ {
			b.WriteString("\n")
		}
	}

	content := b.String()

	return tea.NewView(content)
}

func highlightMatches(s string, indexes []int, matchStyle, normalStyle lipgloss.Style) string {
	if len(indexes) == 0 {
		return normalStyle.Render(s)
	}

	matchSet := make(map[int]bool, len(indexes))
	for _, idx := range indexes {
		matchSet[idx] = true
	}

	var result strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		ch := string(r)
		if matchSet[i] {
			result.WriteString(matchStyle.Render(ch))
		} else {
			result.WriteString(normalStyle.Render(ch))
		}
	}
	return result.String()
}

func (m Model) Chosen() string { return m.chosen }
func (m Model) Quit() bool     { return m.quit }
func (m Model) LoadErr() error { return m.loadErr }
func (m Model) Loading() bool  { return m.loading }
