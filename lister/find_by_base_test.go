package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
)

func TestFindTmuxSessionByBase(t *testing.T) {
	newListerWith := func(names ...string) *RealLister {
		sessions := make([]*model.TmuxSession, 0, len(names))
		for _, n := range names {
			sessions = append(sessions, &model.TmuxSession{Name: n, Path: "/p/" + n})
		}
		mockTmux := new(tmux.MockTmux)
		mockTmux.On("ListSessions").Return(sessions, nil)
		return &RealLister{tmux: mockTmux}
	}

	t.Run("exact match wins", func(t *testing.T) {
		l := newListerWith("sesh", "other")
		got, ok := l.FindTmuxSessionByBase("sesh")
		assert.True(t, ok)
		assert.Equal(t, "sesh", got.Name)
	})

	t.Run("matches enriched prefix", func(t *testing.T) {
		l := newListerWith("400-status — warm the cache", "other")
		got, ok := l.FindTmuxSessionByBase("400-status")
		assert.True(t, ok)
		assert.Equal(t, "400-status — warm the cache", got.Name)
	})

	t.Run("prefers exact over enriched prefix", func(t *testing.T) {
		l := newListerWith("400-status — warm the cache", "400-status")
		got, ok := l.FindTmuxSessionByBase("400-status")
		assert.True(t, ok)
		assert.Equal(t, "400-status", got.Name)
	})

	t.Run("does not match a bare prefix without the separator", func(t *testing.T) {
		l := newListerWith("sesh-ui")
		_, ok := l.FindTmuxSessionByBase("sesh")
		assert.False(t, ok)
	})

	t.Run("matches base containing a slash", func(t *testing.T) {
		l := newListerWith("w/400 — warm the cache")
		got, ok := l.FindTmuxSessionByBase("w/400")
		assert.True(t, ok)
		assert.Equal(t, "w/400 — warm the cache", got.Name)
	})

	t.Run("no match returns false", func(t *testing.T) {
		l := newListerWith("alpha", "beta")
		_, ok := l.FindTmuxSessionByBase("gamma")
		assert.False(t, ok)
	})
}
