package wezterm

import (
	"encoding/json"
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

// weztermListItem represents a single pane from `wezterm cli list --format json`.
type weztermListItem struct {
	WindowID    int    `json:"window_id"`
	TabID       int    `json:"tab_id"`
	PaneID      int    `json:"pane_id"`
	Workspace   string `json:"workspace"`
	Cwd         string `json:"cwd"`
	Title       string `json:"title"`
	IsActive    bool   `json:"is_active"`
	IsZoomed    bool   `json:"is_zoomed"`
	CursorX     int    `json:"cursor_x"`
	CursorY     int    `json:"cursor_y"`
	TabTitle    string `json:"tab_title"`
	PaneTitle   string `json:"pane_title"`
	CursorShape string `json:"cursor_shape"`
}

func (w *RealWezterm) ListAllPanes() ([]model.WeztermPane, error) {
	output, err := w.shell.Cmd("wezterm", "cli", "list", "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("couldn't list wezterm panes: %w", err)
	}

	var items []weztermListItem
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return nil, fmt.Errorf("couldn't parse wezterm list output: %w", err)
	}

	panes := make([]model.WeztermPane, len(items))
	for i, item := range items {
		panes[i] = model.WeztermPane{
			PaneID:    item.PaneID,
			TabID:     item.TabID,
			WindowID:  item.WindowID,
			Workspace: item.Workspace,
			Cwd:       item.Cwd,
			Title:     item.Title,
			IsActive:  item.IsActive,
		}
	}
	return panes, nil
}

func (w *RealWezterm) ListWorkspaces() ([]*model.WeztermWorkspace, error) {
	panes, err := w.ListAllPanes()
	if err != nil {
		return nil, err
	}

	workspaceMap := make(map[string]*model.WeztermWorkspace)
	var orderedNames []string

	for _, pane := range panes {
		ws, exists := workspaceMap[pane.Workspace]
		if !exists {
			ws = &model.WeztermWorkspace{
				Name:  pane.Workspace,
				Panes: []model.WeztermPane{},
			}
			workspaceMap[pane.Workspace] = ws
			orderedNames = append(orderedNames, pane.Workspace)
		}
		ws.Panes = append(ws.Panes, pane)
	}

	workspaces := make([]*model.WeztermWorkspace, len(orderedNames))
	for i, name := range orderedNames {
		workspaces[i] = workspaceMap[name]
	}
	return workspaces, nil
}
