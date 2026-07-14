package seshcli

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
)

// tmuxRenamer embeds the generated tmux mock (so it satisfies the full
// tmux.Tmux interface) and overrides RenameSession to capture its arguments.
type tmuxRenamer struct {
	*tmux.MockTmux
	called     bool
	gotTarget  string
	gotNewName string
}

func newTmuxRenamer() *tmuxRenamer {
	return &tmuxRenamer{MockTmux: new(tmux.MockTmux)}
}

func (t *tmuxRenamer) RenameSession(target string, newName string) (string, error) {
	t.called = true
	t.gotTarget = target
	t.gotNewName = newName
	return "", nil
}

func TestEnrichedName(t *testing.T) {
	t.Run("appends sanitized issue title to base", func(t *testing.T) {
		mockNamer := new(namer.MockNamer)
		mockGithub := new(github.MockGithub)
		mockNamer.On("Name", "/p").Return("400-status", nil)
		mockGithub.On("Issue", "/p").
			Return(github.Issue{Number: 400, Title: "warm: the cache", State: "OPEN"}, true, nil)

		deps := &Deps{Namer: mockNamer}
		deps.Github = mockGithub

		assert.Equal(t, "400-status — warm the cache", enrichedName(deps, "/p"))
	})

	t.Run("returns bare base when no issue", func(t *testing.T) {
		mockNamer := new(namer.MockNamer)
		mockGithub := new(github.MockGithub)
		mockNamer.On("Name", "/p").Return("400-status", nil)
		mockGithub.On("Issue", "/p").Return(github.Issue{}, false, nil)

		deps := &Deps{Namer: mockNamer}
		deps.Github = mockGithub

		assert.Equal(t, "400-status", enrichedName(deps, "/p"))
	})

	t.Run("returns empty when namer fails", func(t *testing.T) {
		mockNamer := new(namer.MockNamer)
		mockNamer.On("Name", "/p").Return("", assert.AnError)

		deps := &Deps{Namer: mockNamer}

		assert.Equal(t, "", enrichedName(deps, "/p"))
	})
}

func TestRenameTarget(t *testing.T) {
	t.Run("uses the named session when an arg is given", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockLister.On("FindTmuxSession", "foo").
			Return(model.SeshSession{Name: "foo", Path: "/p/foo"}, true)

		deps := &Deps{Lister: mockLister}
		got, ok := renameTarget(deps, []string{"foo"})

		assert.True(t, ok)
		assert.Equal(t, "foo", got.Name)
	})

	t.Run("falls back to the attached session", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockLister.On("GetAttachedTmuxSession").
			Return(model.SeshSession{Name: "bar", Path: "/p/bar"}, true)

		deps := &Deps{Lister: mockLister}
		got, ok := renameTarget(deps, nil)

		assert.True(t, ok)
		assert.Equal(t, "bar", got.Name)
	})
}

func TestRunEnrich(t *testing.T) {
	t.Run("renames when the enriched name differs", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockNamer := new(namer.MockNamer)
		mockGithub := new(github.MockGithub)
		mockTmux := newTmuxRenamer()
		mockLister.On("GetAttachedTmuxSession").
			Return(model.SeshSession{Name: "400-status", Path: "/p"}, true)
		mockNamer.On("Name", "/p").Return("400-status", nil)
		mockGithub.On("Issue", "/p").
			Return(github.Issue{Number: 400, Title: "warm the cache", State: "OPEN"}, true, nil)

		deps := &Deps{Lister: mockLister, Namer: mockNamer}
		deps.Github = mockGithub
		deps.Tmux = mockTmux

		err := runEnrich(deps, nil)

		assert.NoError(t, err)
		assert.Equal(t, "400-status", mockTmux.gotTarget)
		assert.Equal(t, "400-status — warm the cache", mockTmux.gotNewName)
	})

	t.Run("does not rename when the name is unchanged (idempotent)", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockNamer := new(namer.MockNamer)
		mockGithub := new(github.MockGithub)
		mockTmux := newTmuxRenamer()
		mockLister.On("GetAttachedTmuxSession").
			Return(model.SeshSession{Name: "400-status — warm the cache", Path: "/p"}, true)
		mockNamer.On("Name", "/p").Return("400-status", nil)
		mockGithub.On("Issue", "/p").
			Return(github.Issue{Number: 400, Title: "warm the cache", State: "OPEN"}, true, nil)

		deps := &Deps{Lister: mockLister, Namer: mockNamer}
		deps.Github = mockGithub
		deps.Tmux = mockTmux

		err := runEnrich(deps, nil)

		assert.NoError(t, err)
		assert.False(t, mockTmux.called)
	})
}
