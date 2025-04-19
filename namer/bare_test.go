package namer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitRoot(t *testing.T) {
	t.Run("run should find worktree root", func(t *testing.T) {
		workTreeList := `
/Users/hansolo/code/project/sesh             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`
		bareRoot, err := parseBareFromWorkTreeList(workTreeList)
		assert.Nil(t, err)
		assert.Equal(t, "/Users/hansolo/code/project/sesh", bareRoot)
	})

	t.Run("run should find worktree root with .bare folder convention", func(t *testing.T) {
		workTreeList := `
/Users/hansolo/code/project/sesh/.bare        (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`
		bareRoot, err := parseBareFromWorkTreeList(workTreeList)
		assert.Nil(t, err)
		assert.Equal(t, "/Users/hansolo/code/project/sesh", bareRoot)
	})

	t.Run("run should find non-worktree root", func(t *testing.T) {
		workTreeList := `
/Users/hansolo/.dotfiles        ba04ca494 [5.x]
`
		root, err := parseBareFromWorkTreeList(workTreeList)
		assert.Nil(t, err)
		assert.Equal(t, "/Users/hansolo/.dotfiles", root)
	})
}
