package previewer

import (
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/shell"
)

type ConfigPreviewStrategy struct {
	lister lister.Lister
	shell  shell.Shell
}

func NewConfigStrategy(lister lister.Lister, shell shell.Shell) *ConfigPreviewStrategy {
	return &ConfigPreviewStrategy{lister: lister, shell: shell}
}

func (s *ConfigPreviewStrategy) Execute(name string) (string, error) {
	session, configExists := s.lister.FindConfigSession(name)
	if !configExists {
		return "", nil
	}

	if session.PreviewCommand == "" {
		return "", nil
	}

	replacements := map[string]string{
		"{}": session.Path,
	}
	cmdParts, err := s.shell.PrepareCmd(session.PreviewCommand, replacements)
	if err != nil {
		return "", err
	}

	cmdOutput, err := s.shell.Cmd(cmdParts[0], cmdParts[1:]...)
	if err != nil {
		return "", err
	}

	return cmdOutput, nil
}
