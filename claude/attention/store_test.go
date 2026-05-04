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

func TestReconcile_TriggersNewFlag(t *testing.T) {
	t0 := time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)
	clk := &frozenClock{t: t0}
	s, _ := newStore(t, clk)

	err := s.Reconcile(map[string]Signal{
		"figma-design": {Reason: "mcp_auth", TriggerPID: 100},
	}, nil)
	require.NoError(t, err)

	flags := s.Load()
	require.Len(t, flags, 1)
	assert.Equal(t, t0, flags["figma-design"].FirstAt)
	assert.Equal(t, "mcp_auth", flags["figma-design"].Reason)
	assert.Equal(t, 100, flags["figma-design"].TriggerPID)
}

func TestReconcile_StickyAfterSignalGone(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]Signal{
		"bg-research": {Reason: "permission", TriggerPID: 200},
	}, nil))

	// 信号消失，但 flag 应保留
	require.NoError(t, s.Reconcile(map[string]Signal{}, nil))

	flags := s.Load()
	assert.Len(t, flags, 1)
	assert.Equal(t, "permission", flags["bg-research"].Reason)
}

func TestReconcile_PreservesFirstAtOnReTrigger(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]Signal{
		"x": {Reason: "auth_url"},
	}, nil))
	first := s.Load()["x"].FirstAt

	clk.t = clk.t.Add(5 * time.Minute)
	require.NoError(t, s.Reconcile(map[string]Signal{
		"x": {Reason: "auth_url"},
	}, nil))

	got := s.Load()["x"]
	assert.Equal(t, first, got.FirstAt, "FirstAt 必须保留首次触发时刻")
}

func TestReconcile_UpgradesReason(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]Signal{
		"x": {Reason: "pending", TriggerPID: 1},
	}, nil))
	require.NoError(t, s.Reconcile(map[string]Signal{
		"x": {Reason: "auth_url", TriggerPID: 2},
	}, nil))

	got := s.Load()["x"]
	assert.Equal(t, "auth_url", got.Reason)
	assert.Equal(t, 2, got.TriggerPID)
}

func TestAck_RemovesFlag(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]Signal{
		"a": {Reason: "auth_url"},
		"b": {Reason: "permission"},
	}, nil))

	require.NoError(t, s.Ack("a"))

	flags := s.Load()
	assert.NotContains(t, flags, "a")
	assert.Contains(t, flags, "b")
}

func TestAck_NoOpWhenAbsent(t *testing.T) {
	s, _ := newStore(t, &frozenClock{t: time.Now()})
	assert.NoError(t, s.Ack("nonexistent"))
}

func TestReconcile_GarbageCollectsDeadSessions(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]Signal{
		"alive": {Reason: "auth_url"},
		"dead":  {Reason: "permission"},
	}, nil))

	// dead 被移出 active 列表
	require.NoError(t, s.Reconcile(nil, []string{"alive"}))

	flags := s.Load()
	assert.Contains(t, flags, "alive")
	assert.NotContains(t, flags, "dead")
}

func TestReconcile_NilActiveListSkipsGC(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, _ := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]Signal{
		"a": {Reason: "auth_url"},
	}, nil))
	require.NoError(t, s.Reconcile(nil, nil)) // active=nil 不做 GC

	assert.Contains(t, s.Load(), "a")
}

func TestPersistence_RoundTrip(t *testing.T) {
	clk := &frozenClock{t: time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)}
	s, path := newStore(t, clk)

	require.NoError(t, s.Reconcile(map[string]Signal{
		"x": {Reason: "auth_url", TriggerPID: 42},
	}, nil))

	// 文件已落盘
	_, err := os.Stat(path)
	require.NoError(t, err)

	// 新 store 实例读同一文件
	s2 := New(path)
	flags := s2.Load()
	require.Contains(t, flags, "x")
	assert.Equal(t, "auth_url", flags["x"].Reason)
	assert.Equal(t, 42, flags["x"].TriggerPID)
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

	require.NoError(t, s.Reconcile(map[string]Signal{
		"a": {Reason: "auth_url"},
		"b": {Reason: "permission"},
	}, nil))

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
