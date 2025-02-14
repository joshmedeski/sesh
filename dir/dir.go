package dir

import (
	"strings"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
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
	isGitBare, absPath := gitBareRootDir(d, path)
	if isGitBare {
		return true, absPath
	}
	isGit, absPath := gitRootDir(d, path)
	if isGit {
		return true, absPath
	}
	return false, ""
}

func gitBareRootDir(d *RealDir, path string) (hasRootDir bool, absPath string) {
	isGitBare, commonDir, _ := d.git.GitCommonDir(path)
	if isGitBare && strings.HasSuffix(commonDir, "/.bare") {
		topLevelDir := strings.TrimSuffix(commonDir, "/.bare")
		relativePath := strings.TrimPrefix(path, topLevelDir)
		firstDir := strings.Split(relativePath, string("/"))[1]
		name, err := d.path.Abs(topLevelDir + "/" + firstDir)
		if err != nil {
			return false, ""
		}
		return true, name
	} else {
		return false, ""
	}
}

func gitRootDir(d *RealDir, path string) (hasDir bool, absPath string) {
	isGit, topLevelDir, _ := d.git.ShowTopLevel(path)
	if isGit && topLevelDir != "" {
		return true, topLevelDir
	} else {
		return false, ""
	}
}
