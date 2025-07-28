package namer

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func TestGitNamer(t *testing.T) {
	t.Run("should find git name when on top level", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := NewNamer(mockPathwrap, mockGit, mockHome, config)
		mockGit.On("ShowTopLevel", "/Users/hansolo/code/project/sesh").Return(true, "/Users/hansolo/code/project/sesh", nil)
		mockPathwrap.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")
		name, err := gitName(n.(*RealNamer), "/Users/hansolo/code/project/sesh")
		assert.NoError(t, err)
		assert.Equal(t, "sesh", name)
	})

	t.Run("should find git name when nested in git repo", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := NewNamer(mockPathwrap, mockGit, mockHome, config)
		mockGit.On("ShowTopLevel", "/Users/hansolo/code/project/sesh/namer").Return(true, "/Users/hansolo/code/project/sesh", nil)
		mockPathwrap.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")
		name, err := gitName(n.(*RealNamer), "/Users/hansolo/code/project/sesh/namer")
		assert.NoError(t, err)
		assert.Equal(t, "sesh/namer", name)
	})

	t.Run("should return empty string when not in a git repo", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := NewNamer(mockPathwrap, mockGit, mockHome, config)
		mockGit.On("ShowTopLevel", "/Users/hansolo/code/project/sesh").Return(false, "", nil)
		name, err := gitName(n.(*RealNamer), "/Users/hansolo/code/project/sesh")
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})
}
