package seshcli

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/stretchr/testify/assert"
)

func TestFormatStatus(t *testing.T) {
	t.Run("open issue gets a green OPEN badge", func(t *testing.T) {
		got := formatStatus(github.Issue{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"})
		assert.Equal(t, "#[fg=green,bold]OPEN#[default] #400 Dynamic tmux status bar", got)
	})

	t.Run("closed issue gets a red CLOSED badge", func(t *testing.T) {
		got := formatStatus(github.Issue{Number: 400, Title: "Dynamic tmux status bar", State: "CLOSED"})
		assert.Equal(t, "#[fg=red,bold]CLOSED#[default] #400 Dynamic tmux status bar", got)
	})

	t.Run("any non-OPEN state is treated as red", func(t *testing.T) {
		got := formatStatus(github.Issue{Number: 7, Title: "x", State: "MERGED"})
		assert.Equal(t, "#[fg=red,bold]MERGED#[default] #7 x", got)
	})
}
