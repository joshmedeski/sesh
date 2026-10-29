package namer

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/jj"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func newJjNamer(cfg model.Config) (*RealNamer, *jj.MockJj) {
	mockJj := new(jj.MockJj)
	return &RealNamer{pathwrap: pathwrap.NewPath(), jj: mockJj, config: cfg}, mockJj
}

func TestJjName(t *testing.T) {
	tests := []struct {
		name      string
		cfg       model.Config
		path      string
		workspace string
		repo      string
		want      string
		wantRoot  string
	}{
		{"repo root", model.Config{}, "/code/proj", "/code/proj", "/code/proj", "proj", "proj"},
		{"nested subdir", model.Config{}, "/code/proj/src/pkg", "/code/proj", "/code/proj", "proj/src/pkg", "proj"},
		{"workspace inside repo", model.Config{}, "/code/proj/ws/feat/src", "/code/proj/ws/feat", "/code/proj", "proj/ws/feat/src", "proj/ws/feat"},
		{"workspace root option", model.Config{GitNamerUseWorktreeRoot: true}, "/code/feat/src", "/code/feat", "/code/proj", "feat/src", "feat"},
		{"git_dir_length", model.Config{GitDirLength: 2}, "/code/proj/src", "/code/proj", "/code/proj", "code/proj/src", "code/proj"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, mj := newJjNamer(tt.cfg)
			mj.On("Roots", tt.path).Return(true, tt.workspace, tt.repo)

			name, err := jjName(n, tt.path)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, name)

			root, err := jjRootName(n, tt.path)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantRoot, root)
		})
	}

	t.Run("not a jj repo", func(t *testing.T) {
		n, mj := newJjNamer(model.Config{})
		mj.On("Roots", "/tmp/x").Return(false, "", "")
		name, _ := jjName(n, "/tmp/x")
		root, _ := jjRootName(n, "/tmp/x")
		assert.Empty(t, name)
		assert.Empty(t, root)
	})
}
