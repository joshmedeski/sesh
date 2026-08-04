package picker

import (
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sahilm/fuzzy"

	"github.com/joshmedeski/sesh/v2/model"
)

// WorktreeFetchFunc loads the worktree entries asynchronously. It is called in a
// goroutine by Init(), because resolving titles can reach the network on the
// numbers whose cache entry has expired. refresh asks for every title to be
// refetched rather than read from the cache.
type WorktreeFetchFunc func(refresh bool) ([]model.WorktreeEntry, error)

// worktreesLoadedMsg carries the result of the async fetch back to Update().
type worktreesLoadedMsg struct {
	entries []model.WorktreeEntry
	err     error
	// refreshed marks the result of a ctrl+r refetch rather than the initial
	// load, which is rendered differently: the list stayed on screen throughout,
	// so the cursor is kept rather than reset.
	refreshed bool
}

// WorktreeOptions configures the worktree picker model. Zero values are valid
// and mean "unset", except Prompt and Placeholder which are passed through
// as-is.
type WorktreeOptions struct {
	Prompt      string
	Placeholder string
	// ShowIcons declares that the terminal's font has the nerd font glyphs, so
	// the state badges are rounded off into pills rather than squared off with
	// filled spaces.
	ShowIcons bool
	// Query pre-fills the filter input, as if it had just been typed, so it can
	// be backspaced away.
	Query string
	// AutoRefresh refetches every title once the cached rows are on screen, so
	// an issue renamed since the last listing shows its new title without the
	// user having to know that ctrl+r exists. It is what makes the long issue
	// cache TTL tolerable: the cache is what makes the picker open instantly,
	// and this is what keeps it from being wrong for a day afterwards.
	AutoRefresh bool
}

// worktreeItem is one row of the picker: a worktree directory, plus the strings
// derived from it that filtering and rendering need.
type worktreeItem struct {
	entry  model.WorktreeEntry
	number string // the issue number as displayed
	// searchText is "<number> <title>", so a query matches on either. It is the
	// one string handed to the fuzzy matcher, which keeps a query spanning both
	// — "409 preview" — working as a single match rather than needing the
	// columns to be searched separately.
	searchText string
}

// worktreeItems implements fuzzy.Source for fuzzy matching.
type worktreeItems []worktreeItem

func (w worktreeItems) String(i int) string { return w[i].searchText }
func (w worktreeItems) Len() int            { return len(w) }

type filteredWorktree struct {
	item worktreeItem
	// numberIndexes and titleIndexes are the matched rune positions split into
	// the two columns they are rendered in, since the badge sits between them.
	numberIndexes []int
	titleIndexes  []int
}

// Issue and PR states as GitHub's GraphQL API spells them. A PR also reports
// MERGED; an issue only ever reports the first two.
const (
	stateOpen   = "OPEN"
	stateClosed = "CLOSED"
	stateMerged = "MERGED"
)

// badgeCellWidth is the width the badge column is padded to: the longest label,
// the two cells rounding it off either side, and one plain space after it so the
// pill never butts up against the title. Padding the column rather than the pill
// keeps every pill hugging its own text, so the badges read as one column of
// tags rather than bars of differing length.
const badgeCellWidth = len("MERGED") + badgeCapWidth*2 + 1

// badgeCapWidth is the display width of each end cap, a half circle or the space
// standing in for one. Taken as a constant rather than measured so a font that
// reports the glyph as double width cannot pull the column out of alignment —
// the same single cell the session picker's alias chip assumes.
const badgeCapWidth = 1

// badgeColor returns the ANSI color a state's badge is filled with, and whether
// the state is one that gets a badge at all. The coding follows GitHub's own:
// green for open, purple for merged, red for closed.
func badgeColor(state string) (lipgloss.ANSIColor, bool) {
	switch strings.ToUpper(state) {
	case stateOpen:
		return lipgloss.ANSIColor(2), true
	case stateMerged:
		return lipgloss.ANSIColor(5), true
	case stateClosed:
		return lipgloss.ANSIColor(1), true
	default:
		// Includes the empty state of a worktree whose issue never resolved.
		// Its title is missing too, so the row is left deliberately quiet
		// rather than badged as something the fetch never confirmed.
		return lipgloss.ANSIColor(0), false
	}
}

// stateBadge renders a state as a pill, e.g. `OPEN` rounded off with a half
// circle either side, padded out to the badge column. An unrecognized or absent
// state renders as blanks of the same width, so the titles stay aligned.
//
// The label is a foreground color under reverse video rather than an explicit
// fg/bg pair: that fills it with the state color and paints the text in the
// terminal's own background color, which stays legible under any color scheme.
// The half circles take the same color as a plain foreground, which is exactly
// the color the reversed label is filled with, so the whole thing reads as one
// pill. It is the alias chip's reasoning, with the fill color made explicit
// because a state has a color to say rather than just a shape.
func stateBadge(state string, showIcons bool) string {
	color, ok := badgeColor(state)
	if !ok {
		return padding(badgeCellWidth)
	}

	label := strings.ToUpper(state)
	fill := lipgloss.NewStyle().Foreground(color)
	filled := fill.Reverse(true)

	var pill string
	if showIcons {
		pill = fill.Render(chipLeftGlyph) + filled.Render(label) + fill.Render(chipRightGlyph)
	} else {
		// Without nerd fonts the half circles render as tofu, so the pill is
		// squared off with filled spaces instead. Both forms are the same width,
		// so the column lines up either way.
		pill = filled.Render(" " + label + " ")
	}
	return pill + padding(badgeCellWidth-len(label)-badgeCapWidth*2)
}

// buildWorktreeItems derives the filterable rows from the fetched entries,
// preserving the order they were listed in — by number, ascending.
func buildWorktreeItems(entries []model.WorktreeEntry) worktreeItems {
	items := make(worktreeItems, 0, len(entries))
	for _, entry := range entries {
		number := strconv.Itoa(entry.Number)
		searchText := number
		if entry.Title != "" {
			searchText += " " + entry.Title
		}
		items = append(items, worktreeItem{
			entry:      entry,
			number:     number,
			searchText: searchText,
		})
	}
	return items
}

// WorktreeModel is the picker for worktrees. It is a separate model from the
// session picker rather than a mode of it: the rows are columns of typed fields
// instead of a single name, and none of the session picker's aliases, icons or
// previews have a meaning here.
type WorktreeModel struct {
	allItems    worktreeItems
	filtered    []filteredWorktree
	filterInput textinput.Model
	cursor      int
	offset      int
	width       int
	height      int
	chosen      model.WorktreeEntry
	picked      bool
	quit        bool
	focusCmd    tea.Cmd
	loading     bool
	fetchFunc   WorktreeFetchFunc
	loadErr     error
	// refreshing is a ctrl+r refetch in flight. It is kept apart from loading
	// because the list already on screen stays usable while it runs.
	refreshing bool
	// refreshErr is the last failed refresh. Unlike a failed initial load it
	// does not end the picker — the rows fetched before it are still good — so
	// it is reported in the header and cleared by the next refresh.
	refreshErr error
	// autoRefresh is the option: refetch once after the initial load lands.
	autoRefresh bool
	// refreshIsAuto marks the in-flight refresh as the automatic one rather than
	// a ctrl+r, which only changes how its failure is treated. One bool is
	// enough because refreshing already admits a single refresh at a time.
	refreshIsAuto bool
	showIcons     bool
	// numberWidth is the width of the widest number, which the column is
	// right-aligned to so the numbers read as a column.
	numberWidth int
}

func NewWorktreeModel(fetchFunc WorktreeFetchFunc, opts WorktreeOptions) WorktreeModel {
	ti := textinput.New()
	ti.Placeholder = opts.Placeholder
	ti.Prompt = opts.Prompt
	// SetValue leaves the cursor at the end, so the pre-filled query reads as
	// something just typed.
	ti.SetValue(opts.Query)

	m := WorktreeModel{
		filterInput: ti,
		loading:     true,
		fetchFunc:   fetchFunc,
		showIcons:   opts.ShowIcons,
		autoRefresh: opts.AutoRefresh,
	}
	m.focusCmd = m.filterInput.Focus()
	return m
}

func (m WorktreeModel) Init() tea.Cmd {
	return tea.Batch(m.focusCmd, m.fetchWorktrees(false))
}

func (m WorktreeModel) fetchWorktrees(refresh bool) tea.Cmd {
	fetchFunc := m.fetchFunc
	return func() tea.Msg {
		entries, err := fetchFunc(refresh)
		return worktreesLoadedMsg{entries: entries, err: err, refreshed: refresh}
	}
}

func (m *WorktreeModel) syncInputWidth() {
	m.filterInput.SetWidth(m.width - 4)
}

func (m WorktreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case worktreesLoadedMsg:
		wasAuto := m.refreshIsAuto
		if msg.refreshed {
			m.refreshing = false
			m.refreshIsAuto = false
		}
		if msg.err != nil {
			// A refresh that failed leaves the list it was refreshing intact,
			// so it is reported rather than fatal.
			if msg.refreshed {
				// Except when nobody asked for it: an automatic refresh that
				// cannot reach GitHub is not news, and reporting it would put an
				// error on screen every time the picker opens offline. The
				// cached titles it was trying to improve are still on screen.
				if !wasAuto {
					m.refreshErr = msg.err
				}
				return m, nil
			}
			m.loadErr = msg.err
			return m, tea.Quit
		}
		m.loading = false
		m.allItems = buildWorktreeItems(msg.entries)
		m.numberWidth = widestNumber(m.allItems)
		m.applyFilter()
		// A refresh can drop rows out from under the cursor — a worktree
		// removed since the load — so it is pulled back onto the list.
		m.clampCursor()

		// With the cached rows on screen, go get the current titles. Deferred to
		// here rather than batched into Init so the two fetches cannot race to
		// decide what is rendered: this one starts from a list already drawn and
		// replaces it, which is the same path ctrl+r takes.
		if m.autoRefresh && !msg.refreshed && len(m.allItems) > 0 {
			m.autoRefresh = false // once per picker, not once per fetch
			m.refreshing = true
			m.refreshIsAuto = true
			return m, m.fetchWorktrees(true)
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.syncInputWidth()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if m.loading {
				return m, nil
			}
			if m.cursor < len(m.filtered) {
				m.chosen = m.filtered[m.cursor].item.entry
				m.picked = true
			}
			return m, tea.Quit

		case "esc", "ctrl+c":
			m.quit = true
			return m, tea.Quit

		case "up", "ctrl+k", "ctrl+p":
			m.cursorUp(1)
			return m, nil

		case "down", "ctrl+j", "ctrl+n":
			m.cursorDown(1)
			return m, nil

		case "ctrl+u":
			m.cursorUp(m.visibleCount() / 2)
			return m, nil

		case "ctrl+d":
			m.cursorDown(m.visibleCount() / 2)
			return m, nil

		case "ctrl+r":
			// Refetch every title and state, ignoring the cache. The rows
			// already on screen stay put and are replaced when it lands, so a
			// refresh over a slow network is not a blank list.
			if m.loading || m.refreshing {
				return m, nil
			}
			m.refreshing = true
			m.refreshIsAuto = false
			m.refreshErr = nil
			return m, m.fetchWorktrees(true)
		}
	}

	// Forward to text input
	prevValue := m.filterInput.Value()
	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)

	if m.filterInput.Value() != prevValue {
		if !m.loading {
			m.applyFilter()
		}
		m.cursor = 0
		m.offset = 0
	}

	return m, cmd
}

// widestNumber is the display width of the largest number in the list, used to
// right-align the number column.
func widestNumber(items worktreeItems) int {
	width := 0
	for _, item := range items {
		width = max(width, len(item.number))
	}
	return width
}

// applyFilter narrows the list by the query: everything when it is empty, a
// fuzzy match over "<number> <title>" otherwise.
func (m *WorktreeModel) applyFilter() {
	pattern := m.filterInput.Value()
	if pattern == "" {
		filtered := make([]filteredWorktree, len(m.allItems))
		for i, item := range m.allItems {
			filtered[i] = filteredWorktree{item: item}
		}
		m.filtered = filtered
		return
	}

	matches := rankMatches(fuzzy.FindFromNoSort(pattern, m.allItems))
	filtered := make([]filteredWorktree, len(matches))
	for i, match := range matches {
		item := m.allItems[match.Index]
		numberIndexes, titleIndexes := splitMatchIndexes(match.MatchedIndexes, len([]rune(item.number)))
		filtered[i] = filteredWorktree{
			item:          item,
			numberIndexes: numberIndexes,
			titleIndexes:  titleIndexes,
		}
	}
	m.filtered = filtered
}

// splitMatchIndexes divides the rune positions matched in "<number> <title>"
// into positions within the number and within the title, so each can be
// highlighted in the column it is rendered in. The separating space is at
// numberLen and belongs to neither.
func splitMatchIndexes(indexes []int, numberLen int) (numberIndexes, titleIndexes []int) {
	for _, idx := range indexes {
		switch {
		case idx < numberLen:
			numberIndexes = append(numberIndexes, idx)
		case idx > numberLen:
			titleIndexes = append(titleIndexes, idx-numberLen-1)
		}
	}
	return numberIndexes, titleIndexes
}

// clampCursor pulls the cursor and the scroll offset back onto a list that has
// shrunk under them.
func (m *WorktreeModel) clampCursor() {
	if last := len(m.filtered) - 1; m.cursor > last {
		m.cursor = max(last, 0)
	}
	if m.offset > m.cursor {
		m.offset = m.cursor
	}
}

func (m *WorktreeModel) cursorUp(n int) {
	m.cursor -= n
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
}

func (m *WorktreeModel) cursorDown(n int) {
	m.cursor += n
	if last := len(m.filtered) - 1; m.cursor > last {
		m.cursor = max(last, 0)
	}
	visible := m.visibleCount()
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
}

// visibleCount is how many worktree rows fit. The picker runs on the alt
// screen, so it gets every row apart from the filter row and the blank line
// under it.
func (m WorktreeModel) visibleCount() int {
	available := m.height - headerLines
	if available < 1 {
		return fallbackVisibleCount
	}
	return available
}

func (m WorktreeModel) View() tea.View {
	var b strings.Builder

	visible := m.visibleCount()
	faint := lipgloss.NewStyle().Faint(true)

	b.WriteString("  " + m.filterInput.View())
	b.WriteString("\n")
	// The blank line under the filter doubles as the status line: the input is
	// padded out to the full width, so there is no room beside it, and a line of
	// its own would cost a row of worktrees.
	b.WriteString(faint.Render(m.statusLine()))
	b.WriteString("\n")

	switch {
	case m.loading:
		b.WriteString(faint.Render("  Loading worktrees..."))
		b.WriteString("\n")
		// Pad the remaining lines so the layout doesn't jump once loaded.
		for i := 1; i < visible; i++ {
			b.WriteString("\n")
		}

	case len(m.allItems) == 0:
		// Distinguished from a query that matched nothing: there is no worktree
		// to find, so the filter is not what is hiding them.
		b.WriteString(faint.Render("  No worktrees found"))
		b.WriteString("\n")
		for i := 1; i < visible; i++ {
			b.WriteString("\n")
		}

	default:
		end := min(m.offset+visible, len(m.filtered))

		cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)
		matchStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1)).Bold(true)
		normalStyle := lipgloss.NewStyle()

		for i := m.offset; i < end; i++ {
			b.WriteString(m.rowView(m.filtered[i], i == m.cursor, cursorStyle, matchStyle, normalStyle))
			b.WriteString("\n")
		}

		for i := end - m.offset; i < visible; i++ {
			b.WriteString("\n")
		}
	}

	// The last row leaves a trailing newline behind. Dropping it keeps the
	// frame exactly the height of the terminal.
	v := tea.NewView(strings.TrimSuffix(b.String(), "\n"))
	// Full window mode: the picker fills the terminal and hands the user's
	// scrollback back untouched when it quits.
	v.AltScreen = true
	return v
}

// statusLine reports an in-flight or failed ctrl+r refresh, and is empty the
// rest of the time. It is clipped to the terminal width so it cannot wrap and
// throw the row count off.
func (m WorktreeModel) statusLine() string {
	var line string
	switch {
	case m.refreshing:
		line = "  Refreshing worktrees..."
	case m.refreshErr != nil:
		line = "  Refresh failed: " + m.refreshErr.Error()
	default:
		return ""
	}

	if m.width > 0 && lipgloss.Width(line) > m.width {
		runes := []rune(line)
		line = string(runes[:m.width-1]) + "…"
	}
	return line
}

// numberGap separates the number column from the badge column.
const numberGap = "  "

// rowView renders one row: the number, the state badge, then the title, with
// the query's matches highlighted in the number and the title.
func (m WorktreeModel) rowView(row filteredWorktree, atCursor bool, cursorStyle, matchStyle, normalStyle lipgloss.Style) string {
	prefix := "  "
	if atCursor {
		prefix = cursorStyle.Render("> ")
	}

	number := padding(m.numberWidth-len(row.item.number)) +
		highlightMatches(row.item.number, row.numberIndexes, matchStyle, normalStyle)
	badge := stateBadge(row.item.entry.State, m.showIcons)

	used := lipgloss.Width(prefix) + m.numberWidth + len(numberGap) + badgeCellWidth
	title, titleIndexes := m.clipTitle(row.item.entry.Title, row.titleIndexes, used)

	return prefix + number + numberGap + badge +
		highlightMatches(title, titleIndexes, matchStyle, normalStyle)
}

// clipTitle shortens a title that would wrap onto the next line, dropping the
// match highlights that fall past the cut with it. Wrapping is what has to be
// avoided rather than merely tidied: a row spilling over two lines throws off
// the count of rows the cursor and the scroll offset are working from.
func (m WorktreeModel) clipTitle(title string, indexes []int, used int) (string, []int) {
	available := m.width - used
	// A width of zero means the terminal size is not known yet, so there is no
	// budget to clip against.
	if m.width <= 0 || lipgloss.Width(title) <= available {
		return title, indexes
	}
	if available <= 1 {
		return "", nil
	}

	runes := []rune(title)
	kept := available - 1 // room for the ellipsis
	var keptIndexes []int
	for _, idx := range indexes {
		if idx < kept {
			keptIndexes = append(keptIndexes, idx)
		}
	}
	return string(runes[:kept]) + "…", keptIndexes
}

// Chosen returns the selected worktree, and whether one was selected at all —
// false when the picker was quit, or when it was empty.
func (m WorktreeModel) Chosen() (model.WorktreeEntry, bool) { return m.chosen, m.picked }
func (m WorktreeModel) Quit() bool                          { return m.quit }
func (m WorktreeModel) LoadErr() error                      { return m.loadErr }
