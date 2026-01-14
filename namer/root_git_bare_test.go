package namer

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func TestGitBareRootName(t *testing.T) {
	t.Run("should return repo/worktree name for .bare folder convention", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

		path := "/Users/hansolo/code/project/sesh/main"
		worktreeListOutput := `/Users/hansolo/code/project/sesh/.bare             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [main]
/Users/hansolo/code/project/sesh/feature     c1d2e3f45 [feature]
`
		mockGit.On("WorktreeList", path).Return(true, worktreeListOutput, nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/sesh/main", nil)
		mockPathwrap.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")

		name, err := gitBareRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "sesh/main", name)
	})

	t.Run("should return repo/worktree name for .git folder convention", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

		path := "/Users/hansolo/code/project/myrepo/feature-branch"
		// .git suffix is trimmed just like .bare suffix
		worktreeListOutput := `/Users/hansolo/code/project/myrepo/.git             (bare)
/Users/hansolo/code/project/myrepo/main           ba04ca494 [main]
/Users/hansolo/code/project/myrepo/feature-branch c1d2e3f45 [feature-branch]
`
		mockGit.On("WorktreeList", path).Return(true, worktreeListOutput, nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/myrepo/feature-branch", nil)
		mockPathwrap.On("Base", "/Users/hansolo/code/project/myrepo").Return("myrepo")

		name, err := gitBareRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "myrepo/feature-branch", name)
	})

	t.Run("should return repo/worktree for nested worktree path", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

		// Path is a subdirectory within a worktree
		path := "/Users/hansolo/code/project/sesh/main/src/components"
		worktreeListOutput := `/Users/hansolo/code/project/sesh/.bare             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [main]
`
		mockGit.On("WorktreeList", path).Return(true, worktreeListOutput, nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/sesh/main", nil)
		mockPathwrap.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")

		name, err := gitBareRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "sesh/main", name)
	})

	t.Run("should return empty string for non-bare git repo", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

		path := "/Users/hansolo/code/project/regular-repo"
		// No (bare) entry in worktree list - this is a regular git repo
		worktreeListOutput := `/Users/hansolo/code/project/regular-repo     c1d2e3f45 [main]
`
		mockGit.On("WorktreeList", path).Return(true, worktreeListOutput, nil)

		name, err := gitBareRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})

	t.Run("should return empty string for non-git directory", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

		path := "/Users/hansolo/code/not-a-git-repo"
		mockGit.On("WorktreeList", path).Return(false, "", nil)

		name, err := gitBareRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})

	t.Run("should return just repo name when ShowTopLevel fails", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}

		path := "/Users/hansolo/code/project/sesh/main"
		worktreeListOutput := `/Users/hansolo/code/project/sesh/.bare             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [main]
`
		mockGit.On("WorktreeList", path).Return(true, worktreeListOutput, nil)
		mockGit.On("ShowTopLevel", path).Return(false, "", nil)
		mockPathwrap.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")

		name, err := gitBareRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "sesh", name)
	})
}

func TestRootNameWithBareRepo(t *testing.T) {
	t.Run("RootName should use gitBareRootName for .bare convention", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := NewNamer(mockPathwrap, mockGit, mockHome, config)

		path := "/Users/hansolo/code/project/sesh/main"
		worktreeListOutput := `/Users/hansolo/code/project/sesh/.bare             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [main]
`
		mockPathwrap.On("EvalSymlinks", path).Return(path, nil)
		mockGit.On("WorktreeList", path).Return(true, worktreeListOutput, nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/sesh/main", nil)
		mockPathwrap.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")

		name, err := n.RootName(path)
		assert.NoError(t, err)
		assert.Equal(t, "sesh/main", name)
	})

	t.Run("RootName should use gitBareRootName for .git convention", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := NewNamer(mockPathwrap, mockGit, mockHome, config)

		path := "/Users/hansolo/code/myrepo/develop"
		worktreeListOutput := `/Users/hansolo/code/myrepo/.git             (bare)
/Users/hansolo/code/myrepo/main        ba04ca494 [main]
/Users/hansolo/code/myrepo/develop     c1d2e3f45 [develop]
`
		mockPathwrap.On("EvalSymlinks", path).Return(path, nil)
		mockGit.On("WorktreeList", path).Return(true, worktreeListOutput, nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/myrepo/develop", nil)
		mockPathwrap.On("Base", "/Users/hansolo/code/myrepo").Return("myrepo")

		name, err := n.RootName(path)
		assert.NoError(t, err)
		assert.Equal(t, "myrepo/develop", name)
	})

	t.Run("RootName should fall back to gitRootName for regular repos", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := NewNamer(mockPathwrap, mockGit, mockHome, config)

		path := "/Users/hansolo/code/regular-repo/src"
		// No bare entry - regular git repo
		worktreeListOutput := `/Users/hansolo/code/regular-repo     c1d2e3f45 [main]
`
		mockPathwrap.On("EvalSymlinks", path).Return(path, nil)
		mockGit.On("WorktreeList", path).Return(true, worktreeListOutput, nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/regular-repo", nil)
		mockPathwrap.On("Base", "/Users/hansolo/code/regular-repo").Return("regular-repo")

		name, err := n.RootName(path)
		assert.NoError(t, err)
		assert.Equal(t, "regular-repo", name)
	})
}
