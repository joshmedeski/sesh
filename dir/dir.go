package dir

import (
	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/oswrap"
	"github.com/joshmedeski/sesh/pathwrap"
)

type Dir interface {
	Dir(name string) (isDir bool, absPath string)
	RootDir(name string) (hasRootDir bool, absPath string)
}

type RealDir struct {
	os   oswrap.Os
	git  git.Git
	path pathwrap.Path
}

func NewDir(os oswrap.Os, git git.Git, path pathwrap.Path) Dir {
	return &RealDir{os, git, path}
}

func (d *RealDir) Dir(path string) (isDir bool, absPath string) {
	absPath, err := d.path.Abs(path)
	if err != nil {
		return false, ""
	}

	info, err := d.os.Stat(absPath)
	if err != nil {
		return false, ""
	}
	if !info.IsDir() {
		return false, ""
	}

	return true, absPath
}

func (d *RealDir) RootDir(path string) (hasRootDir bool, absPath string) {
	isGit, absPath := gitRootDir(d, path)
	if isGit {
		return true, absPath
	}
	return false, ""
}

func gitRootDir(d *RealDir, path string) (hasDir bool, absPath string) {
	isGit, rootDir, _ := d.git.GitRoot(path)
	if isGit && rootDir != "" {
		return true, rootDir
	} else {
		return false, ""
	}
}
