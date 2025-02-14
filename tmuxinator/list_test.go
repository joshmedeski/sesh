package tmuxinator

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestListConfigs(t *testing.T) {
	t.Run("List Tmuxinator Configs", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		tmuxinator := &RealTmuxinator{shell: mockShell}
		mockShell.EXPECT().ListCmd("tmuxinator", "list", "-n").Return([]string{
			"tmuxinator projects:",
			"dotfiles",
			"sesh",
			"home",
			"",
		}, nil)
		expected := []*model.TmuxinatorConfig{
			{Name: "dotfiles"},
			{Name: "sesh"},
			{Name: "home"},
		}
		actual, err := tmuxinator.List()
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
