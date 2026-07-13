package zoxide

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestListResults(t *testing.T) {
	t.Run("parses zoxide scored output by default", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{})
		mockShell.EXPECT().
			PrepareCmd("zoxide query --list --score", map[string]string{}).
			Return([]string{"zoxide", "query", "--list", "--score"}, nil)
		mockShell.EXPECT().ListCmd("zoxide", "query", "--list", "--score").Return([]string{
			"100.0 /Users/joshmedeski/Downloads",
			" 82.0 /Users/joshmedeski/c/dotfiles/.config/fish",
			" 73.5 /Users/joshmedeski/c/dotfiles/.config/tmux",
			" 56.0 /Users/joshmedeski/c/sesh/v2",
			" 51.5 /Users/joshmedeski/c/dotfiles/.config/sesh",
			" 48.0 /Users/joshmedeski/c/sesh/main",
			"",
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

	t.Run("parses paths-only output from a custom backend with zero score", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{
			ListCommand: "fasd -d -l -R",
		})
		mockShell.EXPECT().
			PrepareCmd("fasd -d -l -R", map[string]string{}).
			Return([]string{"fasd", "-d", "-l", "-R"}, nil)
		mockShell.EXPECT().ListCmd("fasd", "-d", "-l", "-R").Return([]string{
			"/Users/joshmedeski/c/sesh/v2",
			"/Users/joshmedeski/Downloads",
			"",
		}, nil)
		expected := []*model.ZoxideResult{
			{Path: "/Users/joshmedeski/c/sesh/v2", Score: 0},
			{Path: "/Users/joshmedeski/Downloads", Score: 0},
		}
		actual, err := zoxide.ListResults()
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("auto-detects score per line for mixed output", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		zoxide := NewZoxide(mockShell, model.FrecencyConfig{})
		mockShell.EXPECT().
			PrepareCmd("zoxide query --list --score", map[string]string{}).
			Return([]string{"zoxide", "query", "--list", "--score"}, nil)
		mockShell.EXPECT().ListCmd("zoxide", "query", "--list", "--score").Return([]string{
			"42.0 /scored/path",
			"/unscored/path",
		}, nil)
		expected := []*model.ZoxideResult{
			{Path: "/scored/path", Score: 42.0},
			{Path: "/unscored/path", Score: 0},
		}
		actual, err := zoxide.ListResults()
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
