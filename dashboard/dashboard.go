package dashboard

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/dashboard/core"
	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/model"
)

// dashPage is one dashboard tab: rows of panes, each row laid out side by side.
type dashPage struct {
	title string
	rows  [][]Section
}

// Model is the dashboard TUI. Its tabs are one per [[dashboard.page]]; a page
// stacks its rows of shared frames vertically.
type Model struct {
	config model.DashboardConfig
	pages  []dashPage

	// page is the active tab's index into pages.
	page int
	// focus is the focused pane index on the active page, row-major over its
	// rows.
	focus int

	width    int
	height   int
	tooSmall bool
	chosen   string
	quit     bool

	chosenWorktree *model.WorktreeConnectOpts

	contentHeight int
	rowWidths     [][]int
	rowHeights    []int

	lastHoveredSession string

	showHelp bool
}

func New(config model.Config, deps SectionDeps) Model {
	built := BuildPages(config.Dashboard, config.WorktreeConfigs, deps)

	m := Model{
		config: config.Dashboard,
		pages:  built.pages(),
		width:  80,
		height: 24,
	}
	return m.withLayout()
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, sec := range m.dashboardSections() {
		cmds = append(cmds, sec.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.width < 20 || m.height < 5 {
			m.tooSmall = true
			return m, tea.Quit
		}
		m.tooSmall = false
		m = m.withLayout()
		return m.broadcast(msg)

	case tea.KeyPressMsg:
		return m.handleKey(msg)

	case tea.MouseClickMsg:
		return m.handleMouseClick(msg)

	default:
		return m.broadcast(msg)
	}
}

// handleMouseClick maps a left click to the pane under the cursor: the pane is
// focused (dashboard pages) and list sections move their selection to the
// clicked row. Clicks outside any pane are ignored.
func (m Model) handleMouseClick(msg tea.MouseClickMsg) (Model, tea.Cmd) {
	e := msg.Mouse()
	if m.showHelp || e.Button != tea.MouseLeft {
		return m, nil
	}
	idx, sec, row, ok := m.hitTest(e.X, e.Y)
	if !ok {
		return m, nil
	}
	m.focus = idx
	if c, ok := sec.(Clicker); ok {
		c.ClickAt(row)
	}
	return m, nil
}

// hitTest maps terminal cell coordinates to the pane under the cursor,
// returning the flat focus index, the pane's section, and the clicked list
// row (0-based, view-relative). ok is false when the click lands in the
// header, footer, or frame chrome.
func (m Model) hitTest(x, y int) (idx int, sec Section, row int, ok bool) {
	// Content begins below the two-row header.
	cy := y - 2
	top, base := 0, 0
	for r, panes := range m.rows() {
		if r >= len(m.rowHeights) {
			break
		}
		if cy >= top && cy < top+m.rowHeights[r] {
			col := paneCol(x, m.rowWidths[r])
			if col < 0 {
				return 0, nil, 0, false
			}
			return base + col, panes[col], cy - top - 1, true
		}
		top += m.rowHeights[r]
		base += len(panes)
	}
	return 0, nil, 0, false
}

// paneCol returns the index of the pane containing column x, accounting for
// the leading corner and one junction column per pane. -1 when x is chrome or
// out of range.
func paneCol(x int, widths []int) int {
	cursor := 1
	for i, w := range widths {
		if x >= cursor && x < cursor+w {
			return i
		}
		cursor += w + 1
	}
	return -1
}

// broadcast forwards a non-key message to every section (each section ignores
// messages it does not recognize), then refreshes the details widget if the
// hovered session changed.
func (m Model) broadcast(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	for p := range m.pages {
		for r := range m.pages[p].rows {
			for i := range m.pages[p].rows[r] {
				w, c := m.pages[p].rows[r][i].Update(msg)
				m.pages[p].rows[r][i] = w
				if c != nil {
					cmds = append(cmds, c)
				}
			}
		}
	}

	m, syncCmd := m.syncHoveredSession()
	if syncCmd != nil {
		cmds = append(cmds, syncCmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// While the focused pane is filtering, route every key to it (except ctrl+c
	// which always quits) so quit/switch/nav keys don't fire mid-typing.
	if m.focusedPaneFiltering() {
		if msg.String() == "ctrl+c" {
			m.quit = true
			return m, tea.Quit
		}
		return m.routeKey(msg)
	}

	if m.showHelp {
		switch msg.String() {
		case "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "?", "esc", "q":
			m.showHelp = false
		}
		return m, nil
	}

	switch msg.String() {
	case "?":
		m.showHelp = true
		return m, nil
	case "q", "esc", "ctrl+c":
		m.quit = true
		return m, tea.Quit

	case "tab":
		m.page = (m.page + 1) % len(m.pages)
		m.focus = 0
		return m.withLayout(), nil
	case "shift+tab":
		m.page = (m.page - 1 + len(m.pages)) % len(m.pages)
		m.focus = 0
		return m.withLayout(), nil
	}

	// Pane navigation (ctrl+h/l/j/k and the backspace alias). Matched by both
	// the string form ("ctrl+h") and the raw code+modifier the decoder emits,
	// so navigation keeps working regardless of how the terminal/decoder
	// reports the key.
	switch {
	case isBackspaceKey(msg):
		m = m.moveFocus(-1)
		return m, nil
	case isCtrlKey(msg, 'l'):
		m = m.moveFocus(1)
		return m, nil
	case isCtrlKey(msg, 'j'):
		m = m.moveFocusRow(1) // down a row
		return m, nil
	case isCtrlKey(msg, 'k'):
		m = m.moveFocusRow(-1) // up a row
		return m, nil

	case isDigitKey(msg):
		m = m.jumpFocus(int(msg.String()[0] - '0'))
		return m, nil
	}

	return m.routeKey(msg)
}

// focusedSection returns the focused pane on the active page (nil when there
// is none).
func (m Model) focusedSection() Section {
	r, c, ok := m.paneAt(m.focus)
	if !ok {
		return nil
	}
	return m.rows()[r][c]
}

// focusedPaneFiltering reports whether the focused pane is in filter mode.
func (m Model) focusedPaneFiltering() bool {
	if f, ok := m.focusedSection().(Filterer); ok {
		return f.Filtering()
	}
	return false
}

// focusedFilterState returns the filtering state and query of the focused pane.
func (m Model) focusedFilterState() (filtering bool, query string) {
	if f, ok := m.focusedSection().(Filterer); ok {
		return f.Filtering(), f.FilterQuery()
	}
	return false, ""
}

// sortLabel returns the focused pane's sort mode label, or "" when the
// focused pane cannot be sorted.
func (m Model) sortLabel() string {
	if sorter, ok := m.focusedSection().(Sorter); ok {
		return sorter.SortLabel()
	}
	return ""
}

// isCtrlKey reports whether msg is ctrl+<letter>. It matches both the string
// form ("ctrl+h") and the raw Code+Modifier form the decoder emits, so it
// survives decoder variations in how ctrl keys are reported.
func isCtrlKey(msg tea.KeyPressMsg, letter rune) bool {
	return msg.String() == "ctrl+"+string(letter) ||
		(msg.Mod.Contains(tea.ModCtrl) && msg.Code == letter)
}

// isDigitKey reports whether msg is a number key "1".."9". Digits arrive with
// the Text form set (no Code required), so we match on the string form only.
func isDigitKey(msg tea.KeyPressMsg) bool {
	s := msg.String()
	return len(s) == 1 && s[0] >= '1' && s[0] <= '9'
}

// isBackspaceKey reports whether msg is any backspace form: the plain
// backspace key, ctrl+backspace (some decoders report the combo as a modified
// backspace), or ctrl+h (terminals send ctrl+h as BS 0x08). In handleKey it
// moves focus left; in the filter handlers it deletes the last query rune.
// Matching is Code-based so it survives decoder variations in the String()
// form.
func isBackspaceKey(msg tea.KeyPressMsg) bool {
	return msg.Code == tea.KeyBackspace || isCtrlKey(msg, 'h')
}

// routeKey forwards a key to the focused pane and handles enter→chosen.
func (m Model) routeKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated Section

	r, c, ok := m.paneAt(m.focus)
	if !ok {
		return m, nil
	}
	updated, cmd = m.pages[m.page].rows[r][c].Update(msg)
	m.pages[m.page].rows[r][c] = updated

	if msg.String() == "enter" {
		if chosen := updated.Chosen(); chosen != "" {
			m.chosen = chosen
			return m, tea.Quit
		}
		if ws, ok := updated.(*WorktreeSection); ok {
			if opts, ok := ws.ChosenWorktree(); ok {
				m.chosenWorktree = &opts
				return m, tea.Quit
			}
		}
	}

	m, syncCmd := m.syncHoveredSession()
	return m, tea.Batch(cmd, syncCmd)
}

// moveFocus shifts the focused pane by delta across the flat row-major pane
// list (wrapping) on the active page. No-op with a single pane.
func (m Model) moveFocus(delta int) Model {
	n := m.paneCount()
	if n <= 1 {
		m.focus = 0
		return m
	}
	m.focus = ((m.focus+delta)%n + n) % n
	return m
}

// moveFocusRow moves focus to the row above (dir < 0) or below (dir > 0) on a
// dashboard page, keeping the column index clamped to the target row's
// length. No-op when there is no such row.
func (m Model) moveFocusRow(dir int) Model {
	r, c, ok := m.paneAt(m.focus)
	if !ok {
		return m
	}
	rows := m.rows()
	target := r + dir
	if dir == 0 || target < 0 || target >= len(rows) {
		return m
	}
	m.focus = m.flatIndex(target, min(c, len(rows[target])-1))
	return m
}

// jumpFocus sets focus to the pane at the given 1-based flat position on a
// active page, lazygit-style. Out-of-range digits are a no-op.
func (m Model) jumpFocus(digit int) Model {
	if idx := digit - 1; idx >= 0 && idx < m.paneCount() {
		m.focus = idx
	}
	return m
}

// rows returns the active page's rows, or nil when there is no page.
func (m Model) rows() [][]Section {
	if m.page >= len(m.pages) {
		return nil
	}
	return m.pages[m.page].rows
}

func (m Model) paneCount() int {
	n := 0
	for _, row := range m.rows() {
		n += len(row)
	}
	return n
}

// paneAt maps a flat focus index to its row and column on the active page.
func (m Model) paneAt(idx int) (row, col int, ok bool) {
	if idx < 0 {
		return 0, 0, false
	}
	for r, panes := range m.rows() {
		if idx < len(panes) {
			return r, idx, true
		}
		idx -= len(panes)
	}
	return 0, 0, false
}

func (m Model) flatIndex(row, col int) int {
	idx := col
	for _, panes := range m.rows()[:row] {
		idx += len(panes)
	}
	return idx
}

// dashboardSections returns every pane on every dashboard page.
func (m Model) dashboardSections() []Section {
	var out []Section
	for _, p := range m.pages {
		for _, row := range p.rows {
			out = append(out, row...)
		}
	}
	return out
}

func (m Model) pageSections() []Section {
	var out []Section
	for _, row := range m.rows() {
		out = append(out, row...)
	}
	return out
}

// syncHoveredSession keeps a Details widget in sync with the hovered session
// in the first sessions list on the active page.
func (m Model) syncHoveredSession() (Model, tea.Cmd) {
	sessions := firstSessions(m.pageSections())
	if sessions == nil {
		return m, nil
	}
	for r, row := range m.rows() {
		for c, sec := range row {
			if _, ok := sec.(interface{ LayoutRow() int }); !ok {
				continue
			}
			name, path, windows := sessions.HoveredSession()
			if name == m.lastHoveredSession {
				return m, nil
			}
			m.lastHoveredSession = name
			updated, cmd := sec.Update(core.HoveredSessionMsg{Name: name, Path: path, Windows: windows})
			m.pages[m.page].rows[r][c] = updated
			return m, cmd
		}
	}
	return m, nil
}

func firstSessions(sections []Section) *SessionsSection {
	for _, sec := range sections {
		if s, ok := sec.(*SessionsSection); ok {
			return s
		}
	}
	return nil
}

// activeCount is the tmux session count shown in the header, or -1 to hide it
// when no dashboard page has a sessions list.
func (m Model) activeCount() int {
	s := firstSessions(m.dashboardSections())
	if s == nil {
		return -1
	}
	return s.TmuxCount()
}

// withLayout recomputes the layout for the active page: content height
// (header 2 + footer 1), the per-row widths, and the row heights.
func (m Model) withLayout() Model {
	m.contentHeight = max(m.height-3, 1)
	rows := m.rows()
	m.rowWidths = make([][]int, len(rows))
	for r, panes := range rows {
		m.rowWidths[r] = m.computePaneWidths(panes)
	}
	m.rowHeights = splitHeights(m.contentHeight, len(rows))
	return m
}

// splitHeights divides total evenly between n rows, giving any remainder to
// the first rows. Every row gets at least one line.
func splitHeights(total, n int) []int {
	heights := make([]int, n)
	for i := range heights {
		heights[i] = max(total/n, 1)
		if i < total%n {
			heights[i]++
		}
	}
	return heights
}

// computePaneWidths allocates column widths to a single row of panes. Panes
// with a positive Width() fraction get a proportional share (scaled so the
// total of all fractions is <= 1); flex panes (Width() 0, e.g. the sessions
// list) share the remainder equally, with leftover distributed round-robin
// from the right. If all panes are fixed, the leftover goes to the last pane.
func (m Model) computePaneWidths(panes []Section) []int {
	n := len(panes)
	pw := make([]int, n)
	if n == 0 {
		return pw
	}

	// The shared frame consumes (n-1) junction characters (┬/│/┴) between panes
	// plus 2 corner characters (┌┐ on top, └┘ on bottom), i.e. n+1 columns of
	// chrome. Subtract all of it so the frame is exactly m.width wide.
	availableWidth := max(m.width-(n-1)-2, n)

	flex := make([]bool, n)
	flexCount := 0
	totalFraction := 0.0
	for i, p := range panes {
		if w := p.Width(); w > 0 {
			totalFraction += w
		} else {
			flex[i] = true
			flexCount++
		}
	}

	scale := 1.0
	if totalFraction > 1.0 {
		scale = 1.0 / totalFraction
	}

	allocated := 0
	for i, p := range panes {
		if w := p.Width(); w > 0 {
			pw[i] = max(int(float64(availableWidth)*w*scale), 1)
			allocated += pw[i]
		}
	}

	remaining := availableWidth - allocated
	if flexCount > 0 {
		each := remaining / flexCount
		for i := range pw {
			if flex[i] {
				pw[i] = each
				remaining -= each
			}
		}
		for i := n - 1; i >= 0 && remaining > 0; i-- {
			if flex[i] {
				pw[i]++
				remaining--
			}
		}
	} else if remaining > 0 {
		pw[n-1] += remaining
	}

	for i := range pw {
		if pw[i] < 1 {
			pw[i] = 1
		}
	}
	return pw
}

func (m Model) View() tea.View {
	if m.quit {
		return tea.NewView("")
	}
	if m.tooSmall {
		return tea.NewView("Terminal too small for dashboard")
	}

	header := render.RenderHeader(m.page, m.activeCount(), m.width, m.tabTitles()...)
	filtering, query := m.focusedFilterState()
	footer := render.RenderFooter(m.width, m.sortLabel(), filtering, query)

	var content string
	switch {
	case m.showHelp:
		content = render.RenderHelp(m.width, m.contentHeight, m.sortLabel() != "")
	default:
		content = m.viewDashboardPage()
	}

	ui := lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
	finalString := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(ui)

	v := tea.NewView(finalString)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m Model) viewDashboardPage() string {
	frames := make([]string, 0, len(m.rows()))
	base := 0
	for r, panes := range m.rows() {
		frames = append(frames, m.renderRow(panes, m.rowWidths[r], m.rowHeights[r], base))
		base += len(panes)
	}
	return lipgloss.JoinVertical(lipgloss.Top, frames...)
}

// renderRow builds the shared frame for one row of panes. flatOffset is the
// flat focus index of the first pane in the row.
func (m Model) renderRow(panes []Section, widths []int, height int, flatOffset int) string {
	innerHeight := height - 2
	fp := make([]render.FramePane, 0, len(panes))
	for i, s := range panes {
		width := m.width
		if i < len(widths) {
			width = widths[i]
		}
		focused := m.focus == flatOffset+i
		title, content := s.ViewBorderless(width, innerHeight, focused)
		title = fmt.Sprintf("%d %s", flatOffset+i+1, title)
		fp = append(fp, render.FramePane{Title: title, Content: content, Width: width, Focused: focused})
	}
	return render.RenderFrame(fp, height)
}

func (m Model) Chosen() string {
	return m.chosen
}

// ChosenWorktree returns the worktree selected in a worktree section, or nil.
func (m Model) ChosenWorktree() *model.WorktreeConnectOpts {
	return m.chosenWorktree
}

func (m Model) tabTitles() []string {
	titles := make([]string, 0, len(m.pages))
	for _, p := range m.pages {
		titles = append(titles, p.title)
	}
	return titles
}

func (m Model) Quit() bool {
	return m.quit
}
