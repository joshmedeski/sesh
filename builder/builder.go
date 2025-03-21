package builder

import (
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

type Builder interface {
	Build(path model.SeshSession) (string, error)
}

type RealBuilder struct {
	os    oswrap.Os
	shell shell.Shell
}

func NewBuilder(
	os oswrap.Os,
	shell shell.Shell,
) Builder {
	return &RealBuilder{
		os,
		shell,
	}
}
