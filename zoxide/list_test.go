package zoxide

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestListResults(t *testing.T) {
	t.Run("ListResults", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := &RealZoxide{shell: mockShell}
		mockShell.EXPECT().ListCmd("zoxide", "query", "--list", "--score").Return([]string{
			"100.0 /Users/joshmedeski/Downloads",
			" 82.0 /Users/joshmedeski/c/dotfiles/.config/fish",
			" 73.5 /Users/joshmedeski/c/dotfiles/.config/tmux",
			" 56.0 /Users/joshmedeski/c/sesh/v2",
			" 51.5 /Users/joshmedeski/c/dotfiles/.config/sesh",
			" 48.0 /Users/joshmedeski/c/sesh/main",
		}, nil)
		expected := []*model.ZoxideResult{
			{Path: "/Users/joshmedeski/Downloads", Score: 100.0},
			{Path: "/Users/joshmedeski/c/dotfiles/.config/fish", Score: 82.0},
			{Path: "/Users/joshmedeski/c/dotfiles/.config/tmux", Score: 73.5},
			{Path: "/Users/joshmedeski/c/sesh/v2", Score: 56.0},
			{Path: "/Users/joshmedeski/c/dotfiles/.config/sesh", Score: 51.5},
			{Path: "/Users/joshmedeski/c/sesh/main", Score: 48.0},
		}
		actual, err := zoxide.ListResults()
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
