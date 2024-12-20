package previewer

import (
	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/tmux"
)

type TmuxPreviewStrategy struct {
	lister lister.Lister
	tmux   tmux.Tmux
}

func NewTmuxStrategy(lister lister.Lister, tmux tmux.Tmux) *TmuxPreviewStrategy {
	return &TmuxPreviewStrategy{lister: lister, tmux: tmux}
}

func (s *TmuxPreviewStrategy) Execute(name string) (string, error) {
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
