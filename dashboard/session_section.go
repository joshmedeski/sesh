package dashboard

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

type sessionsLoadedMsg struct {
	section  *SessionsSection
	sessions model.SeshSessions
	err      error
}

type currentSessionMsg struct {
	name string
}

type sessionIssuesLoadedMsg struct {
	section *SessionsSection
	entries map[string]model.WorktreeEntry
}

type SessionsSection struct {
	config   model.DashboardSectionConfig
	deps     SectionDeps
	sessions []model.SeshSession
	ListState
	loading       bool
	chosen        string
	totalSessions int
	sortMode      string
	currentName   string
	sortOrder     model.SortOrder
	rank          map[string]int
	columns       []string
	worktrees     []model.WorktreeConfig
	issues        map[string]model.WorktreeEntry
}

func NewSessionsSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &SessionsSection{
		config:   cfg,
		deps:     deps,
		loading:  true,
		sortMode: "name",
	}
}

func NewSourcesSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	s := NewSessionsSection(cfg, deps).(*SessionsSection)
	s.sortOrder = cfg.Sources
	s.columns = resolveColumns(cfg.Columns, slices.Concat(render.OpenColumns, render.IssueColumns), cfg.Type)
	s.sortMode = s.sortModes()[0]
	return s
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

func (s *SessionsSection) TmuxCount() int {
	n := 0
	for _, sess := range s.sessions {
		if sess.Src == "tmux" {
			n++
		}
	}
	return n
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
	opts := s.listOptions()
	list := func() tea.Msg {
		sessions, err := s.deps.Lister.List(opts)
		return sessionsLoadedMsg{section: s, sessions: sessions, err: err}
	}
	if !s.showsIssues() {
		return list
	}
	return tea.Batch(list, s.fetchIssues())
}

func (s *SessionsSection) showsIssues() bool {
	return slices.ContainsFunc(s.columns, func(id string) bool { return slices.Contains(render.IssueColumns, id) })
}

func (s *SessionsSection) fetchIssues() tea.Cmd {
	worktrees, svc := s.worktrees, s.deps.Worktree
	return func() tea.Msg {
		entries := map[string]model.WorktreeEntry{}
		if svc == nil {
			return sessionIssuesLoadedMsg{section: s, entries: entries}
		}
		for _, wc := range worktrees {
			list, err := svc.List(model.WorktreeListOpts{Repo: wc.Repo})
			if err != nil {
				slog.Warn("dashboard: listing worktrees for issue columns", "repo", wc.Repo, "error", err)
				continue
			}
			for _, e := range list {
				entries[filepath.Clean(e.Path)] = e
			}
		}
		return sessionIssuesLoadedMsg{section: s, entries: entries}
	}
}

func (s *SessionsSection) issueFor(path string) render.Issue {
	if path == "" || len(s.issues) == 0 {
		return render.Issue{}
	}
	for p := filepath.Clean(path); ; p = filepath.Dir(p) {
		if e, ok := s.issues[p]; ok {
			return render.Issue{Number: e.Number, Title: e.Title, State: e.State}
		}
		if filepath.Dir(p) == p {
			return render.Issue{}
		}
	}
}

func (s *SessionsSection) listOptions() lister.ListOptions {
	groups := s.sortOrder.SortGroups()
	if len(groups) == 0 {
		return lister.ListOptions{Tmux: true}
	}
	opts := lister.ListOptions{SortOrder: s.sortOrder, HideDuplicates: true}
	for _, group := range groups {
		for _, src := range group {
			switch strings.ToLower(src) {
			case "tmux":
				opts.Tmux = true
			case "config":
				opts.Config = true
			case "tmuxinator":
				opts.Tmuxinator = true
			case "zoxide":
				opts.Zoxide = true
			}
		}
	}
	return opts
}

func (s *SessionsSection) sortModes() []string {
	if len(s.sortOrder.SortGroups()) > 0 {
		return []string{"order", "name", "recent", "created"}
	}
	return []string{"name", "recent", "created"}
}

func rankKey(sess model.SeshSession) string {
	return sess.Src + "\x00" + sess.Name
}

func (s *SessionsSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case sessionsLoadedMsg:
		if msg.section != s || msg.err != nil {
			return s, nil
		}
		s.loading = false
		s.sessions = flattenSessions(msg.sessions)
		s.rank = make(map[string]int, len(s.sessions))
		for i, sess := range s.sessions {
			s.rank[rankKey(sess)] = i
		}
		s.totalSessions = len(msg.sessions.OrderedIndex)
		s.applySort()
		s.applyFilter()
		return s, tea.Batch(fetchBranches(s.deps.Git, s.sessions), fetchStatuses(s.deps.Git, s.sessions), s.fetchCurrentSession())

	case branchLoadedMsg:
		s.applyBranch(msg.path, msg.branch)
		return s, nil

	case statusLoadedMsg:
		s.applyStatus(msg.path, msg.status)
		return s, nil

	case sessionIssuesLoadedMsg:
		if msg.section == s {
			s.issues = msg.entries
			s.applyFilter()
		}
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
	case "r":
		return s, s.Init()
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
	s.ListState.handleFilterKey(msg, s.sessions, s.match, s.selectItem)
	return s, nil
}

// cycleSortMode advances to the next sort mode (order, when a dashboard
// sort_order is set, then name → recent → created) and re-sorts.
func (s *SessionsSection) cycleSortMode() {
	modes := s.sortModes()
	next := 0
	for i, mode := range modes {
		if mode == s.sortMode {
			next = (i + 1) % len(modes)
		}
	}
	s.sortMode = modes[next]
	s.applySort()
	s.applyFilter()
}

// applySort sorts the master list (s.sessions) by the current sortMode.
func (s *SessionsSection) applySort() {
	sort.SliceStable(s.sessions, func(i, j int) bool {
		switch s.sortMode {
		case "order":
			return s.rank[rankKey(s.sessions[i])] < s.rank[rankKey(s.sessions[j])]
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
	s.ListState.applyFilter(s.sessions, s.match)
}

func (s *SessionsSection) match(sess model.SeshSession, q string) bool {
	return sessionsMatch(sess, q) || strings.Contains(strings.ToLower(s.issueFor(sess.Path).Title), q)
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

// flattenSessions returns every listed session as a flat list, in the order
// the lister returned them.
func flattenSessions(sessions model.SeshSessions) []model.SeshSession {
	out := make([]model.SeshSession, 0, len(sessions.OrderedIndex))
	for _, key := range sessions.OrderedIndex {
		out = append(out, sessions.Directory[key])
	}
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

func (s *SessionsSection) applyBranch(path, branch string) {
	applyBranch(s.sessions, path, branch)
	s.applyFilter()
}

func (s *SessionsSection) applyStatus(path, status string) {
	applyStatus(s.sessions, path, status)
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
	if sess.Src != "tmux" {
		return nil
	}
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
	s.viewHeight = max(height-1, 1)

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
	end := min(s.offset+s.viewHeight, len(visible))

	var b strings.Builder
	b.WriteString(render.RenderOpenHeader(width, s.columns, branchColumnWidth(visible), aliasColumnWidth(visible), max(s.deps.IconWidth, 1)))
	b.WriteString("\n")
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
	visible := s.visible()
	sess := visible[i]
	dir := render.CollapseHome(sess.Path, s.deps.HomeDir)
	current := sess.Name == s.currentName && s.currentName != ""
	return render.RenderOpenRowFocused(width, s.columns, branchColumnWidth(visible), aliasColumnWidth(visible), iconCol(s.deps, sess, "", i == s.cursor), i == s.cursor, current, focused, sess.Name, sess.Alias, sess.Attached, sess.Windows, dir, sess.Branch, sess.GitStatus, sess.LastAttached, sess.Alerts, s.issueFor(sess.Path))
}

func branchColumnWidth(sessions []model.SeshSession) int {
	w := 0
	for _, sess := range sessions {
		w = max(w, lipgloss.Width(sess.Branch))
	}
	return w
}

func aliasColumnWidth(sessions []model.SeshSession) int {
	aliases := make([]string, len(sessions))
	for i, sess := range sessions {
		aliases[i] = sess.Alias
	}
	return render.AliasColumn(aliases...)
}
