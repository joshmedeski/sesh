package lister

import (
	"fmt"
	"os"

	"github.com/joshmedeski/sesh/v2/model"
)

func tmuxPaneKey(windowName string, paneID string) string {
	return fmt.Sprintf("tmux-pane:%s/%s", windowName, paneID)
}

func tmuxPaneDisplayName(pane *model.TmuxPane) string {
	hostname, _ := os.Hostname()
	if pane.PaneTitle != "" && pane.PaneTitle != hostname {
		return pane.PaneTitle
	}
	return pane.PaneCommand
}

func listTmuxPanes(l *RealLister) (model.SeshSessions, error) {
	tmuxPanes, err := l.tmux.ListTmuxPanes()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list tmux panes: %q", err)
	}

	// Count raw names to detect duplicates needing .0, .1 suffixes
	type paneEntry struct {
		pane    *model.TmuxPane
		rawName string
	}
	entries := make([]paneEntry, len(tmuxPanes))
	nameCounts := make(map[string]int)
	for i, pane := range tmuxPanes {
		rawName := fmt.Sprintf("%s/%s", pane.WindowName, tmuxPaneDisplayName(pane))
		entries[i] = paneEntry{pane: pane, rawName: rawName}
		nameCounts[rawName]++
	}

	directory := make(map[string]model.SeshSession)
	orderedIndex := []string{}
	nameIndexes := make(map[string]int)

	for _, entry := range entries {
		name := entry.rawName
		if nameCounts[entry.rawName] > 1 {
			idx := nameIndexes[entry.rawName]
			name = fmt.Sprintf("%s.%d", entry.rawName, idx)
			nameIndexes[entry.rawName] = idx + 1
		}

		key := tmuxPaneKey(entry.pane.WindowName, entry.pane.PaneID)
		orderedIndex = append(orderedIndex, key)
		directory[key] = model.SeshSession{
			Src:  "tmux-pane",
			Name: name,
			Path: entry.pane.PanePath,
		}
	}

	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) ListTmuxPanes() (model.SeshSessions, error) {
	return listTmuxPanes(l)
}
