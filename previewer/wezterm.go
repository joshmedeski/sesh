package previewer

import (
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/wezterm"
)

type WeztermPreviewStrategy struct {
	lister  lister.Lister
	wezterm wezterm.Wezterm
}

func NewWeztermStrategy(lister lister.Lister, wezterm wezterm.Wezterm) *WeztermPreviewStrategy {
	return &WeztermPreviewStrategy{lister: lister, wezterm: wezterm}
}

func (s *WeztermPreviewStrategy) Execute(name string) (string, error) {
	session, exists := s.lister.FindWeztermWorkspace(name)
	if !exists {
		return "", nil
	}

	// Find the first pane of this workspace and get its text.
	panes, err := s.wezterm.ListAllPanes()
	if err != nil {
		return "", nil
	}

	for _, p := range panes {
		if p.Workspace == session.Name {
			output, err := s.wezterm.GetText(p.PaneID)
			if err != nil {
				return "", err
			}
			return output, nil
		}
	}

	return "", nil
}
