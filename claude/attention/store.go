// Package attention 维护一个粘性的「session 跑完一轮活了，用户还没看」标记。
//
// 核心语义（v2，跟 Claude 实时状态完全解耦）：
//   - 触发：Reconcile 观察到 session 发生 active(busy/subagent) → idle 的转换 → 写入 flag
//   - 粘性：flag 写入后**不自动清除**
//   - 清除：用户 attach 该 session（Ack）/ session 消失（GC）/ 手动 Clear
//
// 跟 needs-input（OAuth、permission prompt）无关 —— 那个走 live 包的实时统计，
// 不再走粘性提醒。本包只回答："有几个 session 跑完一轮活，等着你 review？"
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

const fileVersion = 2

// Flag 是单个 session 的粘性 attention 标记。
type Flag struct {
	FirstAt time.Time `json:"first_at"`
}

// Clock 抽象时间源，便于单测。
type Clock interface{ Now() time.Time }

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// Store 持久化 flag + 上次观察到的 busy 状态（用于检测 active→idle 转换）。
// 并发安全：一个进程内单实例使用，sync.Mutex 兜底。
type Store struct {
	path  string
	clock Clock

	mu       sync.Mutex
	loaded   bool
	flags    map[string]Flag // session-name -> 已触发的 attention flag
	tracking map[string]bool // session-name -> 上次 Reconcile 观察到 busy/subagent
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

// Load 返回当前所有触发中的 attention flag。读盘失败一律当空 map。
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

// Reconcile 输入当前活实例的 busy 状态，按转换语义更新 flag。
//
// busyByName 的语义：
//   - busyByName[name] = true：该 session 当前有 Claude 在 busy/subagent（在跑活）
//   - busyByName[name] = false 或 missing：该 session 当前没在跑活（idle / 仅 needs-input / 无实例）
//
// 触发规则：
//   - tracking 里记录 true 且本轮 busyByName[name] = false → 触发 flag
//   - 已存在 flag 时不更新 FirstAt（保留首次时刻）
//
// activeSessionNames 用于回收幽灵数据 —— 不在该列表中的 flag/tracking 都被删除。
// 传 nil 表示不做 GC。
func (s *Store) Reconcile(busyByName map[string]bool, activeSessionNames []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureLoaded()

	now := s.clock.Now()
	changed := false

	// 1. 按本轮 busy 状态更新 tracking + 触发 flag
	for name, isBusy := range busyByName {
		wasBusy := s.tracking[name]
		switch {
		case isBusy && !wasBusy:
			s.tracking[name] = true
			changed = true
		case !isBusy && wasBusy:
			// active → idle 转换：触发 flag（如还没有）
			if _, has := s.flags[name]; !has {
				s.flags[name] = Flag{FirstAt: now}
			}
			delete(s.tracking, name)
			changed = true
		}
	}

	// 2. 回收 dead session 的 flag 和 tracking
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
		for name := range s.tracking {
			if _, ok := active[name]; !ok {
				delete(s.tracking, name)
				changed = true
			}
		}
	}

	if !changed {
		return nil
	}
	return s.saveLocked()
}

// Ack 清除指定 session 的 flag 和 tracking。无 flag 时也会清 tracking，避免下次旧状态又触发。
func (s *Store) Ack(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureLoaded()

	_, hadFlag := s.flags[name]
	_, hadTracking := s.tracking[name]
	if !hadFlag && !hadTracking {
		return nil
	}
	delete(s.flags, name)
	delete(s.tracking, name)
	return s.saveLocked()
}

// Clear 清空所有 flag 和 tracking。给 `cc-sesh attention clear` 子命令用。
func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags = map[string]Flag{}
	s.tracking = map[string]bool{}
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
	s.tracking = map[string]bool{}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}

	var f fileShape
	if err := json.Unmarshal(data, &f); err != nil {
		return
	}
	if f.Sessions != nil {
		s.flags = f.Sessions
	}
	if f.Tracking != nil {
		s.tracking = f.Tracking
	}
}

// saveLocked 在 mu 锁内调用：原子写。
func (s *Store) saveLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("mkdir state dir: %w", err)
	}

	// 排序输出以让文件 diff 友好
	flagKeys := make([]string, 0, len(s.flags))
	for k := range s.flags {
		flagKeys = append(flagKeys, k)
	}
	sort.Strings(flagKeys)
	orderedFlags := make(map[string]Flag, len(flagKeys))
	for _, k := range flagKeys {
		orderedFlags[k] = s.flags[k]
	}

	trackKeys := make([]string, 0, len(s.tracking))
	for k := range s.tracking {
		trackKeys = append(trackKeys, k)
	}
	sort.Strings(trackKeys)
	orderedTracking := make(map[string]bool, len(trackKeys))
	for _, k := range trackKeys {
		orderedTracking[k] = s.tracking[k]
	}

	data, err := json.MarshalIndent(fileShape{
		Version:  fileVersion,
		Sessions: orderedFlags,
		Tracking: orderedTracking,
	}, "", "  ")
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
	Tracking map[string]bool `json:"tracking,omitempty"`
}
