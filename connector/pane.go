package connector

import (
	"fmt"
	"strings"

	"github.com/Wingsdh/cc-sesh/v2/model"
)

func isTmuxPaneFormat(name string) bool {
	return strings.Contains(name, "/") && !strings.HasPrefix(name, "/")
}

func tmuxPaneStrategy(c *RealConnector, name string) (model.Connection, error) {
	if !c.tmux.IsAttached() || !isTmuxPaneFormat(name) {
		return model.Connection{Found: false}, nil
	}

	sessions, err := c.lister.ListTmuxPanes()
	if err != nil {
		return model.Connection{Found: false}, nil
	}

	for _, key := range sessions.OrderedIndex {
		session := sessions.Directory[key]
		if session.Name == name {
			return model.Connection{
				Found:       true,
				Session:     session,
				New:         false,
				AddToZoxide: false,
			}, nil
		}
	}

	return model.Connection{Found: false}, nil
}

func connectToTmuxPane(c *RealConnector, connection model.Connection, _ model.ConnectOpts) (string, error) {
	panes, err := c.tmux.ListTmuxPanes()
	if err != nil {
		return "", err
	}

	sessions, err := c.lister.ListTmuxPanes()
	if err != nil {
		return "", err
	}

	var targetPaneKey string
	for _, key := range sessions.OrderedIndex {
		if sessions.Directory[key].Name == connection.Session.Name {
			targetPaneKey = key
			break
		}
	}

	for _, pane := range panes {
		paneKey := fmt.Sprintf("tmux-pane:%s/%s", pane.WindowName, pane.PaneID)
		if paneKey == targetPaneKey {
			return c.tmux.SelectPane(pane.WindowIndex, pane.PaneIndex)
		}
	}

	return "", fmt.Errorf("pane not found: %s", connection.Session.Name)
}
