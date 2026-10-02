package pathwrap

import (
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
	return filepath.ToSlash(filepath.Join(elem...))
}

func (p *RealPath) Abs(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err == nil {
		absPath = filepath.ToSlash(absPath)
	}
	return absPath, err
}

func (p *RealPath) Base(path string) string {
	return filepath.ToSlash(filepath.Base(path))
}

func (p *RealPath) EvalSymlinks(path string) (string, error) {
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err == nil {
		resolvedPath = filepath.ToSlash(resolvedPath)
	}
	return resolvedPath, err
}
