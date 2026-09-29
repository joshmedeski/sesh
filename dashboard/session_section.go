package dashboard

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

type sessionsLoadedMsg struct {
	sessions model.SeshSessions
	err      error
}

type currentSessionMsg struct {
	name string
}

type SessionsSection struct {
	config   model.DashboardSectionConfig
	deps     SectionDeps
	sessions []model.SeshSession
	ListState
	loading       bool
	chosen        string
	totalSessions int
	sortMode      string // "name" | "recent" | "created"
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
// aliases) delete the last rune, j/k and the arrow keys move the cursor
// through the filtered results, enter selects the highlighted filtered item
// and exits filtering, and esc cancels filtering without selecting.
func (s *SessionsSection) handleFilterKey(msg tea.KeyPressMsg) (*SessionsSection, tea.Cmd) {
	s.ListState.handleFilterKey(msg, s.sessions, sessionsMatch, s.selectItem)
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
	s.ListState.applyFilter(s.sessions, sessionsMatch)
}

// sessionsMatch reports whether a session matches the query by name or alias.
func sessionsMatch(sess model.SeshSession, q string) bool {
	return strings.Contains(strings.ToLower(sess.Name), q) || strings.Contains(strings.ToLower(sess.Alias), q)
}

// visible returns the currently displayed list (filtered view while filtering,
// the full sorted list otherwise).
func (s *SessionsSection) visible() []model.SeshSession {
	return s.ListState.visible(s.sessions)
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
			return statusLoadedMsg{path: path, status: render.FormatGitStatus(status)}
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
	s.ListState.clampCursor(len(s.visible()))
}

func (s *SessionsSection) cursorUp(n int) {
	s.ListState.cursorUp(n)
}

func (s *SessionsSection) cursorDown(n int) {
	s.ListState.cursorDown(n, len(s.visible()))
}

// ClickAt moves the cursor to the clicked view row, scrolling to reveal it.
func (s *SessionsSection) ClickAt(row int) {
	s.ListState.ClickAt(row, len(s.visible()))
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
	dir := render.CollapseHome(sess.Path, s.deps.HomeDir)
	current := sess.Name == s.currentName && s.currentName != ""
	return render.RenderOpenRowFocused(width, i == s.cursor, current, focused, sess.Name, sess.Alias, sess.Attached, sess.Windows, dir, sess.Branch, sess.GitStatus, sess.LastAttached, sess.Alerts)
}
