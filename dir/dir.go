package dir

import (
	"github.com/joshmedeski/sesh/oswrap"
	"github.com/joshmedeski/sesh/pathwrap"
)

type Dir interface {
	Dir(name string) (isDir bool, absPath string)
}

type RealDir struct {
	os   oswrap.Os
	path pathwrap.Path
}

func NewDir(os oswrap.Os, path pathwrap.Path) Dir {
	return &RealDir{os, path}
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
