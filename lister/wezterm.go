package lister

import (
	"fmt"
	"net/url"

	"github.com/joshmedeski/sesh/v2/model"
)

func weztermKey(name string) string {
	return fmt.Sprintf("wezterm:%s", name)
}

func listWezterm(l *RealLister) (model.SeshSessions, error) {
	workspaces, err := l.wezterm.ListWorkspaces()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list wezterm workspaces: %q", err)
	}

	directory := make(map[string]model.SeshSession)
	orderedIndex := []string{}

	for _, ws := range workspaces {
		key := weztermKey(ws.Name)
		orderedIndex = append(orderedIndex, key)

		// Use the CWD of the first pane as the workspace path.
		var path string
		if len(ws.Panes) > 0 {
			path = stripFileScheme(ws.Panes[0].Cwd)
		}

		directory[key] = model.SeshSession{
			Src:      "wezterm",
			Name:     ws.Name,
			Path:     path,
			Windows:  len(ws.Panes),
			Attached: attachedCount(ws),
		}
	}

	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func attachedCount(ws *model.WeztermWorkspace) int {
	for _, pane := range ws.Panes {
		if pane.IsActive {
			return 1
		}
	}
	return 0
}

// stripFileScheme removes a "file://" prefix from WezTerm CWD URIs.
func stripFileScheme(s string) string {
	u, err := url.Parse(s)
	if err != nil || u.Scheme != "file" {
		return s
	}
	return u.Path
}

func (l *RealLister) FindWeztermWorkspace(name string) (model.SeshSession, bool) {
	sessions, err := listWezterm(l)
	if err != nil {
		return model.SeshSession{}, false
	}
	key := weztermKey(name)
	if session, exists := sessions.Directory[key]; exists {
		return session, true
	}
	return model.SeshSession{}, false
}

func (l *RealLister) GetActiveWeztermWorkspace() (model.SeshSession, bool) {
	workspaces, err := l.wezterm.ListWorkspaces()
	if err != nil {
		return model.SeshSession{}, false
	}
	for _, ws := range workspaces {
		for _, pane := range ws.Panes {
			if pane.IsActive {
				var path string
				if len(ws.Panes) > 0 {
					path = stripFileScheme(ws.Panes[0].Cwd)
				}
				return model.SeshSession{
					Src:      "wezterm",
					Name:     ws.Name,
					Path:     path,
					Windows:  len(ws.Panes),
					Attached: 1,
				}, true
			}
		}
	}
	return model.SeshSession{}, false
}
