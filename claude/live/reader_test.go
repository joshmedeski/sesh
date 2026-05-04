package live

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeProc struct {
	alive map[int]bool
}

func (f fakeProc) IsAlive(pid int) bool { return f.alive[pid] }

func writeSession(t *testing.T, claudeDir string, pid int, status, kind, cwd string) {
	t.Helper()
	raw := map[string]any{
		"pid":       pid,
		"sessionId": "session-of-" + filepath.Base(cwd),
		"cwd":       cwd,
		"status":    status,
		"kind":      kind,
		"updatedAt": int64(1),
	}
	b, err := json.Marshal(raw)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "sessions", fmt.Sprintf("%d.json", pid)), b, 0o644))
}

func setupHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude", "sessions"), 0o755))
	return home
}

func TestRead_ClassifiesAndAggregatesByCwd(t *testing.T) {
	home := setupHome(t)

	writeSession(t, filepath.Join(home, ".claude"), 100, "busy", "interactive", "/work/proj-a")
	writeSession(t, filepath.Join(home, ".claude"), 101, "running", "subagent", "/work/proj-a")
	writeSession(t, filepath.Join(home, ".claude"), 102, "idle", "interactive", "/work/proj-b")
	writeSession(t, filepath.Join(home, ".claude"), 103, "auth_url", "interactive", "/work/proj-c")

	r := NewReader(home, fakeProc{alive: map[int]bool{100: true, 101: true, 102: true, 103: true}})
	got, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Status{Total: 2, Busy: 1, Subagent: 1}, got["/work/proj-a"])
	assert.Equal(t, Status{Total: 1}, got["/work/proj-b"])
	assert.Equal(t, Status{Total: 1, Needing: 1}, got["/work/proj-c"])
}

func TestRead_SkipsDeadPids(t *testing.T) {
	home := setupHome(t)

	writeSession(t, filepath.Join(home, ".claude"), 200, "busy", "interactive", "/x")
	writeSession(t, filepath.Join(home, ".claude"), 201, "busy", "interactive", "/x")

	r := NewReader(home, fakeProc{alive: map[int]bool{200: true}}) // 201 已死
	got, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Status{Total: 1, Busy: 1}, got["/x"])
}

func TestRead_HandlesMissingSessionsDir(t *testing.T) {
	home := t.TempDir() // 不创建 .claude/sessions
	r := NewReader(home, fakeProc{})
	got, err := r.Read()
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestRead_SkipsCorruptJSON(t *testing.T) {
	home := setupHome(t)
	require.NoError(t, os.WriteFile(filepath.Join(home, ".claude", "sessions", "garbage.json"), []byte("not-json"), 0o644))
	writeSession(t, filepath.Join(home, ".claude"), 300, "busy", "interactive", "/y")

	r := NewReader(home, fakeProc{alive: map[int]bool{300: true}})
	got, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Status{Total: 1, Busy: 1}, got["/y"])
}

func TestStatus_Severity(t *testing.T) {
	assert.Equal(t, LogicalIdle, Status{}.Severity())
	assert.Equal(t, LogicalIdle, Status{Total: 1}.Severity())
	assert.Equal(t, LogicalSubagent, Status{Total: 1, Subagent: 1}.Severity())
	assert.Equal(t, LogicalBusy, Status{Total: 2, Busy: 1, Subagent: 1}.Severity())
	assert.Equal(t, LogicalNeedsInput, Status{Total: 3, Needing: 1, Busy: 1, Subagent: 1}.Severity())
}

func TestStatus_Idle(t *testing.T) {
	assert.Equal(t, 2, Status{Total: 5, Busy: 1, Subagent: 1, Needing: 1}.Idle())
	assert.Equal(t, 0, Status{Total: 0}.Idle())
}

func TestAggregateBySession_MapsInstanceToOwningPaneSession(t *testing.T) {
	instances := []Instance{
		{PID: 100, Cwd: "/work/proj-a", Logical: LogicalBusy},
		{PID: 101, Cwd: "/work/proj-a", Logical: LogicalSubagent},
		{PID: 102, Cwd: "/work/proj-b", Logical: LogicalNeedsInput},
		{PID: 103, Cwd: "/work/proj-c", Logical: LogicalIdle}, // 没 pane 在这 cwd
	}
	panes := []PaneInfo{
		{SessionName: "alpha", Cwd: "/work/proj-a"},
		{SessionName: "alpha", Cwd: "/work/proj-a/subdir"}, // 同 session 的另一个 pane，cwd 不在 instances
		{SessionName: "beta", Cwd: "/work/proj-b"},
	}

	got := AggregateBySession(instances, panes)
	assert.Equal(t, Status{Total: 2, Busy: 1, Subagent: 1}, got["alpha"])
	assert.Equal(t, Status{Total: 1, Needing: 1}, got["beta"])
	assert.NotContains(t, got, "gamma", "无 pane 匹配的 instance 不应出现")
}

func TestAggregateBySession_DedupesSamePaneSessionCombo(t *testing.T) {
	instances := []Instance{
		{PID: 100, Cwd: "/x", Logical: LogicalBusy},
	}
	// 同 session 多个 pane 都 cd 在 /x —— instance 只该计一次
	panes := []PaneInfo{
		{SessionName: "alpha", Cwd: "/x"},
		{SessionName: "alpha", Cwd: "/x"},
	}
	got := AggregateBySession(instances, panes)
	assert.Equal(t, Status{Total: 1, Busy: 1}, got["alpha"])
}

func TestAggregateBySession_OneInstanceMultipleSessions(t *testing.T) {
	instances := []Instance{
		{PID: 100, Cwd: "/x", Logical: LogicalBusy},
	}
	// 不同 session 都 cd 到 /x —— 两个 session 都该看到这个 instance
	panes := []PaneInfo{
		{SessionName: "alpha", Cwd: "/x"},
		{SessionName: "beta", Cwd: "/x"},
	}
	got := AggregateBySession(instances, panes)
	assert.Equal(t, Status{Total: 1, Busy: 1}, got["alpha"])
	assert.Equal(t, Status{Total: 1, Busy: 1}, got["beta"])
}

func TestAggregateBySession_EmptyInputs(t *testing.T) {
	assert.Empty(t, AggregateBySession(nil, nil))
	assert.Empty(t, AggregateBySession([]Instance{{PID: 1, Cwd: "/x"}}, nil))
	assert.Empty(t, AggregateBySession(nil, []PaneInfo{{SessionName: "a", Cwd: "/x"}}))
}

func TestClassify(t *testing.T) {
	tests := []struct {
		status, kind string
		want         Logical
	}{
		{"auth_url", "interactive", LogicalNeedsInput},
		{"pending", "interactive", LogicalNeedsInput},
		{"busy", "interactive", LogicalBusy},
		{"running", "subagent", LogicalSubagent},
		{"async_launched", "subagent", LogicalSubagent},
		{"in_progress", "interactive", LogicalBusy},
		{"compacting", "interactive", LogicalBusy},
		{"idle", "interactive", LogicalIdle},
		{"completed", "interactive", LogicalIdle},
		{"unknown-future-value", "interactive", LogicalIdle},
		{"", "", LogicalIdle},
	}
	for _, tt := range tests {
		t.Run(tt.status+"/"+tt.kind, func(t *testing.T) {
			assert.Equal(t, tt.want, classify(tt.status, tt.kind))
		})
	}
}
