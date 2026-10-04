package pathwrap

import (
	"path/filepath"
)

type Path interface {
	Join(elem ...string) string
	Abs(path string) (string, error)
	Base(path string) string
	Rel(base string, target string) (string, error)
	EvalSymlinks(path string) (string, error)
	FromSlash(path string) string
	ToSlash(path string) string
}

type RealPath struct{}

func NewPath() Path {
	return &RealPath{}
}

func (p *RealPath) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (p *RealPath) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (p *RealPath) Base(path string) string {
	return filepath.Base(path)
}

func (p *RealPath) Rel(base string, target string) (string, error) {
	return filepath.Rel(base, target)
}

func (p *RealPath) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

func (p *RealPath) FromSlash(path string) string {
	return filepath.FromSlash(path)
}

func (p *RealPath) ToSlash(path string) string {
	return filepath.ToSlash(path)
}
