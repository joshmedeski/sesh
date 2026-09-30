package dashboard

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

// configuredLoadedMsg carries the config-source sessions (sorted) plus the set
// of tmux session names currently running (used for the running-state column).
type configuredLoadedMsg struct {
	sessions []model.SeshSession
	running  map[string]bool
	err      error
}

// ConfiguredSection lists pre-configured sessions from the sesh config.
// Selecting a session sets Chosen() to the session name; the CLI connector
// opens it.
type ConfiguredSection struct {
	config   model.DashboardSectionConfig
	deps     SectionDeps
	sessions []model.SeshSession
	ListState
	running map[string]bool
	loading bool
	chosen  string
}

func NewConfiguredSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &ConfiguredSection{
		config:  cfg,
		deps:    deps,
		loading: true,
		running: map[string]bool{},
	}
}

func (s *ConfiguredSection) Width() float64 {
	return s.config.Width
}

func (s *ConfiguredSection) Name() string {
	return s.config.Title
}

func (s *ConfiguredSection) TotalItems() int {
	return len(s.sessions)
}

func (s *ConfiguredSection) Chosen() string {
	return s.chosen
}

// Filtering implements Filterer.
func (s *ConfiguredSection) Filtering() bool {
	return s.filtering
}

// FilterQuery implements Filterer.
func (s *ConfiguredSection) FilterQuery() string {
	return s.filterQuery
}

// Init loads config-source sessions and the running tmux set, then kicks off
// async git branch/status enrichment for the loaded paths.
func (s *ConfiguredSection) Init() tea.Cmd {
	return s.fetch()
}

func (s *ConfiguredSection) fetch() tea.Cmd {
	return func() tea.Msg {
		configSessions, err := s.deps.Lister.List(lister.ListOptions{Config: true})
		if err != nil {
			return configuredLoadedMsg{err: err}
		}

		running := map[string]bool{}
		if tmuxSessions, err := s.deps.Lister.List(lister.ListOptions{Tmux: true}); err == nil {
			for _, key := range tmuxSessions.OrderedIndex {
				running[tmuxSessions.Directory[key].Name] = true
			}
		}

		sessions := make([]model.SeshSession, 0, len(configSessions.OrderedIndex))
		for _, key := range configSessions.OrderedIndex {
			sessions = append(sessions, configSessions.Directory[key])
		}
		sort.Slice(sessions, func(i, j int) bool { return sessions[i].Name < sessions[j].Name })

		return configuredLoadedMsg{sessions: sessions, running: running}
	}
}

func (s *ConfiguredSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case configuredLoadedMsg:
		if msg.err != nil {
			return s, nil
		}
		s.loading = false
		s.sessions = msg.sessions
		s.running = msg.running
		s.clampCursor()
		return s, tea.Batch(fetchBranches(s.deps.Git, s.sessions), fetchStatuses(s.deps.Git, s.sessions))

	case branchLoadedMsg:
		s.applyBranch(msg.path, msg.branch)
		return s, nil

	case statusLoadedMsg:
		s.applyStatus(msg.path, msg.status)
		return s, nil

	case tea.KeyPressMsg:
		return s.handleKey(msg)
	}
	return s, nil
}

func (s *ConfiguredSection) handleKey(msg tea.KeyPressMsg) (Section, tea.Cmd) {
	if s.filtering {
		s.handleFilterKey(msg)
		return s, nil
	}
	switch msg.String() {
	case "j", "down":
		s.cursorDown(1)
	case "k", "up":
		s.cursorUp(1)
	case "enter":
		s.selectItem()
	case "r":
		s.loading = true
		return s, s.fetch()
	case "/":
		s.filtering = true
		s.filterQuery = ""
		s.applyFilter()
	}
	return s, nil
}

// handleFilterKey consumes keys while type-to-filter is active (mirrors
// SessionsSection): printable characters append to the query, backspace (and
// its ctrl+h / ctrl+backspace aliases) delete the last rune, j/k and the arrow
// keys move the cursor through the filtered results, enter selects the
// highlighted filtered item and exits filtering, and esc cancels filtering
// without selecting.
func (s *ConfiguredSection) handleFilterKey(msg tea.KeyPressMsg) {
	s.ListState.handleFilterKey(msg, s.sessions, configuredMatch, s.selectItem)
}

// applyFilter rebuilds the filtered view from the master list and clamps the
// cursor.
func (s *ConfiguredSection) applyFilter() {
	s.ListState.applyFilter(s.sessions, configuredMatch)
}

// configuredMatch reports whether a session matches the query by name.
func configuredMatch(sess model.SeshSession, q string) bool {
	return strings.Contains(strings.ToLower(sess.Name), q)
}

// visible returns the currently displayed list.
func (s *ConfiguredSection) visible() []model.SeshSession {
	return s.ListState.visible(s.sessions)
}

func (s *ConfiguredSection) clampCursor() {
	s.ListState.clampCursor(len(s.visible()))
}

func (s *ConfiguredSection) cursorUp(n int) {
	s.ListState.cursorUp(n)
}

func (s *ConfiguredSection) cursorDown(n int) {
	s.ListState.cursorDown(n, len(s.visible()))
}

// ClickAt moves the cursor to the clicked view row, scrolling to reveal it.
func (s *ConfiguredSection) ClickAt(row int) {
	s.ListState.ClickAt(row, len(s.visible()))
}

func (s *ConfiguredSection) selectItem() {
	if len(s.visible()) == 0 {
		return
	}
	s.chosen = s.visible()[s.cursor].Name
}

func (s *ConfiguredSection) applyBranch(path, branch string) {
	applyBranch(s.sessions, path, branch)
	s.applyFilter()
}

func (s *ConfiguredSection) applyStatus(path, status string) {
	applyStatus(s.sessions, path, status)
	s.applyFilter()
}

// ViewBorderless renders the configured list with columns:
// marker(2) | state(2) | name(24) | alias(longest) | path(fill) | branch(16) | status(12).
func (s *ConfiguredSection) ViewBorderless(width, height int, focused bool) (string, string) {
	s.viewHeight = max(height-1, 1)

	title := s.config.Title
	if title == "" {
		title = "Configured"
	}

	const minWidth = 34
	if width < minWidth {
		msg := fmt.Sprintf("  Enlarge pane to see sessions (need ≥%d cols, have %d)", minWidth, width)
		return title, msg
	}

	if s.loading {
		return title, "  Loading sessions..."
	}
	if len(s.sessions) == 0 {
		return title, "  No sessions configured"
	}

	visible := s.visible()
	end := min(s.offset+s.viewHeight, len(visible))

	var b strings.Builder
	b.WriteString(render.RenderConfiguredHeader(width, aliasColumnWidth(visible), max(s.deps.IconWidth, 1)))
	b.WriteString("\n")
	for i := s.offset; i < end; i++ {
		sess := visible[i]
		path := shortenHome(s.deps, sess.Path)
		b.WriteString(render.RenderConfiguredRowFocused(width, aliasColumnWidth(visible), iconCol(s.deps, sess, "", i == s.cursor), i == s.cursor, focused, sess.Name, sess.Alias, sess.StartupCommand, s.running[sess.Name], path, sess.Branch, sess.GitStatus))
		b.WriteString("\n")
	}

	return title, b.String()
}
