package dashboard

import (
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/model"
)

type worktreesLoadedMsg struct {
	repo    string
	entries []model.WorktreeEntry
	err     error
}

// WorktreeSection lists the worktrees of one [[worktree]] config entry on its
// own tab. Selecting a row sets ChosenWorktree.
type WorktreeSection struct {
	config   model.WorktreeConfig
	deps     SectionDeps
	sessions []model.SeshSession
	entries  map[string]model.WorktreeEntry
	ListState
	loading bool
	err     error
	chosen  int
}

func NewWorktreeSection(cfg model.WorktreeConfig, deps SectionDeps) *WorktreeSection {
	return &WorktreeSection{config: cfg, deps: deps, loading: true}
}

// TabTitle is the repo name without its owner, e.g. "sesh" for
// "joshmedeski/sesh".
func (s *WorktreeSection) TabTitle() string {
	repo := s.config.Repo
	if repo == "" {
		return s.config.Path
	}
	return repo[strings.LastIndex(repo, "/")+1:]
}

func (s *WorktreeSection) Width() float64 { return 0 }

func (s *WorktreeSection) Name() string { return s.config.Repo }

func (s *WorktreeSection) TotalItems() int { return len(s.sessions) }

func (s *WorktreeSection) Chosen() string { return "" }

// ChosenWorktree returns the connect options for the selected worktree.
func (s *WorktreeSection) ChosenWorktree() (model.WorktreeConnectOpts, bool) {
	if s.chosen == 0 {
		return model.WorktreeConnectOpts{}, false
	}
	return model.WorktreeConnectOpts{Number: s.chosen, Repo: s.config.Repo}, true
}

func (s *WorktreeSection) Filtering() bool { return s.filtering }

func (s *WorktreeSection) FilterQuery() string { return s.filterQuery }

func (s *WorktreeSection) Init() tea.Cmd { return s.fetch() }

func (s *WorktreeSection) fetch() tea.Cmd {
	repo := s.config.Repo
	return func() tea.Msg {
		if s.deps.Worktree == nil {
			return worktreesLoadedMsg{repo: repo}
		}
		entries, err := s.deps.Worktree.List(model.WorktreeListOpts{Repo: repo})
		return worktreesLoadedMsg{repo: repo, entries: entries, err: err}
	}
}

func (s *WorktreeSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case worktreesLoadedMsg:
		if msg.repo != s.config.Repo {
			return s, nil
		}
		s.loading = false
		s.err = msg.err
		s.sessions = make([]model.SeshSession, 0, len(msg.entries))
		s.entries = make(map[string]model.WorktreeEntry, len(msg.entries))
		for _, e := range msg.entries {
			s.sessions = append(s.sessions, model.SeshSession{Src: "worktree", Name: e.Title, Path: e.Path})
			s.entries[e.Path] = e
		}
		s.applyFilter()
		return s, tea.Batch(fetchBranches(s.deps.Git, s.sessions), fetchStatuses(s.deps.Git, s.sessions))

	case branchLoadedMsg:
		applyBranch(s.sessions, msg.path, msg.branch)
		s.applyFilter()
		return s, nil

	case statusLoadedMsg:
		applyStatus(s.sessions, msg.path, msg.status)
		s.applyFilter()
		return s, nil

	case tea.KeyPressMsg:
		return s.handleKey(msg)
	}
	return s, nil
}

func (s *WorktreeSection) handleKey(msg tea.KeyPressMsg) (Section, tea.Cmd) {
	if s.filtering {
		s.ListState.handleFilterKey(msg, s.sessions, s.match, s.selectItem)
		return s, nil
	}
	switch msg.String() {
	case "j", "down":
		s.ListState.cursorDown(1, len(s.visible()))
	case "k", "up":
		s.ListState.cursorUp(1)
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

func (s *WorktreeSection) match(sess model.SeshSession, q string) bool {
	number := strconv.Itoa(s.entries[sess.Path].Number)
	return strings.Contains(number, q) || strings.Contains(strings.ToLower(sess.Name), q)
}

func (s *WorktreeSection) applyFilter() {
	s.ListState.applyFilter(s.sessions, s.match)
}

func (s *WorktreeSection) visible() []model.SeshSession {
	return s.ListState.visible(s.sessions)
}

func (s *WorktreeSection) ClickAt(row int) {
	s.ListState.ClickAt(row, len(s.visible()))
}

func (s *WorktreeSection) selectItem() {
	visible := s.visible()
	if len(visible) == 0 {
		return
	}
	s.chosen = s.entries[visible[s.cursor].Path].Number
}

func (s *WorktreeSection) ViewBorderless(width, height int, focused bool) (string, string) {
	s.viewHeight = height
	title := s.Name()

	switch {
	case s.loading:
		return title, "  Loading worktrees..."
	case s.err != nil:
		return title, "  " + s.err.Error()
	case len(s.sessions) == 0:
		return title, "  No worktrees"
	}

	visible := s.visible()
	end := min(s.offset+max(height, 1), len(visible))
	branchCol := branchColumnWidth(visible)

	var b strings.Builder
	for i := s.offset; i < end; i++ {
		sess := visible[i]
		e := s.entries[sess.Path]
		b.WriteString(render.RenderWorktreeRowFocused(width, branchCol, i == s.cursor, focused, e.Number, sess.Name, e.State == "CLOSED", sess.Branch, sess.GitStatus))
		b.WriteString("\n")
	}
	return title, b.String()
}
