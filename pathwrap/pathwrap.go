package pathwrap

import "path"

type Path interface {
	Join(elem ...string) string
}

type RealPath struct{}

func NewPath() Path {
	return &RealPath{}
}

func (p *RealPath) Join(elem ...string) string {
	return path.Join(elem...)
}
