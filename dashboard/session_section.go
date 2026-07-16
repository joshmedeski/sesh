package dashboard

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

type sessionsLoadedMsg struct {
	sessions model.SeshSessions
	err      error
}

type branchLoadedMsg struct {
	path   string
	branch string
}

type statusLoadedMsg struct {
	path   string
	status string
}

type group struct {
	name      string
	patterns  []string
	sessions  []model.SeshSession
	collapsed bool
}

type flatItem struct {
	isGroup  bool
	groupIdx int
	sessIdx  int
}

type SessionsSection struct {
	config        model.DashboardSectionConfig
	deps          SectionDeps
	groups        []*group
	items         []flatItem
	cursor        int
	offset        int
	loading       bool
	chosen        string
	totalSessions int
	viewHeight    int
}

func NewSessionsSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &SessionsSection{
		config:  cfg,
		deps:    deps,
		loading: true,
	}
}

func (s *SessionsSection) Width() float64 {
	return s.config.Width
}

// name of the section
func (s *SessionsSection) Name() string {
	return s.config.Title
}

// number of items in the section
func (s *SessionsSection) TotalItems() int {
	return s.totalSessions
}

// fetch tmux sessions
func (s *SessionsSection) Init() tea.Cmd {
	return func() tea.Msg {
		sessions, err := s.deps.Lister.List(lister.ListOptions{Tmux: true})
		return sessionsLoadedMsg{sessions: sessions, err: err}
	}
}

func (s *SessionsSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case sessionsLoadedMsg:
		if msg.err != nil {
			return s, nil
		}
		s.loading = false
		s.groupSessions(msg.sessions)
		s.totalSessions = len(msg.sessions.OrderedIndex)
		s.rebuildItems()
		return s, tea.Batch(s.fetchBranches(msg.sessions), s.fetchStatuses(msg.sessions))

	case branchLoadedMsg:
		s.applyBranch(msg.path, msg.branch)
		return s, nil

	case statusLoadedMsg:
		s.applyStatus(msg.path, msg.status)
		return s, nil

	case tea.KeyPressMsg:
		s, cmd := s.handleKey(msg)
		return s, cmd
	}
	return s, nil
}

func (s *SessionsSection) Chosen() string {
	return s.chosen
}

func (s *SessionsSection) handleKey(msg tea.KeyPressMsg) (*SessionsSection, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		s.cursorDown(1)
	case "k", "up":
		s.cursorUp(1)
	case "t":
		s.toggleGroup()
	case "enter":
		s.selectItem()
	case "ctrl+d":
		return s, s.killSession()
	}
	return s, nil
}

func (s *SessionsSection) groupSessions(sessions model.SeshSessions) {
	expanded := make([]*group, 0, len(s.config.Groups))
	for i := range s.config.Groups {
		g := &group{
			name:      s.config.Groups[i].Name,
			patterns:  s.config.Groups[i].Patterns,
			collapsed: false,
		}
		expanded = append(expanded, g)
	}

	other := &group{name: "Other", collapsed: true}

	for _, key := range sessions.OrderedIndex {
		sess := sessions.Directory[key]
		matched := false
		for _, g := range expanded {
			for _, pattern := range g.patterns {
				p := pattern
				if strings.HasPrefix(p, "~/") {
					p = filepath.Join(s.deps.HomeDir, p[2:])
				}
				if ok, _ := filepath.Match(p, sess.Path); ok {
					g.sessions = append(g.sessions, sess)
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}
		if !matched {
			other.sessions = append(other.sessions, sess)
		}
	}

	s.groups = expanded
	if len(other.sessions) > 0 {
		s.groups = append(s.groups, other)
	}
}

func (s *SessionsSection) rebuildItems() {
	s.items = nil
	for gi, g := range s.groups {
		if len(g.sessions) == 0 {
			continue
		}
		s.items = append(s.items, flatItem{isGroup: true, groupIdx: gi})
		if !g.collapsed {
			for si := range g.sessions {
				s.items = append(s.items, flatItem{isGroup: false, groupIdx: gi, sessIdx: si})
			}
		}
	}
	if s.cursor >= len(s.items) {
		s.cursor = max(len(s.items)-1, 0)
	}
	if s.offset >= len(s.items) {
		s.offset = 0
	}
}

func (s *SessionsSection) fetchBranches(sessions model.SeshSessions) tea.Cmd {
	paths := make(map[string]bool)
	for _, key := range sessions.OrderedIndex {
		p := sessions.Directory[key].Path
		if p != "" {
			paths[p] = true
		}
	}
	cmds := make([]tea.Cmd, 0, len(paths))
	for p := range paths {
		path := p
		cmds = append(cmds, func() tea.Msg {
			found, branch, err := s.deps.Git.CurrentBranch(path)
			if err != nil || !found {
				return branchLoadedMsg{path: path, branch: ""}
			}
			return branchLoadedMsg{path: path, branch: strings.TrimSpace(branch)}
		})
	}
	return tea.Batch(cmds...)
}

func (s *SessionsSection) applyBranch(path, branch string) {
	for _, g := range s.groups {
		for i := range g.sessions {
			if g.sessions[i].Path == path {
				g.sessions[i].Branch = branch
			}
		}
	}
}

func (s *SessionsSection) fetchStatuses(sessions model.SeshSessions) tea.Cmd {
	paths := make(map[string]bool)
	for _, key := range sessions.OrderedIndex {
		p := sessions.Directory[key].Path
		if p != "" {
			paths[p] = true
		}
	}
	cmds := make([]tea.Cmd, 0, len(paths))
	for p := range paths {
		path := p
		cmds = append(cmds, func() tea.Msg {
			status, err := s.deps.Git.StatusSummary(path)
			if err != nil {
				return statusLoadedMsg{path: path, status: ""}
			}
			parts := make([]string, 0, 4)
			if status.Staged > 0 {
				parts = append(parts, fmt.Sprintf("+%d", status.Staged))
			}
			if status.Unstaged > 0 {
				parts = append(parts, fmt.Sprintf("~%d", status.Unstaged))
			}
			if status.Deleted > 0 {
				parts = append(parts, fmt.Sprintf("-%d", status.Deleted))
			}
			if status.Untracked > 0 {
				parts = append(parts, fmt.Sprintf("!%d", status.Untracked))
			}
			return statusLoadedMsg{path: path, status: strings.Join(parts, " ")}
		})
	}
	return tea.Batch(cmds...)
}

func (s *SessionsSection) applyStatus(path, status string) {
	for _, g := range s.groups {
		for i := range g.sessions {
			if g.sessions[i].Path == path {
				g.sessions[i].GitStatus = status
			}
		}
	}
}

func (s *SessionsSection) cursorUp(n int) {
	s.cursor -= n
	if s.cursor < 0 {
		s.cursor = 0
	}
	if s.cursor < s.offset {
		s.offset = s.cursor
	}
}

func (s *SessionsSection) cursorDown(n int) {
	s.cursor += n
	max := len(s.items) - 1
	if max < 0 {
		max = 0
	}
	if s.cursor > max {
		s.cursor = max
	}
	visible := s.visibleCount()
	if s.cursor >= s.offset+visible {
		s.offset = s.cursor - visible + 1
	}
}

func (s *SessionsSection) visibleCount() int {
	if s.viewHeight <= 0 {
		return 20
	}
	available := s.viewHeight - 2
	if available < 1 {
		return 1
	}
	return available
}

func (s *SessionsSection) toggleGroup() {
	if len(s.items) == 0 {
		return
	}
	item := s.items[s.cursor]
	if !item.isGroup {
		return
	}
	g := s.groups[item.groupIdx]
	g.collapsed = !g.collapsed
	s.rebuildItems()
}

func (s *SessionsSection) killSession() tea.Cmd {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[s.cursor]
	if item.isGroup {
		return nil
	}
	g := s.groups[item.groupIdx]
	if _, err := s.deps.Tmux.KillSession(g.sessions[item.sessIdx].Name); err != nil {
		slog.Error("failed to kill session", "name", g.sessions[item.sessIdx].Name, "error", err)
	}
	return s.Init()
}

func (s *SessionsSection) selectItem() {
	if len(s.items) == 0 {
		return
	}
	item := s.items[s.cursor]
	if item.isGroup {
		s.toggleGroup()
		return
	}
	g := s.groups[item.groupIdx]
	s.chosen = g.sessions[item.sessIdx].Name
}

// HoveredSession returns the name and path of the session under the cursor.
// Returns empty strings if cursor is on a group or items are empty.
func (s *SessionsSection) HoveredSession() (name, path string, windows int) {
	if len(s.items) == 0 {
		return "", "", 0
	}
	item := s.items[s.cursor]
	if item.isGroup {
		return "", "", 0
	}
	g := s.groups[item.groupIdx]
	sess := g.sessions[item.sessIdx]
	name = sess.Name
	path = sess.Path
	if after, ok := strings.CutPrefix(path, s.deps.HomeDir); ok {
		path = filepath.Join("~", after)
	}
	windows = sess.Windows
	return name, path, windows
}

func (s *SessionsSection) View(width, height int, focused bool) string {
	s.viewHeight = height

	// Guard: Minimum layout size checks
	const minWidth = 34
	if width < minWidth {
		msg := fmt.Sprintf("  Enlarge pane to see sessions (need ≥%d cols, have %d)", minWidth, width)
		return lipgloss.NewStyle().Faint(true).Width(width).Height(height).Render(msg)
	}

	// State Guards: Loading or Empty List
	if s.loading {
		return lipgloss.NewStyle().Faint(true).Width(width).Height(height).Render("  Loading sessions...")
	}
	if len(s.items) == 0 {
		return lipgloss.NewStyle().Faint(true).Width(width).Height(height).Render("  No sessions found")
	}

	// Calculate active available viewing rows
	chrome := 4 // Accounts for title header line space
	available := height - chrome
	if available < 1 {
		available = 5
	}

	end := min(s.offset+available, len(s.items))

	var b strings.Builder

	// Style Definitions
	groupStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)

	// Render the section title
	groupName := GroupNameRender(s.config.Title, width)
	b.WriteString(groupName.Render(s.config.Title))
	b.WriteString("\n\n")

	for i := s.offset; i < end; i++ {
		item := s.items[i]
		var prefix string
		if i == s.cursor {
			prefix = cursorStyle.Render("▸ ")
		} else {
			prefix = "  "
		}

		if item.isGroup {
			g := s.groups[item.groupIdx]
			groupLine := prefix + groupStyle.Render("") + g.name + groupStyle.Render(fmt.Sprintf(" (%d)", len(g.sessions)))
			b.WriteString(lipgloss.NewStyle().Width(width).Render(groupLine))
			b.WriteString("\n")
			continue
		}

		g := s.groups[item.groupIdx]
		sess := g.sessions[item.sessIdx]

		windNum := fmt.Sprintf("%4d.", item.sessIdx+1)
		name := sess.Name

		branch := ""
		if sess.Branch != "" {
			branch = fmt.Sprintf("[%s]", sess.Branch)
		}

		gitStatus := ""
		if sess.GitStatus != "" {
			gitStatus = fmt.Sprintf("[%s]", sess.GitStatus)
		}

		meta := ""
		if sess.Attached > 0 {
			meta = "∗"
		}

		// Join columns together line by line
		line := SessionLineRender(windNum, name, branch, gitStatus, meta, width)
		b.WriteString(prefix)
		b.WriteString(line)
		b.WriteString("\n")
	}

	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder())
	if focused {
		style = style.BorderForeground(lipgloss.Color("14"))
	}
	return style.Render(b.String())
}
