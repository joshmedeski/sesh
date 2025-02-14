package previewer

import (
	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/ls"
)

type DirectoryPreviewStrategy struct {
	home home.Home
	dir  dir.Dir
	ls   ls.Ls
}

func NewDirectoryStrategy(home home.Home, dir dir.Dir, ls ls.Ls) *DirectoryPreviewStrategy {
	return &DirectoryPreviewStrategy{home: home, dir: dir, ls: ls}
}

func (s *DirectoryPreviewStrategy) Execute(name string) (string, error) {
	path, _ := s.home.ExpandHome(name)
	isDir, absPath := s.dir.Dir(path)

	if isDir {
		output, err := s.ls.ListDirectory(absPath)
		if err != nil {
			return "", err
		}

		return output, nil
	}

	return "", nil
}
