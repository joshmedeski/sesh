package dashboard

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

// wmStatus mirrors the JSON emitted by `workmux status --json --git`.
// Parsing is lenient: unknown fields are ignored.
type wmStatus struct {
	Agents []wmAgent `json:"agents"`
}

type wmAgent struct {
	Worktree    string  `json:"worktree"`
	Branch      string  `json:"branch"`
	Status      string  `json:"status"` // "working" | "waiting" | "done" | "-"
	ElapsedSecs *uint64 `json:"elapsed_secs"`
	Title       *string `json:"title"`
	Session     *string `json:"session"` // tmux session name (nil for window-mode agents)
	AgentKind   string  `json:"agent_kind"`
	Git         *wmGit  `json:"git"`
}

type wmGit struct {
	HasStaged          bool `json:"has_staged"`
	HasUnstaged        bool `json:"has_unstaged"`
	HasUnmergedCommits bool `json:"has_unmerged_commits"`
}

// workmuxLoadedMsg carries the result of a `workmux status` call. errMsg is
// non-empty when the command failed or the JSON could not be parsed.
type workmuxLoadedMsg struct {
	agents []wmAgent
	errMsg string
}

type WorkmuxSection struct {
	config     model.DashboardSectionConfig
	deps       SectionDeps
	agents     []wmAgent
	cursor     int
	offset     int
	loading    bool
	chosen     string
	errorMsg   string
	viewHeight int
}

func NewWorkmuxSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &WorkmuxSection{
		config:  cfg,
		deps:    deps,
		loading: true,
	}
}

func (s *WorkmuxSection) Name() string    { return s.config.Title }
func (s *WorkmuxSection) Width() float64  { return s.config.Width }
func (s *WorkmuxSection) Chosen() string  { return s.chosen }
func (s *WorkmuxSection) TotalItems() int { return len(s.agents) }

func (s *WorkmuxSection) Init() tea.Cmd {
	return s.fetch()
}

func (s *WorkmuxSection) fetch() tea.Cmd {
	return func() tea.Msg {
		out, err := runCommand("workmux", "status", "--json", "--git")
		if err != nil {
			return workmuxLoadedMsg{errMsg: fmt.Sprintf("workmux not found or errored: %v", err)}
		}
		agents, err := parseWorkmuxStatus(out)
		if err != nil {
			return workmuxLoadedMsg{errMsg: "Failed to parse workmux status"}
		}
		return workmuxLoadedMsg{agents: agents}
	}
}

// parseWorkmuxStatus decodes the `workmux status --json --git` output. Unknown
// fields are ignored by encoding/json.
func parseWorkmuxStatus(out string) ([]wmAgent, error) {
	var status wmStatus
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		return nil, err
	}
	return status.Agents, nil
}

func (s *WorkmuxSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case workmuxLoadedMsg:
		s.loading = false
		s.errorMsg = msg.errMsg
		s.agents = msg.agents
		s.clampCursor()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			s.cursorDown(1)
		case "k", "up":
			s.cursorUp(1)
		case "enter":
			s.selectItem()
		case "r":
			s.loading = true
			s.errorMsg = ""
			return s, s.fetch()
		}
	}
	return s, nil
}

func (s *WorkmuxSection) clampCursor() {
	n := len(s.agents)
	if s.cursor >= n {
		s.cursor = max(n-1, 0)
	}
	if s.offset >= n {
		s.offset = 0
	}
}

func (s *WorkmuxSection) cursorUp(n int) {
	s.cursor -= n
	if s.cursor < 0 {
		s.cursor = 0
	}
	if s.cursor < s.offset {
		s.offset = s.cursor
	}
}

func (s *WorkmuxSection) cursorDown(n int) {
	s.cursor += n
	if maxIdx := len(s.agents) - 1; s.cursor > maxIdx {
		if maxIdx < 0 {
			s.cursor = 0
		} else {
			s.cursor = maxIdx
		}
	}
	visible := s.visibleCount()
	if s.cursor >= s.offset+visible {
		s.offset = s.cursor - visible + 1
	}
}

func (s *WorkmuxSection) visibleCount() int {
	if s.viewHeight <= 0 {
		return 20
	}
	return max(s.viewHeight, 1)
}

// ClickAt moves the cursor to the clicked view row, scrolling to reveal it.
func (s *WorkmuxSection) ClickAt(row int) {
	n := len(s.agents)
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

// selectItem sets Chosen() to the agent's tmux session. Window-mode agents
// (nil Session) are a no-op.
func (s *WorkmuxSection) selectItem() {
	if len(s.agents) == 0 {
		return
	}
	a := s.agents[s.cursor]
	if a.Session != nil && *a.Session != "" {
		s.chosen = *a.Session
	}
}

func (s *WorkmuxSection) ViewBorderless(width, height int, focused bool) (string, string) {
	s.viewHeight = height

	title := s.config.Title
	if title == "" {
		title = "Workmux"
	}

	const minWidth = 20
	if width < minWidth {
		msg := fmt.Sprintf("  Enlarge pane to see agents (need ≥%d cols, have %d)", minWidth, width)
		return title, msg
	}

	if s.loading {
		return title, "  Loading..."
	}
	if s.errorMsg != "" {
		return title, "  " + s.errorMsg
	}
	if len(s.agents) == 0 {
		return title, "  No agents found"
	}

	available := max(height, 1)
	end := min(s.offset+available, len(s.agents))

	var b strings.Builder
	for i := s.offset; i < end; i++ {
		b.WriteString(renderWorkmuxRow(width, i == s.cursor, s.agents[i]))
		b.WriteString("\n")
	}

	return title, b.String()
}

// renderWorkmuxRow renders a single agent row with columns:
// marker(2) | state(2) | kind(10) | branch(10) | elapsed(6) | title(fill).
// Progressive degradation keeps branch+elapsed visible as long as possible:
// the title drops first (below 50 cols), then the kind column shrinks to
// absorb the width. Only below 27 cols do branch+elapsed drop, leaving
// state+kind.
func renderWorkmuxRow(width int, selected bool, a wmAgent) string {
	const (
		stateW, kindW, branchW, elapsedW = 2, 10, 15, 5
		minTitleW                        = 16
	)

	// state+branch+elapsed fit from 27 cols up:
	// marker(2) + state(2) + kind(≥4) + branch(10) + elapsed(6) + seps(3) = 27.
	if width >= 27 {
		// With a title (5 cols, 4 seps) the fixed part is
		// marker(2)+state(2)+kind(10)+branch(10)+elapsed(6) = 30, so the title
		// gets width-34. Below that the title drops and kind absorbs the slack
		// (4 cols, 3 seps → kind = width-23).
		// titleW := width - 34
		kind := kindW
		// var title string
		// if titleW >= minTitleW {
		// 	if a.Title != nil {
		// 		title = *a.Title
		// 	}
		// } else {
		// 	kind = min(max(width-23, 4), 24)
		// }

		cols := []col{
			{text: wmStateGlyph(a.Status), width: stateW},
			{text: truncateRight(a.AgentKind, 24), width: kind, style: textStyle()},
			{text: truncateRight(paren(a.Branch), 24), width: branchW, style: branchStyle()},
			{text: wmElapsed(a.ElapsedSecs), width: elapsedW, style: dimmedStyle(), align: lipgloss.Left},
		}
		// if title != "" {
		// 	cols = append(cols, col{text: truncateRight(title, titleW), width: titleW, style: dimmedStyle()})
		// }
		return renderRow(rowMarker(selected), cols, selected)
	}

	// Too narrow for branch+elapsed: state+kind only (2 cols, 1 sep).
	kind := max(width-5, 1)
	cols := []col{
		{text: wmStateGlyph(a.Status), width: stateW},
		{text: truncateRight(a.AgentKind, 24), width: kind, style: textStyle()},
	}
	return renderRow(rowMarker(selected), cols, selected)
}

// wmStateGlyph returns the styled state glyph for an agent status:
// working → ● yellow, waiting → ○ blue, done → ✓ green, unknown → dimmed -.
func wmStateGlyph(status string) string {
	switch status {
	case "working":
		return lipgloss.NewStyle().Render("🟠")
	case "waiting":
		return lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(6)).Render("⏸️")
	case "done":
		return lipgloss.NewStyle().Render("✅")
	default:
		return dimmedStyle().Render("-")
	}
}

// wmGitCell renders the git status flags joined without separators:
// + staged (green), ~ unstaged (yellow), u unmerged (red). Blank when git is
// nil or no flags are set.
func wmGitCell(g *wmGit) string {
	if g == nil {
		return ""
	}
	var b strings.Builder
	if g.HasStaged {
		b.WriteString(successStyle().Render("+"))
	}
	if g.HasUnstaged {
		b.WriteString(warningStyle().Render("~"))
	}
	if g.HasUnmergedCommits {
		b.WriteString(lipgloss.NewStyle().Foreground(colorDeleted).Render("u"))
	}
	return b.String()
}

// wmElapsed renders ElapsedSecs as a compact duration ("42s"/"2m"/"1h"/"3d"),
// or "" when nil.
func wmElapsed(secs *uint64) string {
	if secs == nil {
		return ""
	}
	s := *secs
	switch {
	case s < 60:
		return fmt.Sprintf("%ds", s)
	case s < 60*60:
		return fmt.Sprintf("%dm", s/60)
	case s < 24*60*60:
		return fmt.Sprintf("%dh", s/(60*60))
	default:
		return fmt.Sprintf("%dd", s/(24*60*60))
	}
}
