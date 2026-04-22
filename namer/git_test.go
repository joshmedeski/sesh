package namer

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func TestDetermineGitRootPath(t *testing.T) {
	t.Run("bare clone without suffix", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/sesh             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`
		assert.Equal(t, "/Users/hansolo/code/project/sesh", determineGitRootPath(out))
	})

	t.Run("bare clone with .bare suffix is trimmed", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/sesh/.bare             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`
		assert.Equal(t, "/Users/hansolo/code/project/sesh", determineGitRootPath(out))
	})

	t.Run("bare clone with .git suffix is trimmed", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/sesh/.git             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`
		assert.Equal(t, "/Users/hansolo/code/project/sesh", determineGitRootPath(out))
	})

	t.Run("regular clone uses first entry as-is", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/nu            bb976dcdc [main]
/Users/hansolo/code/project/nu/.wk/5969   f31c5985c [jam/5969-something]
`
		assert.Equal(t, "/Users/hansolo/code/project/nu", determineGitRootPath(out))
	})

	t.Run("standalone clone with single entry", func(t *testing.T) {
		out := "/Users/hansolo/code/project/regular-repo     c1d2e3f45 [main]\n"
		assert.Equal(t, "/Users/hansolo/code/project/regular-repo", determineGitRootPath(out))
	})

	t.Run("bare only, no .bare or .git suffix", func(t *testing.T) {
		out := "/Users/hansolo/code/project/repo.git (bare)"
		assert.Equal(t, "/Users/hansolo/code/project/repo.git", determineGitRootPath(out))
	})

	t.Run("empty output", func(t *testing.T) {
		assert.Equal(t, "", determineGitRootPath(""))
	})

	t.Run("whitespace-only output", func(t *testing.T) {
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
		list := "/Users/hansolo/code/project/nu          bb976dcdc [main]\n"
		mg.On("WorktreeList", path).Return(true, list, nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu", name)
	})

	t.Run("regular clone, nested subdir", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/server"
		list := "/Users/hansolo/code/project/nu          bb976dcdc [main]\n"
		mg.On("WorktreeList", path).Return(true, list, nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu/server", name)
	})

	t.Run("regular clone, linked worktree", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/.wk/5969"
		list := `/Users/hansolo/code/project/nu              bb976dcdc [main]
/Users/hansolo/code/project/nu/.wk/5969     f31c5985c [jam/5969-something]
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
		list := `/Users/hansolo/code/project/sesh/.bare             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [main]
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
		list := `/Users/hansolo/code/project/sesh             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [main]
`
		mg.On("WorktreeList", path).Return(true, list, nil)
		mp.On("Base", "/Users/hansolo/code/project/sesh").Return("sesh")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "sesh/main", name)
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
		list := "/Users/hansolo/code/project/nu          bb976dcdc [main]\n"
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
		list := `/Users/hansolo/code/project/nu              bb976dcdc [main]
/Users/hansolo/code/project/nu/.wk/5969     f31c5985c [jam/5969]
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
		list := `/Users/hansolo/code/project/sesh/.bare             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [main]
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
		list := `/Users/hansolo/code/myrepo/.git             (bare)
/Users/hansolo/code/myrepo/main        ba04ca494 [main]
/Users/hansolo/code/myrepo/develop     c1d2e3f45 [develop]
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
		list := `/Users/hansolo/code/project/sesh/.bare             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [main]
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
