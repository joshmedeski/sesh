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

	"github.com/Wingsdh/cc-sesh/v2/icon"
	"github.com/Wingsdh/cc-sesh/v2/model"
)

type sessionItem struct {
	session    model.SeshSession
	name       string
	searchName string
	src        string
	decoration Decoration
}

// sessionItems 实现 fuzzy.Source，让 fuzzy 匹配只看 searchName。
type sessionItems []sessionItem

func (s sessionItems) String(i int) string { return s[i].searchName }
func (s sessionItems) Len() int            { return len(s) }

type filteredItem struct {
	item           sessionItem
	matchedIndexes []int
}

// FetchFunc 在 Init 与 mode 切换时被调用。mode 由 picker 内部按 ctrl+a/t/g/x/f 切换；
// 调用方根据 mode 决定数据源（all/tmux/config/zoxide/find）。
type FetchFunc func(mode string) (model.SeshSessions, Decorator, error)

// 五种 fetch mode 常量。
const (
	ModeAll     = "all"
	ModeTmux    = "tmux"
	ModeConfig  = "config"
	ModeZoxide  = "zoxide"
	ModeFind    = "find"
)

type sessionsLoadedMsg struct {
	sessions  model.SeshSessions
	decorator Decorator
	err       error
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
	separatorAware bool
	focusCmd       tea.Cmd
	loading        bool
	fetchFunc      FetchFunc
	loadErr        error
	decorator      Decorator
	dismisser      Dismisser
	killer         Killer
	now            func() time.Time
	mode           string // 当前 fetch mode：all/tmux/config/zoxide/find
}

// srcIcon 返回 sesh 原本的来源 icon + ANSI 颜色。
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

var separatorReplacer = strings.NewReplacer("-", " ", "_", " ", "/", " ", "\\", " ")

func normalizeSeparators(s string) string {
	return separatorReplacer.Replace(s)
}

func buildItems(sessions model.SeshSessions, dec Decorator, separatorAware bool) sessionItems {
	if dec == nil {
		dec = NoDecoration{}
	}
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
			decoration: dec.Decorate(s),
		})
	}
	return items
}

func New(fetchFunc FetchFunc, dec Decorator, dis Dismisser, kil Killer, showIcons, separatorAware bool, prompt, placeholder string) Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Prompt = prompt

	if dec == nil {
		dec = NoDecoration{}
	}

	m := Model{
		filterInput:    ti,
		showIcons:      showIcons,
		separatorAware: separatorAware,
		loading:        true,
		fetchFunc:      fetchFunc,
		decorator:      dec,
		dismisser:      dis,
		killer:         kil,
		now:            time.Now,
		mode:           ModeAll,
	}
	m.focusCmd = m.filterInput.Focus()
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.focusCmd, m.fetchSessions())
}

func (m Model) fetchSessions() tea.Cmd {
	mode := m.mode
	return func() tea.Msg {
		sessions, dec, err := m.fetchFunc(mode)
		return sessionsLoadedMsg{sessions: sessions, decorator: dec, err: err}
	}
}

// switchMode 切换数据源并触发异步 fetch。
func (m *Model) switchMode(mode string) tea.Cmd {
	if m.mode == mode {
		return nil
	}
	m.mode = mode
	m.loading = true
	m.cursor = 0
	m.offset = 0
	m.filterInput.SetValue("")
	return m.fetchSessions()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sessionsLoadedMsg:
		if msg.err != nil {
			m.loadErr = msg.err
			return m, tea.Quit
		}
		m.loading = false
		if msg.decorator != nil {
			m.decorator = msg.decorator
		}
		m.allItems = buildItems(msg.sessions, m.decorator, m.separatorAware)
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

		case "up", "ctrl+k", "shift+tab":
			m.cursorUp(1)
			return m, nil

		case "down", "ctrl+j", "tab":
			m.cursorDown(1)
			return m, nil

		case "ctrl+u":
			m.cursorUp(m.visibleCount() / 2)
			return m, nil

		case "pgdown":
			m.cursorDown(m.visibleCount() / 2)
			return m, nil

		case "pgup":
			m.cursorUp(m.visibleCount() / 2)
			return m, nil

		case "ctrl+a":
			return m, m.switchMode(ModeAll)

		case "ctrl+t":
			return m, m.switchMode(ModeTmux)

		case "ctrl+g":
			return m, m.switchMode(ModeConfig)

		case "ctrl+x":
			return m, m.switchMode(ModeZoxide)

		case "ctrl+f":
			return m, m.switchMode(ModeFind)

		case "ctrl+d":
			// kill 当前 cursor 所指 tmux session（与 fzf-tmux 习惯一致）。
			// session 不存在后 attention 也会被自然 GC，无需额外清理。
			m.killCurrent()
			return m, nil

		case "alt+d":
			// 手动 dismiss 当前 attention 行。alt+d 避开和搜索字符 'd' 冲突。
			m.dismissCurrent()
			return m, nil
		}
	}

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

// killCurrent kill 当前 cursor 所指的 tmux session。仅对 src=tmux 有效；
// 其他类型（zoxide / config / tmuxinator template）没真实 session 可 kill，no-op。
func (m *Model) killCurrent() {
	if m.killer == nil || len(m.filtered) == 0 {
		return
	}
	cur := m.filtered[m.cursor].item
	if cur.session.Src != "tmux" {
		return
	}
	if err := m.killer.Kill(cur.name); err != nil {
		return
	}
	// 从 allItems 移除被 kill 的项，让列表立即同步
	newItems := make(sessionItems, 0, len(m.allItems))
	for _, it := range m.allItems {
		if it.name == cur.name && it.session.Src == "tmux" {
			continue
		}
		newItems = append(newItems, it)
	}
	m.allItems = newItems
	m.applyFilter()
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
}

// dismissCurrent 在 cursor 落在 attention 行时清除该行的 flag，并重新装饰所有项。
func (m *Model) dismissCurrent() {
	if m.dismisser == nil || len(m.filtered) == 0 {
		return
	}
	cur := m.filtered[m.cursor].item
	if !cur.decoration.Attention.Triggered {
		return
	}
	if err := m.dismisser.Dismiss(cur.name); err != nil {
		return
	}
	// 重新装饰：对所有 allItems 再调一次 decorator.Decorate
	for i := range m.allItems {
		m.allItems[i].decoration = m.decorator.Decorate(m.allItems[i].session)
	}
	m.applyFilter()
}

func (m *Model) applyFilter() {
	pattern := m.filterInput.Value()

	var matches []fuzzy.Match
	if pattern != "" {
		searchPat := pattern
		if m.separatorAware {
			searchPat = normalizeSeparators(pattern)
		}
		matches = fuzzy.FindFrom(searchPat, m.allItems)
	}

	if pattern == "" {
		m.filtered = make([]filteredItem, len(m.allItems))
		for i, item := range m.allItems {
			m.filtered[i] = filteredItem{item: item}
		}
	} else {
		m.filtered = make([]filteredItem, len(matches))
		for i, match := range matches {
			m.filtered[i] = filteredItem{
				item:           m.allItems[match.Index],
				matchedIndexes: match.MatchedIndexes,
			}
		}
	}

	// attention 项稳定排序到前面
	sort.SliceStable(m.filtered, func(i, j int) bool {
		ai := m.filtered[i].item.decoration.Attention.Triggered
		aj := m.filtered[j].item.decoration.Attention.Triggered
		if ai != aj {
			return ai
		}
		return false
	})

	if m.cursor >= len(m.filtered) {
		m.cursor = 0
		m.offset = 0
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
	// chrome: filter(1) + header(1) + 2 blank + section overhead(~5)
	chrome := 10
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

// 颜色调色：尽量用 ANSI 数字，确保跨终端一致。
var (
	colorCursor   = lipgloss.ANSIColor(2)  // green
	colorAttenFg  = lipgloss.ANSIColor(9)  // bright red
	colorAttenBg  = lipgloss.ANSIColor(52) // 深红 256-color；不支持时退化为默认
	colorBusy     = lipgloss.ANSIColor(12) // bright blue
	colorSubagent = lipgloss.ANSIColor(11) // bright yellow
	colorNeeding  = lipgloss.ANSIColor(9)  // bright red
	colorIdle     = lipgloss.ANSIColor(8)  // bright black / dim
	colorMatch    = lipgloss.ANSIColor(1)  // red
	colorTail     = lipgloss.ANSIColor(8)
	colorHeader   = lipgloss.ANSIColor(8)
)

func (m Model) View() tea.View {
	var b strings.Builder

	b.WriteString("  " + m.filterInput.View())
	b.WriteString("\n")
	b.WriteString(renderHotkeyHeader(m.mode))
	b.WriteString("\n")

	visible := m.visibleCount()

	if m.loading {
		loadingStyle := lipgloss.NewStyle().Faint(true)
		b.WriteString(loadingStyle.Render("  Loading sessions..."))
		b.WriteString("\n")
		for i := 1; i < visible; i++ {
			b.WriteString("\n")
		}
	} else {
		end := m.offset + visible
		if end > len(m.filtered) {
			end = len(m.filtered)
		}

		// 计算 needs-you 数；> 0 时第一行渲染分组标题（如果在可视范围内）。
		needsCount := 0
		for _, fi := range m.filtered {
			if fi.item.decoration.Attention.Triggered {
				needsCount++
			}
		}

		linesPrinted := 0
		printedHeader := needsCount == 0 // 没 attention 时不显示
		printedDivider := false

		for i := m.offset; i < end; i++ {
			fi := m.filtered[i]

			// 在第一条 attention 行前插 "needs you" header
			if !printedHeader && fi.item.decoration.Attention.Triggered {
				b.WriteString(renderHeader(fmt.Sprintf("needs you (%d)", needsCount)))
				b.WriteString("\n")
				printedHeader = true
				linesPrinted++
			}
			// 在第一条非 attention 行前插分割线（前提是 attention 区有内容）
			if needsCount > 0 && !printedDivider && !fi.item.decoration.Attention.Triggered {
				b.WriteString(renderDivider())
				b.WriteString("\n")
				printedDivider = true
				linesPrinted++
			}

			b.WriteString(m.renderRow(fi, i == m.cursor))
			b.WriteString("\n")
			linesPrinted++
		}

		for i := linesPrinted; i < visible; i++ {
			b.WriteString("\n")
		}
	}

	return tea.NewView(b.String())
}

func renderHeader(label string) string {
	return lipgloss.NewStyle().Foreground(colorHeader).Faint(true).Render("  ─── " + label + " ──")
}

// renderHotkeyHeader 在搜索框下方显示模式切换的 hotkey 提示，当前 mode 高亮。
func renderHotkeyHeader(mode string) string {
	dim := lipgloss.NewStyle().Foreground(colorHeader).Faint(true)
	hi := lipgloss.NewStyle().Foreground(colorCursor).Bold(true)
	pick := func(name, label string) string {
		if mode == name {
			return hi.Render(label)
		}
		return dim.Render(label)
	}
	parts := []string{
		pick(ModeAll, "^a all"),
		pick(ModeTmux, "^t tmux"),
		pick(ModeConfig, "^g configs"),
		pick(ModeZoxide, "^x zoxide"),
		pick(ModeFind, "^f find"),
		dim.Render("^d kill"),
	}
	return "  " + strings.Join(parts, dim.Render("  "))
}

func renderDivider() string {
	return lipgloss.NewStyle().Foreground(colorHeader).Faint(true).Render("  ───────────────")
}

func (m Model) renderRow(fi filteredItem, isCursor bool) string {
	item := fi.item
	dec := item.decoration

	// 1. cursor 列
	cursorPrefix := "  "
	if isCursor {
		cursorPrefix = lipgloss.NewStyle().Foreground(colorCursor).Bold(true).Render("> ")
	}

	// 2. attention 列：⚠ + 空格，或两个空格
	attenCol := "  "
	if dec.Attention.Triggered {
		attenCol = lipgloss.NewStyle().Foreground(colorAttenFg).Bold(true).Render("⚠ ")
	}

	// 3. src icon 列
	var srcCol string
	if m.showIcons {
		icn, clr := srcIcon(item.src)
		srcCol = lipgloss.NewStyle().Foreground(clr).Render(icn)
	}

	// 4. live badge 列
	liveCol := renderLiveBadge(dec.Live)

	// 5. name 列（fuzzy 高亮）
	nameStyle := lipgloss.NewStyle()
	matchStyle := lipgloss.NewStyle().Foreground(colorMatch).Bold(true)
	name := highlightMatches(item.name, fi.matchedIndexes, matchStyle, nameStyle)

	// 6. tail 列
	tail := m.renderTail(dec)

	body := cursorPrefix + attenCol + srcCol + liveCol + name
	if tail != "" {
		body += "  " + tail
	}

	// attention 行整行染暗红色，区别于普通行
	if dec.Attention.Triggered {
		body = lipgloss.NewStyle().Foreground(colorAttenFg).Render(body)
	}

	return body
}

// renderLiveBadge 输出 8 字符宽的列：`● 2/2   ` / `○ 1     ` / `        `
func renderLiveBadge(b LiveBadge) string {
	if b.IsEmpty() {
		return strings.Repeat(" ", 8)
	}

	var glyph rune
	var clr color.Color
	switch b.Severity() {
	case SevNeeding:
		glyph = '⚠'
		clr = colorNeeding
	case SevBusy:
		glyph = '●'
		clr = colorBusy
	case SevSubagent:
		glyph = '◐'
		clr = colorSubagent
	default:
		glyph = '○'
		clr = colorIdle
	}

	// 计数显示：idle 时只显示总数；其他时显示 active/total
	var counter string
	switch b.Severity() {
	case SevIdle:
		counter = fmt.Sprintf("%d", b.Total)
	case SevSubagent:
		counter = fmt.Sprintf("%d/%d", b.Subagent, b.Total)
	case SevBusy:
		counter = fmt.Sprintf("%d/%d", b.Busy, b.Total)
	case SevNeeding:
		counter = fmt.Sprintf("%d/%d", b.Needing, b.Total)
	}

	plain := fmt.Sprintf("%c %s", glyph, counter)
	if width := lipgloss.Width(plain); width < 8 {
		plain += strings.Repeat(" ", 8-width)
	}
	return lipgloss.NewStyle().Foreground(clr).Render(plain)
}

func (m Model) renderTail(dec Decoration) string {
	if !dec.Attention.Triggered {
		return ""
	}
	dur := durationShort(m.now().Sub(dec.Attention.FirstAt))
	reason := dec.Attention.Reason
	tailStyle := lipgloss.NewStyle().Foreground(colorTail).Faint(true)

	if dec.Live.IsEmpty() || dec.Live.Severity() == SevIdle {
		// "幽灵" 行：信号已消失但 attention 还在
		return tailStyle.Render(fmt.Sprintf("was: %s · %s · resolved", reason, dur))
	}
	return tailStyle.Render(fmt.Sprintf("%s · %s", reason, dur))
}

func durationShort(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
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
