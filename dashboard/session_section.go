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
		s.selectItem()
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
			collapsed: true, // start with group collapsed. user can toggle to expand.
		}
		expanded = append(expanded, g)
	}

	homeDir, _ := os.UserHomeDir()
	other := &group{name: "Other", collapsed: true} // start with group collapsed

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

// TODO: review this
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
	_, _ = s.deps.Tmux.KillSession(g.sessions[item.sessIdx].Name)
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
func (s *SessionsSection) HoveredSession() (name, path string) {
	if len(s.items) == 0 {
		return "", ""
	}
	item := s.items[s.cursor]
	if item.isGroup {
		return "", ""
	}
	g := s.groups[item.groupIdx]
	sess := g.sessions[item.sessIdx]
	name = sess.Name
	path = sess.Path
	if after, ok := strings.CutPrefix(path, s.deps.HomeDir); ok {
		path = filepath.Join("~", after)
	}
	return name, path
}

func (s *SessionsSection) View(width, height int) string {
	s.viewHeight = height
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

	const minWidth = 34
	if width < minWidth {
		var bb strings.Builder
		msg := lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("  Enlarge pane to see sessions (need ≥%d cols, have %d)", minWidth, width))
		bb.WriteString(msg)
		for i := 1; i < height; i++ {
			bb.WriteString("\n")
		}
		return bb.String()
	}

	chrome := 2
	available := height - chrome
	if available < 1 {
		available = 5
	}

	end := min(s.offset+available, len(s.items))

	var b strings.Builder

	// TODO: make the colors configurable
	sectionStyle := lipgloss.NewStyle().Bold(true)
	groupStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)
	sessionStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15))
	branchStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(12))
	gitStatusStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(5))
	// metaStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8))
	attachedStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2))
	numStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8))
	// cursorBg := lipgloss.NewStyle().Background(lipgloss.ANSIColor(25))

	b.WriteString(sectionStyle.Render("" + s.config.Title))
	b.WriteString("\n\n")

	// Content-aware column widths based on actual session data
	numW := 5
	overhead := 8 // "   "(3) + numW(5) + gaps(3*1)
	colSpace := width - overhead

	maxName := 5
	maxBranch := 0
	maxGit := 0
	for _, g := range s.groups {
		for _, sess := range g.sessions {
			if l := len(sess.Name); l > maxName {
				maxName = l
			}
			if l := len(sess.Branch); l > maxBranch {
				maxBranch = l
			}
			if l := len(sess.GitStatus); l > maxGit {
				maxGit = l
			}
		}
	}

	// Ideal widths from content (capped)
	nameW := max(int(float64(colSpace)*0.25), 15)
	branchW := max(int(float64(colSpace)*0.25), 10)
	gitW := max(int(float64(colSpace)*0.45), 15)
	metaW := max(int(float64(colSpace)*0.25), 25)

	// Shrink iteratively until everything fits in colSpace
	for nameW+branchW+gitW+metaW > colSpace {
		if metaW > 2 {
			metaW--
			continue
		}
		if gitW > 2 {
			gitW--
			continue
		}
		if branchW > 3 {
			branchW--
			continue
		}
		if nameW > 3 {
			nameW--
			continue
		}
		break
	}

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
			b.WriteString(lipgloss.NewStyle().Width(width).MaxWidth(width).Render(groupLine))
			b.WriteString("\n")

		} else {
			g := s.groups[item.groupIdx]
			sess := g.sessions[item.sessIdx]

			colStyle := lipgloss.NewStyle().Width(numW)
			num := colStyle.Render(numStyle.Render(fmt.Sprintf("%4d.", item.sessIdx+1)))

			colStyle = lipgloss.NewStyle().Width(nameW).MaxWidth(nameW)
			name := colStyle.Render(sessionStyle.Render(sess.Name))
			// if len(name) > nameW-2 {
			// 	name = name[:nameW-2] + "…"
			// }

			colStyle = lipgloss.NewStyle().Width(branchW).MaxWidth(branchW)
			branchContent := sess.Branch
			// if len(branchContent) > branchW-2 {
			// 	branchContent = branchContent[:max(branchW-3, 0)] + "…"
			// }
			// if no branch, don't show the brackets
			branch := ""
			if sess.Branch == "" {
				branchContent = ""
				branch = colStyle.Render("")
			} else {
				branch = colStyle.Render(branchStyle.Render("[" + branchContent + "]"))
			}

			colStyle = lipgloss.NewStyle().Width(gitW).MaxWidth(gitW)
			gitContent := sess.GitStatus
			// if len(gitContent) > gitW-2 {
			// 	gitContent = gitContent[:max(gitW-3, 0)] + "…"
			// }
			// if no status, don't show the brackets
			gitStatus := ""
			if sess.GitStatus == "" {
				gitContent = ""
				gitStatus = colStyle.Render("")
			} else {
				gitStatus = colStyle.Render(gitStatusStyle.Render("[" + gitContent + "]"))
			}

			right := ""
			if sess.Attached > 0 {
				right += attachedStyle.Render("∗")
			}
			colStyle = lipgloss.NewStyle().Width(metaW).MaxWidth(metaW)
			meta := colStyle.Render(right)

			// if i == s.cursor {
			// 	num = cursorBg.Render(num)
			// 	name = cursorBg.Render(name)
			// 	branch = cursorBg.Render(branch)
			// 	gitStatus = cursorBg.Render(gitStatus)
			// 	meta = cursorBg.Render(meta)
			// }

			line := fmt.Sprintf("%s%s %s %s %s", num, name, branch, gitStatus, meta)

			// if i == s.cursor {
			// 	lineWidth := lipgloss.Width(line) // Correctly ignores hidden ANSI style codes
			//
			// 	if lineWidth < width {
			// 		// Create the padding spaces
			// 		padding := strings.Repeat(" ", width-lineWidth)
			// 		// Render BOTH the text and the spaces inside the background color
			// 		line = cursorBg.Render(line + padding)
			// 	} else {
			// 		line = cursorBg.Render(line)
			// 	}
			// }

			// b.WriteString(prefix)
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	// Pad to fill allocated height
	lines := strings.Count(b.String(), "\n")
	for i := lines; i < height; i++ {
		b.WriteString("\n")
	}

	return b.String()
}
