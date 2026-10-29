package dir

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/jj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRootDir(t *testing.T) {
	notGit := errors.New("not a git repository")

	t.Run("uses the git top level", func(t *testing.T) {
		mockGit := new(git.MockGit)
		mockJj := new(jj.MockJj)
		d := NewDir(nil, mockGit, mockJj, nil)
		mockGit.On("GitCommonDir", "/code/proj/src").Return(true, "/code/proj/.git", nil)
		mockGit.On("ShowTopLevel", "/code/proj/src").Return(true, "/code/proj", nil)

		ok, root := d.RootDir("/code/proj/src")

		assert.True(t, ok)
		assert.Equal(t, "/code/proj", root)
		mockJj.AssertNotCalled(t, "Roots", mock.Anything)
	})

	t.Run("falls back to the jj workspace root", func(t *testing.T) {
		mockGit := new(git.MockGit)
		mockJj := new(jj.MockJj)
		d := NewDir(nil, mockGit, mockJj, nil)
		mockGit.On("GitCommonDir", "/code/proj/w/feat/src").Return(false, "", notGit)
		mockGit.On("ShowTopLevel", "/code/proj/w/feat/src").Return(false, "", notGit)
		mockJj.On("Roots", "/code/proj/w/feat/src").Return(true, "/code/proj/w/feat", "/code/proj")

		ok, root := d.RootDir("/code/proj/w/feat/src")

		assert.True(t, ok)
		assert.Equal(t, "/code/proj/w/feat", root)
	})

	t.Run("no root outside git and jj", func(t *testing.T) {
		mockGit := new(git.MockGit)
		mockJj := new(jj.MockJj)
		d := NewDir(nil, mockGit, mockJj, nil)
		mockGit.On("GitCommonDir", "/tmp/x").Return(false, "", notGit)
		mockGit.On("ShowTopLevel", "/tmp/x").Return(false, "", notGit)
		mockJj.On("Roots", "/tmp/x").Return(false, "", "")

		ok, root := d.RootDir("/tmp/x")

		assert.False(t, ok)
		assert.Empty(t, root)
	})
}
