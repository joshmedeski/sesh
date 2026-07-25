package picker

import (
	"fmt"
	"image/color"
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
}

// sessionItems implements fuzzy.Source for fuzzy matching.
type sessionItems []sessionItem

func (s sessionItems) String(i int) string { return s[i].searchName }
func (s sessionItems) Len() int            { return len(s) }

type filteredItem struct {
	item           sessionItem
	matchedIndexes []int
}

// FetchFunc loads sessions asynchronously. It is called in a goroutine by Init().
type FetchFunc func() (model.SeshSessions, error)

// sessionsLoadedMsg carries the result of the async fetch back to Update().
type sessionsLoadedMsg struct {
	sessions model.SeshSessions
	err      error
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
	// Aliases is keyed by lowercased alias.
	Aliases map[string]Alias
	// AliasAutoConnectDelay is the grace period before auto-connect fires,
	// leaving room to finish typing a longer alias that shares a prefix.
	AliasAutoConnectDelay time.Duration
	// DisableAliasAutoConnect suppresses auto-connect for this invocation.
	DisableAliasAutoConnect bool
}

type Model struct {
	allItems       sessionItems
	filtered       []filteredItem
	filterInput    textinput.Model
	cursor         int
	offset         int
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

	// aliases is keyed by lowercased alias; aliasByName is keyed by session
	// name so rows can be tagged while rendering.
	aliases                 map[string]Alias
	aliasByName             map[string]string
	aliasAutoConnectDelay   time.Duration
	disableAliasAutoConnect bool
	// aliasSeq increments on every filter change so a pending auto-connect
	// tick can tell whether it is stale.
	aliasSeq int
}

// srcIcon returns the nerd font icon and color for a session source.
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

// Powerline half circles used to round off the alias chip. They are nerd font
// glyphs, so they are only used when icons are enabled.
const (
	chipLeftGlyph  = ""
	chipRightGlyph = ""
)

// aliasChip renders the alias of a session as a chip that sits before the
// session name, e.g. ` wp ` with nerd fonts or `[wp] ` without. It returns ""
// for sessions that have no alias.
//
// The label uses reverse video rather than an explicit fg/bg pair so the text
// is always the terminal's background color on its foreground color, which
// stays legible under any color scheme. The half circles are left unstyled for
// the same reason: their default foreground is exactly the color the reversed
// label fills with, so the chip reads as one shape.
func (m Model) aliasChip(name string) string {
	alias, ok := m.aliasByName[name]
	if !ok {
		return ""
	}

	label := lipgloss.NewStyle().Reverse(true)
	if !m.showIcons {
		// Without nerd fonts the half circles render as tofu, so fall back to
		// brackets that still read as a chip.
		return label.Render("["+alias+"]") + " "
	}
	return chipLeftGlyph + label.Render(alias) + chipRightGlyph + " "
}

var separatorReplacer = strings.NewReplacer("-", " ", "_", " ", "/", " ", "\\", " ")

func normalizeSeparators(s string) string {
	return separatorReplacer.Replace(s)
}

func buildItems(sessions model.SeshSessions, separatorAware bool) sessionItems {
	items := make(sessionItems, 0, len(sessions.OrderedIndex))
	for _, key := range sessions.OrderedIndex {
		s := sessions.Directory[key]
		searchName := s.Name
		if separatorAware {
			searchName = normalizeSeparators(s.Name)
		}
		items = append(items, sessionItem{
			session:    s,
			name:       s.Name,
			searchName: searchName,
			src:        s.Src,
		})
	}
	return items
}

func New(fetchFunc FetchFunc, opts Options) Model {
	ti := textinput.New()
	ti.Placeholder = opts.Placeholder
	ti.Prompt = opts.Prompt

	aliasByName := make(map[string]string, len(opts.Aliases))
	for _, alias := range opts.Aliases {
		aliasByName[alias.Target] = alias.Alias
	}

	m := Model{
		filterInput:             ti,
		showIcons:               opts.ShowIcons,
		showWindows:             opts.ShowWindows,
		separatorAware:          opts.SeparatorAware,
		loading:                 true,
		fetchFunc:               fetchFunc,
		aliases:                 opts.Aliases,
		aliasByName:             aliasByName,
		aliasAutoConnectDelay:   opts.AliasAutoConnectDelay,
		disableAliasAutoConnect: opts.DisableAliasAutoConnect,
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case aliasAutoConnectMsg:
		// The tick can't be cancelled, so a stale one still arrives: ignore it
		// unless it is the latest and the filter still reads as that alias.
		if msg.seq != m.aliasSeq {
			return m, nil
		}
		alias, ok := m.aliases[strings.ToLower(m.filterInput.Value())]
		if !ok || alias.Alias != msg.alias || !alias.AutoConnect || m.disableAliasAutoConnect {
			return m, nil
		}
		m.chosen = alias.Target
		return m, tea.Quit

	case sessionsLoadedMsg:
		if msg.err != nil {
			m.loadErr = msg.err
			return m, tea.Quit
		}
		m.loading = false
		m.allItems = buildItems(msg.sessions, m.separatorAware)
		m.applyFilter()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.filterInput.SetWidth(m.contentWidth() - 4)
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
		if tick := m.scheduleAliasAutoConnect(); tick != nil {
			return m, tea.Batch(cmd, tick)
		}
	}

	return m, cmd
}

// scheduleAliasAutoConnect invalidates any pending auto-connect and, when the
// filter now reads as an alias that opted into auto-connect, arms a fresh one.
// It returns nil when there is nothing to schedule. Auto-connect deliberately
// works while sessions are still loading: the target comes from the config, not
// from the list, so there is no reason to make the user wait.
func (m *Model) scheduleAliasAutoConnect() tea.Cmd {
	m.aliasSeq++

	if m.disableAliasAutoConnect {
		return nil
	}
	// Matched against the raw input: separator-aware normalization must not
	// turn `w p` into a match for the alias `w-p`.
	alias, ok := m.aliases[strings.ToLower(m.filterInput.Value())]
	if !ok || !alias.AutoConnect {
		return nil
	}

	msg := aliasAutoConnectMsg{seq: m.aliasSeq, alias: alias.Alias}
	if m.aliasAutoConnectDelay <= 0 {
		return func() tea.Msg { return msg }
	}
	return tea.Tick(m.aliasAutoConnectDelay, func(time.Time) tea.Msg { return msg })
}

func (m *Model) applyFilter() {
	pattern := m.filterInput.Value()
	if pattern == "" {
		m.filtered = make([]filteredItem, len(m.allItems))
		for i, item := range m.allItems {
			m.filtered[i] = filteredItem{item: item}
		}
		return
	}

	if m.separatorAware {
		pattern = normalizeSeparators(pattern)
	}

	matches := fuzzy.FindFrom(pattern, m.allItems)
	m.filtered = make([]filteredItem, len(matches))
	for i, match := range matches {
		m.filtered[i] = filteredItem{
			item:           m.allItems[match.Index],
			matchedIndexes: match.MatchedIndexes,
		}
	}
}

func (m *Model) cursorUp(n int) {
	m.cursor -= n
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
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
	visible := m.visibleCount()
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
}

func (m Model) visibleCount() int {
	// border(2) + title(1) + blank(1) + filter(1) + blank(1) + counter(1) + help(1) + blank before counter(1)
	chrome := 9
	available := m.height - chrome
	if available < 1 {
		available = 5
	}
	if available > 15 {
		available = 15
	}
	return available
}

func (m Model) contentWidth() int {
	w := m.width
	if w < 30 {
		w = 40
	}
	if w > 60 {
		w = 60
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
	} else {
		// Session list
		end := m.offset + visible
		if end > len(m.filtered) {
			end = len(m.filtered)
		}

		cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)
		matchStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1)).Bold(true)
		normalStyle := lipgloss.NewStyle()
		windowStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8)).Faint(true)

		for i := m.offset; i < end; i++ {
			item := m.filtered[i]
			prefix := "  "
			if i == m.cursor {
				prefix = cursorStyle.Render("> ")
			}

			var tag string
			if m.showIcons {
				icn, clr := srcIcon(item.item.src)
				iconStyle := lipgloss.NewStyle().Foreground(clr)
				tag = iconStyle.Render(icn)
			}
			chip := m.aliasChip(item.item.name)
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

	content := b.String()

	return tea.NewView(content)
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

	matchSet := make(map[int]bool, len(indexes))
	for _, idx := range indexes {
		matchSet[idx] = true
	}

	var result strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		ch := string(r)
		if matchSet[i] {
			result.WriteString(matchStyle.Render(ch))
		} else {
			result.WriteString(normalStyle.Render(ch))
		}
	}
	return result.String()
}

func (m Model) Chosen() string { return m.chosen }
func (m Model) Quit() bool     { return m.quit }
func (m Model) LoadErr() error { return m.loadErr }
func (m Model) Loading() bool  { return m.loading }
