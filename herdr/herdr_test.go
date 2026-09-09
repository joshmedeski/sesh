package herdr

import (
	"encoding/json"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// liveSnapshot is a `session.snapshot` result recorded from herdr 0.8.2, with a
// second workspace added so resolution has something to choose between.
const liveSnapshot = `{"snapshot":{"agents":[],"focused_pane_id":"wJ:p1","focused_tab_id":"wJ:t1","focused_workspace_id":"wJ",
"panes":[
 {"agent_status":"unknown","cwd":"/Users/dev","focused":true,"foreground_cwd":"/Users/dev","pane_id":"wJ:p1","tab_id":"wJ:t1","workspace_id":"wJ"},
 {"agent_status":"unknown","cwd":"/Users/dev/c/sesh","focused":false,"foreground_cwd":"/Users/dev/c/sesh/lister","pane_id":"wK:p1","tab_id":"wK:t1","workspace_id":"wK"},
 {"agent_status":"unknown","cwd":"/Users/dev/c/sesh","focused":false,"foreground_cwd":"","pane_id":"wK:p2","tab_id":"wK:t2","workspace_id":"wK"}
],
"protocol":20,
"tabs":[
 {"focused":true,"label":"1","number":1,"pane_count":1,"tab_id":"wJ:t1","workspace_id":"wJ"},
 {"focused":false,"label":"editor","number":1,"pane_count":1,"tab_id":"wK:t1","workspace_id":"wK"},
 {"focused":false,"label":"logs","number":2,"pane_count":1,"tab_id":"wK:t2","workspace_id":"wK"}
],
"version":"0.8.2",
"workspaces":[
 {"active_tab_id":"wJ:t1","focused":true,"label":"~","number":1,"pane_count":1,"tab_count":1,"workspace_id":"wJ"},
 {"active_tab_id":"wK:t2","focused":false,"label":"sesh","number":2,"pane_count":2,"tab_count":2,"workspace_id":"wK"}
]}}`

func newTestHerdr(t *testing.T) (*Herdr, *MockClient, *oswrap.MockOs, *shell.MockShell) {
	client := NewMockClient(t)
	os := oswrap.NewMockOs(t)
	sh := shell.NewMockShell(t)
	os.EXPECT().Getenv("HERDR_BIN_PATH").Return("").Maybe()
	client.EXPECT().Call("session.snapshot", nil).Return(json.RawMessage(liveSnapshot), nil).Maybe()
	return NewHerdr(os, sh, client, ""), client, os, sh
}

func ok() json.RawMessage { return json.RawMessage(`{"type":"ok"}`) }

func TestListSessions(t *testing.T) {
	h, _, _, _ := newTestHerdr(t)

	sessions, err := h.ListSessions()

	require.NoError(t, err)
	require.Len(t, sessions, 2)
	assert.Equal(t, &model.TmuxSession{ID: "wJ", Name: "~", Path: "/Users/dev", Attached: 1, Windows: 1}, sessions[0])
	// Not focused, so the path comes from a pane on the active tab, which has
	// no foreground process and falls back to the pane's cwd.
	assert.Equal(t, &model.TmuxSession{ID: "wK", Name: "sesh", Path: "/Users/dev/c/sesh", Attached: 0, Windows: 2}, sessions[1])
}

func TestIsAttached(t *testing.T) {
	h, _, os, _ := newTestHerdr(t)
	os.EXPECT().Getenv("HERDR_ENV").Return("").Once()
	os.EXPECT().Getenv("HERDR_WORKSPACE_ID").Return("").Once()
	assert.False(t, h.IsAttached())

	os.EXPECT().Getenv("HERDR_ENV").Return("1").Once()
	assert.True(t, h.IsAttached())
}

func TestNewSessionReturnsWorkspaceID(t *testing.T) {
	h, client, _, _ := newTestHerdr(t)
	client.EXPECT().Call("workspace.create", map[string]any{"cwd": "/Users/dev/c/sesh", "label": "sesh", "focus": false}).
		Return(json.RawMessage(`{"type":"workspace_created","workspace":{"workspace_id":"wL","label":"sesh"},"tab":{"tab_id":"wL:t1"},"root_pane":{"pane_id":"wL:p1"}}`), nil)

	id, err := h.NewSession("sesh", "/Users/dev/c/sesh")

	require.NoError(t, err)
	assert.Equal(t, "wL", id)
}

func paneIs(id string) any {
	return mock.MatchedBy(func(p map[string]any) bool { return p["pane_id"] == id })
}

func TestSendKeysTypesIntoTheActivePane(t *testing.T) {
	t.Run("workspace by label lands on its active tab", func(t *testing.T) {
		h, client, _, _ := newTestHerdr(t)
		client.EXPECT().Call("pane.send_text", map[string]any{"pane_id": "wK:p2", "text": "nvim\n"}).
			Return(ok(), nil)

		_, err := h.SendKeys("sesh", "nvim")

		require.NoError(t, err)
	})

	t.Run("workspace by ID", func(t *testing.T) {
		h, client, _, _ := newTestHerdr(t)
		client.EXPECT().Call("pane.send_text", paneIs("wJ:p1")).Return(ok(), nil)

		_, err := h.SendKeys("wJ", "ls")

		require.NoError(t, err)
	})

	t.Run("window target by number", func(t *testing.T) {
		h, client, _, _ := newTestHerdr(t)
		client.EXPECT().Call("pane.send_text", paneIs("wK:p1")).Return(ok(), nil)

		_, err := h.SendKeys("sesh:1", "ls")

		require.NoError(t, err)
	})

	t.Run("tab ID", func(t *testing.T) {
		h, client, _, _ := newTestHerdr(t)
		client.EXPECT().Call("pane.send_text", paneIs("wK:p1")).Return(ok(), nil)

		_, err := h.SendKeys("wK:t1", "ls")

		require.NoError(t, err)
	})

	t.Run("unknown target", func(t *testing.T) {
		h, _, _, _ := newTestHerdr(t)

		_, err := h.SendKeys("nope", "ls")

		require.Error(t, err)
	})
}

func TestResolveWorkspaceFallsBackToEnrichedLabel(t *testing.T) {
	client := NewMockClient(t)
	os := oswrap.NewMockOs(t)
	os.EXPECT().Getenv("HERDR_BIN_PATH").Return("")
	enriched := `{"snapshot":{"workspaces":[{"workspace_id":"wZ","label":"460` + model.SessionNameSeparator + `herdr support","active_tab_id":"wZ:t1"}],"tabs":[],"panes":[]}}`
	client.EXPECT().Call("session.snapshot", nil).Return(json.RawMessage(enriched), nil)
	h := NewHerdr(os, shell.NewMockShell(t), client, "")

	w, err := h.resolveWorkspace("460")

	require.NoError(t, err)
	assert.Equal(t, "wZ", w.ID)
}

func TestCapturePaneReadsVisibleAnsi(t *testing.T) {
	h, client, _, _ := newTestHerdr(t)
	client.EXPECT().Call("pane.read", map[string]any{"pane_id": "wJ:p1", "source": "visible", "format": "ansi", "strip_ansi": false}).
		Return(json.RawMessage(`{"type":"pane_read","read":{"pane_id":"wJ:p1","text":"\u001b[1m~\u001b[0m $ "}}`), nil)

	out, err := h.CapturePane("~")

	require.NoError(t, err)
	assert.Equal(t, "\x1b[1m~\x1b[0m $ ", out)
}

func TestListWindows(t *testing.T) {
	h, _, _, _ := newTestHerdr(t)

	windows, err := h.ListWindows("sesh")

	require.NoError(t, err)
	assert.Equal(t, []*model.TmuxWindow{
		{Name: "editor", Path: "/Users/dev/c/sesh/lister", Index: 1, Active: false},
		{Name: "logs", Path: "/Users/dev/c/sesh", Index: 2, Active: true},
	}, windows)
}

func TestListAllWindowNames(t *testing.T) {
	h, _, _, _ := newTestHerdr(t)

	names, err := h.ListAllWindowNames("#{window_name}")

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"~": {"1"}, "sesh": {"editor", "logs"}}, names)
}

func TestNewWindowInSession(t *testing.T) {
	h, client, _, _ := newTestHerdr(t)
	client.EXPECT().Call("tab.create", map[string]any{"workspace_id": "wK", "cwd": "/Users/dev/notes", "label": "notes", "focus": false}).
		Return(json.RawMessage(`{"type":"tab_created","tab":{"tab_id":"wK:t3","label":"notes"},"root_pane":{"pane_id":"wK:p3"}}`), nil)
	client.EXPECT().Call("pane.send_text", map[string]any{"pane_id": "wK:p3", "text": "nvim .\n"}).
		Return(ok(), nil)

	target, err := h.NewWindowInSession(model.TmuxWindowOpts{
		Name: "notes", StartDir: "/Users/dev/notes", TargetSession: "sesh", Command: "nvim .", Background: true,
	})

	require.NoError(t, err)
	assert.Equal(t, "wK:t3", target)
}

func TestNextWindowInSessionWraps(t *testing.T) {
	h, client, _, _ := newTestHerdr(t)
	// sesh's active tab is its last, so "next" is the first.
	client.EXPECT().Call("tab.focus", map[string]any{"tab_id": "wK:t1"}).Return(ok(), nil)

	_, err := h.NextWindowInSession("sesh")

	require.NoError(t, err)
}

func TestSelectWindow(t *testing.T) {
	h, client, _, _ := newTestHerdr(t)
	client.EXPECT().Call("tab.focus", map[string]any{"tab_id": "wK:t2"}).Return(ok(), nil)

	_, err := h.SelectWindow("sesh:2")

	require.NoError(t, err)
}

func TestRenameSession(t *testing.T) {
	h, client, _, _ := newTestHerdr(t)
	client.EXPECT().Call("workspace.rename", map[string]any{"workspace_id": "wK", "label": "sesh — herdr"}).Return(ok(), nil)

	_, err := h.RenameSession("sesh", "sesh — herdr")

	require.NoError(t, err)
}

func TestSwitchOrAttach(t *testing.T) {
	t.Run("inside herdr only focuses", func(t *testing.T) {
		h, client, os, _ := newTestHerdr(t)
		os.EXPECT().Getenv("HERDR_ENV").Return("1")
		client.EXPECT().Call("workspace.focus", map[string]any{"workspace_id": "wK"}).Return(ok(), nil)

		out, err := h.SwitchOrAttach("sesh", model.ConnectOpts{})

		require.NoError(t, err)
		assert.Equal(t, "switching to herdr workspace: sesh", out)
	})

	t.Run("outside herdr focuses then attaches a client", func(t *testing.T) {
		h, client, os, sh := newTestHerdr(t)
		os.EXPECT().Getenv("HERDR_ENV").Return("")
		os.EXPECT().Getenv("HERDR_WORKSPACE_ID").Return("")
		os.EXPECT().Getenv("HERDR_SESSION").Return("work")
		client.EXPECT().Call("workspace.focus", map[string]any{"workspace_id": "wK"}).Return(ok(), nil)
		sh.EXPECT().CmdWithOutput("herdr", "--session", "work").Return("", nil)

		out, err := h.SwitchOrAttach("sesh", model.ConnectOpts{})

		require.NoError(t, err)
		assert.Equal(t, "attaching to herdr workspace: sesh", out)
	})

	t.Run("switch flag never attaches", func(t *testing.T) {
		h, client, _, _ := newTestHerdr(t)
		client.EXPECT().Call("workspace.focus", map[string]any{"workspace_id": "wJ"}).Return(ok(), nil)

		_, err := h.SwitchOrAttach("wJ", model.ConnectOpts{Switch: true})

		require.NoError(t, err)
	})
}

func TestDecodeResponse(t *testing.T) {
	t.Run("unwraps result", func(t *testing.T) {
		raw, err := decodeResponse("ping", []byte(`{"id":"sesh","result":{"type":"pong"}}`+"\n"))
		require.NoError(t, err)
		assert.JSONEq(t, `{"type":"pong"}`, string(raw))
	})

	t.Run("surfaces protocol errors", func(t *testing.T) {
		_, err := decodeResponse("workspace.focus", []byte(`{"id":"","error":{"code":"not_found","message":"workspace wQ not found"}}`))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not_found")
		assert.Contains(t, err.Error(), "wQ")
	})
}

func TestSocketPath(t *testing.T) {
	t.Run("explicit path wins", func(t *testing.T) {
		os := oswrap.NewMockOs(t)
		os.EXPECT().Getenv("HERDR_SOCKET_PATH").Return("/tmp/h.sock")
		assert.Equal(t, "/tmp/h.sock", SocketPath(os))
	})

	t.Run("named session", func(t *testing.T) {
		os := oswrap.NewMockOs(t)
		os.EXPECT().Getenv("HERDR_SOCKET_PATH").Return("")
		os.EXPECT().UserHomeDir().Return("/Users/dev", nil)
		os.EXPECT().Getenv("HERDR_SESSION").Return("work")
		assert.Equal(t, "/Users/dev/.config/herdr/sessions/work/herdr.sock", SocketPath(os))
	})

	t.Run("default session", func(t *testing.T) {
		os := oswrap.NewMockOs(t)
		os.EXPECT().Getenv("HERDR_SOCKET_PATH").Return("")
		os.EXPECT().UserHomeDir().Return("/Users/dev", nil)
		os.EXPECT().Getenv("HERDR_SESSION").Return("")
		assert.Equal(t, "/Users/dev/.config/herdr/herdr.sock", SocketPath(os))
	})
}
