package namer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineBareWorktreePath(t *testing.T) {
	t.Run("should find the bare path for standard bare clone", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/sesh             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`
		barePath := determineBareWorktreePath(out)
		assert.Equal(t, "/Users/hansolo/code/project/sesh", barePath)
	})

	t.Run("should find the bare path for bare clone to .bare folder and trim .bare suffix", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/sesh/.bare             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`
		barePath := determineBareWorktreePath(out)
		assert.Equal(t, "/Users/hansolo/code/project/sesh", barePath)
	})

	t.Run("should return empty string when no bare repository exists", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
/Users/hansolo/code/project/sesh/feature     c1d2e3f45 [feature-branch]
`
		barePath := determineBareWorktreePath(out)
		assert.Equal(t, "", barePath)
	})

	t.Run("should handle empty output", func(t *testing.T) {
		out := ""
		barePath := determineBareWorktreePath(out)
		assert.Equal(t, "", barePath)
	})

	t.Run("should handle whitespace-only output", func(t *testing.T) {
		out := "   \n  \n   "
		barePath := determineBareWorktreePath(out)
		assert.Equal(t, "", barePath)
	})

	t.Run("should find bare when it's the only entry", func(t *testing.T) {
		out := "/Users/hansolo/code/project/repo.git (bare)"
		barePath := determineBareWorktreePath(out)
		assert.Equal(t, "/Users/hansolo/code/project/repo.git", barePath)
	})

	t.Run("should find the bare path for bare clone to .git folder and trim .git suffix", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/sesh/.git             (bare)
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`
		barePath := determineBareWorktreePath(out)
		assert.Equal(t, "/Users/hansolo/code/project/sesh", barePath)
	})

	t.Run("should handle multiple lines with bare on second line", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
/Users/hansolo/code/project/sesh             (bare)
`
		barePath := determineBareWorktreePath(out)
		assert.Equal(t, "/Users/hansolo/code/project/sesh", barePath)
	})

	t.Run("should ignore malformed lines that don't match expected format", func(t *testing.T) {
		out := `
/Users/hansolo/code/project/sesh             (bare)
some malformed line without proper format
/Users/hansolo/code/project/sesh/main        ba04ca494 [5.x]
`
		barePath := determineBareWorktreePath(out)
		assert.Equal(t, "/Users/hansolo/code/project/sesh", barePath)
	})
}
