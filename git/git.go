package git

import (
	"github.com/joshmedeski/sesh/v2/shell"
)

type Git interface {
	ShowTopLevel(name string) (bool, string, error)
	GitCommonDir(name string) (bool, string, error)
	Clone(url string, cmdDir string, dir string) (string, error)
}

type RealGit struct {
	shell shell.Shell
}

func NewGit(shell shell.Shell) Git {
	return &RealGit{shell}
}

func (g *RealGit) ShowTopLevel(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "rev-parse", "--show-toplevel")
	if err != nil {
		return false, "", err
	}
	return true, out, nil
}

func (g *RealGit) GitCommonDir(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "rev-parse", "--git-common-dir")
	if err != nil {
		return false, "", err
	}
	return true, out, nil
}

func (g *RealGit) Clone(url string, cmdDir string, dir string) (string, error) {
	var out string
	var err error

	args := []string{"clone", url}
	if cmdDir != "" {
		args = append([]string{"-C", cmdDir}, args...)
	}
	if dir != "" {
		args = append(args, dir)
	}

	out, err = g.shell.Cmd("git", args...)
	if err != nil {
		return "", err
	}
	return out, nil
}
