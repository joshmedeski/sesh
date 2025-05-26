package previewer

import (
	"strings"
	
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type TmuxPreviewStrategy struct {
	lister lister.Lister
	tmux   tmux.Tmux
}

func NewTmuxStrategy(lister lister.Lister, tmux tmux.Tmux) *TmuxPreviewStrategy {
	return &TmuxPreviewStrategy{lister: lister, tmux: tmux}
}

func (s *TmuxPreviewStrategy) Execute(name string) (string, error) {
	// Check if this is a marked window format: "session:window_number"
	if strings.Contains(name, ":") {
		// For marked windows, capture the specific window
		output, err := s.tmux.CapturePane(name)
		if err != nil {
			return "", err
		}
		return output, nil
	}
	
	// For regular sessions, use the session lookup
	session, sessionExists := s.lister.FindTmuxSession(name)
	if sessionExists {
		output, err := s.tmux.CapturePane(session.Name)
		if err != nil {
			return "", err
		}
		return output, nil
	}

	return "", nil
}
