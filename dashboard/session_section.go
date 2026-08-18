package dashboard

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

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

type currentSessionMsg struct {
	name string
}

type SessionsSection struct {
	config        model.DashboardSectionConfig
	deps          SectionDeps
	sessions      []model.SeshSession
	filtered      []model.SeshSession // filtered view when filtering
	cursor        int
	offset        int
	loading       bool
	chosen        string
	totalSessions int
	viewHeight    int
	sortMode      string // "name" | "recent" | "created"
	filtering     bool
	filterQuery   string
	currentName   string
}

func NewSessionsSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &SessionsSection{
		config:   cfg,
		deps:     deps,
		loading:  true,
		sortMode: "name",
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

// SortLabel implements Sorter.
func (s *SessionsSection) SortLabel() string {
	return s.sortMode
}

// Filtering implements Filterer.
func (s *SessionsSection) Filtering() bool {
	return s.filtering
}

// FilterQuery implements Filterer.
func (s *SessionsSection) FilterQuery() string {
	return s.filterQuery
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
		s.sessions = flattenSessions(msg.sessions)
		s.totalSessions = len(msg.sessions.OrderedIndex)
		s.applySort()
		s.applyFilter()
		return s, tea.Batch(s.fetchBranches(), s.fetchStatuses(), s.fetchCurrentSession())

	case branchLoadedMsg:
		s.applyBranch(msg.path, msg.branch)
		return s, nil

	case statusLoadedMsg:
		s.applyStatus(msg.path, msg.status)
		return s, nil

	case currentSessionMsg:
		s.currentName = msg.name
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
	if s.filtering {
		return s.handleFilterKey(msg)
	}
	switch msg.String() {
	case "j", "down":
		s.cursorDown(1)
	case "k", "up":
		s.cursorUp(1)
	case "enter":
		s.selectItem()
	case "ctrl+d":
		return s, s.killSession()
	case "s":
		s.cycleSortMode()
	case "/":
		s.filtering = true
		s.filterQuery = ""
		s.applyFilter()
	}
	return s, nil
}

// handleFilterKey consumes keys while type-to-filter is active: printable
// characters append to the query, backspace (and its ctrl+h / ctrl+backspace
// aliases) delete the last rune, esc/enter exit and clear the filter.
func (s *SessionsSection) handleFilterKey(msg tea.KeyPressMsg) (*SessionsSection, tea.Cmd) {
	if isBackspaceKey(msg) {
		if s.filterQuery != "" {
			r := []rune(s.filterQuery)
			s.filterQuery = string(r[:len(r)-1])
		}
		s.applyFilter()
		return s, nil
	}
	switch msg.String() {
	case "esc", "enter":
		s.filtering = false
		s.filterQuery = ""
		s.applyFilter()
	default:
		if msg.Text != "" {
			s.filterQuery += msg.Text
			s.applyFilter()
		}
	}
	return s, nil
}

// cycleSortMode advances sortMode name → recent → created → name and re-sorts.
func (s *SessionsSection) cycleSortMode() {
	switch s.sortMode {
	case "name":
		s.sortMode = "recent"
	case "recent":
		s.sortMode = "created"
	default:
		s.sortMode = "name"
	}
	s.applySort()
	s.applyFilter()
}

// applySort sorts the master list (s.sessions) by the current sortMode.
func (s *SessionsSection) applySort() {
	sort.SliceStable(s.sessions, func(i, j int) bool {
		switch s.sortMode {
		case "recent":
			ti := timeOrZero(s.sessions[i].LastAttached)
			tj := timeOrZero(s.sessions[j].LastAttached)
			if !ti.Equal(tj) {
				return ti.After(tj)
			}
		case "created":
			ti := timeOrZero(s.sessions[i].Created)
			tj := timeOrZero(s.sessions[j].Created)
			if !ti.Equal(tj) {
				return ti.After(tj)
			}
		default:
			return s.sessions[i].Name < s.sessions[j].Name
		}
		return s.sessions[i].Name < s.sessions[j].Name
	})
}

// applyFilter rebuilds the filtered view from the master list and clamps the
// cursor.
func (s *SessionsSection) applyFilter() {
	if !s.filtering || s.filterQuery == "" {
		s.filtered = nil
		s.clampCursor()
		return
	}
	q := strings.ToLower(s.filterQuery)
	out := make([]model.SeshSession, 0, len(s.sessions))
	for _, sess := range s.sessions {
		if strings.Contains(strings.ToLower(sess.Name), q) {
			out = append(out, sess)
		}
	}
	s.filtered = out
	s.clampCursor()
}

// visible returns the currently displayed list (filtered view while filtering,
// the full sorted list otherwise).
func (s *SessionsSection) visible() []model.SeshSession {
	if s.filtering && s.filtered != nil {
		return s.filtered
	}
	return s.sessions
}

// timeOrZero returns t as a non-pointer time.Time, treating nil as the zero
// time (so nil times sort before real ones).
func timeOrZero(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// flattenSessions returns every tmux session as a flat list, sorted
// alphabetically by name (stable order).
func flattenSessions(sessions model.SeshSessions) []model.SeshSession {
	out := make([]model.SeshSession, 0, len(sessions.OrderedIndex))
	for _, key := range sessions.OrderedIndex {
		out = append(out, sessions.Directory[key])
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (s *SessionsSection) fetchCurrentSession() tea.Cmd {
	return func() tea.Msg {
		sess, ok := s.deps.Lister.GetAttachedTmuxSession()
		if !ok {
			return currentSessionMsg{}
		}
		return currentSessionMsg{name: sess.Name}
	}
}

func (s *SessionsSection) fetchBranches() tea.Cmd {
	paths := make(map[string]bool)
	for _, sess := range s.sessions {
		if sess.Path != "" {
			paths[sess.Path] = true
		}
	}
	cmds := make([]tea.Cmd, 0, len(paths))
	for path := range paths {
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
	for i := range s.sessions {
		if s.sessions[i].Path == path {
			s.sessions[i].Branch = branch
		}
	}
	s.applyFilter()
}

func (s *SessionsSection) fetchStatuses() tea.Cmd {
	paths := make(map[string]bool)
	for _, sess := range s.sessions {
		if sess.Path != "" {
			paths[sess.Path] = true
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
			return statusLoadedMsg{path: path, status: formatGitStatus(status)}
		})
	}
	return tea.Batch(cmds...)
}

func (s *SessionsSection) applyStatus(path, status string) {
	for i := range s.sessions {
		if s.sessions[i].Path == path {
			s.sessions[i].GitStatus = status
		}
	}
	s.applyFilter()
}

func (s *SessionsSection) clampCursor() {
	n := len(s.visible())
	if s.cursor >= n {
		s.cursor = max(n-1, 0)
	}
	if s.offset >= n {
		s.offset = 0
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
	maxIdx := max(len(s.visible())-1, 0)
	if s.cursor > maxIdx {
		s.cursor = maxIdx
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
	return max(s.viewHeight, 1)
}

// ClickAt moves the cursor to the clicked view row, scrolling to reveal it.
func (s *SessionsSection) ClickAt(row int) {
	n := len(s.visible())
	if n == 0 {
		return
	}
	s.cursor = min(max(s.offset+row, 0), n-1)
	if s.cursor < s.offset {
		s.offset = s.cursor
	}
	if visible := s.visibleCount(); s.cursor >= s.offset+visible {
		s.offset = s.cursor - visible + 1
	}
}

func (s *SessionsSection) killSession() tea.Cmd {
	if len(s.visible()) == 0 {
		return nil
	}
	sess := s.visible()[s.cursor]
	if _, err := s.deps.Tmux.KillSession(sess.Name); err != nil {
		slog.Error("failed to kill session", "name", sess.Name, "error", err)
	}
	return s.Init()
}

func (s *SessionsSection) selectItem() {
	if len(s.visible()) == 0 {
		return
	}
	s.chosen = s.visible()[s.cursor].Name
}

// HoveredSession returns the name and path of the session under the cursor.
// Returns empty strings if the list is empty.
func (s *SessionsSection) HoveredSession() (name, path string, windows int) {
	if len(s.visible()) == 0 {
		return "", "", 0
	}
	sess := s.visible()[s.cursor]
	name = sess.Name
	path = sess.Path
	if after, ok := strings.CutPrefix(path, s.deps.HomeDir); ok {
		path = filepath.Join("~", after)
	}
	windows = sess.Windows
	return name, path, windows
}

func (s *SessionsSection) ViewBorderless(width, height int, focused bool) (string, string) {
	s.viewHeight = height

	title := s.config.Title
	if title == "" {
		title = "Sessions"
	}

	// Guard: Minimum layout size checks
	const minWidth = 34
	if width < minWidth {
		msg := fmt.Sprintf("  Enlarge pane to see sessions (need ≥%d cols, have %d)", minWidth, width)
		return title, msg
	}

	// State Guards: Loading or Empty List
	if s.loading {
		return title, "  Loading sessions..."
	}
	if len(s.sessions) == 0 {
		return title, "  No sessions found"
	}

	// Calculate active available viewing rows (the pane content area, already
	// reduced by the shared frame's top/bottom borders).
	visible := s.visible()
	available := max(height, 1)
	end := min(s.offset+available, len(visible))

	var b strings.Builder
	for i := s.offset; i < end; i++ {
		b.WriteString(s.renderItemFocused(i, width, focused))
		b.WriteString("\n")
	}

	return title, b.String()
}

// renderItem renders a single flat session row for Tab 1.
func (s *SessionsSection) renderItem(i, width int) string {
	return s.renderItemFocused(i, width, true)
}

// renderItemFocused is renderItem with an explicit focused flag, so unfocused
// panes render a dimmed selection highlight.
func (s *SessionsSection) renderItemFocused(i, width int, focused bool) string {
	sess := s.visible()[i]
	dir := collapseHome(sess.Path, s.deps.HomeDir)
	current := sess.Name == s.currentName && s.currentName != ""
	return renderOpenRowFocused(width, i == s.cursor, current, focused, sess.Name, sess.Attached, sess.Windows, dir, sess.Branch, sess.GitStatus, sess.Created, sess.Alerts)
}
