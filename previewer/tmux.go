package previewer

import (
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/mux"
)

type TmuxPreviewStrategy struct {
	lister lister.Lister
	mux    mux.Mux
}

func NewTmuxStrategy(lister lister.Lister, m mux.Mux) *TmuxPreviewStrategy {
	return &TmuxPreviewStrategy{lister: lister, mux: m}
}

func (s *TmuxPreviewStrategy) Execute(name string) (string, error) {
	session, sessionExists := s.lister.FindTmuxSession(name)

	if sessionExists {
		output, err := s.mux.CapturePane(session.Name)
		if err != nil {
			return "", err
		}

		return output, nil
	}

	return "", nil
}
