package git

import (
	"strings"

	"github.com/joshmedeski/sesh/shell"
)

type Git interface {
	GitRoot(name string) (bool, string, error)
	Clone(name string) (string, error)
}

type RealGit struct {
	shell shell.Shell
}

func NewGit(shell shell.Shell) Git {
	return &RealGit{shell}
}

func (g *RealGit) GitRoot(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "worktree", "list")
	if err != nil {
		return false, "", err
	}
	main := strings.Fields(out)[0]
	return true, main, nil
}

func (g *RealGit) Clone(name string) (string, error) {
	out, err := g.shell.Cmd("git", "clone", name)
	if err != nil {
		return "", err
	}
	return out, nil
}
