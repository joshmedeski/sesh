package git

import (
	"errors"
	"strings"

	"github.com/joshmedeski/sesh/v2/shell"
)

type Git interface {
	GitRoot(name string) (bool, string, error)
	Clone(url string, cmdDir string, dir string) (string, error)
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

	fields := strings.Fields(out)
	if len(fields) == 0 {
		return false, "", errors.New("error parsing git worktree fields")
	}

	root := fields[0]
	if strings.HasSuffix(root, "/.bare") {
		root = strings.TrimSuffix(root, "/.bare")
	}

	return true, root, nil
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
