package ls

import (
	"strings"

	"github.com/joshmedeski/sesh/model"
	"github.com/joshmedeski/sesh/shell"
)

type Ls interface {
	ListDirectory(name string) (string, error)
}

type RealLs struct {
	config model.Config
	shell  shell.Shell
}

func NewLs(config model.Config, shell shell.Shell) Ls {
	return &RealLs{config, shell}
}

func (g *RealLs) ListDirectory(path string) (string, error) {
	command := g.config.DefaultSessionConfig.LsCommand
	if command == "" {
		command = "ls"
	}

	cmd := strings.Split(command, " ")
	cmd = append(cmd, path)

	out, err := g.shell.Cmd(cmd[0], cmd[1:]...)
	if err != nil {
		return "", err
	}
	return out, nil
}
