package git

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/shell"
)

type StatusSummary struct {
	Staged    int
	Unstaged  int
	Untracked int
	Deleted   int
}

type Git interface {
	ShowTopLevel(name string) (bool, string, error)
	GitCommonDir(name string) (bool, string, error)
	Clone(url string, cmdDir string, dir string) (string, error)
	WorktreeList(name string) (bool, string, error)
	CurrentBranch(path string) (bool, string, error)
	StatusSummary(path string) (StatusSummary, error)
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
	args := []string{"clone", url}
	if cmdDir != "" {
		args = append([]string{"-C", cmdDir}, args...)
	}
	if dir != "" {
		args = append(args, dir)
	}

	_, err := g.shell.CmdWithOutput("git", args...)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (g *RealGit) CurrentBranch(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return false, "", err
	}
	return true, out, nil
}

func (g *RealGit) WorktreeList(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "worktree", "list", "--porcelain")
	if err != nil {
		return false, "", err
	}
	return true, out, nil
}

func (g *RealGit) StatusSummary(path string) (StatusSummary, error) {
	out, err := g.shell.Cmd("git", "-C", path, "status", "--porcelain")
	if err != nil {
		return StatusSummary{}, err
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return StatusSummary{}, nil
	}
	var s StatusSummary
	for _, line := range lines {
		if strings.HasPrefix(line, "?? ") {
			s.Untracked++
			continue
		}
		first := line[0]
		second := line[1]
		if first != ' ' {
			s.Staged++
		}
		if second == 'M' {
			s.Unstaged++
		}
		if first == 'D' || second == 'D' {
			s.Deleted++
		}
	}
	return s, nil
}
