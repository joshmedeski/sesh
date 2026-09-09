package herdr

// The subset of herdr's API objects sesh reads. Field names follow the
// protocol-20 schema (`herdr api schema --json`).

type workspace struct {
	ID          string `json:"workspace_id"`
	Number      int    `json:"number"`
	Label       string `json:"label"`
	Focused     bool   `json:"focused"`
	PaneCount   int    `json:"pane_count"`
	TabCount    int    `json:"tab_count"`
	ActiveTabID string `json:"active_tab_id"`
}

type tab struct {
	ID          string `json:"tab_id"`
	WorkspaceID string `json:"workspace_id"`
	Number      int    `json:"number"`
	Label       string `json:"label"`
	Focused     bool   `json:"focused"`
}

type pane struct {
	ID            string `json:"pane_id"`
	WorkspaceID   string `json:"workspace_id"`
	TabID         string `json:"tab_id"`
	Focused       bool   `json:"focused"`
	CWD           string `json:"cwd"`
	ForegroundCWD string `json:"foreground_cwd"`
}

// path is the directory sesh should show for a pane: what the foreground
// process is in when herdr knows it, else where the shell was started.
func (p pane) path() string {
	if p.ForegroundCWD != "" {
		return p.ForegroundCWD
	}
	return p.CWD
}

type snapshot struct {
	FocusedWorkspaceID string      `json:"focused_workspace_id"`
	Workspaces         []workspace `json:"workspaces"`
	Tabs               []tab       `json:"tabs"`
	Panes              []pane      `json:"panes"`
}

// workspacePath picks the directory a workspace is "in": its focused pane,
// else a pane on its active tab, else any of its panes.
func (s snapshot) workspacePath(w workspace) string {
	var onActiveTab, any string
	for _, p := range s.Panes {
		if p.WorkspaceID != w.ID || p.path() == "" {
			continue
		}
		if p.Focused {
			return p.path()
		}
		if p.TabID == w.ActiveTabID && onActiveTab == "" {
			onActiveTab = p.path()
		}
		if any == "" {
			any = p.path()
		}
	}
	if onActiveTab != "" {
		return onActiveTab
	}
	return any
}

// activePane is the pane input to a workspace should land in: the focused
// pane of its active tab, else the first pane on that tab, else any pane.
func (s snapshot) activePane(w workspace) (pane, bool) {
	var onActiveTab, any *pane
	for i := range s.Panes {
		p := &s.Panes[i]
		if p.WorkspaceID != w.ID {
			continue
		}
		if p.TabID == w.ActiveTabID {
			if p.Focused {
				return *p, true
			}
			if onActiveTab == nil {
				onActiveTab = p
			}
		}
		if any == nil {
			any = p
		}
	}
	if onActiveTab != nil {
		return *onActiveTab, true
	}
	if any != nil {
		return *any, true
	}
	return pane{}, false
}

// tabPane is the pane input to a tab should land in.
func (s snapshot) tabPane(tabID string) (pane, bool) {
	var first *pane
	for i := range s.Panes {
		p := &s.Panes[i]
		if p.TabID != tabID {
			continue
		}
		if p.Focused {
			return *p, true
		}
		if first == nil {
			first = p
		}
	}
	if first != nil {
		return *first, true
	}
	return pane{}, false
}
