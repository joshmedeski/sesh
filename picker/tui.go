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
	// chrome: filter(1) + hotkey(1) + table-top(1) + header(1) + 2 blank + section overhead(~5)
	chrome := 12
	available := m.height - chrome
	if available < 1 {
		available = 3
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
	colorCursor   = lipgloss.ANSIColor(2)   // green
	colorAttn     = lipgloss.ANSIColor(208) // 256-color 橙；ATTN 列圆点
	colorBusy     = lipgloss.ANSIColor(12)  // bright blue
	colorSubagent = lipgloss.ANSIColor(11)  // bright yellow
	colorNeeding  = lipgloss.ANSIColor(11)  // bright yellow（WAIT 列）
	colorRun      = lipgloss.ANSIColor(14)  // bright cyan（RUN 列）
	colorIdle     = lipgloss.ANSIColor(8)   // bright black / dim
	colorMatch    = lipgloss.ANSIColor(1)   // red
	colorTail     = lipgloss.ANSIColor(8)
	colorHeader   = lipgloss.ANSIColor(8)
)

func (m Model) View() tea.View {
	var b strings.Builder

	b.WriteString("  " + m.filterInput.View())
	b.WriteString("\n")
	b.WriteString(renderHotkeyHeader(m.mode))
	b.WriteString("\n")
	b.WriteString(renderTableTop(m.showIcons, m.contentWidth()))
	b.WriteString("\n")
	b.WriteString(renderColumnHeaders(m.showIcons))
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

// 表格列宽常量。每行的 ATTN/IDLE/RUN/WAIT 4 列，前 3 列宽 5（含右侧分隔空格），
// 最后一列宽 4。Header 字符串与行内 cell 必须使用相同宽度，确保竖向对齐。
const (
	colCellWidth     = 5
	colLastCellWidth = 4
	colsTotalWidth   = colCellWidth*3 + colLastCellWidth // 19
)

// 表格列布局：每个数字 cell 4 字符宽（居中对齐），后跟 1 字符分隔（最后列无尾分隔）。
const colNumWidth = 4

// renderTableTop 输出表格上方的横线，把表格区与 hotkey 区视觉分隔开。
// 横线长度 = contentWidth - leftPad，覆盖整张表格（含 name 区）。
func renderTableTop(showIcons bool, contentWidth int) string {
	leftPad := 2
	if showIcons {
		leftPad += 2
	}
	lineLen := contentWidth - leftPad
	if lineLen < colsTotalWidth {
		lineLen = colsTotalWidth
	}
	style := lipgloss.NewStyle().Foreground(colorHeader).Faint(true)
	return strings.Repeat(" ", leftPad) + style.Render(strings.Repeat("─", lineLen))
}

// renderColumnHeaders 输出表格列标题，与每行的 4 列徽章列竖向对齐。
//
// 左侧 padding = cursor(2) + src_icon(showIcons ? 2 : 0)，刚好对齐到行内 ATTN 列起点。
// 标题用 bold 默认前景色，比 dim 更醒目。
func renderColumnHeaders(showIcons bool) string {
	leftPad := 2
	if showIcons {
		leftPad += 2
	}
	style := lipgloss.NewStyle().Bold(true)
	cell := func(label string, last bool) string {
		s := style.Width(colNumWidth).Align(lipgloss.Center).Render(label)
		if !last {
			s += " "
		}
		return s
	}
	return strings.Repeat(" ", leftPad) +
		cell("ATTN", false) +
		cell("IDLE", false) +
		cell("RUN", false) +
		cell("WAIT", true)
}

// renderRowCounts 渲染单行的 4 列徽章数字（19 字符宽）。
//
//   - ATTN：橙底白字 ⚠ 色块（粘性提醒）；未触发为空白
//   - IDLE / RUN / WAIT：右对齐数字。0 dim、非零彩色加粗。
//   - 整行无 Claude（src 非 tmux 等）→ 全列空白对齐
func renderRowCounts(dec Decoration) string {
	// 整体空白：非 tmux session（无 live、无 attention）
	if dec.Live.IsEmpty() && !dec.Attention.Triggered {
		return strings.Repeat(" ", colsTotalWidth)
	}

	// ATTN 列：橙色圆点居中（与 IDLE/RUN/WAIT 同款 cell 形状）
	var attnCell string
	if dec.Attention.Triggered {
		dot := lipgloss.NewStyle().
			Foreground(colorAttn).
			Bold(true).
			Width(colNumWidth).
			Align(lipgloss.Center).
			Render("●")
		attnCell = dot + " "
	} else {
		attnCell = strings.Repeat(" ", colCellWidth)
	}

	// IDLE/RUN/WAIT 列：数字居中对齐 4 宽 + 右侧 1 空格分隔（最后列无尾空格）
	cell := func(n int, clr lipgloss.ANSIColor, last bool) string {
		width := colCellWidth
		if last {
			width = colLastCellWidth
		}
		if dec.Live.IsEmpty() {
			return strings.Repeat(" ", width)
		}
		txt := fmt.Sprintf("%d", n)
		st := lipgloss.NewStyle().Width(colNumWidth).Align(lipgloss.Center)
		if n == 0 {
			st = st.Foreground(colorIdle).Faint(true)
		} else {
			st = st.Foreground(clr).Bold(true)
		}
		styled := st.Render(txt)
		if !last {
			styled += " "
		}
		return styled
	}

	return attnCell +
		cell(dec.Live.Idle(), colorIdle, false) +
		cell(dec.Live.Busy+dec.Live.Subagent, colorRun, false) +
		cell(dec.Live.Needing, colorNeeding, true)
}

func (m Model) renderRow(fi filteredItem, isCursor bool) string {
	item := fi.item
	dec := item.decoration

	// 1. cursor 列（2 字符）
	cursorPrefix := "  "
	if isCursor {
		cursorPrefix = lipgloss.NewStyle().Foreground(colorCursor).Bold(true).Render("> ")
	}

	// 2. src icon 列（2 字符；showIcons=false 时省略）
	var srcCol string
	if m.showIcons {
		icn, clr := srcIcon(item.src)
		srcCol = lipgloss.NewStyle().Foreground(clr).Render(icn)
	}

	// 3. ATTN/IDLE/RUN/WAIT 4 列徽章（19 字符）
	countsCol := renderRowCounts(dec)

	// 4. name 列（fuzzy 高亮）
	nameStyle := lipgloss.NewStyle()
	matchStyle := lipgloss.NewStyle().Foreground(colorMatch).Bold(true)
	name := highlightMatches(item.name, fi.matchedIndexes, matchStyle, nameStyle)

	// 5. tail 列
	tail := m.renderTail(dec)

	body := cursorPrefix + srcCol + countsCol + " " + name
	if tail != "" {
		body += "  " + tail
	}
	return body
}

// renderTail 在 attention 行末尾显示「完成多久」提示，便于用户判断紧迫度。
func (m Model) renderTail(dec Decoration) string {
	if !dec.Attention.Triggered {
		return ""
	}
	dur := durationShort(m.now().Sub(dec.Attention.FirstAt))
	return lipgloss.NewStyle().Foreground(colorTail).Faint(true).
		Render(fmt.Sprintf("done %s ago", dur))
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
