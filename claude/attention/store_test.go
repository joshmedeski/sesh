package attention

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type frozenClock struct{ t time.Time }

func (f *frozenClock) Now() time.Time { return f.t }

func newStore(t *testing.T, clk Clock) (*Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "cc-sesh", "attention.json")
	s := New(path)
	if clk != nil {
		s.WithClock(clk)
	}
	return s, path
}

// 第一次看到 busy 不触发 flag —— 仅记录 tracking。
func TestReconcile_BusyAloneDoesNotTrigger(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]bool{"a": true}, []string{"a"}))

	assert.Empty(t, s.Load(), "纯 busy 不该触发 flag")
}

// busy → idle 转换才触发 flag。
func TestReconcile_BusyToIdleTriggers(t *testing.T) {
	t0 := time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)
	clk := &frozenClock{t: t0}
	s, _ := newStore(t, clk)

	// 第一轮：a busy
	require.NoError(t, s.Reconcile(map[string]bool{"a": true}, []string{"a"}))
	require.Empty(t, s.Load())

	// 第二轮：a 不再 busy（变 idle）→ 触发
	clk.t = t0.Add(2 * time.Minute)
	require.NoError(t, s.Reconcile(map[string]bool{"a": false}, []string{"a"}))

	flags := s.Load()
	require.Contains(t, flags, "a")
	assert.Equal(t, clk.t, flags["a"].FirstAt)
}

// 从未 busy 直接 idle 不触发（避免冷启动每个 idle 都被标记）。
func TestReconcile_IdleFromColdStartDoesNotTrigger(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]bool{"a": false}, []string{"a"}))

	assert.Empty(t, s.Load())
}

// flag 触发后是粘性的，后续 Reconcile 不会清除。
func TestReconcile_FlagIsSticky(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	// busy → idle：触发
	require.NoError(t, s.Reconcile(map[string]bool{"a": true}, []string{"a"}))
	require.NoError(t, s.Reconcile(map[string]bool{"a": false}, []string{"a"}))
	require.Contains(t, s.Load(), "a")

	// 又一轮 idle：flag 仍在
	require.NoError(t, s.Reconcile(map[string]bool{"a": false}, []string{"a"}))
	assert.Contains(t, s.Load(), "a")
}

// FirstAt 在已存在 flag 时不被覆盖。
func TestReconcile_PreservesFirstAtAcrossRetrigger(t *testing.T) {
	t0 := time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)
	clk := &frozenClock{t: t0}
	s, _ := newStore(t, clk)

	// 第一次转换：触发 t0
	require.NoError(t, s.Reconcile(map[string]bool{"x": true}, []string{"x"}))
	require.NoError(t, s.Reconcile(map[string]bool{"x": false}, []string{"x"}))
	first := s.Load()["x"].FirstAt

	// 又跑一轮再回 idle：flag 仍存在，FirstAt 不变
	clk.t = t0.Add(10 * time.Minute)
	require.NoError(t, s.Reconcile(map[string]bool{"x": true}, []string{"x"}))
	clk.t = t0.Add(15 * time.Minute)
	require.NoError(t, s.Reconcile(map[string]bool{"x": false}, []string{"x"}))

	assert.Equal(t, first, s.Load()["x"].FirstAt)
}

func TestAck_RemovesFlag(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]bool{"a": true, "b": true}, []string{"a", "b"}))
	require.NoError(t, s.Reconcile(map[string]bool{"a": false, "b": false}, []string{"a", "b"}))

	require.NoError(t, s.Ack("a"))
	flags := s.Load()
	assert.NotContains(t, flags, "a")
	assert.Contains(t, flags, "b")
}

func TestAck_NoOpWhenAbsent(t *testing.T) {
	s, _ := newStore(t, &frozenClock{t: time.Now()})
	assert.NoError(t, s.Ack("nonexistent"))
}

// Ack 也清 tracking：避免 ack 后又出现 idle 立即重新触发。
func TestAck_ClearsTrackingToo(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	// busy → idle → 触发 flag
	require.NoError(t, s.Reconcile(map[string]bool{"a": true}, []string{"a"}))
	require.NoError(t, s.Reconcile(map[string]bool{"a": false}, []string{"a"}))
	require.NoError(t, s.Ack("a"))

	// 又一轮 busy → 又一轮 idle：要想重新触发，必须重新走完转换
	require.NoError(t, s.Reconcile(map[string]bool{"a": false}, []string{"a"}))
	assert.NotContains(t, s.Load(), "a", "ack 后再来 idle 不该立即触发")

	// 走完完整 busy→idle 又会触发
	require.NoError(t, s.Reconcile(map[string]bool{"a": true}, []string{"a"}))
	require.NoError(t, s.Reconcile(map[string]bool{"a": false}, []string{"a"}))
	assert.Contains(t, s.Load(), "a")
}

func TestReconcile_GarbageCollectsDeadSessions(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	// 触发两个
	require.NoError(t, s.Reconcile(map[string]bool{"alive": true, "dead": true}, []string{"alive", "dead"}))
	require.NoError(t, s.Reconcile(map[string]bool{"alive": false, "dead": false}, []string{"alive", "dead"}))

	// 下一轮 dead 不在 active 列表
	require.NoError(t, s.Reconcile(nil, []string{"alive"}))

	flags := s.Load()
	assert.Contains(t, flags, "alive")
	assert.NotContains(t, flags, "dead")
}

func TestReconcile_NilActiveListSkipsGC(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]bool{"a": true}, []string{"a"}))
	require.NoError(t, s.Reconcile(map[string]bool{"a": false}, []string{"a"}))
	require.Contains(t, s.Load(), "a")

	require.NoError(t, s.Reconcile(nil, nil)) // 不做 GC
	assert.Contains(t, s.Load(), "a")
}

func TestPersistence_RoundTrip(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, path := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]bool{"x": true}, []string{"x"}))
	require.NoError(t, s.Reconcile(map[string]bool{"x": false}, []string{"x"}))

	_, err := os.Stat(path)
	require.NoError(t, err)

	s2 := New(path)
	flags := s2.Load()
	require.Contains(t, flags, "x")
	assert.Equal(t, clk.t, flags["x"].FirstAt)
}

// tracking 也要持久化：跨进程 busy→idle 转换仍能识别。
func TestPersistence_TrackingSurvivesRestart(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, path := newStore(t, clk)

	// 第一进程：仅看到 busy，写 tracking 但无 flag
	require.NoError(t, s.Reconcile(map[string]bool{"x": true}, []string{"x"}))
	require.Empty(t, s.Load())

	// 第二进程：只看到 idle —— 应该识别出 busy→idle 转换并触发
	s2 := New(path).WithClock(clk)
	require.NoError(t, s2.Reconcile(map[string]bool{"x": false}, []string{"x"}))
	assert.Contains(t, s2.Load(), "x")
}

func TestLoad_MissingFileReturnsEmpty(t *testing.T) {
	s, _ := newStore(t, &frozenClock{t: time.Now()})
	assert.Empty(t, s.Load())
}

func TestLoad_CorruptFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "attention.json")
	require.NoError(t, os.WriteFile(path, []byte("{not json"), 0o644))

	s := New(path)
	assert.Empty(t, s.Load())
}

func TestClear_RemovesAll(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]bool{"a": true, "b": true}, []string{"a", "b"}))
	require.NoError(t, s.Reconcile(map[string]bool{"a": false, "b": false}, []string{"a", "b"}))

	require.NoError(t, s.Clear())
	assert.Empty(t, s.Load())
}

func TestDefaultPath_UsesXDGStateHome(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", "/custom/state")
	got, err := DefaultPath()
	require.NoError(t, err)
	assert.Equal(t, "/custom/state/cc-sesh/attention.json", got)
}

func TestDefaultPath_FallbackHomeDotLocal(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", "")
	t.Setenv("HOME", "/home/test")
	got, err := DefaultPath()
	require.NoError(t, err)
	assert.Equal(t, "/home/test/.local/state/cc-sesh/attention.json", got)
}
