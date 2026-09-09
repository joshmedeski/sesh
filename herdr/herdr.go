// Package herdr drives herdr (https://herdr.dev) as a sesh multiplexer backend.
//
// The mapping is: a sesh session is a herdr workspace, a window is a tab, and
// typing into a session means sending input to the focused pane of the
// workspace's active tab. Every operation goes over the socket API; only
// attaching a terminal, which has to take over the tty, shells out to the
// herdr binary.
//
// Sesh addresses sessions by name, herdr by ID, so every call resolves the
// target against a session snapshot first. Labels aren't unique in herdr; the
// first workspace with a matching label wins, and IDs are accepted anywhere a
// name is so callers that already hold one can skip the lookup.
package herdr

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/mux"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

const Name = "herdr"

type Herdr struct {
	os     oswrap.Os
	shell  shell.Shell
	client Client
	bin    string
	// snap caches one session snapshot per process, so a connect that lists,
	// resolves and types costs one round trip for the reads. Mutations drop it.
	snap *snapshot
}

var _ mux.Mux = (*Herdr)(nil)

// NewHerdr builds the backend. bin is the herdr executable used to attach;
// empty falls back to HERDR_BIN_PATH (what herdr sets for its plugins) and
// then to "herdr" on PATH.
func NewHerdr(os oswrap.Os, shell shell.Shell, client Client, bin string) *Herdr {
	if bin == "" {
		bin = os.Getenv("HERDR_BIN_PATH")
	}
	if bin == "" {
		bin = "herdr"
	}
	return &Herdr{os: os, shell: shell, client: client, bin: bin}
}

func (h *Herdr) Name() string { return Name }

// IsAttached is true inside a herdr pane, where herdr injects HERDR_ENV=1 and
// the pane's coordinates into the environment.
func (h *Herdr) IsAttached() bool {
	return h.os.Getenv("HERDR_ENV") != "" || h.os.Getenv("HERDR_WORKSPACE_ID") != ""
}

func (h *Herdr) snapshot() (*snapshot, error) {
	if h.snap != nil {
		return h.snap, nil
	}
	var result struct {
		Snapshot snapshot `json:"snapshot"`
	}
	if err := call(h.client, "session.snapshot", nil, &result); err != nil {
		return nil, err
	}
	h.snap = &result.Snapshot
	return h.snap, nil
}

func (h *Herdr) invalidate() { h.snap = nil }

func (h *Herdr) ListSessions() ([]*model.TmuxSession, error) {
	snap, err := h.snapshot()
	if err != nil {
		return nil, err
	}
	sessions := make([]*model.TmuxSession, 0, len(snap.Workspaces))
	for _, w := range snap.Workspaces {
		attached := 0
		if w.Focused {
			attached = 1
		}
		sessions = append(sessions, &model.TmuxSession{
			ID:       w.ID,
			Name:     w.Label,
			Path:     snap.workspacePath(w),
			Attached: attached,
			Windows:  w.TabCount,
		})
	}
	return sessions, nil
}

func (h *Herdr) NewSession(sessionName string, startDir string) (string, error) {
	var result struct {
		Workspace workspace `json:"workspace"`
	}
	params := map[string]any{"cwd": startDir, "label": sessionName, "focus": false}
	if err := call(h.client, "workspace.create", params, &result); err != nil {
		return "", err
	}
	h.invalidate()
	return result.Workspace.ID, nil
}

func (h *Herdr) RenameSession(target string, newName string) (string, error) {
	w, err := h.resolveWorkspace(target)
	if err != nil {
		return "", err
	}
	params := map[string]any{"workspace_id": w.ID, "label": newName}
	if err := call(h.client, "workspace.rename", params, nil); err != nil {
		return "", err
	}
	h.invalidate()
	return "", nil
}

func (h *Herdr) SwitchClient(target string) (string, error) {
	w, err := h.resolveWorkspace(target)
	if err != nil {
		return "", err
	}
	if err := call(h.client, "workspace.focus", map[string]any{"workspace_id": w.ID}, nil); err != nil {
		return "", err
	}
	h.invalidate()
	return "", nil
}

// SwitchOrAttach focuses the workspace, then hands the terminal to a herdr
// client when there isn't one yet. Focusing first means the client comes up
// on the right workspace instead of wherever it was last.
func (h *Herdr) SwitchOrAttach(name string, opts model.ConnectOpts) (string, error) {
	if _, err := h.SwitchClient(name); err != nil {
		return "", fmt.Errorf("failed to switch to herdr workspace: %w", err)
	}
	if opts.Switch || h.IsAttached() {
		return fmt.Sprintf("switching to herdr workspace: %s", name), nil
	}
	args := []string{}
	if session := h.os.Getenv("HERDR_SESSION"); session != "" {
		args = append(args, "--session", session)
	}
	if _, err := h.shell.CmdWithOutput(h.bin, args...); err != nil {
		return "", fmt.Errorf("failed to attach to herdr: %w", err)
	}
	return fmt.Sprintf("attaching to herdr workspace: %s", name), nil
}

// SendKeys types command and a newline into the pane behind target, the same
// shape as `tmux send-keys ... Enter`. target may be a workspace, a tab
// ("name:N" or a tab ID) or a pane ID.
func (h *Herdr) SendKeys(target string, command string) (string, error) {
	p, err := h.resolvePane(target)
	if err != nil {
		return "", err
	}
	if err := h.typeInto(p.ID, command); err != nil {
		return "", err
	}
	return "", nil
}

func (h *Herdr) CapturePane(target string) (string, error) {
	p, err := h.resolvePane(target)
	if err != nil {
		return "", err
	}
	var result struct {
		Read struct {
			Text string `json:"text"`
		} `json:"read"`
	}
	params := map[string]any{"pane_id": p.ID, "source": "visible", "format": "ansi", "strip_ansi": false}
	if err := call(h.client, "pane.read", params, &result); err != nil {
		return "", err
	}
	return result.Read.Text, nil
}

func (h *Herdr) ListWindows(target string) ([]*model.TmuxWindow, error) {
	w, err := h.resolveWorkspace(target)
	if err != nil {
		return nil, err
	}
	snap, _ := h.snapshot()
	windows := []*model.TmuxWindow{}
	for _, t := range h.tabsOf(w.ID) {
		path := ""
		if p, ok := snap.tabPane(t.ID); ok {
			path = p.path()
		}
		windows = append(windows, &model.TmuxWindow{
			Name:   t.Label,
			Path:   path,
			Index:  t.Number,
			Active: t.ID == w.ActiveTabID,
		})
	}
	return windows, nil
}

// ListAllWindowNames ignores format: herdr has no format language, tabs are
// listed by label.
func (h *Herdr) ListAllWindowNames(_ string) (map[string][]string, error) {
	snap, err := h.snapshot()
	if err != nil {
		return nil, err
	}
	names := make(map[string][]string, len(snap.Workspaces))
	for _, w := range snap.Workspaces {
		for _, t := range h.tabsOf(w.ID) {
			names[w.Label] = append(names[w.Label], t.Label)
		}
	}
	return names, nil
}

// NewWindowInSession creates a tab and returns its ID, which every other
// window method accepts as a target. A command is typed into the new tab
// rather than becoming its process: herdr tabs always run a shell.
func (h *Herdr) NewWindowInSession(opts model.TmuxWindowOpts) (string, error) {
	w, err := h.resolveWorkspace(opts.TargetSession)
	if err != nil {
		return "", err
	}
	params := map[string]any{"workspace_id": w.ID, "focus": !opts.Background}
	if opts.StartDir != "" {
		params["cwd"] = opts.StartDir
	}
	if opts.Name != "" {
		params["label"] = opts.Name
	}
	var result struct {
		Tab      tab  `json:"tab"`
		RootPane pane `json:"root_pane"`
	}
	if err := call(h.client, "tab.create", params, &result); err != nil {
		return "", err
	}
	h.invalidate()
	if opts.Command != "" && result.RootPane.ID != "" {
		if err := h.typeInto(result.RootPane.ID, opts.Command); err != nil {
			return "", err
		}
	}
	return result.Tab.ID, nil
}

// NextWindowInSession focuses the tab after the workspace's active one,
// wrapping around, which is how startup lands on the first configured window
// after creating the rest.
func (h *Herdr) NextWindowInSession(target string) (string, error) {
	w, err := h.resolveWorkspace(target)
	if err != nil {
		return "", err
	}
	tabs := h.tabsOf(w.ID)
	if len(tabs) < 2 {
		return "", nil
	}
	for i, t := range tabs {
		if t.ID == w.ActiveTabID {
			return h.focusTab(tabs[(i+1)%len(tabs)].ID)
		}
	}
	return h.focusTab(tabs[0].ID)
}

func (h *Herdr) SelectWindow(target string) (string, error) {
	t, err := h.resolveTab(target)
	if err != nil {
		return "", err
	}
	return h.focusTab(t.ID)
}

func (h *Herdr) focusTab(tabID string) (string, error) {
	if err := call(h.client, "tab.focus", map[string]any{"tab_id": tabID}, nil); err != nil {
		return "", err
	}
	h.invalidate()
	return "", nil
}

// tabsOf returns a workspace's tabs in sidebar order.
func (h *Herdr) tabsOf(workspaceID string) []tab {
	snap, err := h.snapshot()
	if err != nil {
		return nil
	}
	tabs := []tab{}
	for _, t := range snap.Tabs {
		if t.WorkspaceID == workspaceID {
			tabs = append(tabs, t)
		}
	}
	sort.SliceStable(tabs, func(i, j int) bool { return tabs[i].Number < tabs[j].Number })
	return tabs
}

var (
	tabIDPattern  = regexp.MustCompile(`^(w[^:\s]+):t\d+$`)
	paneIDPattern = regexp.MustCompile(`^(w[^:\s]+):p\d+$`)
)

// resolveWorkspace accepts a workspace ID or a label. A label match prefers an
// exact hit and falls back to a name sesh has enriched with an issue title, the
// same rule the lister applies to tmux sessions.
func (h *Herdr) resolveWorkspace(target string) (workspace, error) {
	snap, err := h.snapshot()
	if err != nil {
		return workspace{}, err
	}
	for _, w := range snap.Workspaces {
		if w.ID == target || w.Label == target {
			return w, nil
		}
	}
	prefix := target + model.SessionNameSeparator
	for _, w := range snap.Workspaces {
		if strings.HasPrefix(w.Label, prefix) {
			return w, nil
		}
	}
	return workspace{}, fmt.Errorf("herdr workspace %q not found", target)
}

// resolveTab accepts a tab ID or "workspace:N", where N is the tab's number
// as ListWindows reports it.
func (h *Herdr) resolveTab(target string) (tab, error) {
	snap, err := h.snapshot()
	if err != nil {
		return tab{}, err
	}
	if tabIDPattern.MatchString(target) {
		for _, t := range snap.Tabs {
			if t.ID == target {
				return t, nil
			}
		}
		return tab{}, fmt.Errorf("herdr tab %q not found", target)
	}
	sep := strings.LastIndex(target, ":")
	if sep < 0 {
		return tab{}, fmt.Errorf("%q is not a tab target", target)
	}
	number, err := strconv.Atoi(target[sep+1:])
	if err != nil {
		return tab{}, fmt.Errorf("%q is not a tab target", target)
	}
	w, err := h.resolveWorkspace(target[:sep])
	if err != nil {
		return tab{}, err
	}
	for _, t := range h.tabsOf(w.ID) {
		if t.Number == number {
			return t, nil
		}
	}
	return tab{}, fmt.Errorf("herdr tab %d not found in workspace %q", number, w.Label)
}

// resolvePane finds where input for target should go: the pane itself for a
// pane ID, a tab's focused pane for a tab target, and the active tab's focused
// pane for a workspace.
func (h *Herdr) resolvePane(target string) (pane, error) {
	snap, err := h.snapshot()
	if err != nil {
		return pane{}, err
	}
	if paneIDPattern.MatchString(target) {
		for _, p := range snap.Panes {
			if p.ID == target {
				return p, nil
			}
		}
		return pane{}, fmt.Errorf("herdr pane %q not found", target)
	}
	if w, err := h.resolveWorkspace(target); err == nil {
		if p, ok := snap.activePane(w); ok {
			return p, nil
		}
		return pane{}, fmt.Errorf("herdr workspace %q has no panes", w.Label)
	}
	t, err := h.resolveTab(target)
	if err != nil {
		return pane{}, err
	}
	if p, ok := snap.tabPane(t.ID); ok {
		return p, nil
	}
	return pane{}, fmt.Errorf("herdr tab %q has no panes", t.ID)
}

// typeInto sends command followed by a newline as terminal input. Verified
// against herdr 0.8.2: `pane.send_text` with a trailing "\n" runs the command
// in the pane's shell.
func (h *Herdr) typeInto(paneID string, command string) error {
	params := map[string]any{"pane_id": paneID, "text": command + "\n"}
	return call(h.client, "pane.send_text", params, nil)
}
