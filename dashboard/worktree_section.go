package dashboard

import (
	"sort"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/model"
)

type worktreesLoadedMsg struct {
	repo    string
	entries []model.WorktreeEntry
	err     error
}

var worktreeSortModes = []string{"issue", "age", "state", "status"}

var issueStateOrder = map[string]int{"OPEN": 0, "MERGED": 1, "CLOSED": 2}

// WorktreeSection lists the worktrees of one [[worktree]] config entry.
// Selecting a row sets ChosenWorktree.
type WorktreeSection struct {
	config   model.WorktreeConfig
	deps     SectionDeps
	sessions []model.SeshSession
	entries  map[string]model.WorktreeEntry
	ListState
	loading  bool
	err      error
	chosen   int
	sortMode string
	changes  map[string]int
	title    string
	columns  []string
}

func NewWorktreeSection(cfg model.WorktreeConfig, deps SectionDeps) *WorktreeSection {
	return &WorktreeSection{
		config:   cfg,
		deps:     deps,
		loading:  true,
		sortMode: worktreeSortModes[0],
		changes:  map[string]int{},
	}
}

// SortLabel implements Sorter.
func (s *WorktreeSection) SortLabel() string { return s.sortMode }

func (s *WorktreeSection) Width() float64 { return 0 }

func (s *WorktreeSection) Name() string {
	if s.title != "" {
		return s.title
	}
	return s.config.Repo
}

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
		s.applySort()
		return s, tea.Batch(fetchBranches(s.deps.Git, s.sessions), fetchStatuses(s.deps.Git, s.sessions))

	case branchLoadedMsg:
		applyBranch(s.sessions, msg.path, msg.branch)
		s.applyFilter()
		return s, nil

	case statusLoadedMsg:
		applyStatus(s.sessions, msg.path, msg.status)
		if _, ok := s.entries[msg.path]; ok {
			s.changes[msg.path] = msg.changes
		}
		if s.sortMode == "status" {
			s.applySort()
		} else {
			s.applyFilter()
		}
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
	case "s":
		s.cycleSortMode()
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

func (s *WorktreeSection) cycleSortMode() {
	i := 0
	for j, mode := range worktreeSortModes {
		if mode == s.sortMode {
			i = j
		}
	}
	s.sortMode = worktreeSortModes[(i+1)%len(worktreeSortModes)]
	s.applySort()
}

func (s *WorktreeSection) applySort() {
	number := func(i int) int { return s.entries[s.sessions[i].Path].Number }
	sort.SliceStable(s.sessions, func(i, j int) bool {
		pi, pj := s.sessions[i].Path, s.sessions[j].Path
		switch s.sortMode {
		case "age":
			if ti, tj := s.entries[pi].Created, s.entries[pj].Created; !ti.Equal(tj) {
				return ti.After(tj)
			}
		case "state":
			if oi, oj := stateRank(s.entries[pi].State), stateRank(s.entries[pj].State); oi != oj {
				return oi < oj
			}
		case "status":
			if ci, cj := s.changes[pi], s.changes[pj]; ci != cj {
				return ci > cj
			}
		}
		return number(i) < number(j)
	})
	s.applyFilter()
}

func stateRank(state string) int {
	if rank, ok := issueStateOrder[state]; ok {
		return rank
	}
	return len(issueStateOrder)
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

// iconCol uses the worktree's own icon, else the icon of its repo root, else
// the config glyph.
func (s *WorktreeSection) iconCol(sess model.SeshSession, selected bool) render.Col {
	if s.deps.Icon != nil && s.deps.Icon(sess) == "" {
		sess = model.SeshSession{Path: expandHome(s.config.Path, s.deps.HomeDir)}
	}
	return iconCol(s.deps, sess, "config", selected)
}

func expandHome(path, homeDir string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		return homeDir + path[1:]
	}
	return path
}

func (s *WorktreeSection) ViewBorderless(width, height int, focused bool) (string, string) {
	s.viewHeight = max(height-1, 1)
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
	end := min(s.offset+s.viewHeight, len(visible))
	branchCol := branchColumnWidth(visible)

	var b strings.Builder
	b.WriteString(render.RenderWorktreeHeader(width, s.columns, branchCol, max(s.deps.IconWidth, 1)))
	b.WriteString("\n")
	for i := s.offset; i < end; i++ {
		sess := visible[i]
		e := s.entries[sess.Path]
		var created *time.Time
		if !e.Created.IsZero() {
			created = &e.Created
		}
		b.WriteString(render.RenderWorktreeRowFocused(width, s.columns, branchCol, s.iconCol(sess, i == s.cursor), i == s.cursor, focused, e.Number, sess.Name, e.State, sess.Branch, sess.GitStatus, created))
		b.WriteString("\n")
	}
	return title, b.String()
}
