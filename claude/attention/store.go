// Package attention 维护一个粘性的「这个 session 出现过需介入的事，用户还没看过」标记。
//
// 核心语义：
//   - 触发：Reconcile 时发现 session 有 needs-input 信号 → 写入 flag
//   - 粘性：信号消失后 flag **不自动清除**
//   - 清除：用户 attach 该 session 时调 Ack(name)
//
// 与 claude/live 完全解耦：live 算实时状态、不写盘；attention 写盘、跨进程累积。
// 两者都从 ~/.claude/ 读数据，但归类逻辑、产出、生命周期独立。
package attention

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

const fileVersion = 1

// Flag 是单个 session 的粘性标记。Reason / FirstAt 用于 picker 显示「等了多久 / 因何而起」。
type Flag struct {
	FirstAt    time.Time `json:"first_at"`
	Reason     string    `json:"reason"`
	TriggerPID int       `json:"trigger_pid,omitempty"`
}

// Signal 是 Reconcile 的单条输入：表示「现在 session X 有 needs-input 类信号」。
// 由调用方（通常是 picker 主流程）从 live.Status / mcp-cache 等数据源派生。
type Signal struct {
	Reason     string
	TriggerPID int
}

// Clock 抽象时间源，便于单测。
type Clock interface{ Now() time.Time }

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// Store 持久化 flag 到 attention.json。所有方法都做了惰性懒加载 + 原子写。
// 并发安全：一个进程内单实例使用，sync.Mutex 兜底。
type Store struct {
	path  string
	clock Clock

	mu     sync.Mutex
	loaded bool
	flags  map[string]Flag // session-name -> flag
}

// New 用给定路径创建 Store。路径通常通过 DefaultPath() 得到。
func New(path string) *Store {
	return &Store{
		path:  path,
		clock: realClock{},
	}
}

// WithClock 用于单测注入假时钟。
func (s *Store) WithClock(c Clock) *Store {
	s.clock = c
	return s
}

// DefaultPath 返回 attention.json 的标准位置：$XDG_STATE_HOME/cc-sesh/attention.json
// 兜底为 $HOME/.local/state/cc-sesh/attention.json。
func DefaultPath() (string, error) {
	if v := os.Getenv("XDG_STATE_HOME"); v != "" {
		return filepath.Join(v, "cc-sesh", "attention.json"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("locate home dir: %w", err)
	}
	return filepath.Join(home, ".local", "state", "cc-sesh", "attention.json"), nil
}

// Load 返回当前所有 triggered 的 flag。读盘失败一律当空 map，不阻塞 picker。
func (s *Store) Load() map[string]Flag {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureLoaded()
	out := make(map[string]Flag, len(s.flags))
	for k, v := range s.flags {
		out[k] = v
	}
	return out
}

// Reconcile 接收当前 needs-input 信号集，写入新触发的 flag。
//
// 语义（强粘性）：
//   - 对 signals 里出现且 store 里**没有**的 session：新建 flag，FirstAt=now
//   - 对 signals 里出现且 store 里**已有**的 session：保留原 FirstAt，仅可能更新 Reason/TriggerPID
//   - 对 signals 里**未出现**的 session：保留（不清除！）
//
// activeSessionNames 用于回收幽灵 flag —— 不在该列表中的 flag 会被删除（tmux session 已消失）。
// 传 nil 表示不做回收。
func (s *Store) Reconcile(signals map[string]Signal, activeSessionNames []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureLoaded()

	now := s.clock.Now()
	changed := false

	// 1. 触发新 flag / 更新 reason
	for name, sig := range signals {
		existing, ok := s.flags[name]
		if !ok {
			s.flags[name] = Flag{FirstAt: now, Reason: sig.Reason, TriggerPID: sig.TriggerPID}
			changed = true
			continue
		}
		// 已存在：保留 FirstAt，但允许 reason 升级（auth_url 比之前的更紧迫等场景）
		if existing.Reason != sig.Reason || existing.TriggerPID != sig.TriggerPID {
			existing.Reason = sig.Reason
			existing.TriggerPID = sig.TriggerPID
			s.flags[name] = existing
			changed = true
		}
	}

	// 2. 回收幽灵 flag（session 已不存在）
	if activeSessionNames != nil {
		active := make(map[string]struct{}, len(activeSessionNames))
		for _, n := range activeSessionNames {
			active[n] = struct{}{}
		}
		for name := range s.flags {
			if _, ok := active[name]; !ok {
				delete(s.flags, name)
				changed = true
			}
		}
	}

	if !changed {
		return nil
	}
	return s.saveLocked()
}

// Ack 清除指定 session 的 flag。无 flag 时是 no-op。
func (s *Store) Ack(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureLoaded()

	if _, ok := s.flags[name]; !ok {
		return nil
	}
	delete(s.flags, name)
	return s.saveLocked()
}

// Clear 清空所有 flag。给 `cc-sesh attention clear` 子命令用。
func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags = map[string]Flag{}
	s.loaded = true
	return s.saveLocked()
}

// ensureLoaded 在 mu 锁内调用：首次访问时读盘。
func (s *Store) ensureLoaded() {
	if s.loaded {
		return
	}
	s.loaded = true
	s.flags = map[string]Flag{}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return // 不存在/读不了：当空 map
	}

	var f fileShape
	if err := json.Unmarshal(data, &f); err != nil {
		return // 损坏：当空 map，下次写入会覆盖
	}
	if f.Sessions != nil {
		s.flags = f.Sessions
	}
}

// saveLocked 在 mu 锁内调用：原子写。
func (s *Store) saveLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("mkdir state dir: %w", err)
	}

	// 排序输出以让 attention.json 变更对人友好
	keys := make([]string, 0, len(s.flags))
	for k := range s.flags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ordered := make(map[string]Flag, len(keys))
	for _, k := range keys {
		ordered[k] = s.flags[k]
	}

	data, err := json.MarshalIndent(fileShape{Version: fileVersion, Sessions: ordered}, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal attention: %w", err)
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

type fileShape struct {
	Version  int             `json:"version"`
	Sessions map[string]Flag `json:"sessions"`
}
