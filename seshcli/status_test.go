package seshcli

import (
	"os"
	"testing"
	"time"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/statuscache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	t.Run("open issue gets a green OPEN badge and magenta issue label", func(t *testing.T) {
		got := formatStatus(github.Issue{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"})
		assert.Equal(t, "#[fg=green,bold]OPEN#[default] #[fg=magenta]Issue #400#[default] Dynamic tmux status bar", got)
	})

	t.Run("closed issue gets a red CLOSED badge and magenta issue label", func(t *testing.T) {
		got := formatStatus(github.Issue{Number: 400, Title: "Dynamic tmux status bar", State: "CLOSED"})
		assert.Equal(t, "#[fg=red,bold]CLOSED#[default] #[fg=magenta]Issue #400#[default] Dynamic tmux status bar", got)
	})

	t.Run("any non-OPEN state is treated as red", func(t *testing.T) {
		got := formatStatus(github.Issue{Number: 7, Title: "x", State: "MERGED"})
		assert.Equal(t, "#[fg=red,bold]MERGED#[default] #[fg=magenta]Issue #7#[default] x", got)
	})
}

func TestToIssue(t *testing.T) {
	got := toIssue(statuscache.Ref{Number: 400, Title: "x", State: "OPEN"})
	assert.Equal(t, github.Issue{Number: 400, Title: "x", State: "OPEN"}, got)
}

func TestComputeStatus(t *testing.T) {
	path := "/repo"
	ref := github.BranchRef{RepoRoot: "/repo", Branch: "400", Number: 400, HasNumber: true}
	key := statuscache.Key("/repo", "400")

	t.Run("fresh issue hit prints and does not spawn", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		sc.On("Read", key).Return(statuscache.Entry{
			Issue:     &statuscache.Ref{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"},
			Timestamp: time.Now(),
		}, true, nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		out, spawn := computeStatus(deps, 60, path)

		assert.Equal(t, "#[fg=green,bold]OPEN#[default] #[fg=magenta]Issue #400#[default] Dynamic tmux status bar", out)
		assert.False(t, spawn)
	})

	t.Run("fresh negative hit prints nothing and does not spawn", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		sc.On("Read", key).Return(statuscache.Entry{Timestamp: time.Now()}, true, nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		out, spawn := computeStatus(deps, 60, path)

		assert.Equal(t, "", out)
		assert.False(t, spawn)
	})

	t.Run("stale hit prints and spawns", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		sc.On("Read", key).Return(statuscache.Entry{
			Issue:     &statuscache.Ref{Number: 400, Title: "T", State: "OPEN"},
			Timestamp: time.Now().Add(-2 * time.Minute),
		}, true, nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		out, spawn := computeStatus(deps, 60, path)

		assert.NotEqual(t, "", out)
		assert.True(t, spawn)
	})

	t.Run("miss prints nothing and spawns", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		sc.On("Read", key).Return(statuscache.Entry{}, false, nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		out, spawn := computeStatus(deps, 60, path)

		assert.Equal(t, "", out)
		assert.True(t, spawn)
	})

	t.Run("ttl zero uses live Issue and never spawns", func(t *testing.T) {
		gh := new(github.MockGithub)
		gh.On("Issue", path).Return(github.Issue{Number: 400, Title: "T", State: "OPEN"}, true, nil)
		deps := &Deps{}
		deps.Github = gh

		out, spawn := computeStatus(deps, 0, path)

		assert.NotEqual(t, "", out)
		assert.False(t, spawn)
		gh.AssertNotCalled(t, "Resolve", path)
	})
}

func TestRunRefresh(t *testing.T) {
	path := "/repo"
	ref := github.BranchRef{RepoRoot: "/repo", Branch: "400", Number: 400, HasNumber: true}
	key := statuscache.Key("/repo", "400")

	t.Run("writes issue entry on success", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		gh.On("Issue", path).Return(github.Issue{Number: 400, Title: "T", State: "OPEN"}, true, nil)
		sc.On("Write", key, mock.MatchedBy(func(e statuscache.Entry) bool {
			return e.Issue != nil && e.Issue.Number == 400 && e.PR == nil
		})).Return(nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		assert.NoError(t, runRefresh(deps, path))
		sc.AssertExpectations(t)
	})

	t.Run("writes negative entry when no issue", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(github.BranchRef{RepoRoot: "/repo", Branch: "main"}, true)
		sc.On("Write", statuscache.Key("/repo", "main"), mock.MatchedBy(func(e statuscache.Entry) bool {
			return e.Issue == nil && e.PR == nil
		})).Return(nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		assert.NoError(t, runRefresh(deps, path))
		sc.AssertExpectations(t)
	})

	t.Run("not a repo writes nothing", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(github.BranchRef{}, false)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		assert.NoError(t, runRefresh(deps, path))
		sc.AssertNotCalled(t, "Write", mock.Anything, mock.Anything)
	})
}
