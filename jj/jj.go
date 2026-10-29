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
