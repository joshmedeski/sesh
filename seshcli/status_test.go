package seshcli

import (
	"os"
	"testing"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestStatusPath(t *testing.T) {
	t.Run("returns the attached session path", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		session := model.SeshSession{Path: "/Users/josh/c/sesh"}
		mockLister.On("GetAttachedTmuxSession").Return(session, true)

		deps := &Deps{Lister: mockLister}
		got := statusPath(deps)

		assert.Equal(t, "/Users/josh/c/sesh", got)
	})

	t.Run("falls back to cwd when not attached", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockLister.On("GetAttachedTmuxSession").Return(model.SeshSession{}, false)

		deps := &Deps{Lister: mockLister}
		expectedCwd, _ := os.Getwd()
		got := statusPath(deps)

		assert.Equal(t, expectedCwd, got)
	})
}

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
