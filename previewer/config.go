package previewer

import (
	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/ls"
)

type ConfigPreviewStrategy struct {
	lister lister.Lister
	ls     ls.Ls
}

func NewConfigStrategy(lister lister.Lister, ls ls.Ls) *ConfigPreviewStrategy {
	return &ConfigPreviewStrategy{lister: lister, ls: ls}
}

func (s *ConfigPreviewStrategy) Execute(name string) (string, error) {
	session, configExists := s.lister.FindConfigSession(name)
	if configExists {
		output, err := s.ls.ListDirectory(session.Path)
		if err != nil {
			return "", err
		}

		return output, nil
	}

	return "", nil
}
