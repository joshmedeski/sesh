package dashboard

import (
	"fmt"
	"os"
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
}

func NewSessionsSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &SessionsSection{
		config:  cfg,
		deps:    deps,
		loading: true,
	}
}

func (s *SessionsSection) Name() string {
	return s.config.Title
}

func (s *SessionsSection) TotalItems() int {
	return s.totalSessions
}

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
		return s, s.fetchBranches(msg.sessions)

	case branchLoadedMsg:
		s.applyBranch(msg.path, msg.branch)
		return s, nil

	case tea.KeyPressMsg:
		return s.handleKey(msg), nil
	}
	return s, nil
}

func (s *SessionsSection) Chosen() string {
	return s.chosen
}

func (s *SessionsSection) handleKey(msg tea.KeyPressMsg) *SessionsSection {
	switch msg.String() {
	case "j", "down":
		s.cursorDown(1)
	case "k", "up":
		s.cursorUp(1)
	case "t":
		s.toggleGroup()
	case "enter":
		s.selectItem()
	}
	return s
}

func (s *SessionsSection) groupSessions(sessions model.SeshSessions) {
	expanded := make([]*group, 0, len(s.config.Groups))
	for i := range s.config.Groups {
		g := &group{
			name:     s.config.Groups[i].Name,
			patterns: s.config.Groups[i].Patterns,
		}
		expanded = append(expanded, g)
	}

	homeDir, _ := os.UserHomeDir()
	other := &group{name: "Other"}

	for _, key := range sessions.OrderedIndex {
		sess := sessions.Directory[key]
		matched := false
		for _, g := range expanded {
			for _, pattern := range g.patterns {
				p := pattern
				if strings.HasPrefix(p, "~/") {
					p = filepath.Join(homeDir, p[2:])
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
	return 20
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

func (s *SessionsSection) View(width, height int) string {
	if s.loading {
		msg := lipgloss.NewStyle().Faint(true).Render("  Loading sessions...")
		lines := strings.Count(msg, "\n")
		var b strings.Builder
		b.WriteString(msg)
		for i := lines + 1; i < height; i++ {
			b.WriteString("\n")
		}
		return b.String()
	}

	if len(s.items) == 0 {
		msg := lipgloss.NewStyle().Faint(true).Render("  No sessions found")
		lines := strings.Count(msg, "\n")
		var b strings.Builder
		b.WriteString(msg)
		for i := lines + 1; i < height; i++ {
			b.WriteString("\n")
		}
		return b.String()
	}

	chrome := 2
	available := height - chrome
	if available < 1 {
		available = 5
	}

	end := min(s.offset+available, len(s.items))

	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().Bold(true)
	groupStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(33))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)
	sessionStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15))
	branchStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8))
	pathStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8)).Faint(true)
	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8))
	attachedStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2))

	b.WriteString(sectionStyle.Render("  "+s.config.Title) + "\n\n")

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
			if g.collapsed {
				fmt.Fprintf(&b, "%s%s ▶ %s (%d)\n", prefix, groupStyle.Render(""), g.name, len(g.sessions))
			} else {
				fmt.Fprintf(&b, "%s%s ▼ %s\n", prefix, groupStyle.Render(""), g.name)
			}
		} else {
			g := s.groups[item.groupIdx]
			sess := g.sessions[item.sessIdx]

			name := sessionStyle.Render(sess.Name)
			branch := ""
			if sess.Branch != "" {
				branch = branchStyle.Render("[" + sess.Branch + "]")
			}
			path := pathStyle.Render(sess.Path)

			right := ""
			if sess.Windows > 0 {
				right = metaStyle.Render(fmt.Sprintf("%d windows", sess.Windows))
			}
			if sess.Attached > 0 {
				if right != "" {
					right += "  "
				}
				right += attachedStyle.Render("attached")
			}

			line := fmt.Sprintf("  %s  %s  %s  %s", name, branch, path, right)
			b.WriteString(prefix + line + "\n")
		}
	}

	// Pad to fill allocated height
	lines := strings.Count(b.String(), "\n")
	for i := lines; i < height; i++ {
		b.WriteString("\n")
	}

	return b.String()
}
