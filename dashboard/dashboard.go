package dashboard

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/tmux"
)

const (
	pageOpen       = 0
	pageConfigured = 1
)

// Model is the dashboard TUI. It has two permanent tabs (page 0 "Open" and
// page 1 "Configured"). Tab 1 lays panes out in two rows of shared frames:
// row 1 is the sessions list, row 2 is the remaining widgets side by side.
// Tab 2 is the single-pane configured list.
type Model struct {
	config     model.DashboardConfig
	sessions   *SessionsSection
	configured *ConfiguredSection
	widgets    []Section

	// page is the active tab: pageOpen or pageConfigured.
	page int
	// focus is the focused pane index on page 0, row-major over the flat pane
	// list [row1..., row2...]. 0 = sessions list. Ignored on page 1.
	focus int

	width    int
	height   int
	tooSmall bool
	chosen   string
	quit     bool

	contentHeight int
	row1Widths    []int
	row2Widths    []int
	row1Height    int
	row2Height    int

	lastHoveredSession string
}

func New(config model.DashboardConfig, tmux tmux.Tmux, lister lister.Lister, git git.Git, connector connector.Connector, sh shell.Shell, homeDir string) Model {
	deps := SectionDeps{
		Tmux:      tmux,
		Lister:    lister,
		Git:       git,
		Connector: connector,
		Shell:     sh,
		HomeDir:   homeDir,
	}

	built := BuildSections(config, deps)

	m := Model{
		config:     config,
		sessions:   built.Sessions,
		configured: built.Configured,
		widgets:    built.Widgets,
		page:       pageOpen,
		focus:      0,
		width:      80,
		height:     24,
	}
	return m.withLayout()
}

func (m Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.widgets)+2)
	cmds = append(cmds, m.sessions.Init(), m.configured.Init())
	for _, w := range m.widgets {
		cmds = append(cmds, w.Init())
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
// focused (page 0) and list sections move their selection to the clicked row.
// Clicks outside any pane are ignored.
func (m Model) handleMouseClick(msg tea.MouseClickMsg) (Model, tea.Cmd) {
	e := msg.Mouse()
	if e.Button != tea.MouseLeft {
		return m, nil
	}
	idx, sec, row, ok := m.hitTest(e.X, e.Y)
	if !ok {
		return m, nil
	}
	if m.page == pageOpen {
		m.focus = idx
	}
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
	if m.page == pageConfigured {
		if cy < 0 || cy >= m.contentHeight || x < 1 || x >= m.width {
			return 0, nil, 0, false
		}
		return 0, m.configured, cy - 1, true
	}

	panes, widths, base, top := m.rowAt(cy)
	if panes == nil {
		return 0, nil, 0, false
	}
	col := paneCol(x, widths)
	if col < 0 {
		return 0, nil, 0, false
	}
	return base + col, panes[col], cy - top - 1, true
}

// rowAt returns the panes, widths, flat base index, and top offset of the row
// containing content y (row 1 or row 2 on page 0), or nil when y is outside
// both rows.
func (m Model) rowAt(cy int) (panes []Section, widths []int, base, top int) {
	if cy >= 0 && cy < m.row1Height {
		return m.row1Panes(), m.row1Widths, 0, 0
	}
	if m.row2Height > 0 && cy >= m.row1Height && cy < m.row1Height+m.row2Height {
		return m.row2Panes(), m.row2Widths, len(m.row1Panes()), m.row1Height
	}
	return nil, nil, 0, 0
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

	s, c := m.sessions.Update(msg)
	m.sessions = s.(*SessionsSection)
	if c != nil {
		cmds = append(cmds, c)
	}

	cfg, c := m.configured.Update(msg)
	m.configured = cfg.(*ConfiguredSection)
	if c != nil {
		cmds = append(cmds, c)
	}

	for i := range m.widgets {
		var w Section
		w, c = m.widgets[i].Update(msg)
		m.widgets[i] = w
		if c != nil {
			cmds = append(cmds, c)
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

	switch msg.String() {
	case "q", "esc", "ctrl+c":
		m.quit = true
		return m, tea.Quit

	case "tab", "shift+tab":
		m.page = 1 - m.page
		m.focus = 0
		return m, nil
	}

	// Pane navigation (ctrl+h/l/j/k and the backspace alias). Matched by both
	// the string form ("ctrl+h") and the raw code+modifier the decoder emits,
	// so navigation keeps working regardless of how the terminal/decoder
	// reports the key.
	switch {
	case isLeftKey(msg):
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

// focusedSection returns the currently focused pane (sessions list, configured
// list, or a widget).
func (m Model) focusedSection() Section {
	if m.page == pageConfigured {
		return m.configured
	}
	wi := m.flatWidgetIndex(m.focus)
	if wi < 0 {
		return m.sessions
	}
	return m.widgets[wi]
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

// sortLabel returns the sessions list's current sort mode label (the sessions
// list is always a Sorter on page 0). Defaults to "name".
func (m Model) sortLabel() string {
	if m.sessions == nil {
		return "name"
	}
	if label := m.sessions.SortLabel(); label != "" {
		return label
	}
	return "name"
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

// isLeftKey reports whether msg means "move focus left": ctrl+h or the
// backspace alias (some terminals send the backspace key as ctrl+h / BS, and
// some decoders report ctrl+h as ctrl+backspace).
func isLeftKey(msg tea.KeyPressMsg) bool {
	return msg.Code == tea.KeyBackspace || isCtrlKey(msg, 'h')
}

// routeKey forwards a key to the focused pane and handles enter→chosen.
func (m Model) routeKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated Section

	if m.page == pageConfigured {
		updated, cmd = m.configured.Update(msg)
		m.configured = updated.(*ConfiguredSection)
	} else {
		wi := m.flatWidgetIndex(m.focus)
		if wi < 0 {
			updated, cmd = m.sessions.Update(msg)
			m.sessions = updated.(*SessionsSection)
		} else {
			updated, cmd = m.widgets[wi].Update(msg)
			m.widgets[wi] = updated
		}
	}

	if msg.String() == "enter" {
		if chosen := updated.Chosen(); chosen != "" {
			m.chosen = chosen
			return m, tea.Quit
		}
	}

	m, syncCmd := m.syncHoveredSession()
	return m, tea.Batch(cmd, syncCmd)
}

// moveFocus shifts the focused pane by delta across the flat row-major pane
// list (wrapping) on page 0. No-op on page 1 or with a single pane.
func (m Model) moveFocus(delta int) Model {
	if m.page != pageOpen {
		return m
	}
	n := len(m.row1Panes()) + len(m.row2Panes())
	if n <= 1 {
		m.focus = 0
		return m
	}
	m.focus += delta
	m.focus %= n
	if m.focus < 0 {
		m.focus += n
	}
	return m
}

// moveFocusRow moves focus between rows on page 0, keeping the column index
// (clamped to the target row's length). dir > 0 moves down (row 1 → row 2),
// dir < 0 moves up. No-op on page 1, when there is no second row, or when the
// target row does not exist.
func (m Model) moveFocusRow(dir int) Model {
	if m.page != pageOpen {
		return m
	}
	row1Len := len(m.row1Panes())
	row2Len := len(m.row2Panes())
	if row2Len == 0 {
		return m
	}

	if m.focus < row1Len {
		// Currently in row 1.
		if dir <= 0 {
			return m // no row above
		}
		m.focus = row1Len + min(m.focus, row2Len-1)
	} else {
		// Currently in row 2.
		if dir >= 0 {
			return m // no row below
		}
		m.focus = min(m.focus-row1Len, row1Len-1)
	}
	return m
}

// jumpFocus sets focus to the pane at the given 1-based flat position (1 =
// sessions) on page 0, lazygit-style. Out-of-range digits and the configured
// page are a no-op.
func (m Model) jumpFocus(digit int) Model {
	if m.page != pageOpen {
		return m
	}
	n := len(m.row1Panes()) + len(m.row2Panes())
	idx := digit - 1
	if idx >= 0 && idx < n {
		m.focus = idx
	}
	return m
}

// detailsIndex returns the index of the details widget in m.widgets, or -1.
func (m Model) detailsIndex() int {
	for i, w := range m.widgets {
		if _, ok := w.(*DetailsSection); ok {
			return i
		}
	}
	return -1
}

// row1Panes returns row 1: the sessions list, followed by the details widget
// (if configured).
func (m Model) row1Panes() []Section {
	panes := []Section{m.sessions}
	if di := m.detailsIndex(); di >= 0 {
		panes = append(panes, m.widgets[di])
	}
	return panes
}

// row2Order returns the widget indices for row 2 (all non-details widgets in
// config order).
func (m Model) row2Order() []int {
	order := make([]int, 0, len(m.widgets))
	for i, w := range m.widgets {
		if _, ok := w.(*DetailsSection); ok {
			continue
		}
		order = append(order, i)
	}
	return order
}

// row2Panes returns row 2: all non-details widgets in config order.
func (m Model) row2Panes() []Section {
	order := m.row2Order()
	panes := make([]Section, 0, len(order))
	for _, wi := range order {
		panes = append(panes, m.widgets[wi])
	}
	return panes
}

// flatWidgetIndex maps a flat row-major focus index on page 0 to a widget
// index, or -1 for the sessions list.
func (m Model) flatWidgetIndex(idx int) int {
	if idx <= 0 {
		return -1 // sessions
	}
	row1Len := len(m.row1Panes())
	if idx < row1Len {
		return m.detailsIndex() // idx == 1 → details
	}
	order := m.row2Order()
	i := idx - row1Len
	if i < 0 || i >= len(order) {
		return -1
	}
	return order[i]
}

// syncHoveredSession keeps the Details widget in sync with the hovered session
// in the sessions list. It only runs on page 0 (avoids background tmux queries
// from tab 2).
func (m Model) syncHoveredSession() (Model, tea.Cmd) {
	if m.page != pageOpen {
		return m, nil
	}

	dsIdx := -1
	for i, w := range m.widgets {
		if _, ok := w.(*DetailsSection); ok {
			dsIdx = i
			break
		}
	}
	if dsIdx < 0 {
		return m, nil
	}

	name, path, windows := m.sessions.HoveredSession()
	if name == m.lastHoveredSession {
		return m, nil
	}
	m.lastHoveredSession = name

	updated, cmd := m.widgets[dsIdx].Update(hoveredSessionMsg{Name: name, Path: path, Windows: windows})
	m.widgets[dsIdx] = updated
	return m, cmd
}

// withLayout recomputes the layout: content height (header 2 + footer 1), the
// per-row widths, and the row heights.
func (m Model) withLayout() Model {
	m.contentHeight = max(m.height-3, 1)
	m.row1Widths = m.computePaneWidths(m.row1Panes())
	m.row2Widths = m.computePaneWidths(m.row2Panes())
	m.row1Height, m.row2Height = m.computeRowHeights()
	return m
}

// computeRowHeights splits contentHeight between the two rows. Row 1 gets 3/5
// when a second row exists; otherwise row 1 takes the full height. Degenerate
// tiny terminals are clamped so both rows stay usable.
func (m Model) computeRowHeights() (int, int) {
	if len(m.row2Panes()) == 0 {
		return m.contentHeight, 0
	}
	row1 := m.contentHeight * 3 / 5
	row2 := m.contentHeight - row1
	if row2 < 4 {
		row1 = m.contentHeight - 4
		row2 = 4
	}
	if row1 < 1 {
		row1 = 1
	}
	if row2 < 1 {
		row2 = 1
	}
	return row1, row2
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

	header := renderHeader(m.page, m.sessions.totalSessions, m.width)
	filtering, query := m.focusedFilterState()
	footer := renderFooter(m.page, m.width, m.sortLabel(), filtering, query)

	var content string
	if m.page == pageOpen {
		content = m.viewOpenPage()
	} else {
		content = m.viewConfiguredPage()
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

func (m Model) viewOpenPage() string {
	row1 := m.row1Panes()
	row2 := m.row2Panes()

	row1Frame := m.renderRow(row1, m.row1Widths, m.row1Height, 0)
	if len(row2) == 0 {
		return row1Frame
	}

	row2Frame := m.renderRow(row2, m.row2Widths, m.row2Height, len(row1))
	return lipgloss.JoinVertical(lipgloss.Top, row1Frame, row2Frame)
}

// renderRow builds the shared frame for one row of panes. flatOffset is the
// flat focus index of the first pane in the row.
func (m Model) renderRow(panes []Section, widths []int, height int, flatOffset int) string {
	innerHeight := height - 2
	fp := make([]framePane, 0, len(panes))
	for i, s := range panes {
		width := m.width
		if i < len(widths) {
			width = widths[i]
		}
		focused := m.focus == flatOffset+i
		title, content := s.ViewBorderless(width, innerHeight, focused)
		title = fmt.Sprintf("%d %s", flatOffset+i+1, title)
		fp = append(fp, framePane{title: title, content: content, width: width, focused: focused})
	}
	return renderFrame(fp, height)
}

func (m Model) viewConfiguredPage() string {
	innerHeight := m.contentHeight - 2
	// The single pane sits between the frame's two corner columns, so reserve
	// them to keep the frame exactly m.width wide.
	paneWidth := max(m.width-2, 1)
	title, content := m.configured.ViewBorderless(paneWidth, innerHeight, true)
	return renderFrame([]framePane{
		{title: title, content: content, width: paneWidth, focused: true},
	}, m.contentHeight)
}

func (m Model) Chosen() string {
	return m.chosen
}

func (m Model) Quit() bool {
	return m.quit
}
