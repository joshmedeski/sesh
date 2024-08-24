package pathwrap

import (
	"path"
	"path/filepath"
)

type Path interface {
	Join(elem ...string) string
	Abs(path string) (string, error)
	Base(path string) string
	EvalSymlinks(path string) (string, error)
}

type RealPath struct{}

func NewPath() Path {
	return &RealPath{}
}

func (p *RealPath) Join(elem ...string) string {
	return path.Join(elem...)
}

func (p *RealPath) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (p *RealPath) Base(path string) string {
	return filepath.Base(path)
}

func (p *RealPath) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}
