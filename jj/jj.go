package jj

import (
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

type Jj interface {
	// Finds the jj workspace containing path and the root of the repo that
	// workspace belongs to (the same directory for the default workspace)
	Roots(path string) (isJj bool, workspace string, repo string)
	Clone(url string, dest string, flags ...string) error
	Fetch(repo string, flags ...string) (string, error)
	RemoteAdd(repo, name, url string) (string, error)
	WorkspaceAdd(repo, name, rev, target string) (string, error)
	BookmarkCreate(workspace, name, rev string) (string, error)
	New(workspace, rev string) (string, error)
}

type RealJj struct {
	os    oswrap.Os
	shell shell.Shell
}

func NewJj(os oswrap.Os, shell shell.Shell) Jj {
	return &RealJj{os, shell}
}

func (j *RealJj) Roots(path string) (bool, string, string) {
	for dir := path; ; {
		if info, err := j.os.Stat(filepath.Join(dir, ".jj")); err == nil && info.IsDir() {
			return true, dir, j.repoRoot(dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false, "", ""
		}
		dir = parent
	}
}

func (j *RealJj) repoRoot(workspace string) string {
	jjDir := filepath.Join(workspace, ".jj")
	pointer, err := j.os.ReadFile(filepath.Join(jjDir, "repo"))
	if err != nil {
		return workspace
	}
	repoDir := strings.TrimSpace(string(pointer))
	if !filepath.IsAbs(repoDir) {
		repoDir = filepath.Join(jjDir, repoDir)
	}
	return filepath.Dir(filepath.Dir(repoDir))
}

func (j *RealJj) Clone(url string, dest string, flags ...string) error {
	args := append([]string{"git", "clone"}, flags...)
	_, err := j.shell.CmdWithOutput("jj", append(args, url, dest)...)
	return err
}

func (j *RealJj) Fetch(repo string, flags ...string) (string, error) {
	return j.shell.CmdWithOutput("jj", append([]string{"-R", repo, "git", "fetch"}, flags...)...)
}

func (j *RealJj) RemoteAdd(repo, name, url string) (string, error) {
	return j.shell.CmdWithOutput("jj", "-R", repo, "git", "remote", "add", name, url)
}

func (j *RealJj) WorkspaceAdd(repo, name, rev, target string) (string, error) {
	return j.shell.CmdWithOutput("jj", "-R", repo, "workspace", "add", "--name", name, "-r", rev, target)
}

func (j *RealJj) BookmarkCreate(workspace, name, rev string) (string, error) {
	return j.shell.CmdWithOutput("jj", "-R", workspace, "bookmark", "create", name, "-r", rev)
}

func (j *RealJj) New(workspace, rev string) (string, error) {
	return j.shell.CmdWithOutput("jj", "-R", workspace, "new", rev)
}
