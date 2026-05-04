package live

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// rawSession 是 ~/.claude/sessions/<pid>.json 的字段子集。
type rawSession struct {
	PID       int    `json:"pid"`
	SessionID string `json:"sessionId"`
	Cwd       string `json:"cwd"`
	Status    string `json:"status"`
	Kind      string `json:"kind"`
	UpdatedAt int64  `json:"updatedAt"`
}

// ProcessChecker 用于探活 pid，便于单测注入假实现。
type ProcessChecker interface {
	IsAlive(pid int) bool
}

// Instance 是单个 Claude Code 进程经过分类后的视图。
type Instance struct {
	PID       int
	SessionID string
	Cwd       string
	Logical   Logical
}

// Reader 从 ~/.claude/sessions/ 扫出所有活实例并按 cwd 聚合。
type Reader struct {
	sessionsDir string
	proc        ProcessChecker
}

func NewReader(homeDir string, proc ProcessChecker) *Reader {
	return &Reader{
		sessionsDir: filepath.Join(homeDir, ".claude", "sessions"),
		proc:        proc,
	}
}

// Read 扫描、过滤、归类、聚合，返回 cwd -> Status 的映射。
// 适用于 zoxide / config 类 SeshSession（这些 session 还未 attach，只有 path 信息）。
func (r *Reader) Read() (map[string]Status, error) {
	instances, err := r.scan()
	if err != nil {
		return nil, err
	}
	return aggregateByCwd(instances), nil
}

// ReadInstances 返回所有活实例（不聚合），让上层按需关联到 tmux pane / session。
func (r *Reader) ReadInstances() ([]Instance, error) {
	return r.scan()
}

// PaneInfo 是聚合 Claude 实例到 tmux session 时所需的最小输入。
// 由调用方（通常通过 tmux.ListAllPanes）提供。
type PaneInfo struct {
	SessionName string
	Cwd         string
}

// AggregateBySession 把 Claude 实例按 cwd 与 pane 关联，归到对应 tmux session。
// 一个 instance 只归一个 session（按 panes 顺序遇到的第一个 cwd 匹配）。
// 跟 aggregateByCwd 互不冲突：调用方对 tmux 类 session 用本函数，对 zoxide/config 用 Read()。
func AggregateBySession(instances []Instance, panes []PaneInfo) map[string]Status {
	out := make(map[string]Status)
	if len(panes) == 0 {
		return out
	}
	cwdToSession := make(map[string][]string)
	for _, p := range panes {
		cwd := NormalizeCwd(p.Cwd)
		if cwd == "" {
			continue
		}
		cwdToSession[cwd] = append(cwdToSession[cwd], p.SessionName)
	}
	for _, it := range instances {
		sessions, ok := cwdToSession[it.Cwd]
		if !ok {
			continue
		}
		// 同一 cwd 多个 pane / 多个 session：每个 session 各计一份
		// （罕见，但保留正确性）
		seen := make(map[string]struct{}, len(sessions))
		for _, name := range sessions {
			if _, dup := seen[name]; dup {
				continue
			}
			seen[name] = struct{}{}
			s := out[name]
			s.Total++
			switch it.Logical {
			case LogicalNeedsInput:
				s.Needing++
			case LogicalBusy:
				s.Busy++
			case LogicalSubagent:
				s.Subagent++
			}
			out[name] = s
		}
	}
	return out
}

func (r *Reader) scan() ([]Instance, error) {
	entries, err := os.ReadDir(r.sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read sessions dir: %w", err)
	}

	out := make([]Instance, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		fullPath := filepath.Join(r.sessionsDir, e.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		var raw rawSession
		if err := json.Unmarshal(data, &raw); err != nil {
			continue
		}
		if raw.PID == 0 {
			continue
		}
		if !r.proc.IsAlive(raw.PID) {
			continue
		}
		out = append(out, Instance{
			PID:       raw.PID,
			SessionID: raw.SessionID,
			Cwd:       NormalizeCwd(raw.Cwd),
			Logical:   classify(raw.Status, raw.Kind),
		})
	}
	return out, nil
}

func aggregateByCwd(items []Instance) map[string]Status {
	out := make(map[string]Status, len(items))
	for _, it := range items {
		s := out[it.Cwd]
		s.Total++
		switch it.Logical {
		case LogicalNeedsInput:
			s.Needing++
		case LogicalBusy:
			s.Busy++
		case LogicalSubagent:
			s.Subagent++
		}
		out[it.Cwd] = s
	}
	return out
}

// NormalizeCwd 让 lister 侧的 SeshSession.Path 与 Claude 写的 cwd 走同一种归一化。
func NormalizeCwd(p string) string {
	if p == "" {
		return ""
	}
	p = filepath.Clean(p)
	if len(p) > 1 && strings.HasSuffix(p, string(filepath.Separator)) {
		p = strings.TrimRight(p, string(filepath.Separator))
	}
	return p
}
