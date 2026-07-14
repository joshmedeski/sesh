package git

import (
	"github.com/joshmedeski/sesh/v2/shell"
)

type Git interface {
	ShowTopLevel(name string) (bool, string, error)
	GitCommonDir(name string) (bool, string, error)
	Clone(url string, cmdDir string, dir string, gitFlags ...string) (string, error)
	WorktreeList(name string) (bool, string, error)
	Fetch(repoPath string) (string, error)
	WorktreeAdd(repoPath, target, branch, base string) (string, error)
	WorktreeAddDetached(repoPath, target, base string) (string, error)
	Pull(repoPath string) (string, error)
	CurrentBranch(path string) (bool, string, error)
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

func (g *RealGit) Clone(url string, cmdDir string, dir string, gitFlags ...string) (string, error) {
	args := []string{}

	if cmdDir != "" {
		args = append(args, "-C", cmdDir)
	}

	args = append(args, "clone")
	args = append(args, gitFlags...)
	args = append(args, url)

	if dir != "" {
		args = append(args, dir)
	}

	_, err := g.shell.CmdWithOutput("git", args...)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (g *RealGit) WorktreeList(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "worktree", "list", "--porcelain")
	if err != nil {
		return false, "", err
	}
	return true, out, nil
}

func (g *RealGit) Fetch(repoPath string) (string, error) {
	return g.shell.CmdWithOutput("git", "-C", repoPath, "fetch")
}

func (g *RealGit) WorktreeAdd(repoPath, target, branch, base string) (string, error) {
	return g.shell.CmdWithOutput("git", "-C", repoPath, "worktree", "add", target, "-b", branch, "--no-track", base)
}

func (g *RealGit) WorktreeAddDetached(repoPath, target, base string) (string, error) {
	return g.shell.CmdWithOutput("git", "-C", repoPath, "worktree", "add", "--detach", target, base)
}

func (g *RealGit) Pull(repoPath string) (string, error) {
	return g.shell.CmdWithOutput("git", "-C", repoPath, "pull", "--ff-only")
}

func (g *RealGit) CurrentBranch(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return false, "", err
	}
	return true, out, nil
}
