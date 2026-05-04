package picker

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Wingsdh/cc-sesh/v2/model"
)

// TestVisualPreview 是一个**手动预览**测试：默认 skip，仅当 PREVIEW=1 时执行。
// 用法：PREVIEW=1 go test -v -run TestVisualPreview ./picker/...
//
// 直接把一帧 picker UI 输出到 stdout，让人眼看 ANSI 配色与对齐效果。
func TestVisualPreview(t *testing.T) {
	if os.Getenv("PREVIEW") == "" {
		t.Skip("set PREVIEW=1 to render a sample frame")
	}

	m := Model{
		showIcons: true,
		mode:      ModeAll,
		width:     60,
		height:    30,
		now:       time.Now,
	}

	now := time.Now()

	// 模拟 4 种典型 session 状态，覆盖各列组合
	rows := []struct {
		name      string
		src       string
		idle      int
		busy      int
		subagent  int
		needing   int
		attn      bool
		attnFirst time.Time
	}{
		// 真实 tmux + 1 个 idle Claude（最常见）
		{name: "default", src: "tmux", idle: 1},
		// 真实 tmux 但无 Claude
		{name: "bay-translate-extension", src: "tmux"},
		// 跑完一轮活的提醒 + 当前 idle
		{name: "ai-dev-kit", src: "tmux", idle: 1, attn: true, attnFirst: now.Add(-15 * time.Minute)},
		// 在跑活
		{name: "long-running-task", src: "tmux", busy: 2},
		// 等权限（OAuth/permission prompt）
		{name: "oauth-flow", src: "tmux", needing: 1},
		// 同时多种状态
		{name: "mixed", src: "tmux", busy: 1, subagent: 1, idle: 1},
		// zoxide 类（非 tmux 行 → 全列空白对齐）
		{name: "~/Code/backend/athena", src: "zoxide"},
		{name: "~/AI-Workspace/bay-translate", src: "zoxide"},
	}

	items := make(sessionItems, 0, len(rows))
	for _, r := range rows {
		dec := Decoration{
			Live: LiveBadge{
				Total:    r.idle + r.busy + r.subagent + r.needing,
				Busy:     r.busy,
				Subagent: r.subagent,
				Needing:  r.needing,
			},
		}
		if r.attn {
			dec.Attention = AttentionBadge{Triggered: true, FirstAt: r.attnFirst}
		}
		items = append(items, sessionItem{
			session:    model.SeshSession{Name: r.name, Src: r.src},
			name:       r.name,
			searchName: r.name,
			src:        r.src,
			decoration: dec,
		})
	}
	m.allItems = items
	m.filtered = make([]filteredItem, len(items))
	for i, it := range items {
		m.filtered[i] = filteredItem{item: it}
	}

	fmt.Println()
	fmt.Println("  > Filter Sessions...")
	fmt.Println(renderHotkeyHeader(m.mode))
	fmt.Println(renderTableTop(m.showIcons, m.contentWidth()))
	fmt.Println(renderColumnHeaders(m.showIcons))
	for i, fi := range m.filtered {
		fmt.Println(m.renderRow(fi, i == 0)) // 第一行带 cursor
	}
	fmt.Println()
}
