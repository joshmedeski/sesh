package git

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestCurrentBranch(t *testing.T) {
	t.Run("returns the branch name on success", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		g := NewGit(mockShell)
		path := "/Users/josh/c/sesh"
		mockShell.On("Cmd", "git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD").
			Return("400", nil)

		ok, branch, err := g.CurrentBranch(path)

		assert.True(t, ok)
		assert.Equal(t, "400", branch)
		assert.NoError(t, err)
	})

	t.Run("returns false when not a git repo", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		g := NewGit(mockShell)
		path := "/tmp/not-a-repo"
		mockShell.On("Cmd", "git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD").
			Return("", fmt.Errorf("fatal: not a git repository"))

		ok, branch, err := g.CurrentBranch(path)

		assert.False(t, ok)
		assert.Equal(t, "", branch)
		assert.Error(t, err)
	})
}
