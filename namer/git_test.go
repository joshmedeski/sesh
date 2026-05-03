package namer

import (
	"fmt"
	"testing"

	"github.com/Wingsdh/cc-sesh/v2/git"
	"github.com/Wingsdh/cc-sesh/v2/home"
	"github.com/Wingsdh/cc-sesh/v2/model"
	"github.com/Wingsdh/cc-sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func TestDetermineGitRootPath(t *testing.T) {
	t.Run("bare clone without suffix", func(t *testing.T) {
		out := `worktree /Users/hansolo/code/project/sesh
bare

worktree /Users/hansolo/code/project/sesh/main
HEAD ba04ca494
branch refs/heads/main
`
		assert.Equal(t, "/Users/hansolo/code/project/sesh", determineGitRootPath(out))
	})

	t.Run("bare clone with .bare suffix is trimmed", func(t *testing.T) {
		out := `worktree /Users/hansolo/code/project/sesh/.bare
bare

worktree /Users/hansolo/code/project/sesh/main
HEAD ba04ca494
branch refs/heads/main
`
		assert.Equal(t, "/Users/hansolo/code/project/sesh", determineGitRootPath(out))
	})

	t.Run("bare clone with .git suffix is trimmed", func(t *testing.T) {
		out := `worktree /Users/hansolo/code/project/sesh/.git
bare

worktree /Users/hansolo/code/project/sesh/main
HEAD ba04ca494
branch refs/heads/main
`
		assert.Equal(t, "/Users/hansolo/code/project/sesh", determineGitRootPath(out))
	})

	t.Run("regular clone uses first entry as-is", func(t *testing.T) {
		out := `worktree /Users/hansolo/code/project/nu
HEAD bb976dcdc
branch refs/heads/main

worktree /Users/hansolo/code/project/nu/.wk/5969
HEAD f31c5985c
branch refs/heads/jam/5969-something
`
		assert.Equal(t, "/Users/hansolo/code/project/nu", determineGitRootPath(out))
	})

	t.Run("standalone clone with single entry", func(t *testing.T) {
		out := `worktree /Users/hansolo/code/project/regular-repo
HEAD c1d2e3f45
branch refs/heads/main
`
		assert.Equal(t, "/Users/hansolo/code/project/regular-repo", determineGitRootPath(out))
	})

	t.Run("bare only, no .bare or .git suffix", func(t *testing.T) {
		out := `worktree /Users/hansolo/code/project/repo.git
bare
`
		assert.Equal(t, "/Users/hansolo/code/project/repo.git", determineGitRootPath(out))
	})

	t.Run("path containing spaces is preserved", func(t *testing.T) {
		out := `worktree /Users/alice/My Projects/cool repo
HEAD abcdef123
branch refs/heads/main

worktree /Users/alice/My Projects/cool repo/feature branch
HEAD 0123456789
branch refs/heads/feature
`
		assert.Equal(t, "/Users/alice/My Projects/cool repo", determineGitRootPath(out))
	})

	t.Run("bare clone with spaces and .bare suffix", func(t *testing.T) {
		out := `worktree /Users/alice/My Projects/cool repo/.bare
bare

worktree /Users/alice/My Projects/cool repo/main
HEAD abcdef123
branch refs/heads/main
`
		assert.Equal(t, "/Users/alice/My Projects/cool repo", determineGitRootPath(out))
	})

	t.Run("empty output", func(t *testing.T) {
		assert.Equal(t, "", determineGitRootPath(""))
	})

	t.Run("output without worktree lines", func(t *testing.T) {
		assert.Equal(t, "", determineGitRootPath("   \n  \n   "))
	})
}

func newNamer(t *testing.T) (*RealNamer, *pathwrap.MockPath, *git.MockGit) {
	t.Helper()
	mockPathwrap := new(pathwrap.MockPath)
	mockGit := new(git.MockGit)
	mockHome := new(home.MockHome)
	config := model.Config{DirLength: 1}
	return &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}, mockPathwrap, mockGit
}

func TestGitName(t *testing.T) {
	t.Run("regular clone at main tree root", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu"
		list := `worktree /Users/hansolo/code/project/nu
HEAD bb976dcdc
branch refs/heads/main
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu", name)
	})

	t.Run("regular clone, nested subdir", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/server"
		list := `worktree /Users/hansolo/code/project/nu
HEAD bb976dcdc
branch refs/heads/main
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu/server", name)
	})

	t.Run("regular clone, linked worktree", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/.wk/5969"
		list := `worktree /Users/hansolo/code/project/nu
HEAD bb976dcdc
branch refs/heads/main

worktree /Users/hansolo/code/project/nu/.wk/5969
HEAD f31c5985c
branch refs/heads/jam/5969-something
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu/.wk/5969", name)
	})

	t.Run("bare repo, .bare suffix", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/sesh/main"
		list := `worktree /Users/hansolo/code/project/sesh/.bare
bare

worktree /Users/hansolo/code/project/sesh/main
HEAD ba04ca494
branch refs/heads/main
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mp.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "sesh/main", name)
	})

	t.Run("bare repo, no suffix", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/sesh/main"
		list := `worktree /Users/hansolo/code/project/sesh
bare

worktree /Users/hansolo/code/project/sesh/main
HEAD ba04ca494
branch refs/heads/main
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mp.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "sesh/main", name)
	})

	t.Run("path with spaces", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/alice/My Projects/cool repo/src"
		list := `worktree /Users/alice/My Projects/cool repo
HEAD abcdef123
branch refs/heads/main
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mp.On("Base", "/Users/alice/My Projects/cool repo").Return("cool repo")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "cool repo/src", name)
	})

	t.Run("non-git directory, WorktreeList returns false", func(t *testing.T) {
		n, _, mg := newNamer(t)
		path := "/Users/hansolo/.config/nvim"
		mg.On("WorktreeList", path).Return(false, "", nil)

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})

	t.Run("non-git directory, WorktreeList errors", func(t *testing.T) {
		n, _, mg := newNamer(t)
		path := "/Users/hansolo/.config/nvim"
		mg.On("WorktreeList", path).Return(false, "", fmt.Errorf("not a git repository"))

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})
}

func TestGitRootName(t *testing.T) {
	t.Run("regular clone, nested subdir collapses to repo root", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/server/subdir"
		list := `worktree /Users/hansolo/code/project/nu
HEAD bb976dcdc
branch refs/heads/main
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/nu", nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu", name)
	})

	t.Run("regular clone, nested subdir inside linked worktree", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/.wk/5969/src"
		list := `worktree /Users/hansolo/code/project/nu
HEAD bb976dcdc
branch refs/heads/main

worktree /Users/hansolo/code/project/nu/.wk/5969
HEAD f31c5985c
branch refs/heads/jam/5969
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/nu/.wk/5969", nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu/.wk/5969", name)
	})

	t.Run("bare repo, nested subdir in worktree", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/sesh/main/namer"
		list := `worktree /Users/hansolo/code/project/sesh/.bare
bare

worktree /Users/hansolo/code/project/sesh/main
HEAD ba04ca494
branch refs/heads/main
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/sesh/main", nil)
		mp.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "sesh/main", name)
	})

	t.Run("bare repo .git suffix, feature worktree", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/myrepo/develop"
		list := `worktree /Users/hansolo/code/myrepo/.git
bare

worktree /Users/hansolo/code/myrepo/main
HEAD ba04ca494
branch refs/heads/main

worktree /Users/hansolo/code/myrepo/develop
HEAD c1d2e3f45
branch refs/heads/develop
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/myrepo/develop", nil)
		mp.On("Base", "/Users/hansolo/code/myrepo").Return("myrepo")

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "myrepo/develop", name)
	})

	t.Run("returns just repo name when ShowTopLevel returns empty", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/sesh/main"
		list := `worktree /Users/hansolo/code/project/sesh/.bare
bare

worktree /Users/hansolo/code/project/sesh/main
HEAD ba04ca494
branch refs/heads/main
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mg.On("ShowTopLevel", path).Return(false, "", nil)
		mp.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "sesh", name)
	})

	t.Run("non-git directory", func(t *testing.T) {
		n, _, mg := newNamer(t)
		path := "/Users/hansolo/.config/nvim"
		mg.On("WorktreeList", path).Return(false, "", nil)

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})
}
