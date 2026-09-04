package picker

import (
	"fmt"
	"image/color"
	"sort"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sahilm/fuzzy"

	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/model"
)

type sessionItem struct {
	session    model.SeshSession
	name       string // raw session name (no icons/ANSI)
	searchName string // normalized name used for fuzzy matching
	src        string // source type (tmux, config, zoxide, tmuxinator)
	icon       string // custom icon from config, "" to use the source glyph
}

// sessionItems implements fuzzy.Source for fuzzy matching.
type sessionItems []sessionItem

func (s sessionItems) String(i int) string { return s[i].searchName }
func (s sessionItems) Len() int            { return len(s) }

type filteredItem struct {
	item           sessionItem
	matchedIndexes []int
	// chipMatchLen is how many leading runes of the alias chip matched the
	// query, used to highlight the chip in alias-filter mode.
	chipMatchLen int
}

// displayRow is one line of the rendered list: either a session, or the rule
// drawn where one sort_order group ends and the next begins. A separator exists
// only here — the cursor indexes filtered, never rows, so it can neither land
// on one nor select one.
type displayRow struct {
	// index is the position in filtered, or -1 for a separator.
	index     int
	separator bool
}

// FetchFunc loads sessions asynchronously. It is called in a goroutine by Init().
type FetchFunc func() (model.SeshSessions, error)

// IconFunc resolves the icon configured for a session, returning "" when none
// is and the source glyph should be used instead.
type IconFunc func(session model.SeshSession) string

// PreviewFunc renders the preview of a single session. It is called from a
// tea.Cmd, never from Update or View: some strategies shell out, and blocking
// the event loop on every cursor move would make the picker feel sluggish.
type PreviewFunc func(name string) (string, error)

// sessionsLoadedMsg carries the result of the async fetch back to Update().
type sessionsLoadedMsg struct {
	sessions model.SeshSessions
	err      error
}

// previewFetchMsg fires after the debounce window and starts the fetch for the
// session highlighted at the time it was scheduled. The seq identifies that
// cursor position so later movement can invalidate it.
type previewFetchMsg struct {
	seq  int
	name string
}

// previewLoadedMsg carries a rendered preview back to Update(). A result whose
// seq is no longer current is discarded: the cursor has moved on, and the pane
// keeps showing the last preview that did belong to its session.
type previewLoadedMsg struct {
	seq     int
	name    string
	content string
	err     error
}

// aliasAutoConnectMsg fires once the auto-connect grace period has elapsed. The
// seq identifies the filter change that scheduled it so later keystrokes can
// invalidate it.
type aliasAutoConnectMsg struct {
	seq   int
	alias string
}

// Alias is a picker shortcut for a single session, derived from an
// [[session]] config block.
type Alias struct {
	// Alias is the shortcut as written in the config, used for display.
	Alias string
	// Target is the name of the session the alias resolves to.
	Target string
	// AutoConnect connects to Target as soon as the alias is fully typed.
	AutoConnect bool
}

// Options configures the picker model. Zero values are valid and mean
// "off"/"unset", except Prompt and Placeholder which are passed through as-is.
type Options struct {
	ShowIcons      bool
	ShowWindows    bool
	SeparatorAware bool
	Prompt         string
	Placeholder    string
	// Query pre-fills the filter input. It goes in as if it had been typed —
	// the sigils for alias and index mode included — but never triggers alias
	// auto-connect on its own, which stays a reward for an actual keystroke.
	Query string
	// Aliases is keyed by lowercased alias.
	Aliases map[string]Alias
	// AliasFilterPrefix is the sigil that, typed first, narrows the list to
	// aliased sessions. Empty disables the mode.
	AliasFilterPrefix string
	// AliasAutoConnectDelay is the grace period before auto-connect fires,
	// leaving room to finish typing a longer alias that shares a prefix.
	AliasAutoConnectDelay time.Duration
	// DisableAliasAutoConnect suppresses auto-connect for this invocation.
	DisableAliasAutoConnect bool
	// Icon resolves the icon configured for a session. Nil means no session
	// declares one, so every row keeps its source glyph.
	Icon IconFunc
	// IconWidth is the display width the icon column is padded to. Anything
	// below 1 means the single cell the source glyphs need.
	IconWidth int
	// Preview starts the picker with the preview pane showing. It can still be
	// toggled at runtime, and is ignored on a terminal narrower than
	// PreviewMinWidth.
	Preview bool
	// PreviewWidth is the share of the terminal, in percent, guaranteed to the
	// preview pane.
	PreviewWidth int
	// PreviewMinWidth is the narrowest terminal that gets a preview pane.
	PreviewMinWidth int
	// PreviewBorder names the divider drawn between the list and the preview
	// pane, one of the model.PreviewBorder* values. Empty or unrecognized
	// means the default.
	PreviewBorder string
	// PreviewFunc renders previews. Nil disables the pane entirely.
	PreviewFunc PreviewFunc
	// GroupSeparator draws a faint rule between sort_order groups. It is
	// suppressed while the list is filtered, where the groups no longer occupy
	// contiguous ranges.
	GroupSeparator bool
}

type Model struct {
	allItems    sessionItems
	filtered    []filteredItem
	filterInput textinput.Model
	// cursor indexes filtered; offset indexes rows, which is the same list with
	// the group separators laid in, so scrolling counts the lines actually
	// drawn.
	cursor int
	offset int
	// rows is the rendered layout of filtered — see displayRow — and rowOf maps
	// a filtered index back to the row that draws it.
	rows           []displayRow
	rowOf          []int
	groupSeparator bool
	width          int
	height         int
	chosen         string
	quit           bool
	showIcons      bool
	showWindows    bool
	separatorAware bool
	focusCmd       tea.Cmd
	loading        bool
	fetchFunc      FetchFunc
	loadErr        error

	// iconFunc resolves a session's configured icon, and iconWidth is the width
	// the icon column is padded to so wide icons keep the names aligned.
	iconFunc  IconFunc
	iconWidth int

	// aliases is keyed by lowercased alias; aliasByName is keyed by session
	// name so rows can be tagged while rendering.
	aliases                 map[string]Alias
	aliasByName             map[string]string
	aliasFilterPrefix       string
	aliasAutoConnectDelay   time.Duration
	disableAliasAutoConnect bool
	// aliasSeq increments on every filter change so a pending auto-connect
	// tick can tell whether it is stale.
	aliasSeq int

	previewFunc     PreviewFunc
	previewOn       bool
	previewWidthPct int
	previewMinWidth int
	// previewName is the session the current content belongs to, and
	// previewPending the one being fetched. Together they keep the cursor
	// landing back on an already-previewed row from refetching it.
	// previewBorderName is the resolved divider style, so an unset or bogus
	// config value never reaches the renderer.
	previewBorderName string
	previewName       string
	previewPending    string
	previewContent    string
	previewErr        error
	// previewSeq increments on every request so a result that arrives after the
	// cursor moved on can be told apart from the current one.
	previewSeq int
}

// srcIcon returns the nerd font icon and color for a session source. The icon
// carries a trailing space, so it is the whole icon cell for the single-width
// glyphs it comes from; iconCell pads it when a wider icon is in play.
func srcIcon(src string) (string, color.Color) {
	if g, ok := icon.Glyphs[src]; ok {
		var ansi int
		switch {
		case g.ColorCode >= 90 && g.ColorCode <= 97:
			ansi = g.ColorCode - 82
		case g.ColorCode >= 30 && g.ColorCode <= 37:
			ansi = g.ColorCode - 30
		default:
			ansi = g.ColorCode
		}
		return g.Icon + " ", lipgloss.ANSIColor(ansi)
	}
	return "? ", lipgloss.ANSIColor(8)
}

// iconCell renders the icon column for a row: the icon configured for the
// session, or the colored source glyph when it has none, followed by the gap
// before the name. Every cell is padded to the same width so a double-width
// emoji on one row can't push that row's name past the others.
//
// A configured icon is left unstyled — emoji carry their own color, and a nerd
// font glyph inherits the default foreground.
//
// Trailing spaces in a configured icon are the user's own width override: they
// are added to the cell but not counted towards it. Terminals disagree about how
// wide an emoji is — an emoji written with a variation selector (⬆️ is U+2B06
// plus U+FE0F) measures two cells here but is drawn in one by WezTerm and
// others — and nothing in a TUI can ask the terminal which it did. A trailing
// space makes up the difference for that one icon without shifting any other
// row.
func (m Model) iconCell(item sessionItem) string {
	if item.icon == "" {
		glyph, clr := srcIcon(item.src)
		return padIcon(lipgloss.NewStyle().Foreground(clr).Render(glyph), m.iconWidth+1)
	}
	override := strings.TrimRight(item.icon, " ")
	return item.icon + " " + padding(m.iconWidth+1-lipgloss.Width(override+" "))
}

// padIcon right-pads an icon cell to width.
func padIcon(cell string, width int) string {
	return cell + padding(width-lipgloss.Width(cell))
}

func padding(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}

// Powerline half circles used to round off the alias chip. They are nerd font
// glyphs, so they are only used when icons are enabled.
const (
	chipLeftGlyph  = ""
	chipRightGlyph = ""
)

// aliasChip renders the alias of a session as a chip that sits before the
// session name, e.g. ` wp ` with nerd fonts or `[wp] ` without. It returns ""
// for sessions that have no alias. The first matchLen runes of the alias are
// highlighted, for alias-filter mode.
//
// The label uses reverse video rather than an explicit fg/bg pair so the text
// is always the terminal's background color on its foreground color, which
// stays legible under any color scheme. The half circles are left unstyled for
// the same reason: their default foreground is exactly the color the reversed
// label fills with, so the chip reads as one shape.
func (m Model) aliasChip(name string, matchLen int) string {
	alias, ok := m.aliasByName[name]
	if !ok {
		return ""
	}

	label := lipgloss.NewStyle().Reverse(true)
	text := highlightChipPrefix(alias, matchLen, label)
	if !m.showIcons {
		// Without nerd fonts the half circles render as tofu, so fall back to
		// brackets that still read as a chip.
		return label.Render("[") + text + label.Render("]") + " "
	}
	return chipLeftGlyph + text + chipRightGlyph + " "
}

// highlightChipPrefix styles the first matchLen runes of a chip label as
// matched. Reverse video is kept throughout so the chip stays one shape;
// setting a foreground under it paints the matched runes as a colored block,
// which is the same red used to highlight matches in session names.
func highlightChipPrefix(alias string, matchLen int, label lipgloss.Style) string {
	runes := []rune(alias)
	if matchLen <= 0 {
		return label.Render(alias)
	}
	if matchLen > len(runes) {
		matchLen = len(runes)
	}
	matched := label.Foreground(lipgloss.ANSIColor(1)).Render(string(runes[:matchLen]))
	return matched + label.Render(string(runes[matchLen:]))
}

// indexFilterPrefix is the sigil that, typed first, enters index mode: the rows
// are numbered and the next digit jumps straight to one of them.
const indexFilterPrefix = "#"

// maxIndexJump is the highest position index mode can reach. Only single digits
// are in scope, so rows past the ninth are unnumbered.
const maxIndexJump = 9

var separatorReplacer = strings.NewReplacer("-", " ", "_", " ", "/", " ", "\\", " ")

func normalizeSeparators(s string) string {
	return separatorReplacer.Replace(s)
}

func buildItems(sessions model.SeshSessions, separatorAware bool, resolveIcon IconFunc) sessionItems {
	items := make(sessionItems, 0, len(sessions.OrderedIndex))
	for _, key := range sessions.OrderedIndex {
		s := sessions.Directory[key]
		searchName := s.Name
		if separatorAware {
			searchName = normalizeSeparators(s.Name)
		}
		var icn string
		if resolveIcon != nil {
			icn = resolveIcon(s)
		}
		items = append(items, sessionItem{
			session:    s,
			name:       s.Name,
			searchName: searchName,
			src:        s.Src,
			icon:       icn,
		})
	}
	return items
}

func New(fetchFunc FetchFunc, opts Options) Model {
	ti := textinput.New()
	ti.Placeholder = opts.Placeholder
	ti.Prompt = opts.Prompt
	// SetValue leaves the cursor at the end, so the pre-filled query reads as
	// something just typed and can be backspaced away.
	ti.SetValue(opts.Query)

	aliasByName := make(map[string]string, len(opts.Aliases))
	for _, alias := range opts.Aliases {
		aliasByName[alias.Target] = alias.Alias
	}

	m := Model{
		filterInput:             ti,
		showIcons:               opts.ShowIcons,
		showWindows:             opts.ShowWindows,
		separatorAware:          opts.SeparatorAware,
		iconFunc:                opts.Icon,
		iconWidth:               max(opts.IconWidth, 1),
		loading:                 true,
		fetchFunc:               fetchFunc,
		aliases:                 opts.Aliases,
		aliasByName:             aliasByName,
		aliasFilterPrefix:       opts.AliasFilterPrefix,
		aliasAutoConnectDelay:   opts.AliasAutoConnectDelay,
		disableAliasAutoConnect: opts.DisableAliasAutoConnect,
		previewFunc:             opts.PreviewFunc,
		previewOn:               opts.Preview,
		previewWidthPct:         previewWidth(opts.PreviewWidth),
		previewMinWidth:         previewMinWidth(opts.PreviewMinWidth),
		previewBorderName:       previewBorder(opts.PreviewBorder),
		groupSeparator:          opts.GroupSeparator,
	}
	m.focusCmd = m.filterInput.Focus()
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.focusCmd, m.fetchSessions())
}

func (m Model) fetchSessions() tea.Cmd {
	return func() tea.Msg {
		sessions, err := m.fetchFunc()
		return sessionsLoadedMsg{sessions: sessions, err: err}
	}
}

// previewDebounce is how long a cursor position has to hold before its preview
// is fetched. Holding a movement key would otherwise shell out once per row
// passed over; the stale-result guard makes those harmless, but not free.
const previewDebounce = 60 * time.Millisecond

// syncInputWidth fits the filter input to the list column. Toggling the preview
// pane resizes that column, so it is called from there too, not just on resize.
func (m *Model) syncInputWidth() {
	m.filterInput.SetWidth(m.contentWidth() - 4)
}

// highlightedName returns the name of the session under the cursor.
func (m Model) highlightedName() (string, bool) {
	if m.cursor < 0 || m.cursor >= len(m.filtered) {
		return "", false
	}
	return m.filtered[m.cursor].item.name, true
}

// schedulePreview arms a debounced fetch for the highlighted session, and
// returns nil when there is nothing to do — the pane is off, the session is
// already previewed, or its fetch is already in flight.
func (m *Model) schedulePreview() tea.Cmd {
	if m.previewFunc == nil || !m.previewOn {
		return nil
	}

	name, ok := m.highlightedName()
	if !ok {
		// Nothing is highlighted, so there is nothing the old content could
		// still be describing.
		m.previewSeq++
		m.previewName, m.previewPending, m.previewContent, m.previewErr = "", "", "", nil
		return nil
	}
	if name == m.previewName || name == m.previewPending {
		return nil
	}

	m.previewSeq++
	m.previewPending = name
	msg := previewFetchMsg{seq: m.previewSeq, name: name}
	return tea.Tick(previewDebounce, func(time.Time) tea.Msg { return msg })
}

func (m Model) fetchPreview(seq int, name string) tea.Cmd {
	previewFunc := m.previewFunc
	return func() tea.Msg {
		content, err := previewFunc(name)
		return previewLoadedMsg{seq: seq, name: name, content: content, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case aliasAutoConnectMsg:
		// The tick can't be cancelled, so a stale one still arrives: ignore it
		// unless it is the latest and the filter still reads as that alias.
		if msg.seq != m.aliasSeq {
			return m, nil
		}
		alias, ok := m.exactAlias(m.filterInput.Value())
		if !ok || alias.Alias != msg.alias || m.disableAliasAutoConnect {
			return m, nil
		}
		m.chosen = alias.Target
		return m, tea.Quit

	case previewFetchMsg:
		// The tick can't be cancelled, so a stale one still arrives.
		if msg.seq != m.previewSeq {
			return m, nil
		}
		return m, m.fetchPreview(msg.seq, msg.name)

	case previewLoadedMsg:
		if msg.seq != m.previewSeq {
			return m, nil
		}
		m.previewPending = ""
		m.previewName = msg.name
		m.previewContent = msg.content
		m.previewErr = msg.err
		return m, nil

	case sessionsLoadedMsg:
		if msg.err != nil {
			m.loadErr = msg.err
			return m, tea.Quit
		}
		m.loading = false
		m.allItems = buildItems(msg.sessions, m.separatorAware, m.iconFunc)
		m.applyFilter()
		return m, m.schedulePreview()

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
			if len(m.filtered) > 0 {
				selected := m.filtered[m.cursor]
				m.chosen = selected.item.name
			}
			return m, tea.Quit

		case "esc", "ctrl+c":
			m.quit = true
			return m, tea.Quit

		case "up", "ctrl+k", "ctrl+p":
			m.cursorUp(1)
			return m, m.schedulePreview()

		case "down", "ctrl+j", "ctrl+n":
			m.cursorDown(1)
			return m, m.schedulePreview()

		case "ctrl+u":
			m.cursorUp(m.visibleCount() / 2)
			return m, m.schedulePreview()

		case "ctrl+d":
			m.cursorDown(m.visibleCount() / 2)
			return m, m.schedulePreview()

		case "ctrl+o":
			// Toggling stays allowed on a narrow terminal: the pane is gated on
			// width at render time, so it appears once the window has room.
			if m.previewFunc == nil {
				return m, nil
			}
			m.previewOn = !m.previewOn
			m.syncInputWidth()
			if !m.previewOn {
				// Drop anything in flight so switching back on always refetches
				// the highlighted row rather than waiting on a request made
				// before the pane was hidden.
				m.previewSeq++
				m.previewPending = ""
				return m, nil
			}
			return m, m.schedulePreview()

		default:
			// Index mode owns the digits while it is active, so an out-of-range
			// one is swallowed rather than filtering the list.
			if name, handled := m.indexJump(msg.String()); handled {
				if name == "" {
					return m, nil
				}
				m.chosen = name
				return m, tea.Quit
			}
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
		cmds := []tea.Cmd{cmd}
		if tick := m.scheduleAliasAutoConnect(); tick != nil {
			cmds = append(cmds, tick)
		}
		if preview := m.schedulePreview(); preview != nil {
			cmds = append(cmds, preview)
		}
		if len(cmds) > 1 {
			return m, tea.Batch(cmds...)
		}
	}

	return m, cmd
}

// exactAlias resolves the input to an alias when one is typed exactly, in
// either normal or alias-filter mode, and reports whether that alias should
// auto-connect.
//
// In alias-filter mode the per-session alias_auto_connect opt-in is bypassed:
// reaching for the sigil is itself an explicit statement that the next thing
// typed is an alias to jump to. In normal mode it still gates, so an alias
// typed incidentally while filtering doesn't connect on its own.
//
// The input is matched raw, never separator-normalized, so an alias containing
// `-` or `_` can't be matched by typing a space instead.
func (m *Model) exactAlias(raw string) (Alias, bool) {
	if query, ok := m.aliasFilterQuery(raw); ok {
		alias, found := m.aliases[strings.ToLower(query)]
		return alias, found
	}
	alias, found := m.aliases[strings.ToLower(raw)]
	if !found || !alias.AutoConnect {
		return Alias{}, false
	}
	return alias, true
}

// scheduleAliasAutoConnect invalidates any pending auto-connect and, when the
// filter now reads as an alias that should auto-connect, arms a fresh one.
// It returns nil when there is nothing to schedule. Auto-connect deliberately
// works while sessions are still loading: the target comes from the config, not
// from the list, so there is no reason to make the user wait.
func (m *Model) scheduleAliasAutoConnect() tea.Cmd {
	m.aliasSeq++

	if m.disableAliasAutoConnect {
		return nil
	}
	alias, ok := m.exactAlias(m.filterInput.Value())
	if !ok {
		return nil
	}

	msg := aliasAutoConnectMsg{seq: m.aliasSeq, alias: alias.Alias}
	if m.aliasAutoConnectDelay <= 0 {
		return func() tea.Msg { return msg }
	}
	return tea.Tick(m.aliasAutoConnectDelay, func(time.Time) tea.Msg { return msg })
}

// aliasMatch resolves an exactly-typed alias to the session it points at. It
// takes the raw input, never the separator-normalized pattern, so an alias
// containing `-` or `_` can't be matched by typing a space instead.
//
// The session is looked up in the loaded list so its source and window names
// come along, and falls back to a config-sourced item when the target isn't
// listed — an alias always names a [[session]], so connecting still works even
// if that session is filtered out or hasn't loaded yet.
func (m *Model) aliasMatch(raw string) (sessionItem, bool) {
	alias, ok := m.aliases[strings.ToLower(raw)]
	if !ok {
		return sessionItem{}, false
	}
	for _, item := range m.allItems {
		if item.name == alias.Target {
			return item, true
		}
	}
	return m.configItem(alias.Target), true
}

// configItem builds the row for a [[session]] that isn't in the loaded list.
// Only its name is known, so the icon is resolved from the name alone.
func (m *Model) configItem(name string) sessionItem {
	item := sessionItem{name: name, searchName: name, src: "config"}
	if m.iconFunc != nil {
		item.icon = m.iconFunc(model.SeshSession{Src: "config", Name: name})
	}
	return item
}

// aliasFilterQuery reports whether the raw input opens alias-filter mode and,
// if so, returns whatever was typed after the sigil. Only a leading sigil
// counts, so slashes inside a path-like query are left alone.
func (m *Model) aliasFilterQuery(raw string) (string, bool) {
	if m.aliasFilterPrefix == "" || !strings.HasPrefix(raw, m.aliasFilterPrefix) {
		return "", false
	}
	return strings.TrimPrefix(raw, m.aliasFilterPrefix), true
}

// aliasCandidates returns every aliased session as a row, in list order.
// Aliases whose target isn't in the loaded list still appear — an alias always
// names a [[session]], so it must remain reachable even when that session is
// filtered out or hasn't started yet. Those trail the listed ones, sorted by
// alias so the order doesn't shift between runs.
func (m *Model) aliasCandidates() []sessionItem {
	candidates := make([]sessionItem, 0, len(m.aliases))
	listed := make(map[string]bool, len(m.aliases))

	for _, item := range m.allItems {
		alias, ok := m.aliasByName[item.name]
		if !ok {
			continue
		}
		listed[strings.ToLower(alias)] = true
		candidates = append(candidates, item)
	}

	missing := make([]string, 0, len(m.aliases))
	for key := range m.aliases {
		if !listed[key] {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)
	for _, key := range missing {
		candidates = append(candidates, m.configItem(m.aliases[key].Target))
	}

	return candidates
}

// filterAliases narrows the aliased sessions by query, in two tiers: aliases
// the query prefixes come first, then aliases whose session name contains it.
// Session names are multi-word, so a prefix rule there would make "config"
// miss "tmux config"; the alias tier stays a prefix match because aliases are
// short and should narrow deterministically as they're typed.
func (m *Model) filterAliases(query string) []filteredItem {
	candidates := m.aliasCandidates()
	q := strings.ToLower(query)
	matchLen := len([]rune(query))

	byAlias := make([]filteredItem, 0, len(candidates))
	byName := make([]filteredItem, 0, len(candidates))
	for _, item := range candidates {
		alias := m.aliasByName[item.name]
		switch {
		case strings.HasPrefix(strings.ToLower(alias), q):
			byAlias = append(byAlias, filteredItem{item: item, chipMatchLen: matchLen})
		case q != "" && strings.Contains(strings.ToLower(item.name), q):
			byName = append(byName, filteredItem{item: item})
		}
	}
	return append(byAlias, byName...)
}

// indexFilterQuery reports whether the raw input opens index mode and, if so,
// returns whatever was typed after the sigil. Only a leading sigil counts, so a
// `#` inside a branch-like query is left alone. The alias sigil is resolved
// first, so configuring it as `#` keeps alias mode and leaves index mode
// unreachable rather than making the two fight over the same keystroke.
func (m *Model) indexFilterQuery(raw string) (string, bool) {
	if _, ok := m.aliasFilterQuery(raw); ok {
		return "", false
	}
	if !strings.HasPrefix(raw, indexFilterPrefix) {
		return "", false
	}
	return strings.TrimPrefix(raw, indexFilterPrefix), true
}

// indexJump resolves a keypress in index mode to the session it selects. The
// second return reports that the key belongs to the mode: it is true for every
// digit 1-9 while the mode is active, even when no row sits at that position, so
// a digit past the end of the list is a no-op rather than a literal in the
// filter. Digits are counted against the visible list, so the numbering stays
// meaningful when a query narrowed it.
func (m Model) indexJump(key string) (string, bool) {
	if _, ok := m.indexFilterQuery(m.filterInput.Value()); !ok {
		return "", false
	}
	if len(key) != 1 || key[0] < '1' || key[0] > '9' {
		return "", false
	}
	if n := int(key[0] - '0'); n <= len(m.filtered) {
		return m.filtered[n-1].item.name, true
	}
	return "", true
}

func (m *Model) applyFilter() {
	m.filtered = m.filterFor(m.filterInput.Value())
	m.buildRows()
}

// filterFor narrows the loaded list for the raw input, dispatching on the sigil
// it was typed with.
func (m *Model) filterFor(raw string) []filteredItem {
	if query, ok := m.aliasFilterQuery(raw); ok {
		return m.filterAliases(query)
	}
	// Index mode numbers whatever is displayed, so anything typed after the
	// sigil filters exactly as it would on its own.
	if query, ok := m.indexFilterQuery(raw); ok {
		return m.filterSessions(query)
	}
	return m.filterSessions(raw)
}

// buildRows lays the filtered sessions out as the lines the list renders,
// inserting a rule wherever the sort_order group changes.
//
// Separators are suppressed as soon as anything is typed: a query reorders the
// results by match quality, so the groups stop being contiguous ranges and a
// rule would be drawing a boundary that isn't there.
func (m *Model) buildRows() {
	m.rows = make([]displayRow, 0, len(m.filtered))
	m.rowOf = make([]int, len(m.filtered))
	separate := m.groupSeparator && m.filterInput.Value() == ""
	for i, item := range m.filtered {
		if separate && i > 0 && item.item.session.Group != m.filtered[i-1].item.session.Group {
			m.rows = append(m.rows, displayRow{index: -1, separator: true})
		}
		m.rowOf[i] = len(m.rows)
		m.rows = append(m.rows, displayRow{index: i})
	}
}

// filterSessions narrows the loaded list by pattern: everything when it's
// empty, the single target when it is an alias typed exactly, and a fuzzy match
// otherwise.
func (m *Model) filterSessions(pattern string) []filteredItem {
	if pattern == "" {
		filtered := make([]filteredItem, len(m.allItems))
		for i, item := range m.allItems {
			filtered[i] = filteredItem{item: item}
		}
		return filtered
	}

	if item, ok := m.aliasMatch(pattern); ok {
		return []filteredItem{{item: item}}
	}

	if m.separatorAware {
		pattern = normalizeSeparators(pattern)
	}

	matches := rankMatches(fuzzy.FindFromNoSort(pattern, m.allItems))
	filtered := make([]filteredItem, len(matches))
	for i, match := range matches {
		filtered[i] = filteredItem{
			item:           m.allItems[match.Index],
			matchedIndexes: match.MatchedIndexes,
		}
	}
	return filtered
}

// maxUnmatchedCharPenalty caps how much a name can be docked for the characters
// the query didn't match. It is deliberately small — the same size as the fuzzy
// library's adjacent-match bonus — so length can only ever break a tie between
// matches the scorer already rates the same.
const maxUnmatchedCharPenalty = 5

// rankMatches scores unsorted fuzzy matches with the library's length penalty
// capped, and orders them best first.
//
// sahilm/fuzzy docks a match one point for every character the query didn't
// match, with no floor, so a long name loses to a short one even when it is
// plainly the better match: typing "sesh" put a tmux session named
// "sesh/w/423 — Add opt-in preview pane to the picker TUI" below a dozen zoxide
// paths, despite matching at its very first character. Capping the penalty keeps
// length as a tiebreaker without letting it outweigh where the match landed.
//
// Matches must come in unsorted (fuzzy.FindFromNoSort), so they arrive in
// session-list order. The sort is stable, so equally scored matches keep that
// order — which is the order the picker shows with an empty query, tmux
// sessions first.
func rankMatches(matches fuzzy.Matches) fuzzy.Matches {
	ranked := make([]struct {
		match fuzzy.Match
		score int
	}, len(matches))
	for i, match := range matches {
		penalty := len(match.MatchedIndexes) - len(match.Str)
		ranked[i].match = match
		ranked[i].score = match.Score - penalty + max(penalty, -maxUnmatchedCharPenalty)
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	for i, r := range ranked {
		matches[i] = r.match
	}
	return matches
}

func (m *Model) cursorUp(n int) {
	m.cursor -= n
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.scrollToCursor()
}

func (m *Model) cursorDown(n int) {
	m.cursor += n
	max := len(m.filtered) - 1
	if max < 0 {
		max = 0
	}
	if m.cursor > max {
		m.cursor = max
	}
	m.scrollToCursor()
}

// scrollToCursor moves the viewport the least it can to keep the highlighted
// session on screen. It works in row space so the separators between it and the
// cursor are counted as the lines they are.
func (m *Model) scrollToCursor() {
	if m.cursor < 0 || m.cursor >= len(m.rowOf) {
		m.offset = 0
		return
	}
	row := m.rowOf[m.cursor]
	if row < m.offset {
		m.offset = row
	}
	if visible := m.visibleCount(); row >= m.offset+visible {
		m.offset = row - visible + 1
	}
}

// fallbackVisibleCount is used for the frame that can render before the
// terminal size is known.
const fallbackVisibleCount = 5

// visibleCount is how many session rows fit. The picker runs on the alt screen,
// so it gets every row the terminal has apart from the filter row and the blank
// line under it.
func (m Model) visibleCount() int {
	available := m.height - headerLines
	if available < 1 {
		return fallbackVisibleCount
	}
	return available
}

const (
	// maxListWidth is the widest the session list gets, split or not.
	maxListWidth = 60
	// minListWidth is the narrowest list worth splitting off; the preview gives
	// columns back to stay above it.
	minListWidth = 40
	// previewPadding is the gap between the divider — or the list itself, with
	// the divider off — and the preview text. It sits inside the pane's width.
	previewPadding = 1
)

// previewBorderStyle returns the lipgloss border for a resolved divider name,
// and whether a divider is drawn at all.
func previewBorderStyle(name string) (lipgloss.Border, bool) {
	switch name {
	case model.PreviewBorderNone:
		return lipgloss.Border{}, false
	case model.PreviewBorderThick:
		return lipgloss.ThickBorder(), true
	case model.PreviewBorderDouble:
		return lipgloss.DoubleBorder(), true
	default:
		return lipgloss.NormalBorder(), true
	}
}

// previewChrome is how much of the pane's width goes to the divider and the
// padding after it, leaving the text that much less room. Without a divider
// only the padding is charged, so disabling it gives that column to the text.
func (m Model) previewChrome() int {
	if _, drawn := previewBorderStyle(m.previewBorderName); drawn {
		return previewPadding + 1
	}
	return previewPadding
}

// splitActive reports whether this frame renders a preview pane. Width is
// checked here rather than at toggle time so growing the window is enough to
// bring the pane back.
func (m Model) splitActive() bool {
	return m.previewOn && m.previewFunc != nil && m.width >= m.previewMinWidth
}

// previewCols is how many columns the preview pane occupies, divider included.
// The configured percent is a floor: the list is capped at maxListWidth, and
// everything past that goes to the preview rather than being left blank. The
// list is never squeezed below minListWidth.
func (m Model) previewCols() int {
	if !m.splitActive() {
		return 0
	}
	cols := m.width * m.previewWidthPct / 100
	if rest := m.width - maxListWidth; rest > cols {
		cols = rest
	}
	if max := m.width - minListWidth; cols > max {
		cols = max
	}
	if cols <= m.previewChrome() {
		return 0
	}
	return cols
}

func (m Model) contentWidth() int {
	w := m.width - m.previewCols()
	if w < 30 {
		w = minListWidth
	}
	if w > maxListWidth {
		w = maxListWidth
	}
	return w
}

func (m Model) View() tea.View {
	var b strings.Builder

	// Filter input
	b.WriteString("  " + m.filterInput.View())
	b.WriteString("\n\n")

	visible := m.visibleCount()

	if m.loading {
		loadingStyle := lipgloss.NewStyle().Faint(true)
		b.WriteString(loadingStyle.Render("  Loading sessions..."))
		b.WriteString("\n")
		// Pad remaining visible lines to prevent layout jump
		for i := 1; i < visible; i++ {
			b.WriteString("\n")
		}
	} else if _, aliasMode := m.aliasFilterQuery(m.filterInput.Value()); aliasMode && len(m.aliases) == 0 {
		// Without this the sigil looks broken rather than simply unconfigured.
		emptyStyle := lipgloss.NewStyle().Faint(true)
		b.WriteString(emptyStyle.Render("  No aliases configured"))
		b.WriteString("\n")
		for i := 1; i < visible; i++ {
			b.WriteString("\n")
		}
	} else {
		// Session list
		end := m.offset + visible
		if end > len(m.rows) {
			end = len(m.rows)
		}

		cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)
		matchStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1)).Bold(true)
		normalStyle := lipgloss.NewStyle()
		windowStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8)).Faint(true)
		indexStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(4))

		// Index mode numbers the rows so the jump target can be read off the
		// list instead of counted.
		_, indexMode := m.indexFilterQuery(m.filterInput.Value())

		for r := m.offset; r < end; r++ {
			if m.rows[r].separator {
				b.WriteString(m.separatorRule())
				b.WriteString("\n")
				continue
			}
			i := m.rows[r].index
			item := m.filtered[i]
			prefix := "  "
			if i == m.cursor {
				prefix = cursorStyle.Render("> ")
			}
			if indexMode {
				prefix += indexGutter(i, indexStyle)
			}

			var tag string
			if m.showIcons {
				tag = m.iconCell(item.item)
			}
			chip := m.aliasChip(item.item.name, item.chipMatchLen)
			name := highlightMatches(item.item.name, item.matchedIndexes, matchStyle, normalStyle)

			var windows string
			if m.showWindows {
				// Window names are display-only: they are never highlighted as
				// matches and never become part of the selected value.
				used := lipgloss.Width(prefix) + lipgloss.Width(tag) + lipgloss.Width(chip) + lipgloss.Width(item.item.name)
				if text := windowsText(item.item.session.WindowNames, m.contentWidth()-used); text != "" {
					windows = windowStyle.Render(text)
				}
			}

			b.WriteString(fmt.Sprintf("%s%s%s%s%s\n", prefix, tag, chip, name, windows))
		}

		// Pad remaining visible lines
		for i := end - m.offset; i < visible; i++ {
			b.WriteString("\n")
		}
	}

	// The last row leaves a trailing newline behind. Dropping it keeps the frame
	// exactly the height of the terminal, and makes both halves of a split the
	// same height so the divider spans the whole block.
	content := strings.TrimSuffix(b.String(), "\n")

	if cols := m.previewCols(); cols > 0 {
		list := lipgloss.NewStyle().
			Width(m.contentWidth()).
			MaxWidth(m.contentWidth()).
			Render(content)
		content = lipgloss.JoinHorizontal(lipgloss.Top, list, m.previewView(cols, visible))
	}

	v := tea.NewView(content)
	// Full window mode: the picker fills the terminal and hands the user's
	// scrollback back untouched when it quits.
	v.AltScreen = true
	return v
}

// separatorRule draws the boundary between two sort_order groups: a faint rule
// across the list column, so a glance says whether the row under the cursor is
// a live tmux session or somewhere to go.
func (m Model) separatorRule() string {
	width := m.contentWidth() - 2
	if width < 1 {
		width = 1
	}
	return lipgloss.NewStyle().Faint(true).Render("  " + strings.Repeat("─", width))
}

// indexGutter renders the jump number for the row at position i, or blanks of
// the same width once the numbers run out, so every name stays aligned.
func indexGutter(i int, style lipgloss.Style) string {
	if i >= maxIndexJump {
		return "  "
	}
	return style.Render(fmt.Sprintf("%d ", i+1))
}

// headerLines is the filter row plus the blank line under it, which the
// preview pane skips so its first line sits level with the first session row.
const headerLines = 2

// previewView renders the preview pane: a left divider, then the content of the
// highlighted session clipped to the pane.
func (m Model) previewView(cols, rows int) string {
	faint := lipgloss.NewStyle().Faint(true)

	var body string
	switch {
	case m.previewErr != nil:
		// A broken preview_command shouldn't take the picker down with it.
		body = faint.Render(fmt.Sprintf("Preview unavailable: %v", m.previewErr))
	case m.previewName == "":
		body = faint.Render("Loading preview...")
	case strings.TrimSpace(m.previewContent) == "":
		body = faint.Render("No preview")
	default:
		body = m.previewContent
	}

	// The divider and padding live inside Width, so the text gets less room.
	body = clipLines(body, cols-m.previewChrome(), rows)

	border, drawn := previewBorderStyle(m.previewBorderName)

	return lipgloss.NewStyle().
		Border(border, false, false, false, drawn).
		BorderForeground(lipgloss.ANSIColor(8)).
		PaddingLeft(previewPadding).
		PaddingTop(headerLines).
		Width(cols).
		Height(headerLines + rows).
		Render(body)
}

// clipLines cuts text to at most rows lines of width columns each. Truncation
// is ANSI-aware because tmux previews come from `capture-pane -e`, and each
// kept line is reset so a color left open can't bleed into the pane below it.
func clipLines(text string, width, rows int) string {
	if width < 1 || rows < 1 {
		return ""
	}

	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if len(lines) > rows {
		lines = lines[:rows]
	}

	clip := lipgloss.NewStyle().MaxWidth(width)
	for i, line := range lines {
		line = clip.Render(line)
		if strings.Contains(line, "\x1b") {
			line += "\x1b[0m"
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

// windowGap separates the session name from the window names, and each window
// name from the next.
const windowGap = " "

// windowsText renders window names inline within budget columns, eliding the
// ones that don't fit as "+N". It returns "" when there are no names or not
// even one name fits.
func windowsText(names []string, budget int) string {
	if len(names) == 0 || budget <= 0 {
		return ""
	}

	var b strings.Builder
	used, shown := 0, 0
	for i, name := range names {
		entry := windowGap + name
		need := used + lipgloss.Width(entry)
		// Reserve room for the elision label in case a later name doesn't fit.
		if remaining := len(names) - i - 1; remaining > 0 {
			need += lipgloss.Width(elisionLabel(remaining))
		}
		if need > budget {
			break
		}
		b.WriteString(entry)
		used += lipgloss.Width(entry)
		shown++
	}

	if shown == len(names) {
		return b.String()
	}

	label := elisionLabel(len(names) - shown)
	if used+lipgloss.Width(label) <= budget {
		return b.String() + label
	}
	return b.String()
}

func elisionLabel(hidden int) string {
	return fmt.Sprintf("%s+%d", windowGap, hidden)
}

func highlightMatches(s string, indexes []int, matchStyle, normalStyle lipgloss.Style) string {
	if len(indexes) == 0 {
		return normalStyle.Render(s)
	}

	runes := []rune(s)
	matchSet := make([]bool, len(runes))
	for _, idx := range indexes {
		if idx >= 0 && idx < len(runes) {
			matchSet[idx] = true
		}
	}

	// Render one lipgloss call per run of same-styled runes rather than per
	// rune: styling is by far the most expensive part of a keystroke, and a
	// query usually produces only a handful of runs per name.
	var result strings.Builder
	for start := 0; start < len(runes); {
		matched := matchSet[start]
		end := start + 1
		for end < len(runes) && matchSet[end] == matched {
			end++
		}
		run := string(runes[start:end])
		if matched {
			result.WriteString(matchStyle.Render(run))
		} else {
			result.WriteString(normalStyle.Render(run))
		}
		start = end
	}
	return result.String()
}

func (m Model) Chosen() string { return m.chosen }
func (m Model) Quit() bool     { return m.quit }
func (m Model) LoadErr() error { return m.loadErr }
func (m Model) Loading() bool  { return m.loading }
