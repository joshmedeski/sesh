package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/formatter"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActiveWindowNameFormat(t *testing.T) {
	assert.Equal(t, "#{?window_active,#{window_name},}", activeWindowNameFormat)
}

func TestHasTmuxSessions(t *testing.T) {
	t.Run("true when tmux session is present", func(t *testing.T) {
		sessions := model.SeshSessions{
			OrderedIndex: []string{"config:notes", "tmux:work"},
			Directory: model.SeshSessionMap{
				"config:notes": {Src: "config", Name: "notes"},
				"tmux:work":    {Src: "tmux", Name: "work"},
			},
		}
		assert.True(t, hasTmuxSessions(sessions))
	})

	t.Run("false when no tmux session is present", func(t *testing.T) {
		sessions := model.SeshSessions{
			OrderedIndex: []string{"config:notes", "zoxide:work"},
			Directory: model.SeshSessionMap{
				"config:notes": {Src: "config", Name: "notes"},
				"zoxide:work":  {Src: "zoxide", Name: "work"},
			},
		}
		assert.False(t, hasTmuxSessions(sessions))
	})
}

func TestFirstActiveWindowNameBySession(t *testing.T) {
	t.Run("uses first returned name", func(t *testing.T) {
		got := firstActiveWindowNameBySession(map[string][]string{
			"work":     {"nvim", "shell"},
			"dotfiles": {"zsh"},
			"empty":    {},
		})
		assert.Equal(t, map[string]string{
			"work":     "nvim",
			"dotfiles": "zsh",
		}, got)
	})

	t.Run("nil for no window names", func(t *testing.T) {
		assert.Nil(t, firstActiveWindowNameBySession(nil))
		assert.Nil(t, firstActiveWindowNameBySession(map[string][]string{}))
	})
}

func TestFormatSessions(t *testing.T) {
	newLister := func(mockTmux *tmux.MockTmux) Lister {
		config := model.Config{}
		return NewLister(
			config,
			new(home.MockHome),
			mockTmux,
			new(zoxide.MockZoxide),
			new(tmuxinator.MockTmuxinator),
			formatter.NewFormatter(icon.NewIcon(config)),
		)
	}

	t.Run("formats a cloned directory", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		l := newLister(mockTmux)
		sessions := model.SeshSessions{
			OrderedIndex: []string{"tmux:work"},
			Directory: model.SeshSessionMap{
				"tmux:work": {Src: "tmux", Name: "work", Path: "/work"},
			},
		}
		template := "{source}:{name}"

		got, err := l.Format(sessions, ListOptions{
			Format:    template,
			FormatSet: true,
			Icons:     true,
			NoColor:   true,
		})
		require.NoError(t, err)
		assert.Equal(t, "tmux: work", got.Directory["tmux:work"].Name)
		assert.Equal(t, "work", sessions.Directory["tmux:work"].Name)
	})

	t.Run("fetches active window names only when requested", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		mockTmux.On("ListAllWindowNames", activeWindowNameFormat).Return(map[string][]string{
			"work": {"nvim"},
		}, nil).Once()
		l := newLister(mockTmux)
		sessions := model.SeshSessions{
			OrderedIndex: []string{"tmux:work"},
			Directory: model.SeshSessionMap{
				"tmux:work": {Src: "tmux", Name: "work"},
			},
		}

		got, err := l.Format(sessions, ListOptions{
			Format:    "{active_window_name_prefix}{name}",
			FormatSet: true,
		})
		require.NoError(t, err)
		assert.Equal(t, "nvim work", got.Directory["tmux:work"].Name)
		mockTmux.AssertExpectations(t)
	})

	t.Run("does not fetch active window names for unrelated formats", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		l := newLister(mockTmux)
		sessions := model.SeshSessions{
			OrderedIndex: []string{"tmux:work"},
			Directory: model.SeshSessionMap{
				"tmux:work": {Src: "tmux", Name: "work"},
			},
		}

		_, err := l.Format(sessions, ListOptions{Format: "{name}", FormatSet: true})
		require.NoError(t, err)
		mockTmux.AssertNotCalled(t, "ListAllWindowNames")
	})
}
